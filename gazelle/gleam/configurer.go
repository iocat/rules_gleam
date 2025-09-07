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

	// Note: the -proto flag does not set the ModeExplicit flag. We want to
	// be able to switch to DisableMode in vendor directories, even when
	// this is set for compatibility with older versions.
	// fs.Var(&modeFlag{&pc.Mode}, "proto", "default: generates a proto_library rule for one package\n\tpackage: generates a proto_library rule for for each package\n\tdisable: does not touch proto rules\n\tdisable_global: does not touch proto rules and does not use special cases for protos in dependency resolution")
	// fs.StringVar(&pc.groupOption, "proto_group", "", "option name used to group .proto files into proto_library rules")
	// fs.StringVar(&pc.ImportPrefix, "proto_import_prefix", "", "When set, .proto source files in the srcs attribute of the rule are accessible at their path with this prefix appended on.")
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
