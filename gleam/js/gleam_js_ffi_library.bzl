load("//gleam/js:provider.bzl", "GleamJsPackageInfo")
load("@aspect_rules_js//js:providers.bzl", "js_info")

def _gleam_js_ffi_library_impl(ctx):
    """Implementation for the gleam_js_ffi_library rule."""
    transitive_js_modules = [dep[GleamJsPackageInfo].js_module for dep in ctx.attr.deps]
    all_js_files = depset(direct = ctx.files.srcs, transitive = transitive_js_modules)

    transitive_gleam_cache = [dep[GleamJsPackageInfo].gleam_cache for dep in ctx.attr.deps]

    # Accumulate runfiles.
    runfiles = ctx.runfiles(files = ctx.files.srcs + ctx.files.data)
    runfiles = runfiles.merge_all([dep[DefaultInfo].default_runfiles for dep in ctx.attr.deps])

    return [
        DefaultInfo(files = depset(ctx.files.srcs), runfiles = runfiles),
        GleamJsPackageInfo(
            js_module = all_js_files,
            module_names = depset(),  # No gleam modules produced here
            gleam_cache = depset(transitive = transitive_gleam_cache),
            strip_src_prefix = "",  # Not applicable
            target_name = ctx.label.name,
        ),
        js_info(
            target = ctx.label,
            sources = all_js_files,
        )
    ]

gleam_js_ffi_library = rule(
    implementation = _gleam_js_ffi_library_impl,
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = [".mjs"],
            doc = "The list of JavaScript module files for the FFI.",
        ),
        "deps": attr.label_list(
            doc = "The list of dependent gleam ffi modules.",
            providers = [GleamJsPackageInfo],
            default = [],
        ),
        "data": attr.label_list(
            allow_files = True,
            doc = "Runtime dependencies for this library.",
        ),
    },
)
