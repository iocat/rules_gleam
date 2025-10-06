package gleam

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	bzl "github.com/bazelbuild/buildtools/build"

	"github.com/google/go-cmp/cmp"
	"github.com/lithammer/dedent"
)

type buildFile struct {
	pkg, content string
}

type resolveTestCase struct {
	desc      string
	index     []buildFile
	skipIndex bool
	old       buildFile
	want      string
}

type mapResolver map[string]resolve.Resolver

func (mr mapResolver) Resolver(r *rule.Rule, f string) resolve.Resolver {
	return mr[r.Kind()]
}

var testCases = []resolveTestCase{
	{
		desc:      "self import",
		index:     []buildFile{},
		skipIndex: true,
		old: buildFile{
			pkg: "foo",
			content: `
				gleam_library(
					name = "foo",
					srcs = [
						"foo.gleam",
						"bar.gleam"
					],
					_gazelle_imports = ["foo/foo"],
				)
`,
		},
		want: `
				gleam_library(
					name = "foo",
					srcs = [
						"bar.gleam",
						"foo.gleam",
					],
				)
`,
	},
	{
		desc: "import from the same package",
		old: buildFile{
			pkg: "foo",
			content: `
				gleam_library(
					name = "bar",
					srcs = [
						"bar.gleam"
					],
				)

				gleam_binary(
					name = "foo",
					srcs = [
						"foo.gleam",
					],
					_gazelle_imports = ["foo/bar"],
				)
`,
		},
		want: `
				gleam_library(
					name = "bar",
					srcs = [
						"bar.gleam",
					],
				)

				gleam_binary(
					name = "foo",
					srcs = [
						"foo.gleam",
					],
					deps = [":bar"],
				)
`,
	},
	{
		desc: "import from another package",
		index: []buildFile{
			{
				pkg: "foo/bar",
				content: `
					gleam_library(
						name = "bar",
						srcs = [
							"bar.gleam"
						],
					)
`,
			},
			{
				pkg: "test",
				content: `
					gleam_library(
						name = "test",
						srcs = [
							"a.gleam",
							"b.gleam",
						],
						deps = ["//foo/bar"],
					)
`,
			},
		},
		old: buildFile{
			pkg: "foo",
			content: `
					gleam_library(
						name = "bar",
						srcs = [
							"bar.gleam"
						],
						visibility = ["//visibility:public"],
						_gazelle_imports = ["foo/bar/bar", "test/a", "test/b"],
					)

					gleam_library(
						name = "foo",
						srcs = [
							"foo.gleam"
						],
						visibility = ["//visibility:private"],
						_gazelle_imports = ["test/a", "test/b"],
					)
	`,
		},
		want: `
					gleam_library(
						name = "bar",
						srcs = [
							"bar.gleam",
						],
						visibility = ["//visibility:public"],
						deps = [
							"//foo/bar",
							"//test",
						],
					)

					gleam_library(
						name = "foo",
						srcs = [
							"foo.gleam",
						],
						visibility = ["//visibility:private"],
						deps = ["//test"],
					)
`,
	},
	{
		desc:  "import from hex repo",
		index: []buildFile{
			{
				pkg: "foo/bar",
				content: `
					gleam_library(
						name = "bar",
						srcs = [
							"bar.gleam"
						],
					)
`,
			},
		},
		old: buildFile{
			pkg: "foo",
			content: `
					gleam_library(
						name = "foo",
						srcs = [
							"foo.gleam"
						],
						visibility = ["//visibility:private"],
						_gazelle_imports = [
							"foo/bar/bar",
							"gleam/io",
							"gleeunit/gleeunit",
							"gleam/int"
						],
					)
	`,
		},
		want: `
					gleam_library(
						name = "foo",
						srcs = [
							"foo.gleam",
						],
						visibility = ["//visibility:private"],
						deps = [
							"//foo/bar",
							"@hex_gleam_stdlib//gleam",
							"@hex_gleeunit//gleeunit",
						],
					)
`,
	},
	{
		desc: "mixed concept",
		index: []buildFile{
			{
				pkg: "too/many/level/deep",
				content: `
					gleam_library(
						name = "deep",
						srcs = [
							"bar.gleam",
							"level1.gleam",
							"level2.gleam",
							"level3.gleam",
						],
						visibility = ["//visibility:public"],
					)
				`,
			},
			{
				pkg: "test/internal",
				content: `
					gleam_library(
						name = "internal",
						srcs = [
							"a.gleam",
							"b.gleam",
						],
						visibility = ["//test:__subpackages__"],
						deps = ["//too/many/level/deep"],
					)
				`,
			},
		},
		old: buildFile{
			pkg: "test",
			content: `
					gleam_library(
						name = "test_internal",
						srcs = [
							"internal.gleam",
						],
					)

					gleam_binary(
						name = "test",
						srcs = [
							"c.gleam",
							"d.gleam",
						],
						main_module = "c",
						visibility = ["//visibility:private"],
						_gazelle_imports = [
							"test/internal",
							"too/many/level/deep/level1",
							"too/many/level/deep/level2",
							"too/many/level/deep/level3",
							"test/internal/a",
							"test/internal/b",
							"fl",
							"gleam/io",
						],
					)
			`,
		},
		want: `
					gleam_library(
						name = "test_internal",
						srcs = [
							"internal.gleam",
						],
					)

					gleam_binary(
						name = "test",
						srcs = [
							"c.gleam",
							"d.gleam",
						],
						main_module = "c",
						visibility = ["//visibility:private"],
						deps = [
							":test_internal",
							"//test/internal",
							"//too/many/level/deep",
							"@hex_fl//:lib",
							"@hex_gleam_stdlib//gleam",
						],
					)
			`,
	},
}

var testRepos = []repo.Repo{
	repo.Repo{
		Name:     "hex_gleam_stdlib",
		GoPrefix: "gleam/io",
	},
	repo.Repo{
		Name:     "hex_gleeunit",
		GoPrefix: "gleeunit/gleeunit",
	},
	repo.Repo{
		Name:     "hex_gleam_stdlib",
		GoPrefix: "gleam/int",
	},
	repo.Repo{
		Name:     "hex_fl",
		GoPrefix: "fl",
	},
}

func TestResolveGleam(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			c, langs, cexts := testConfig(
				t, fmt.Sprintf("-index=%v", !testCase.skipIndex))
			mrslv := make(mapResolver)
			exts := make([]interface{}, 0, len(langs))
			for _, lang := range langs {
				for kind := range lang.Kinds() {
					mrslv[kind] = lang
				}
				exts = append(exts, lang)
			}
			ix := resolve.NewRuleIndex(mrslv.Resolver, exts...)

			gc := GetGleamConfig(c).clone()
			gc.repos = testRepos
			c.Exts[languageName] = gc

			for _, bf := range testCase.index {
				buildPath := filepath.Join(filepath.FromSlash(bf.pkg), "BUILD")
				f, err := rule.LoadData(buildPath, bf.pkg, []byte(dedent.Dedent(bf.content)))
				if err != nil {
					t.Fatal(err)
				}
				if bf.pkg == "" {
					for _, cext := range cexts {
						cext.Configure(c, "", f)
					}
				}
				for _, r := range f.Rules {
					ix.AddRule(c, r, f)
				}
			}
			buildPath := filepath.Join(filepath.FromSlash(testCase.old.pkg), "BUILD")
			f, err := rule.LoadData(buildPath, testCase.old.pkg, []byte(dedent.Dedent(testCase.old.content)))
			if err != nil {
				t.Fatal(err)
			}
			imports := make([]interface{}, len(f.Rules))
			for i, r := range f.Rules {
				imports[i] = convertImportsAttr(r)
				ix.AddRule(c, r, f)
			}
			ix.Finish()
			for i, r := range f.Rules {
				mrslv.Resolver(r, "").Resolve(c, ix, nil, r, imports[i], label.New("", testCase.old.pkg, r.Name()))
			}
			f.Sync()
			got := strings.TrimSpace(string(bzl.Format(f.File)))
			want := strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(testCase.want)), "\t", "    ")
			diff := cmp.Diff(want, got)
			if diff != "" {
				t.Errorf("(-want, +got):%s", diff)
			}
		})
	}
}

func convertImportsAttr(r *rule.Rule) []string {
	kind := r.Kind()
	value := r.AttrStrings(config.GazelleImportsKey)
	r.DelAttr(config.GazelleImportsKey)
	if _, ok := gleamKinds[kind]; ok {
		return value
	} else {
		return []string{}
	}
}

func testRemoteCache(knownRepos []repo.Repo) *repo.RemoteCache {
	rc, _ := repo.NewRemoteCache(knownRepos)
	return rc
	// rc.RepoRootForImportPath = stubRepoRootForImportPath
	//
	//	rc.HeadCmd = func(_, _ string) (string, error) {
	//		return "", fmt.Errorf("HeadCmd not supported in test")
	//	}
	//
	// rc.ModInfo = stubModInfo
	// return rc
}

// func stubRepoRootForImportPath(importPath string, verbose bool) (*vcs.RepoRoot, error) {
// 	if pathtools.HasPrefix(importPath, "example.com/repo.git") {
// 		return &vcs.RepoRoot{
// 			VCS:  vcs.ByCmd("git"),
// 			Repo: "https://example.com/repo.git",
// 			Root: "example.com/repo.git",
// 		}, nil
// 	}

// 	if pathtools.HasPrefix(importPath, "example.com/repo") {
// 		return &vcs.RepoRoot{
// 			VCS:  vcs.ByCmd("git"),
// 			Repo: "https://example.com/repo.git",
// 			Root: "example.com/repo",
// 		}, nil
// 	}

// 	if pathtools.HasPrefix(importPath, "example.com") {
// 		return &vcs.RepoRoot{
// 			VCS:  vcs.ByCmd("git"),
// 			Repo: "https://example.com",
// 			Root: "example.com",
// 		}, nil
// 	}

// 	return nil, fmt.Errorf("could not resolve import path: %q", importPath)
// }

// func stubModInfo(importPath string) (string, error) {
// 	if pathtools.HasPrefix(importPath, "example.com/repo/v2") {
// 		return "example.com/repo/v2", nil
// 	}
// 	if pathtools.HasPrefix(importPath, "example.com/repo") {
// 		return "example.com/repo", nil
// 	}
// 	return "", fmt.Errorf("could not find module for import path: %q", importPath)
// }
