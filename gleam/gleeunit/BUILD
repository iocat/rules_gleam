load("@rules_gleam//gleam:defs.bzl", "gleam_library", "gleam_erl_library")

exports_files([
    "gleam.toml",
    "manifest.toml",
])

gleam_erl_library(
    name = "gleeunit_ffi",
    srcs = [
        "gleam_gleeunit_progress_ffi.erl",
        "gleam_gleeunit_panic_ffi.erl",
    ],
    visibility = ["//visibility:public"],
)

gleam_library(
    name = "gleeunit",
    srcs = glob(["*.gleam"]),
    deps = [
        ":gleeunit_ffi",
        "@hex_gleam_stdlib//gleam",
        "@hex_gleam_stdlib//gleam/dynamic",
    ],
    strip_src_prefix = "external/rules_gleam+",
    visibility = ["//visibility:public"],
)