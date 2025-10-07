load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")
load("@gleam_hex_repositories_config//:defs.bzl", "DEPS")
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
    data = DEPS + [
        "@gleam_hex_repositories_config//:BUILD",  # keep
    ],
    importpath = "github.com/iocat/rules_gleam/gazelle/gleam",
    deps = [
        "//gazelle/gleam/parser",
        "@com_github_bazelbuild_buildtools//build",
        "@com_github_burntsushi_toml//:toml",
        "@com_github_kr_pretty//:pretty",
        "@gazelle//config",
        "@gazelle//label",
        "@gazelle//language",
        "@gazelle//pathtools",
        "@gazelle//repo",
        "@gazelle//resolve",
        "@gazelle//rule",
        "@rules_go//go/runfiles",
    ],
)

gazelle_binary(
    name = "gazelle_binary",
    languages = [":gleam"],
    visibility = ["//gazelle:__subpackages__"],
)

go_test(
    name = "gleam_test",
    srcs = [
        "config_test.go",
        "language_generate_rules_test.go",
        "resolver_test.go",
    ],
    data = glob(["gentestdata/**"]) + DEPS + [
        "@gleam_hex_repositories_config//:BUILD",  # keep
    ],
    embed = [":gleam"],
    deps = [
        "@com_github_bazelbuild_buildtools//build",
        "@com_github_google_go_cmp//cmp",
        "@com_github_lithammer_dedent//:dedent",
        "@gazelle//config",
        "@gazelle//label",
        "@gazelle//language",
        "@gazelle//merger",
        "@gazelle//repo",
        "@gazelle//resolve",
        "@gazelle//rule",
        "@gazelle//testtools",
        "@gazelle//walk",
    ],
)
