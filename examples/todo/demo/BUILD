load("@gazelle//:def.bzl", "gazelle")
load("@rules_gleam//gleam:defs.bzl", "gleam_binary", "gleam_erl_library", "gleam_test")

package(
    default_visibility = ["//visibility:public"],
)

gazelle(
    name = "gazelle",
    gazelle = "@rules_gleam//gazelle",
)

# Main todo library
gleam_binary(
    name = "todo_lib",
    srcs = ["todo.gleam"],
    deps = [
        "@hex_gleam_stdlib//gleam",
    ],
    visibility = ["//visibility:public"],
)

# Demo binary
gleam_binary(
    name = "demo",
    srcs = ["demo/demo.gleam"],
    deps = [
        ":todo_lib",
        "@hex_gleam_stdlib//gleam",
    ],
    visibility = ["//visibility:public"],
)

# Test target
gleam_test(
    name = "todo_test",
    srcs = ["todo_test.gleam"],
    deps = [
        ":todo_lib",
        "@hex_gleeunit//gleam",
    ],
)
