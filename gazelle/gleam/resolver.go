package gleam

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	_ "github.com/kr/pretty"
)

var (
	gleamTestExt = "_test.gleam"
	gleamExt = ".gleam"
	erlExt   = ".erl"

	errSkipImport    errorType = "skip"
	errNotFound      errorType = "not found"
	errMultipleFound errorType = "multiple found"

	// https://www.erlang.org/doc/apps/stdlib/api-reference.html
	// Global Erlang interop modules
	erlangStdlibModules = map[string]bool{
		"argparse":           true,
		"array":              true,
		"base64":             true,
		"beam_lib":           true,
		"binary":             true,
		"c":                  true,
		"calendar":           true,
		"dets":               true,
		"dict":               true,
		"digraph":            true,
		"digraph_utils":      true,
		"edlin":              true,
		"edlin_expand":       true,
		"epp":                true,
		"erlang":             true,
		"erl_anno":           true,
		"erl_error":          true,
		"erl_eval":           true,
		"erl_expand_records": true,
		"erl_features":       true,
		"erl_id_trans":       true,
		"erl_internal":       true,
		"erl_lint":           true,
		"erl_parse":          true,
		"erl_pp":             true,
		"erl_scan":           true,
		"erl_tar":            true,
		"escript":            true,
		"ets":                true,
		"file_sorter":        true,
		"filelib":            true,
		"filename":           true,
		"gb_sets":            true,
		"gb_trees":           true,
		"gen_event":          true,
		"gen_fsm":            true,
		"gen_server":         true,
		"gen_statem":         true,
		"io":                 true,
		"io_lib":             true,
		"json":               true,
		"lists":              true,
		"log_mf_h":           true,
		"maps":               true,
		"math":               true,
		"ms_transform":       true,
		"orddict":            true,
		"ordsets":            true,
		"peer":               true,
		"pool":               true,
		"proc_lib":           true,
		"proplists":          true,
		"qlc":                true,
		"queue":              true,
		"rand":               true,
		"random":             true,
		"re":                 true,
		"sets":               true,
		"shell":              true,
		"shell_default":      true,
		"shell_docs":         true,
		"slave":              true,
		"sofs":               true,
		"string":             true,
		"supervisor":         true,
		"supervisor_bridge":  true,
		"sys":                true,
		"timer":              true,
		"unicode":            true,
		"uri_string":         true,
		"win32reg":           true,
		"zip":                true,
		"zstd":               true,
	}
)

type errorType string

type gleamGazelleError struct {
	msg       string
	errorType errorType
}

func (gge *gleamGazelleError) ErrorType() errorType {
	return gge.errorType
}

func (gge *gleamGazelleError) Error() string {
	return gge.msg
}

func (*gleamLanguage) Name() string { return "gleam" }

// Returns the Gleam specific import path for the rule.
// These are all of the Gleam modules, declared in srcs.
func (g *gleamLanguage) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	if !isGleamLibrary(r) {
		return nil
	}

	imports := []resolve.ImportSpec{}
	for _, src := range r.AttrStrings("srcs") {
		if path.Ext(src) == gleamExt {
			imports = append(imports, resolve.ImportSpec{Lang: g.Name(), Imp: path.Join(f.Pkg, strings.TrimSuffix(src, gleamExt))})
		} else if path.Ext(src) == erlExt {
			imports = append(imports, resolve.ImportSpec{Lang: g.Name(), Imp: strings.Join([]string{"erl", strings.TrimSuffix(src, erlExt)}, ":")})
		}
	}

	return imports
}

func (g *gleamLanguage) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

// Resolve adds deps to the given rule based on the importRaws.
func (g *gleamLanguage) Resolve(c *config.Config, ix *resolve.RuleIndex, _rc *repo.RemoteCache, r *rule.Rule, importRaws interface{}, from label.Label) {
	// If no imports are given, bail early.
	if importRaws == nil {
		return
	}

	var err error
	gleamConfig := GetGleamConfig(c)
	rc, cleanup := repo.NewRemoteCache(gleamConfig.repos)
	defer func() {
		if cerr := cleanup(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	imports := importRaws.([]string)
	r.DelAttr("deps")

	// Create a set of dependencies so we can avoid duplicates.
	depSet := make(map[string]bool)
	for _, imp := range imports {
		depLabel, err := g.resolveGleam(c, ix, rc, r, imp, from)
		if err != nil && err.ErrorType() == errSkipImport {
			// If resolveGleam returns errSkipImport, skip this import.
			continue
		} else if err != nil {
			// If resolveGleam has any other error, log it.
			log.Print(err.msg)
		} else {
			var label label.Label
			if depLabel.Pkg == from.Pkg && depLabel.Repo == from.Repo {
				label = depLabel.Rel(depLabel.Repo, depLabel.Pkg)
			} else {
				label = depLabel.Abs(depLabel.Repo, depLabel.Pkg)
			}
			depSet[label.String()] = true
		}
	}

	if len(depSet) != 0 {
		// If there are dependencies, set the deps attribute.
		deps := make([]string, 0, len(depSet))
		for dep := range depSet {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}

// For gleamlibrary rule that does self import modules in srcs, we don't need labels for these.
func isSelfImport(r *rule.Rule, f label.Label, imp string) bool {
	localImports := asSet(mapper(r.AttrStrings("srcs"), func(src string) string {
		return path.Join(f.Pkg, strings.TrimSuffix(src, gleamExt))
	}))
	return localImports[imp]
}

func (g *gleamLanguage) resolveGleam(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imp string, from label.Label) (label.Label, *gleamGazelleError) {
	if erlangStdlibModules[strings.TrimPrefix(imp, "erl:")] {
		return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("erlang stdlib module: %s", imp), errorType: errSkipImport}
	}
	if isSelfImport(r, from, imp) {
		return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("self import: %s", imp), errorType: errSkipImport}
	}
	results := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: g.Name(), Imp: imp}, g.Name())

	if len(results) == 0 {
		l, err := g.tryResolveExternalDeps(c, ix, rc, r, imp, from)
		if err != nil {
			return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("no rule may be imported with %q from package %s: %v", imp, from, err), errorType: errNotFound}
		}
		return l, nil
	} else if len(results) > 1 {
		return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("multiple rules (%s and %s) may be imported with %q from %s", results[0].Label, results[1].Label, imp, from), errorType: errMultipleFound}
	}
	return results[0].Label, nil
}

func (g *gleamLanguage) tryResolveExternalDeps(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	imp string,
	from label.Label,
) (label.Label, error) {
	pkg, module, err := rc.Root(imp)
	if err != nil {
		return label.NoLabel, err
	}
	depPkg := filepath.Dir(pkg)
	if depPkg == "." {
		depPkg = ""
	}
	depMod := ""
	if depPkg == "" {
		depMod = "gleam_lib"
	} else {
		depMod = filepath.Base(depPkg)
	}
	l := label.New(module, depPkg, depMod)
	return l, nil
}
