load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")
load("@rules_go//go:def.bzl", "go_library", "go_test")

package(
    default_visibility = ["//visibility:public"],
)

gazelle(name = "gazelle")

go_library(
    name = "gleam",
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
        "@com_github_kr_pretty//:pretty",
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
    name = "gazelle_binary",
    languages = [":gleam"],
    visibility = ["//gazelle:__subpackages__"],
)

go_test(
    name = "gleam_test",
    srcs = ["language_generate_rules_test.go"],
    data = glob(["testdata/**"]),
    embed = [":gleam"],
    deps = [
        "@com_github_bazelbuild_buildtools//build",
        "@com_github_google_go_cmp//cmp",
        "@gazelle//config",
        "@gazelle//language",
        "@gazelle//merger",
        "@gazelle//resolve",
        "@gazelle//rule",
        "@gazelle//testtools",
        "@gazelle//walk",
    ],
)
