package main

import (
	// List of indirect dependencies of the deps from this monorepo,
	// as a hack to get the internal tool to build since the current
	// build is without go_modules enabled and everything has to be
	// linked flatly.
	_ "github.com/bmatcuk/doublestar/v4"
	_ "github.com/kr/text"
	_ "github.com/rogpeppe/go-internal/fmtsort"
	_ "golang.org/x/mod/modfile"
	_ "golang.org/x/sys/unix"
	_ "golang.org/x/tools/go/vcs"
)