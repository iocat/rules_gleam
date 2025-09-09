package gleam

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	lang "github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/pathtools"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/iocat/rules_gleam/gazelle/gleam/parser"
	_ "github.com/kr/pretty"

	"path"
	"path/filepath"
)

// gleamModuleInfo contains information about a Gleam module.
type gleamModuleInfo struct {
	moduleParents []string
	moduleName    string
	file          string

	imports   []string
	hasMainFn bool
}

type ruleKind string

var (
	ruleKindLib    ruleKind = "gleam_library"
	ruleKindBin    ruleKind = "gleam_binary"
	ruleKindErlLib ruleKind = "gleam_erl_library"
)

type gleamModuleBundle struct {
	// default to the package name, but if root will be set to "lib" or "bin"
	name string
	kind ruleKind
	// Maps to fully qualified module names.
	modules        map[string]gleamModuleInfo
	mainModuleName string

	rel string
	c   *config.Config
}

func (gmb *gleamModuleBundle) imports(filterModule func(string) bool) []string {
	if gmb == nil {
		return []string{}
	}
	imports := make(map[string]bool)
	for _, module := range gmb.modules {
		if !filterModule(module.moduleName) {
			continue
		}
		if len(module.imports) == 0 {
			continue
		}
		for _, imp := range module.imports {
			imports[imp] = true
		}
	}

	var importStrings []string
	for imp := range imports {
		importStrings = append(importStrings, imp)
	}
	sort.Strings(importStrings)
	return importStrings
}

func (gmb *gleamModuleBundle) nonInternalImports() []string {
	return gmb.imports(func(m string) bool {
		return m != "internal"
	})
}

func (gmb *gleamModuleBundle) internalImports() []string {
	return gmb.imports(func(m string) bool {
		return m == "internal"
	})
}

func (gmb *gleamModuleBundle) sources() []string {
	files := []string{}
	for _, module := range gmb.modules {
		files = append(files, module.file)
	}
	return files
}

func (gmb *gleamModuleBundle) relsToIndex() []string {
	imports := gmb.imports(func(m string) bool {
		return true
	})
	rel := []string{}
	for _, imp := range imports {
		// stdlib does not need to be indexed.
		if strings.HasPrefix(imp, "gleam/") {
			continue
		}
		rel = append(rel, path.Base(imp))
	}
	return rel
}

// Gets the visiblity for the rule.
func (gmb *gleamModuleBundle) nonInternalVisibility() []string {
	// SEE https://github.com/bazel-contrib/bazel-gazelle/blob/master/language/go/generate.go#L844

	relIndex := pathtools.Index(gmb.rel, "internal")
	configVisibility := GetGleamConfig(gmb.c).gleamVisibility

	visibility := configVisibility
	// Currently processing an internal module (in a internal directory.)
	if relIndex >= 0 {
		parent := strings.TrimSuffix(gmb.rel[:relIndex], "/")
		visibility = append(visibility, fmt.Sprintf("//%s:__subpackages__", parent))
	} else if len(configVisibility) == 0 {
		return []string{"//visibility:public"}
	} else {
		return configVisibility
	}

	return visibility
}

func (gmb *gleamModuleBundle) internalModule() *gleamModuleInfo {
	for _, module := range gmb.modules {
		if module.moduleName == "internal" {
			return &module
		}
	}
	return nil
}

func (gmb *gleamModuleBundle) generateRules() []*rule.Rule {
	internalModule := gmb.internalModule()
	rules := []*rule.Rule{}
	r := rule.NewRule(string(gmb.kind), gmb.name)
	r.SetAttr("srcs", filter(gmb.sources(), func(m string) bool {
		return internalModule == nil || m != internalModule.file
	}))
	switch gmb.kind {
	case ruleKindBin:
		r.SetAttr("visibility", []string{"//visibility:private"})
		// If bin has multiple sources, add a main_module to identify 
		// which module declare the main function.
		if gmb.mainModuleName != "" && len(r.AttrStrings("srcs")) > 1 {
			r.SetAttr("main_module", gmb.mainModuleName)
		}
	case ruleKindLib, ruleKindErlLib:
		r.SetAttr("visibility", gmb.nonInternalVisibility())
	}

	rules = append(rules, r)
	if internalModule != nil {
		internalR := rule.NewRule("gleam_library", gmb.name+"_internal")
		internalR.SetAttr("srcs", []string{internalModule.file})
		internalR.SetAttr("visibility", []string{fmt.Sprintf("//%s:__subpackages__", gmb.rel)})

		rules = append(rules, internalR)
	}

	imports := gmb.generateImports()
	for i, r := range rules {
		// Like go implementation, we set this private useable for testing.
		// After merging phase, this attribute will be removed.
		r.SetPrivateAttr(config.GazelleImportsKey, imports[i].([]string))
	}
	return rules
}

/** Returns import, must be of the same size as generate rules returned. */
func (gmb *gleamModuleBundle) generateImports() []any {
	imports := []any{gmb.nonInternalImports()}
	im := gmb.internalModule()
	if im != nil {
		imports = append(imports, gmb.internalImports())
	}
	return imports
}

func (g *gleamLanguage) GenerateRules(args lang.GenerateArgs) lang.GenerateResult {
	bundles := []*gleamModuleBundle{}
	name := path.Base(args.Rel)
	if len(name) == 0 {
		name = "lib"
	}

	var gleamBundle, gleamFFIBundle *gleamModuleBundle
	// For each of the Gleam file in the directory. Create a
	for _, file := range args.RegularFiles {
		ext := path.Ext(file)

		switch ext {
		case gleamExt:
			if gleamBundle == nil {
				gleamBundle = &gleamModuleBundle{kind: ruleKindLib, name: name, modules: make(map[string]gleamModuleInfo), c: args.Config, rel: args.Rel}
			}
			module, err := getGleamModuleInfo(args.Dir, file, args.Rel)
			if err != nil {
				log.Print(err)
				return lang.GenerateResult{}
			}
			gleamBundle.modules[module.moduleName] = *module
			if module.hasMainFn {
				if len(name) == 0 {
					gleamBundle.name = "bin"
				}
				gleamBundle.mainModuleName = module.moduleName
				gleamBundle.kind = ruleKindBin
			}
		case erlExt:
			if gleamFFIBundle == nil {
				gleamFFIBundle = &gleamModuleBundle{kind: ruleKindErlLib, name: fmt.Sprintf("%s_ffi", name), modules: make(map[string]gleamModuleInfo), c: args.Config, rel: args.Rel}
			}

			nonNsModule := strings.TrimSuffix(filepath.Base(file), erlExt)
			if _, ok := gleamFFIBundle.modules[nonNsModule]; ok {
				log.Printf("Duplicate Erl FFI module name, please consider rename, Gleam FFI implementation doesn't have dir-like namespace: %s", nonNsModule)
				if len(gleamFFIBundle.modules) == 0 {
					gleamFFIBundle = nil
					continue
				}
			}
			gleamFFIBundle.modules[nonNsModule] = gleamModuleInfo{
				moduleName:    nonNsModule,
				file:          file,
				moduleParents: []string{},
				imports:       []string{}, // No imports for Erl FFI
			}
		}
	}

	if gleamBundle != nil {
		bundles = append(bundles, gleamBundle)
	}
	if gleamFFIBundle != nil {
		bundles = append(bundles, gleamFFIBundle)
	}

	importsList := []any{}
	rulesList := []*rule.Rule{}
	for _, bundle := range bundles {
		importsList = append(importsList, bundle.generateImports()...)
		rulesList = append(rulesList, bundle.generateRules()...)
	}
	if len(rulesList) != len(importsList) {
		panic("Rules and imports should be of the same size. Check implementation.")
	}
	// We don't add/remove rules.
	return lang.GenerateResult{
		Imports:     importsList,
		Gen:         rulesList,
		Empty:       []*rule.Rule{},
		RelsToIndex: gleamBundle.relsToIndex(),
	}
}

func getGleamModuleInfo(dir, file string, rel string) (*gleamModuleInfo, error) {
	filePath := path.Clean(path.Join(dir, file))

	imports := map[string]bool{}
	hasMainFunction := false
	parseTree, err := parser.ParseFile(filePath, parser.Debug(false))
	if err != nil {
		log.Printf("failed to parse file %s: %v", filePath, err)
		return nil, nil
	}
	if parseTree != nil {
		for _, stmt := range parseTree.(parser.SourceFile).Statements {
			switch s := stmt.(type) {
			case parser.Import:
				imports[s.Module] = true
			case parser.Function:
				if s.Name == "main" && s.Public && len(s.Parameters) == 0 {
					hasMainFunction = true
				}
				if len(s.ExternalAttributes) > 0 {
					erlImports := filter(s.ExternalAttributes, func(a parser.ExternalAttribute) bool { return a.TargetLang == "erlang" })
					if len(erlImports) > 0 {
						erlImport := erlImports[0]
						imports[fmt.Sprintf("erl:%s", erlImport.Module)] = true
					}
				}
			}
		}
	}

	moduleParents := filepath.SplitList(rel)
	moduleName := strings.TrimSuffix(file, gleamExt)
	return &gleamModuleInfo{imports: collect(imports), moduleParents: moduleParents, moduleName: moduleName, hasMainFn: hasMainFunction, file: file}, nil
}
