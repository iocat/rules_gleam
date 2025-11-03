load("@rules_gleam//gleam:defs.bzl", "gleam_library", "gleam_test")

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

gleam_test(
    name = "demo_test",
    srcs = [
        "demo_test.gleam",
        "internal_test.gleam",
    ],
    deps = [
        ":demo",
        ":internal",
    ],
)
