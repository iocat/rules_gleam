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

type gleamModuleBundle struct {
	// default to the package name, but if root will be set to "lib" or "bin"
	name  string
	isBin bool
	// Maps to fully qualified module names.
	modules map[string]gleamModuleInfo

	rel string
	c   *config.Config
}

func (gmb *gleamModuleBundle) imports(filterModule func(string) bool) []string {
	imports := map[string]bool{}
	for _, module := range gmb.modules {
		if !filterModule(module.moduleName) {
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
	imports := gmb.imports(func (m string) bool {
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
	ruleKind := "gleam_library"
	if gmb.isBin {
		ruleKind = "gleam_binary"
	}

	internalModule := gmb.internalModule()
	rules := []*rule.Rule{}
	r := rule.NewRule(ruleKind, gmb.name)
	if gmb.isBin {
		r.SetAttr("visibility", []string{"//visibility:private"})
	} else {
		r.SetAttr("visibility", gmb.nonInternalVisibility())
	}
	r.SetAttr("srcs", filter(gmb.sources(), func(m string) bool {
		return internalModule == nil || m != internalModule.file
	}))
	rules = append(rules, r)

	if internalModule != nil {
		internalR := rule.NewRule("gleam_library", gmb.name+"_internal")
		internalR.SetAttr("srcs", []string{internalModule.file})
		internalR.SetAttr("visibility", []string{fmt.Sprintf("//%s:__subpackages__", gmb.rel)})
		rules = append(rules, internalR)
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
	gleamBundle := &gleamModuleBundle{modules: make(map[string]gleamModuleInfo), c: args.Config, rel: args.Rel}
	name := path.Base(args.Rel)
	if len(name) == 0 {
		name = "lib"
	}
	isBin := false
	// For each of the Gleam file in the directory. Create a
	for _, file := range args.RegularFiles {
		if path.Ext(file) != gleamExt {
			continue
		}
		module, err := getGleamModuleInfo(args.Dir, file, args.Rel)
		if err != nil {
			log.Print(err)
			return lang.GenerateResult{}
		}
		gleamBundle.modules[module.moduleName] = *module
		if module.hasMainFn {
			name = "bin"
			isBin = true
		}
	}

	gleamBundle.name = name
	gleamBundle.isBin = isBin
	importsList := gleamBundle.generateImports()
	rulesList := gleamBundle.generateRules()
	if len(rulesList) != len(importsList) {
		panic("Rules and imports should be of the same size. Check implementation.")
	}
	// We don't add/remove rules.
	return lang.GenerateResult{
		Imports:     gleamBundle.generateImports(),
		Gen:         gleamBundle.generateRules(),
		Empty:       []*rule.Rule{},
		RelsToIndex: gleamBundle.relsToIndex(),
	}
}

func getGleamModuleInfo(dir, file string, rel string) (*gleamModuleInfo, error) {
	filePath := path.Clean(path.Join(dir, file))

	imports := []string{}
	hasMainFunction := false
	parseTree, err := parser.ParseFile(filePath)
	if err != nil {
		log.Printf("failed to parse file %s: %v", filePath, err)
		return nil, nil
	}
	if parseTree != nil {
		for _, stmt := range parseTree.(parser.SourceFile).Statements {
			switch s := stmt.(type) {
			case parser.Import:
				imports = append(imports, s.Module)
			case parser.Function:
				if s.Name == "main" && s.Public && len(s.Parameters) == 0 {
					hasMainFunction = true
				}
			}
		}
	}

	moduleParents := filepath.SplitList(rel)
	moduleName := strings.TrimSuffix(file, gleamExt)
	return &gleamModuleInfo{imports: imports, moduleParents: moduleParents, moduleName: moduleName, hasMainFn: hasMainFunction, file: file}, nil
}
