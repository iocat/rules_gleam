load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "hack_for_transitive_deps_lib",
    srcs = ["hack_keep_indirect_deps.go"],
    importpath = "github.com/iocat/rules_gleam/internal/tools/hack_for_transitive_deps",
    visibility = ["//:__subpackages__"],
    deps = [
        "@com_github_bazelbuild_buildtools//build",
        "@com_github_bmatcuk_doublestar_v4//:doublestar",
        "@com_github_kr_text//:text",
        "@com_github_rogpeppe_go_internal//fmtsort",
        "@org_golang_x_mod//modfile",
        "@org_golang_x_sync//errgroup",
        "@org_golang_x_sys//unix",
        "@org_golang_x_tools_go_vcs//:vcs",
    ],
)
