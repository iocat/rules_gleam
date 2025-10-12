package gleam

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/rules_go/go/runfiles"
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
	repos := make([]repo.Repo, len(c.repos))
	copy(repos, c.repos)
	copy(visibility, c.gleamVisibility)
	return &GleamConfig{
		gleamVisibility: visibility,
		externalRepo:    c.externalRepo,
		repos:           repos,
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
		if err := maybePopulateRemoteCacheFromGleamTomlForExternalRepo(c, &gc.repos); err != nil {
			return err
		}
	} else {
		if err := maybePopulateRemoteCacheFromBzlMod(c, &gc.repos); err != nil {
			return err
		}
	}

	return nil
}

type GleamToml struct {
	Name         string            `toml:"name"`
	Dependencies map[string]string `toml:"dependencies"`
}

func maybePopulateRemoteCacheFromGleamTomlForExternalRepo(c *config.Config, repos *[]repo.Repo) error {
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

		foundRepos, err := walkDirForRepos(c, externalPath, module)
		*repos = append(*repos, foundRepos...)
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
	configModuleName := c.ModuleToApparentName("gleam_hex_repositories_config")
	if configModuleName == "" {
		configModuleName = "gleam_hex_repositories_config"
	}
	rf, _ := runfiles.New()
	buildFile, err := rf.Rlocation(fmt.Sprintf("%s/BUILD", configModuleName))
	if err != nil {
		return err
	}
	content, err := os.ReadFile(buildFile)
	if err != nil {
		return err
	}
	buildFileParsed, err := build.ParseBuild(fmt.Sprintf("%s:BUILD", configModuleName), content)
	if err != nil {
		return err
	}
	configModuleDirName := filepath.Base(filepath.Dir(buildFile))

	for _, gleamRepo := range buildFileParsed.Rules("gleam_repository") {
		module := gleamRepo.AttrString("module_name")
		moduleDirName := strings.ReplaceAll(configModuleDirName, configModuleName, module)

		var mu sync.Mutex
		var wg sync.WaitGroup
		wg.Go(func() {
			parallelAppendRepos(c, rf, &mu, module, moduleDirName, repos)
		})
		wg.Wait()

	}
	return nil
}

func walkDirForRepos(c *config.Config, dir string, bazelModule string) (repos []repo.Repo, err error) {
	gleamModules := []string{}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".gleam") {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			modulePath := strings.TrimSuffix(relPath, ".gleam")
			modulePath = strings.ReplaceAll(modulePath, string(os.PathSeparator), "/")
			gleamModules = append(gleamModules, modulePath)
		} else if strings.HasSuffix(path, ".erl") {
			relPath := filepath.Base(path)
			importPath := fmt.Sprintf("%s%s", "erl:", strings.TrimSuffix(relPath, ".erl"))
			gleamModules = append(gleamModules, importPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, gleamModule := range gleamModules {
		repos = append(repos, repo.Repo{
			Name:     bazelModule,
			GoPrefix: gleamModule, // Using GoPrefix for now, will need to adjust for Gleam specific prefix if any
		})
	}
	return repos, nil

}

func parallelAppendRepos(c *config.Config, rf *runfiles.Runfiles, mu *sync.Mutex, module string, moduleDirName string, repos *[]repo.Repo) {
	moduleBuild, err := rf.Rlocation(fmt.Sprintf("%s/BUILD", moduleDirName))
	if err != nil {
		log.Printf("Could not find module directory for %s: %v", module, err)
		return
	}
	if moduleBuild == "" {
		log.Printf("Could not find module directory for %s", module)
		return
	}

	foundRepos, _ := walkDirForRepos(c, filepath.Dir(moduleBuild), module)
	mu.Lock()
	defer mu.Unlock()
	*repos = append(*repos, foundRepos...)
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
