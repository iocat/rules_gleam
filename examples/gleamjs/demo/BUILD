load("@rules_gleam//gleam/js:defs.bzl", "gleam_js_library")

package(
    default_visibility = ["//visibility:public"],
)

gleam_js_library(
    name = "demo",
    srcs = ["demo.gleam"],
    visibility = ["//visibility:public"],
    deps = [":internal"],
)

gleam_js_library(
    name = "internal",
    srcs = ["internal.gleam"],
    visibility = ["//demo:__subpackages__"],
)
