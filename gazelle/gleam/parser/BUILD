load("@rules_go//go:def.bzl", "go_library", "go_test")

package(
    default_visibility = ["//gazelle/gleam:__subpackages__"],
)

go_library(
    name = "parser",
    srcs = ["parser.go"],
    importpath = "github.com/iocat/rules_gleam/gazelle/gleam/parser",
)

go_test(
    name = "parser_test",
    srcs = ["parser_test.go"],
    embed = [":parser"],
    deps = ["@com_github_google_go_cmp//cmp"],
)
