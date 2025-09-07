package gleam

import (
	"errors"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

var (
	gleamExt = ".gleam"

	errNotFound = errors.New("rule not found")
	errSkipImport = errors.New("skip import")
)

func (*gleamLanguage) Name() string { return "gleam" }

// Returns the Gleam specific import path for the rule.
// These are all of the Gleam modules, declared in srcs.
func (g *gleamLanguage) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	if !isGleamLibrary(r) {
		return nil
	}
	srcs := r.AttrStrings("srcs")
	imports := make([]resolve.ImportSpec, len(srcs))
	for i, src := range r.AttrStrings("srcs") {
		if path.Ext(src) != ".gleam" {
			continue
		}
		imports[i] = resolve.ImportSpec{Lang: g.Name(), Imp: path.Join(f.Pkg, strings.TrimSuffix(src, gleamExt))}
	}
	return imports
}

func (g *gleamLanguage) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

// Resolve adds deps to the given rule based on the importRaws.
func (g *gleamLanguage) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, importRaws interface{}, from label.Label) {
	// If no imports are given, bail early.
	if importRaws == nil {
		return
	}

	imports := importRaws.([]string)
	r.DelAttr("deps")

	// Create a set of dependencies so we can avoid duplicates.
	depSet := make(map[string]bool)
	for _, imp := range imports {
		depLabel, err := g.resolveGleam(c, ix, rc, r, imp, from)
		if err == errSkipImport {
			// If resolveGleam returns errSkipImport, skip this import.
			continue
		}else if err != nil {
			// If resolveGleam has any other error, log it.
			log.Print(err)
		}else {
			var label label.Label
			if depLabel.Pkg == from.Pkg && depLabel.Repo == from.Repo {
				label = depLabel.Rel(depLabel.Repo, depLabel.Pkg)
			}else {
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

func (g *gleamLanguage) resolveGleam(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imp string, from label.Label) (label.Label, error) {

	if(isSelfImport(r, from, imp)) {
		return label.NoLabel, errSkipImport
	}
	results := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: g.Name(), Imp: imp}, g.Name())
	if len(results) == 0 {
		return label.NoLabel, errNotFound
	} else if len(results) > 1 {
		return label.NoLabel, fmt.Errorf("multiple rules (%s and %s) may be imported with %q from %s", results[0].Label, results[1].Label, imp, from)
	}
	return results[0].Label, nil
}
