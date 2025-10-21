load("@gazelle//:def.bzl", "gazelle")
load("@rules_gleam//gleam:defs.bzl", "gleam_binary", "gleam_erl_library", "gleam_test")

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
        ":example_ffi_ffi",
        "//demo",
        "@hex_gleam_json//gleam:json",
        "@hex_gleam_stdlib//gleam:int",
        "@hex_gleam_stdlib//gleam:io",
        "@hex_gleam_stdlib//gleam:option",
    ],
)

gleam_erl_library(
    name = "example_ffi_ffi",
    srcs = ["example_ffi.erl"],
    visibility = ["//visibility:public"],
)

gleam_test(
    name = "gleam_lib_test",
    srcs = ["example_test.gleam"],
    deps = [
        ":example",
        "@hex_gleam_stdlib//gleam:option",
    ],
)
