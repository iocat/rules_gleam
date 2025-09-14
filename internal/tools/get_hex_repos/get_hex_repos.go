package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	manifestFlag = flag.String("manifest", "", "Path to manifest.toml")
)

func logAndExit(format string, params ...any) {
	fmt.Fprintf(os.Stderr, format, params...)
	os.Exit(1)
}

func logResult(format string, params ...any) {
	fmt.Fprintf(os.Stdout, format, params...)
	os.Exit(0)
}


type Package struct {
	Name string `toml:"name"`
	Version string `toml:"version"`
	BuildTools []string `toml:"build_tools"`
	Requirements []string `toml:"requirements"`
	OtpApp string `toml:"otp_app"`
	Source string `toml:"source"`
	OuterChecksum string `toml:"outer_checksum"`
}

type Manifest struct {
	Packages []Package `toml:"packages"`
}

type GleamRepo struct {
	ModuleName string `json:"module_name"`
	Version string `json:"version"`
	Checksum string `json:"checksum"`
	OtpApp string `json:"otp_app"`
}

type GleamRepoes struct {
	Repos []GleamRepo`json:"repos"`
}

func main() {
	flag.Parse()

	manifestFile := *manifestFlag
	if manifestFile == "" {
		panic(fmt.Errorf("manifest flag is required"))
	}

	var (
		content []byte
		manifest Manifest
	 	err error
	)
	if content, err = os.ReadFile(manifestFile); err != nil {
		logAndExit("failed to read manifest file: %v", err)
	}
	if _, err := toml.Decode(string(content), &manifest); err != nil {
		logAndExit("failed to decode manifest: %v", err)
	}

	repos := GleamRepoes{}
	for _, pkg := range manifest.Packages {
		repos.Repos = append(repos.Repos, GleamRepo{
			ModuleName: pkg.Name,
			Version: pkg.Version,
			Checksum: pkg.OuterChecksum,
			OtpApp: pkg.OtpApp,
		})
	}
	str := strings.Builder{}
	encoder := json.NewEncoder(&str)
	encoder.SetIndent("", "  ")
	encoder.Encode(repos)
	logResult("%s", str.String())
}