load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_binary(
    name = "gazelle",
    embed = [":gazelle_lib"],
    importpath = "github.com/iocat/rules_gleam/internal/tools/gazelle",
    visibility = ["//visibility:private"],
)

go_test(
    name = "gazelle_test",
    srcs = [
        "diff_test.go",
        "fix_test.go",
        "integration_gazelle_tool_test.go",
        "integration_test.go",
        "profiler_test.go",
    ],
    embed = [":gazelle_lib"],
    deps = [
        "//internal/tools/gazelle/wspace",
        "@com_github_google_go_cmp//cmp",
        "@gazelle//config",
        "@gazelle//testtools",
    ],
)

go_library(
    name = "gazelle_lib",
    srcs = [
        "diff.go",
        "fix.go",
        "fix-update.go",
        "langs.go",
        "main.go",
        "metaresolver.go",
        "print.go",
        "profiler.go",
        "update-repos.go",
    ],
    importpath = "github.com/iocat/rules_gleam/internal/tools/gazelle",
    visibility = ["//visibility:private"],
    deps = [
        "//gazelle/gleam",
        "//internal/tools/gazelle/wspace",
        "@com_github_bazelbuild_buildtools//build",
        "@com_github_pmezard_go_difflib//difflib",
        "@gazelle//config",
        "@gazelle//flag",
        "@gazelle//label",
        "@gazelle//language",
        "@gazelle//language/go",
        "@gazelle//language/proto",
        "@gazelle//merger",
        "@gazelle//repo",
        "@gazelle//resolve",
        "@gazelle//rule",
        "@gazelle//walk",
    ],
)
