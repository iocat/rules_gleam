load("@gazelle//:def.bzl", "gazelle")
load("@rules_gleam//gleam:defs.bzl", "gleam_binary")

package(
    default_visibility = ["//visibility:public"],
)

gazelle(
    name = "gazelle",
    gazelle = "@rules_gleam//gazelle/gleam",
)

gleam_binary(
    name = "example",
    srcs = ["example.gleam"],
    visibility = ["//visibility:private"],
    deps = ["//demo"],
)
