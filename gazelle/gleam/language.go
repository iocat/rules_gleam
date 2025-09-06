package gleam

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/config"
	lang "github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const languageName = "gleam"
type gleamLanguage struct{}

var gleamKinds = map[string]rule.KindInfo{
	"gleam_library": {
		MatchAttrs:    []string{"srcs"},
		NonEmptyAttrs: map[string]bool{"srcs": true},
		MergeableAttrs: map[string]bool{
			"srcs":                true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
	"gleam_binary": {
		MatchAttrs:    []string{"srcs"},
		NonEmptyAttrs: map[string]bool{"srcs": true},
		MergeableAttrs: map[string]bool{
			"srcs":                true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
}

func (g *gleamLanguage) Kinds() map[string]rule.KindInfo {
	return gleamKinds
}


func (g *gleamLanguage) ApparentLoads(moduleToApparentName func(string) string) []rule.LoadInfo {
	rulesGleam := moduleToApparentName("rules_gleam")
	if rulesGleam == "" {
		rulesGleam = "rules_gleam"
	}
	return []rule.LoadInfo{
		{
			Name: fmt.Sprintf("@%s//gleam:defs.bzl", rulesGleam),
			Symbols: []string{
				"gleam_library",
				"gleam_binary",
			},
		},
	}
}

// Deprecated.
func (g *gleamLanguage) Loads() []rule.LoadInfo {
	panic("ApparentLoads should be called instead")
}

// We don't fix anything, (yet :).
func (g *gleamLanguage) Fix(c *config.Config, f *rule.File) {

}

func NewLanguage() lang.Language {
	return &gleamLanguage{}
}