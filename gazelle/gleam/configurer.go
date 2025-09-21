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

	// Whether we're generates for an external Gleam (Hex) repository
	externalRepo bool
}

func (c *GleamConfig) clone() *GleamConfig {
	visibility := make([]string, len(c.gleamVisibility))
	copy(visibility, c.gleamVisibility)
	return &GleamConfig{
		gleamVisibility: visibility,
		externalRepo:    c.externalRepo,
	}
}

func (g *gleamLanguage) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	pc := &GleamConfig{}
	c.Exts[languageName] = pc

	fs.BoolVar(&pc.externalRepo, "gleam_external_repo", //
		false, "Whether we're setting up an external Gleam repository")
}

func (g *gleamLanguage) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	return nil
}

func (g *gleamLanguage) KnownDirectives() []string {
	return []string{
		"gleam_visibility",
	}
}

// Configure implements Language.Configure.
//
// It reads the "gleam_visibility" directive, which specifies the visibility
// of the target. Multiple values are allowed.
//
// This is called per directory, child directory inherits config from the parent's.
func (g *gleamLanguage) Configure(c *config.Config, rel string, f *rule.File) {
	var config *GleamConfig
	if c, ok := c.Exts[languageName]; !ok {
		config = &GleamConfig{}
	} else {
		config = c.(*GleamConfig).clone()
	}
	c.Exts[languageName] = config

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
