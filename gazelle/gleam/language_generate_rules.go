package gleam

import (
	"log"
	"sort"
	"strings"

	lang "github.com/bazelbuild/bazel-gazelle/language"
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
}

func (gmb *gleamModuleBundle) imports() []string {
	imports := map[string]bool{}
	for _, module := range gmb.modules {
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
	imports := gmb.imports()
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

func (g *gleamLanguage) GenerateRules(args lang.GenerateArgs) lang.GenerateResult {

	gleamBundle := &gleamModuleBundle{modules: make(map[string]gleamModuleInfo)}
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
	ruleKind := "gleam_library"
	if gleamBundle.isBin {
		ruleKind = "gleam_binary"
	}

	r := rule.NewRule(ruleKind, name)
	if gleamBundle.isBin {
		r.SetAttr("visibility", []string{"//visibility:private"})
	}
	r.SetAttr("srcs", gleamBundle.sources())

	// We don't add/remove rules.
	return lang.GenerateResult{
		Imports: []any{
			gleamBundle.imports(),
		},
		Gen: []*rule.Rule{
			r,
		},
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

	moduleParents := filepath.SplitList(rel)
	moduleName := strings.TrimSuffix(file, gleamExt)
	return &gleamModuleInfo{imports: imports, moduleParents: moduleParents, moduleName: moduleName, hasMainFn: hasMainFunction, file: file}, nil
}
