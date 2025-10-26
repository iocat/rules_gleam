"""Bazel BUILD file for Gleam JS rules."""

package(default_visibility = ["//visibility:public"])
load("//gleam/js:defs.bzl", "gleam_js_binary", "gleam_js_library", "gleam_js_test")

# Example of a Gleam JS library
gleam_js_library(
    name = "my_js_library",
    srcs = ["my_library.gleam"],
    visibility = ["//visibility:public"],
)

# Example of a Gleam JS binary
gleam_js_binary(
    name = "my_js_binary",
    srcs = ["main.gleam"],
    deps = [":my_js_library"],
    visibility = ["//visibility:public"],
)

# Example of a Gleam JS test
gleam_js_test(
    name = "my_js_test",
    srcs = ["my_test.gleam"],
    deps = [":my_js_library"],
    visibility = ["//visibility:public"],
)
