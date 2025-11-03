load("@rules_gleam//gleam/js:defs.bzl", "gleam_js_binary", "gleam_js_ffi_library", "gleam_js_test")

package(
    default_visibility = ["//visibility:public"],
)

gleam_js_binary(
    name = "exec",
    srcs = ["example.gleam"],
    visibility = ["//visibility:private"],
    deps = [
        ":example_ffi",
        "//demo",
    ],
)

gleam_js_test(
    name = "example_test",
    srcs = ["example_test.gleam"],
    visibility = ["//visibility:private"],
)

gleam_js_ffi_library(
    name = "example_ffi",
    srcs = ["example_ffi.mjs"],
    deps = [
    ],
)
