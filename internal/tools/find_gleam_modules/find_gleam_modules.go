package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	repoDir = flag.String("repo_dir", "", "Path to repo directory")
)

func logAndExit(format string, params ...any) {
	fmt.Fprintf(os.Stderr, format, params...)
	os.Exit(1)
}

func logResult(format string, params ...any) {
	fmt.Fprintf(os.Stdout, format, params...)
	os.Exit(0)
}

type Result struct {
	GleamModules []string `json:"gleam_modules"`
}

func findModules(repoDir string) (*Result, error) {
	stat, err := os.Stat(repoDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("repo_dir does not exist: dir '%s'", repoDir)
		}
		return nil, fmt.Errorf("failed to stat repo_dir: %s: %v", repoDir, err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("repo_dir is not a directory: %s", repoDir)
	}

	result := &Result{
		GleamModules: []string{},
	}
	err = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".gleam") {
			relPath, err := filepath.Rel(repoDir, path)
			if err != nil {
				return err
			}
			modulePath := strings.TrimSuffix(relPath, ".gleam")
			modulePath = strings.ReplaceAll(modulePath, string(os.PathSeparator), "/")
			result.GleamModules = append(result.GleamModules, modulePath)
		}
		return nil
	})
	sort.Strings(result.GleamModules)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	flag.Parse()

	repoDir := *repoDir
	if repoDir == "" {
		panic(fmt.Errorf("repo directory is required"))
	}

	result, err := findModules(repoDir)
	if err != nil {
		logAndExit(err.Error())
	}

	str := strings.Builder{}
	encoder := json.NewEncoder(&str)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
	logResult("%s", str.String())
}
