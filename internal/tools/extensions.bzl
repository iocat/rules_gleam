
load(
    "//internal/tools:go_cache_repositories.bzl",
    "go_repository_cache",
)
load(
    "//internal/tools:gleam_repository_tools.bzl",
    "gleam_repository_tools",
)
load(
    "@go_host_compatible_sdk_label//:defs.bzl",
    "HOST_COMPATIBLE_SDK",
)
load("//internal:common.bzl", "extension_metadata")

visibility("//")



def _tools_impl(module_ctx):
    go_repository_cache(
        name = "rules_gleam_go_repository_cache",
        go_env = {},
        go_sdk_name = "@" + HOST_COMPATIBLE_SDK.repo_name,
    )
    gleam_repository_tools(
        name = "rules_gleam_internal_tools",
        go_cache = Label("@rules_gleam_go_repository_cache//:go.env")
    )

    return extension_metadata(module_ctx, reproducible = True)

tools = module_extension(
    _tools_impl,
)
