package gleam

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/bazelbuild/bazel-gazelle/walk"
	bzl "github.com/bazelbuild/buildtools/build"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateRules(t *testing.T) {
	testDir := "testdata"
	c, langs, cexts := testConfig(t, "-build_file_name=BUILD.old", "-repo_root="+testDir)

	f, err := rule.LoadData(filepath.FromSlash("BUILD.config"), "config", []byte(`
# gazelle:follow **
`))

	if err != nil {
		t.Fatal(err)
	}
	for _, cext := range cexts {
		cext.Configure(c, "", f)
	}

	var loads []rule.LoadInfo
	for _, lang := range langs {
		loads = append(loads, lang.(language.ModuleAwareLanguage).ApparentLoads(func(string) string { return "" })...)
	}
	var testsFound int
	walk.Walk(c, cexts, []string{testDir}, walk.VisitAllUpdateSubdirsMode, func(dir, rel string, c *config.Config, update bool, oldFile *rule.File, subdirs, regularFiles, genFiles []string) {
		t.Run(rel, func(t *testing.T) {
			var empty, gen []*rule.Rule
			for _, lang := range langs {
				res := lang.GenerateRules(language.GenerateArgs{
					Config:       c,
					Dir:          dir,
					Rel:          rel,
					File:         oldFile,
					Subdirs:      subdirs,
					RegularFiles: regularFiles,
					GenFiles:     genFiles,
					OtherEmpty:   empty,
					OtherGen:     gen,
				})
				empty = append(empty, res.Empty...)
				gen = append(gen, res.Gen...)
			}
			isTest := false
			for _, name := range regularFiles {
				if name == "BUILD.want" {
					isTest = true
					break
				}
			}
			if !isTest {
				// GenerateRules may have side effects, so we need to run it, even if
				// there's no test.
				return
			}
			testsFound += 1
			f := rule.EmptyFile("test", "")
			for _, r := range gen {
				r.Insert(f)
			}
			convertImportsAttrs(f)
			merger.FixLoads(f, loads)
			f.Sync()
			got := string(bzl.Format(f.File))
			wantPath := filepath.Join(dir, "BUILD.want")
			wantBytes, err := os.ReadFile(wantPath)
			if err != nil {
				t.Fatalf("error reading %s: %v", wantPath, err)
			}
			want := string(wantBytes)
			want = strings.ReplaceAll(want, "\r\n", "\n")
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got): %s", diff)
			}
		})
	})
	// Avoid spurious success if we fail to find any tests.
	if testsFound == 0 {
		t.Error("No rule generation tests were found")
	}
}

func convertImportsAttrs(f *rule.File) {
	for _, r := range f.Rules {
		v := r.PrivateAttr(config.GazelleImportsKey)
		if v != nil {
			r.SetAttr(config.GazelleImportsKey, v)
		}
	}
}

func testConfig(t *testing.T, args ...string) (*config.Config, []language.Language, []config.Configurer) {
	haveRoot := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "-repo_root") {
			haveRoot = true
			break
		}
	}
	if !haveRoot {
		args = append(args, "-repo_root=.")
	}

	cexts := []config.Configurer{
		&config.CommonConfigurer{},
		&walk.Configurer{},
		&resolve.Configurer{},
	}
	langs := []language.Language{NewLanguage()}
	c := testtools.NewTestConfig(t, cexts, langs, args)
	for _, lang := range langs {
		cexts = append(cexts, lang)
	}
	return c, langs, cexts
}
