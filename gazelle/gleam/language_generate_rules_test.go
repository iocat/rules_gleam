package gleam

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/bazel-gazelle/walk"
	bzl "github.com/bazelbuild/buildtools/build"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateRules(t *testing.T) {
	testDir := "gentestdata"
	c, langs, cexts := testConfig(t, "-build_file_name=BUILD.old", "-repo_root="+testDir)

	configFile, err := rule.LoadData(filepath.FromSlash("BUILD.config"), "config", []byte(`
# gazelle:follow **
`))
	if err != nil {
		t.Fatal(err)
	}
	for _, cext := range cexts {
		cext.Configure(c, "", configFile)
	}

	var loads []rule.LoadInfo
	for _, lang := range langs {
		loads = append(loads, lang.(language.ModuleAwareLanguage).ApparentLoads(func(string) string { return "" })...)
	}
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
				return
			}
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
}

func convertImportsAttrs(f *rule.File) {
	for _, r := range f.Rules {
		v := r.PrivateAttr(config.GazelleImportsKey)
		if v != nil {
			r.SetAttr(config.GazelleImportsKey, v)
		}
	}
}
