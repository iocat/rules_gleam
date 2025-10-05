load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "get_hex_repos_lib",
    srcs = ["get_hex_repos.go"],
    importpath = "github.com/iocat/rules_gleam/internal/tools/get_hex_repos",
    visibility = ["//visibility:private"],
    deps = ["@com_github_burntsushi_toml//:toml"],
)

go_binary(
    name = "get_hex_repos",
    embed = [":get_hex_repos_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "get_hex_repos_test",
    srcs = ["get_hex_repos_test.go"],
    embed = [":get_hex_repos_lib"],
)
