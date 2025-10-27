"""JS library rule for Gleam."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR")
load("//gleam/js:build.bzl", "COMMON_ATTRS", "declare_inputs", "declare_libs", "declare_outputs", "get_gleam_compiler", "get_js_prelude")
load("//gleam/js:provider.bzl", "GleamJsPackageInfo")
load("@aspect_rules_js//js:providers.bzl", "js_info")

def _gleam_js_library_impl(ctx):
    """Implementation for the gleam_js_library rule."""
    inputs = declare_inputs(ctx, ctx.files.srcs)
    lib_inputs, lib_path = declare_libs(ctx, ctx.attr.deps)
    outputs = declare_outputs(ctx, ctx.files.srcs)
    working_root = paths.dirname(inputs.toml_file.path)
    gleam_compiler = get_gleam_compiler(ctx)
    if len(outputs.all_files):
        ctx.actions.run_shell(
            inputs = inputs.sources + lib_inputs + [gleam_compiler],
            outputs = outputs.all_files,
            use_default_shell_env = True,
            mnemonic = "GleamJsLibraryCompile",
            command = " && ".join([
                "COMPILER=$(pwd)/%s" % gleam_compiler.path,
                "cd %s" % working_root,
                "$COMPILER compile-package --target javascript --package '.' --out '.' --lib %s --javascript-prelude %s" % (lib_path, get_js_prelude(ctx).path),
                # Move cache outside
                "mv ./%s/* ." % GLEAM_ARTEFACTS_DIR,
            ] + ([
                # Move mjs outside
                "mv ./%s/* ." % ctx.label.package,
            ] if ctx.label.package else [])),
            env = {
                "FORCE_COLOR": "true",
            },
        )

    # Accumulate runfiles.
    runfiles = ctx.runfiles(files = ctx.files.data)
    runfiles = runfiles.merge_all([dep[DefaultInfo].default_runfiles for dep in ctx.attr.deps])
    js_module_depset = depset(direct = outputs.js_files, transitive = [dep[GleamJsPackageInfo].js_module for dep in ctx.attr.deps])

    return [
        DefaultInfo(files = depset(outputs.js_files), runfiles = runfiles),
        GleamJsPackageInfo(
            module_names = outputs.module_names,
            js_module = js_module_depset,
            gleam_cache = depset(direct = outputs.cache_files, transitive = [dep[GleamJsPackageInfo].gleam_cache for dep in ctx.attr.deps]),
            strip_src_prefix = ctx.attr.strip_src_prefix,
            target_name = ctx.label.name,
        ),
        js_info(
            target = ctx.label,
            sources = js_module_depset,
        )
    ]

gleam_js_library = rule(
    implementation = _gleam_js_library_impl,
    attrs = dict(
        COMMON_ATTRS,
        srcs = attr.label_list(
            mandatory = True,
            allow_files = [".gleam"],
        ),
        deps = attr.label_list(
            doc = "The list of dependent gleam modules.",
            providers = [GleamJsPackageInfo],
        ),
    ),
    toolchains = [
        "//gleam_tools:toolchain_type",
    ],
)
