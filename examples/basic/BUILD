load("@gazelle//:def.bzl", "gazelle")
load("@rules_gleam//gleam:defs.bzl", "gleam_binary", "gleam_erl_library")

package(
    default_visibility = ["//visibility:public"],
)

gazelle(
    name = "gazelle",
    gazelle = "@rules_gleam//gazelle",
)

gleam_binary(
    name = "example",
    srcs = ["example.gleam"],
    visibility = ["//visibility:private"],
    deps = [
        ":example_ffi",
        "//demo",
        "@gleam_stdlib//gleam",
        "@gleam_stdlib//gleam/dynamic",
    ],
)

gleam_erl_library(
    name = "example_ffi",
    srcs = ["example_ffi.erl"],
    visibility = ["//visibility:public"],
)
