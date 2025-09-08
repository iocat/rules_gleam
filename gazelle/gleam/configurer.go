package gleam

import (
	"flag"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type GleamConfig struct {
	// For directive gleam_visibility
	gleamVisibility []string
}

func (g *gleamLanguage) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	pc := &GleamConfig{}
	c.Exts[languageName] = pc
}

func (g *gleamLanguage) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

func (g *gleamLanguage) KnownDirectives() []string {
	return []string{
		"gleam_visibility",
	}
}

func (g *gleamLanguage) Configure(c *config.Config, rel string, f *rule.File) {
	config := GetGleamConfig(c)
	if f != nil {
		for _, d := range f.Directives {
			switch d.Key {
			case "gleam_visibility":
				config.gleamVisibility = append(config.gleamVisibility, strings.TrimSpace(d.Value))
			}
		}
	}
}

func GetGleamConfig(c *config.Config) *GleamConfig {
	return c.Exts[languageName].(*GleamConfig)
}
