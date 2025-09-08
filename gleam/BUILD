load("//gleam/private:gleam_bin_for_host.bzl", "gleam_bin_for_host")

package(default_visibility = ["//visibility:public"])

# Runs with
# sh
#    bazel @rules_gleam//gleam -- <args>
alias(
    name = "gleam",
    actual = ":gleam_bin",
    visibility = ["//visibility:public"],
)

gleam_bin_for_host(
    name = "gleam_bin",
    visibility = ["//visibility:private"],
)
