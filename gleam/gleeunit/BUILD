load("@rules_gleam//gleam:defs.bzl", "gleam_erl_library", "gleam_library")

exports_files([
    "gleam.toml",
    "manifest.toml",
])

gleam_erl_library(
    name = "gleeunit_ffi",
    srcs = [
        "gleam_gleeunit_panic_ffi.erl",
        "gleam_gleeunit_progress_ffi.erl",
    ],
    visibility = ["//visibility:public"],
)

gleam_library(
    name = "gleeunit",
    srcs = glob(["*.gleam"]),
    strip_src_prefix = "external/rules_gleam+",
    visibility = ["//visibility:public"],
    deps = [
        ":gleeunit_ffi",
        "@hex_gleam_stdlib//gleam:bit_array",
        "@hex_gleam_stdlib//gleam:dynamic",
        "@hex_gleam_stdlib//gleam:int",
        "@hex_gleam_stdlib//gleam:io",
        "@hex_gleam_stdlib//gleam:list",
        "@hex_gleam_stdlib//gleam:option",
        "@hex_gleam_stdlib//gleam:result",
        "@hex_gleam_stdlib//gleam:string",
    ],
)
