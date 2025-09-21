load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

# go_binary(
#     name = "gazelle",
#     srcs = [
#         "main.go",  # keep
#     ],
#     embed = [":gazelle_lib"],
#     importpath = "github.com/bazelbuild/bazel-gazelle/cmd/gazelle",
#     tags = ["manual"],
#     visibility = ["//visibility:public"],
# )

go_binary(
    name = "gazelle",
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
        "@gazelle//internal/wspace",
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

go_test(
    name = "gazelle_test",
    srcs = [
        "diff_test.go",
        "fix_test.go",
        "integration_test.go",
        "profiler_test.go",
    ],
    embed = [":gazelle_lib"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@gazelle//config",
        "@gazelle//internal/wspace",
        "@gazelle//testtools",
        "@rules_go//go/runfiles",
    ],
)
