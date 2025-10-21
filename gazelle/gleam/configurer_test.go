package gleam

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/google/go-cmp/cmp"
)

func TestMaybePopulateRemoteCacheFromGleamTomlForExternalRepo(t *testing.T) {
	testCases := []struct {
		desc                string
		gleamTomlContent    string
		manifestTomlContent string
		compilerPath        string
		existingRepos       []repo.Repo
		wantRepos           []repo.Repo
		wantErr             bool
	}{
		{
			desc:         "no gleam.toml",
			compilerPath: "echo", // Provide a default compiler path
			wantErr:      false,
		},
		{
			desc: "empty dependencies",
			gleamTomlContent: `
name = "test_package"
version = "0.1.0"

[dependencies]
`,
			wantErr: false,
		},
		{
			desc: "valid dependencies",
			gleamTomlContent: `
name = "test_package"
version = "0.1.0"

[dependencies]
  gleam_stdlib = "0.20.4"
`,
			manifestTomlContent: `
[[packages]]
name = "gleam_stdlib"
version = "0.20.4"
build_tools = []
requirements = []
otp_app = "gleam_stdlib"
source = "hex"
outer_checksum = "valid_checksum"
`,
			compilerPath: "echo", // Mock compiler path
			wantRepos: []repo.Repo{
				{
					Name:     "hex_gleam_stdlib",
					GoPrefix: "gleam_stdlib",
				},
			},
			wantErr: false,
		},
		{
			desc: "multiple dependencies",
			gleamTomlContent: `
name = "test_package"
version = "0.1.0"

[dependencies]
  gleam_stdlib = "0.20.4"
  gleam_otp = "0.5.0"
`,
			manifestTomlContent: `
[[packages]]
name = "gleam_stdlib"
version = "0.20.4"
build_tools = []
requirements = []
otp_app = "gleam_stdlib"
source = "hex"
outer_checksum = "valid_checksum"

[[packages]]
name = "gleam_otp"
version = "0.5.0"
build_tools = []
requirements = []
otp_app = "gleam_otp"
source = "hex"
outer_checksum = "valid_checksum"
`,
			compilerPath: "echo", // Mock compiler path
			wantRepos: []repo.Repo{
				{
					Name:     "hex_gleam_stdlib",
					GoPrefix: "gleam_stdlib",
				},
				{
					Name:     "hex_gleam_otp",
					GoPrefix: "gleam_otp",
				},
				{
					Name:     "hex_gleam_stdlib",
					GoPrefix: "erl:gleam_stdlib_ffi",
				},
				{
					Name:     "hex_gleam_otp",
					GoPrefix: "erl:gleam_otp_ffi",
				},
			},
			wantErr: false,
		},
		{
			desc: "erlang dependencies",
			gleamTomlContent: `
name = "test_package"
version = "0.1.0"

[dependencies]
  gleam_stdlib = "0.20.4"
`,
			manifestTomlContent: `
[[packages]]
name = "gleam_stdlib"
version = "0.20.4"
build_tools = []
requirements = []
otp_app = "gleam_stdlib"
source = "hex"
outer_checksum = "valid_checksum"
`,
			compilerPath: "echo", // Mock compiler path
			wantRepos: []repo.Repo{
				{
					Name:     "hex_gleam_stdlib",
					GoPrefix: "gleam_stdlib",
				},
				{
					Name:     "hex_gleam_stdlib",
					GoPrefix: "erl:gleam_stdlib_ffi",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repoRoot, err := os.MkdirTemp("", "testrepo")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(repoRoot)

			if tc.gleamTomlContent != "" {
				gleamTomlPath := filepath.Join(repoRoot, "gleam.toml")
				if err := os.WriteFile(gleamTomlPath, []byte(tc.gleamTomlContent), 0644); err != nil {
					t.Fatalf("Failed to write gleam.toml: %v", err)
				}
			}

			if tc.manifestTomlContent != "" {
				manifestTomlPath := filepath.Join(repoRoot, "manifest.toml")
				if err := os.WriteFile(manifestTomlPath, []byte(tc.manifestTomlContent), 0644); err != nil {
					t.Fatalf("Failed to write manifest.toml: %v", err)
				}
			}

			// Create build/packages directory and mock source files
			if len(tc.wantRepos) > 0 {
				buildPackagesPath := filepath.Join(repoRoot, "build", "packages")
				if err := os.MkdirAll(buildPackagesPath, 0755); err != nil {
					t.Fatalf("Failed to create build/packages directory: %v", err)
				}

				for _, wantRepo := range tc.wantRepos {
					packageName := strings.TrimPrefix(wantRepo.Name, "hex_")
					packageSrcPath := filepath.Join(buildPackagesPath, packageName, "src")
					if err := os.MkdirAll(packageSrcPath, 0755); err != nil {
						t.Fatalf("Failed to create package src directory: %v", err)
					}

					if strings.HasPrefix(wantRepo.GoPrefix, "erl:") {
						erlFfiPath := filepath.Join(packageSrcPath, packageName+"_ffi.erl")
						if err := os.WriteFile(erlFfiPath, []byte(""), 0644); err != nil {
							t.Fatalf("Failed to write mock .gleam file: %v", err)
						}
					} else {
						// Create a mock .gleam file
						gleamFilePath := filepath.Join(packageSrcPath, packageName+".gleam")
						if err := os.WriteFile(gleamFilePath, []byte(""), 0644); err != nil {
							t.Fatalf("Failed to write mock .gleam file: %v", err)
						}
					}
				}
			}

			c := &config.Config{
				RepoRoot: repoRoot,
				Exts:     map[string]interface{}{"gleam": &GleamConfig{}},
			}
			gc := &GleamConfig{
				gleamCompilerPath: tc.compilerPath,
				externalRepo:      true,
			}

			var repos []repo.Repo
			repos = append(repos, tc.existingRepos...)

			err = maybePopulateRemoteCacheFromGleamTomlForExternalRepo(c, gc, &repos)

			if (err != nil) != tc.wantErr {
				t.Fatalf("maybePopulateRemoteCacheFromGleamTomlForExternalRepo() error = %v, wantErr %v", err, tc.wantErr)
			}

			slices.SortFunc(tc.wantRepos, sortFunc)
			slices.SortFunc(repos, sortFunc)
			if diff := cmp.Diff(tc.wantRepos, repos); diff != "" {
				t.Errorf("maybePopulateRemoteCacheFromGleamTomlForExternalRepo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func sortFunc (a, b repo.Repo) int {
	return strings.Compare(a.GoPrefix, b.GoPrefix)
}