package gleam

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type GleamConfig struct {
	// For directive gleam_visibility
	gleamVisibility []string

	// Whether we're generates for an external Gleam (Hex) repository
	externalRepo bool
	// Cache of remote repositories from gleam.toml
	// Reusing Go repo.Repo for now.
	repos []repo.Repo
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
	gc := GetGleamConfig(c).clone()
	c.Exts[languageName] = gc

	if gc.externalRepo {
		if err := maybePopulateRemoteCacheFromGleamToml(c, &gc.repos); err != nil {
			return err
		}
	}else {
		if err := maybePopulateRemoteCacheFromBzlMod(c, &gc.repos); err != nil {
			return err
		}
	}
	for _, r := range gc.repos {
		fmt.Println(r)
	}

	return nil
}

type GleamToml struct {
	Name         string            `toml:"name"`
	Dependencies map[string]string `toml:"dependencies"`
}

func maybePopulateRemoteCacheFromGleamToml(c *config.Config, repos *[]repo.Repo) error {
	haveGleam := false
	for name := range c.Exts {
		if name == "gleam" {
			haveGleam = true
			break
		}
	}
	if !haveGleam {
		return nil
	}

	gleamTomlPath := filepath.Join(c.RepoRoot, "gleam.toml")
	if _, err := os.Stat(gleamTomlPath); err != nil {
		return nil
	}

	data, err := os.ReadFile(gleamTomlPath)
	if err != nil {
		return err
	}

	var gleamToml GleamToml
	if err := toml.Unmarshal(data, &gleamToml); err != nil {
		return err
	}

	repoName := getRepoNameFromPath(c.RepoRoot)
	for gleamPackage := range gleamToml.Dependencies {
		module := fmt.Sprintf("hex_%s", gleamPackage)
		externalPath := strings.Replace(c.RepoRoot, repoName, module, 1)
		_, err := os.Stat(externalPath)
		if err != nil {
			return err
		}

		err = filepath.Walk(externalPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".gleam") {
				relPath, err := filepath.Rel(externalPath, path)
				if err != nil {
					return err
				}
				importPath := strings.TrimSuffix(relPath, ".gleam")
				importPath = strings.ReplaceAll(importPath, string(os.PathSeparator), "/")
				*repos = append(*repos, repo.Repo{
					Name:     module,
					GoPrefix: importPath, // Using GoPrefix for now, will need to adjust for Gleam specific prefix if any
				})
			} else if strings.HasSuffix(path, ".erl") {
				relPath := filepath.Base(path)
				importPath := fmt.Sprintf("%s%s", "erl:", strings.TrimSuffix(relPath, ".erl"))
				*repos = append(*repos, repo.Repo{
					Name:     module,
					GoPrefix: importPath, // Using GoPrefix for now, will need to adjust for Erlang specific prefix if any
				})
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Given: /blah/blah/external/rules_gleam++gleam+gleam_stdlib
// Extract gleam_stdlib
func getRepoNameFromPath(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	lastPart := parts[len(parts)-1]
	repoComponents := strings.Split(lastPart, "+")
	return repoComponents[len(repoComponents)-1]
}


func maybePopulateRemoteCacheFromBzlMod(c *config.Config, repos *[]repo.Repo) error {
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
