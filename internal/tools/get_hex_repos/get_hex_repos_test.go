package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetHexRepos(t *testing.T) {
	manifestContent := `
packages = [
  { name = "gleam_stdlib", version = "0.34.0", otp_app = "stdlib", outer_checksum = "outersum_stdlib" },
  { name = "gleeunit", version = "1.0.0", otp_app = "gleeunit_app", requirements = ["gleam_stdlib"], outer_checksum = "outersum_gleeunit" },
  { name = "gleam_erlang", version = "0.25.0", otp_app = "erlang", requirements = ["gleam_stdlib"], outer_checksum = "outersum_erlang" },
]
`
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "manifest.toml")
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("failed to write manifest file: %v", err)
	}

	var stdout bytes.Buffer
	if err := run(manifestPath, &stdout); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	var result GleamRepoes
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	if len(result.Repos) != 3 {
		t.Fatalf("expected 3 repos, got %d", len(result.Repos))
	}

	// Check for topological sort. gleam_stdlib should be first.
	if result.Repos[0].ModuleName != "gleam_stdlib" {
		t.Fatalf("expected first repo to be gleam_stdlib, got %s", result.Repos[0].ModuleName)
	}

	// Verify content of one of the packages
	var gleeunitRepo GleamRepo
	found := false
	for _, repo := range result.Repos {
		if repo.ModuleName == "gleeunit" {
			gleeunitRepo = repo
			found = true
			break
		}
	}
	if !found {
		t.Fatal("gleeunit repo not found")
	}

	if gleeunitRepo.ModuleName != "gleeunit" {
		t.Errorf("expected module name gleeunit, got %s", gleeunitRepo.ModuleName)
	}
	if gleeunitRepo.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", gleeunitRepo.Version)
	}
	if gleeunitRepo.Checksum != "outersum_gleeunit" {
		t.Errorf("expected checksum outersum_gleeunit, got %s", gleeunitRepo.Checksum)
	}
	if gleeunitRepo.OtpApp != "gleeunit_app" {
		t.Errorf("expected otp_app gleeunit_app, got %s", gleeunitRepo.OtpApp)
	}
}
