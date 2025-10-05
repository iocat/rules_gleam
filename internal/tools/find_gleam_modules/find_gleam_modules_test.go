package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestFindModules(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	files := []string{
		"module_a.gleam",
		"sub/module_b.gleam",
		"sub/deep/module_c.gleam",
		"not_gleam.txt",
		"sub/another.erl",
		"sub/deep/internal/a.gleam",
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create parent dirs for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	expected := &Result{
		GleamModules: []string{
			"module_a",
			"sub/module_b",
			"sub/deep/module_c",
			"sub/deep/internal/a",
		},
	}

	result, err := findModules(tmpDir)
	if err != nil {
		t.Fatalf("findModules failed: %v", err)
	}

	sort.Strings(result.GleamModules)
	sort.Strings(expected.GleamModules)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
