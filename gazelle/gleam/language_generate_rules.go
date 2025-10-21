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
	"github.com/kr/pretty"

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
	ruleKindTest   ruleKind = "gleam_test"
	ruleKindErlLib ruleKind = "gleam_erl_library"
)

type gleamModuleBundle struct {
	// default to the package name, but if root will be set to "gleam_lib" or "gleam_bin"
	name string
	kind ruleKind
	// Maps to fully qualified module names.
	modules map[string]gleamModuleInfo
	// mainModuleName string

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
		// We do not enforce strict "internal" module fora  hex repository.
		// In Gleam, internal modules are recommendations only.
		if GetGleamConfig(gmb.c).externalRepo {
			return []string{"//visibility:public"}
		}
		parent := strings.TrimSuffix(gmb.rel[:relIndex], "/")
		visibility = append(visibility, fmt.Sprintf("//%s:__subpackages__", parent))
	} else if len(configVisibility) == 0 {
		return []string{"//visibility:public"}
	} else {
		return configVisibility
	}

	return visibility
}

func (gmb *gleamModuleBundle) isInternalModule() bool {
	for _, module := range gmb.modules {
		if module.moduleName == "internal" {
			return true
		}
	}
	return false
}

func (gmb *gleamModuleBundle) setSourcePrefix(r *rule.Rule) {
	if GetGleamConfig(gmb.c).externalRepo {
		externalIndex := strings.LastIndex(gmb.c.RepoRoot, "external/")
		if externalIndex < 0 {
			return
		}
		externalIndex += len("external/")
		r.SetAttr("strip_src_prefix", fmt.Sprintf("%s%s", "external/", gmb.c.RepoRoot[externalIndex:]))
	}
}

func (gmb *gleamModuleBundle) generateRules() []*rule.Rule {
	isInternalModule := gmb.isInternalModule()
	rules := []*rule.Rule{}
	r := rule.NewRule(string(gmb.kind), gmb.name)
	r.SetAttr("srcs", filter(gmb.sources(), func(m string) bool {
		return true
	}))

	gmb.setSourcePrefix(r)
	switch gmb.kind {
	case ruleKindBin:
		r.SetAttr("visibility", []string{"//visibility:private"})
		// If bin has multiple sources, add a main_module to identify
		// which module declare the main function.
		mainModule := ""
		for _, m := range gmb.modules {
			if m.hasMainFn {
				mainModule = m.moduleName
				break
			}
		}
		if mainModule != "" && len(r.AttrStrings("srcs")) > 1 {
			r.SetAttr("main_module", mainModule)
		}
	case ruleKindLib, ruleKindErlLib:
		if isInternalModule {
			// r.SetAttr("visibility", gmb.nonInternalVisibility())
			r.SetAttr("visibility", []string{fmt.Sprintf("//%s:__subpackages__", gmb.rel)})
		} else {
			r.SetAttr("visibility", gmb.nonInternalVisibility())
		}
	}

	if len(r.AttrStrings("srcs")) != 0 {
		rules = append(rules, r)
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
	imports := []any{gmb.imports(func(m string) bool { return true })}
	return imports
}

func (g *gleamLanguage) GenerateRules(args lang.GenerateArgs) lang.GenerateResult {
	bundles := []*gleamModuleBundle{}
	// Do not generate for build directory at root as this is gleam add output directory.
	if args.Rel != "." && (strings.HasPrefix(filepath.Dir(args.Rel), "build/") || //
		filepath.Dir(args.Rel) == "build") {
		return lang.GenerateResult{}
	}
	name := path.Base(args.Rel)
	if len(name) == 0 || name == "." {
		name = "gleam_lib"
	}

	var gleamBundle, gleamTestBundle *gleamModuleBundle
	var ffiBundles []*gleamModuleBundle
	// For each of the Gleam file in the directory. Create a
	for _, file := range args.RegularFiles {
		ext := path.Ext(file)
		filename := file
		if strings.HasSuffix(filename, "_test.gleam") || strings.HasSuffix(filename, "_tests.gleam") {
			ext = gleamTestExt
		}

		switch ext {
		case gleamTestExt:
			if gleamTestBundle == nil {
				gleamTestBundle = &gleamModuleBundle{kind: ruleKindTest, name: fmt.Sprintf("%s_test", name), modules: make(map[string]gleamModuleInfo), c: args.Config, rel: args.Rel}
			}
			module, err := getGleamModuleInfo(args.Dir, file, args.Rel)
			if err != nil {
				log.Print(err)
				return lang.GenerateResult{}
			}
			gleamTestBundle.modules[module.moduleName] = *module
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
				gleamBundle.kind = ruleKindBin
			}
		case erlExt:
			nonNsModule := strings.TrimSuffix(filepath.Base(file), erlExt)
			// Precompiled Gleam module with "@" means these files are part of Gleam compilation
			// with the provided namespace separator, so we exclude them.
			if strings.Contains(nonNsModule, "@") {
				continue
			}
			ffiBundle := &gleamModuleBundle{
				kind:    ruleKindErlLib,
				name:    fmt.Sprintf("%s_ffi", nonNsModule),
				modules: make(map[string]gleamModuleInfo),
				c:       args.Config,
				rel:     args.Rel,
			}
			ffiBundle.modules[nonNsModule] = gleamModuleInfo{
				moduleName:    nonNsModule,
				file:          file,
				moduleParents: []string{},
				imports:       []string{}, // No imports for Erl FFI
			}
			ffiBundles = append(ffiBundles, ffiBundle)
		}
	}

	if gleamBundle != nil {
		if gleamBundle.kind == ruleKindLib {
			for modName, moduleInfo := range gleamBundle.modules {
				libBundle := &gleamModuleBundle{
					kind:    ruleKindLib,
					name:    modName,
					modules: map[string]gleamModuleInfo{modName: moduleInfo},
					rel:     gleamBundle.rel,
					c:       gleamBundle.c,
				}
				bundles = append(bundles, libBundle)
			}
		} else { // ruleKindBin
			for modName, moduleInfo := range gleamBundle.modules {
				if moduleInfo.hasMainFn {
					modules := map[string]gleamModuleInfo{}
					modules[modName] = moduleInfo
					binBundle := &gleamModuleBundle{
						kind:    ruleKindBin,
						name:    modName,
						modules: modules,
						rel:     gleamBundle.rel,
						c:       gleamBundle.c,
					}
					bundles = append(bundles, binBundle)
				}else  {
					libBundle := &gleamModuleBundle{
						kind:    ruleKindLib,
						name:    modName,
						modules: map[string]gleamModuleInfo{modName: moduleInfo},
						rel:     gleamBundle.rel,
						c:       gleamBundle.c,
					}
					bundles = append(bundles, libBundle)
				}
			}
		}
	}
	if len(ffiBundles) > 0 {
		bundles = append(bundles, ffiBundles...)
	}
	if gleamTestBundle != nil {
		bundles = append(bundles, gleamTestBundle)
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
	sortFunc := func(i, j int) bool {
		r1 := rulesList[i]
		r2 := rulesList[i]
		return r1.Name() < r2.Name()
	}
	sort.Slice(importsList, sortFunc)
	sort.Slice(rulesList, sortFunc)
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
	if err != nil || parseTree == nil {
		log.Printf("failed to parse file %s: %v", filePath, err)
		return nil, nil
	}
	printTree := false
	if parseTree != nil {
		for _, stmt := range parseTree.(parser.SourceFile).Statements {
			switch s := stmt.(type) {
			case parser.Import:
				if s.Module == "gleam/erlang/actor" {
					fmt.Println(s, dir, file, rel)
					printTree = true
				}
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
	if printTree {
		pretty.Println(parseTree)
	}

	moduleParents := filepath.SplitList(rel)
	moduleName := strings.TrimSuffix(file, gleamExt)
	return &gleamModuleInfo{imports: collect(imports), moduleParents: moduleParents, moduleName: moduleName, hasMainFn: hasMainFunction, file: file}, nil
}
