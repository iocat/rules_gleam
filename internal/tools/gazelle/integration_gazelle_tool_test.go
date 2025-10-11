package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
)

func TestGazelleToolCallIntegration(t *testing.T) {
	dir, cleanup := testtools.CreateFiles(t, []testtools.FileSpec{
		{Path: "WORKSPACE", Content: ""},
		{Path: "go.mod", Content: "module example.com/test\n\ngo 1.20\n"},
		{Path: "main.go", Content: "package main\nimport \"example.com/test/pkg\"\nfunc main() { pkg.Hello() }"},
		{Path: "pkg/pkg.go", Content: "package pkg\nimport \"fmt\"\nfunc Hello() { fmt.Println(\"Hello from pkg\") }"},
	})
	defer cleanup()

	// Call gazelle tool (simulate integration)
	err := runGazelle(dir, []string{"update", "-go_prefix", "example.com/test", "-repo_root", dir})
	if err != nil {
		t.Fatalf("gazelle tool call failed: %v", err)
	}

	// Check that BUILD or BUILD.bazel file was generated
	buildPath := filepath.Join(dir, "BUILD")
	if _, err := os.Stat(buildPath); err != nil {
		// Try BUILD.bazel if BUILD is missing
		buildBazelPath := filepath.Join(dir, "BUILD.bazel")
		if _, errBazel := os.Stat(buildBazelPath); errBazel != nil {
			// Log directory contents for debugging
			dirs, listErr := os.ReadDir(dir)
			if listErr == nil {
				t.Logf("Directory contents after gazelle run:")
				for _, entry := range dirs {
					info, infoErr := entry.Info()
					if infoErr == nil {
						t.Logf("- %s (size: %d bytes)", entry.Name(), info.Size())
					} else {
						t.Logf("- %s (error getting info)", entry.Name())
					}
				}
			} else {
				t.Logf("Error reading directory contents: %v", listErr)
			}
			t.Fatalf("Neither BUILD nor BUILD.bazel file generated: BUILD error: %v, BUILD.bazel error: %v", err, errBazel)
		}
	}
}
