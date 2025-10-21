package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	manifestFlag = flag.String("manifest", "", "Path to manifest.toml")
)

func logAndExit(w io.Writer, format string, params ...any) error {
	return fmt.Errorf(format, params...)
}

func logResult(w io.Writer, format string, params ...any) {
	fmt.Fprintf(w, format, params...)
}

type Package struct {
	Name          string   `toml:"name"`
	Version       string   `toml:"version"`
	BuildTools    []string `toml:"build_tools"`
	Requirements  []string `toml:"requirements"`
	OtpApp        string   `toml:"otp_app"`
	Source        string   `toml:"source"`
	OuterChecksum string   `toml:"outer_checksum"`
}

type Manifest struct {
	Packages []Package `toml:"packages"`
}

type Graph struct {
	Nodes map[string]*GleamRepo
	Edges map[string][]string
}

type GleamRepo struct {
	ModuleName string   `json:"module_name"`
	Version    string   `json:"version"`
	Checksum   string   `json:"checksum"`
	OtpApp     string   `json:"otp_app"`
	Deps       []string `json:"deps"`
}

type GleamRepoes struct {
	Repos []GleamRepo `json:"repos"`
}

func main() {
	flag.Parse()

	manifestFile := *manifestFlag
	if manifestFile == "" {
		fmt.Fprintln(os.Stderr, "manifest flag is required")
		os.Exit(1)
	}

	if err := run(manifestFile, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(manifestFile string, w io.Writer) error {
	var (
		content  []byte
		manifest Manifest
		err      error
	)
	if content, err = os.ReadFile(manifestFile); err != nil {
		return logAndExit(w, "failed to read manifest file: %v", err)
	}
	if _, err := toml.Decode(string(content), &manifest); err != nil {
		return logAndExit(w, "failed to decode manifest: %v", err)
	}

	repos := GleamRepoes{}
	for _, pkg := range manifest.Packages {
		repos.Repos = append(repos.Repos, GleamRepo{
			ModuleName: pkg.Name,
			Version:    pkg.Version,
			Checksum:   pkg.OuterChecksum,
			OtpApp:     pkg.OtpApp,
			Deps:       pkg.Requirements,
		})
	}

	// Topologically sort repos
	g := Graph{
		Nodes: make(map[string]*GleamRepo),
		Edges: make(map[string][]string),
	}
	for _, repo := range repos.Repos {
		g.Nodes[repo.ModuleName] = &repo
	}
	for _, pkg := range manifest.Packages {
		for _, dep := range pkg.Requirements {
			if edges, ok := g.Edges[dep]; !ok {
				g.Edges[dep] = []string{pkg.Name}
			} else {
				g.Edges[dep] = append(edges, pkg.Name)
			}
		}
	}
	repos.Repos = make([]GleamRepo, 0)
	for node := range g.Topologically() {
		repos.Repos = append(repos.Repos, *node)
	}

	str := strings.Builder{}
	encoder := json.NewEncoder(&str)
	encoder.SetIndent("", "  ")
	encoder.Encode(repos)
	logResult(w, "%s", str.String())
	return nil
}

func (g *Graph) Topologically() iter.Seq[*GleamRepo] {
	return func(yield func(*GleamRepo) bool) {
		visited := make(map[string]bool)
		stack := make([]*GleamRepo, 0, len(g.Nodes))
		for node := range g.Nodes {
			if !visited[node] {
				g.dfs(node, visited, &stack)
			}
		}
		for len(stack) > 0 {
			if !yield(stack[len(stack)-1]) {
				return
			}
			stack = stack[:len(stack)-1]
		}
	}
}
func (g *Graph) dfs(node string, visited map[string]bool, stack *[]*GleamRepo) {
	visited[node] = true
	for _, neighbor := range g.Edges[node] {
		if !visited[neighbor] {
			g.dfs(neighbor, visited, stack)
		}
	}
	*stack = append(*stack, g.Nodes[node])
}
