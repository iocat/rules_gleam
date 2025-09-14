package main

import (
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/iocat/rules_gleam/gazelle/gleam"
)

// This is fed to the gazelle lib declared at @gazelle//cmd/gazelle
var languages = []language.Language{
	gleam.NewLanguage(),
}