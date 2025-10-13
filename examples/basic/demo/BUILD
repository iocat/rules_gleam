load("@rules_gleam//gleam:defs.bzl", "gleam_library")

package(
    default_visibility = ["//visibility:public"],
)

gleam_library(
    name = "demo",
    srcs = ["demo.gleam"],
    visibility = ["//visibility:public"],
    deps = [":internal"],
)

gleam_library(
    name = "internal",
    srcs = ["internal.gleam"],
    visibility = ["//demo:__subpackages__"],
)
