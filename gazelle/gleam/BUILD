load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")
load("@rules_go//go:def.bzl", "go_library")

package(
    default_visibility = ["//visibility:public"],
)

gazelle(name = "gazelle")

go_library(
    name = "gleam_library",
    srcs = [
        "configurer.go",
        "language.go",
        "language_generate_rules.go",
        "resolver.go",
        "utils.go",
    ],
    importpath = "github.com/iocat/rules_gleam/gazelle/gleam",
    deps = [
        "//gazelle/gleam/parser",
        "@gazelle//config",
        "@gazelle//label",
        "@gazelle//language",
        "@gazelle//pathtools",
        "@gazelle//repo",
        "@gazelle//resolve",
        "@gazelle//rule",
    ],
)

gazelle_binary(
    name = "gleam",
    languages = [":gleam_library"],
    visibility = ["//visibility:public"],
)
