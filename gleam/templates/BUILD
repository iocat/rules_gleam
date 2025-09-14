load("@rules_gleam//gleam:defs.bzl", "gleam_erl_library")

exports_files([
    "[[main_module]].app",
    "[[main_module]].sh",
    "[[main_module]]@@main.erl",
])

gleam_erl_library(
    name = "templates_ffi",
    srcs = ["[[main_module]]@@main.erl"],
    visibility = ["//visibility:public"],
)
