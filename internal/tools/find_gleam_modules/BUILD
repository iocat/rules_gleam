load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "find_gleam_modules_lib",
    srcs = ["find_gleam_modules.go"],
    importpath = "github.com/iocat/rules_gleam/internal/tools/find_gleam_modules",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "find_gleam_modules",
    embed = [":find_gleam_modules_lib"],
    visibility = ["//:__subpackages__"],
)

go_test(
    name = "find_gleam_modules_test",
    srcs = ["find_gleam_modules_test.go"],
    embed = [":find_gleam_modules_lib"],
)
