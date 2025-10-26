"""Helper functions for Gleam JS rules."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam_tools:toolchain.bzl", "gleam_toolchain_info")

COMMON_ATTRS = {
    "data": attr.label_list(
        allow_files = True,
        doc = "Extra files to be added to the runfiles.",
    ),
    "strip_src_prefix": attr.string(
        doc = "The prefix to strip from the source file paths when compiling.",
        default = "",
    ),
}

GLEAM_ARTEFACTS_DIR = "build"

def get_gleam_compiler(ctx):
    toolchain = ctx.toolchains["//gleam_tools:toolchain_type"]
    return toolchain[gleam_toolchain_info].gleam_compiler

def declare_inputs(ctx, srcs, is_binary = False, main_module = "", main_template = None):
    toml_file = None
    for f in ctx.files.srcs:
        if f.basename == "gleam.toml":
            toml_file = f
            break
    if not toml_file:
        fail("gleam.toml not found in srcs")

    return struct(
        sources = srcs,
        toml_file = toml_file,
    )

def declare_outputs(ctx, srcs, is_binary = False, main_module = ""):
    js_files = []
    cache_files = []
    module_names = []

    for src in srcs:
        module_name = paths.replace_extension(src.path, "").replace("/", "@")
        js_files.append(ctx.actions.declare_file(module_name + ".mjs"))
        cache_files.append(ctx.actions.declare_file(module_name + ".cache"))
        module_names.append(module_name)

    all_files = js_files + cache_files
    all_files_include_binary = all_files

    return struct(
        js_files = js_files,
        cache_files = cache_files,
        module_names = module_names,
        all_files = all_files,
        all_files_include_binary = all_files_include_binary,
    )

def declare_lib_files_for_dep(ctx, deps):
    lib_files = []
    lib_paths = []
    for dep in deps:
        if GleamJsPackageInfo in dep:
            for js_file in dep[GleamJsPackageInfo].js_module.to_list():
                lib_files.append(js_file)
            for cache_file in dep[GleamJsPackageInfo].gleam_cache.to_list():
                lib_files.append(cache_file)
            lib_paths.append(paths.dirname(dep[GleamJsPackageInfo].js_module.to_list()[0].path))
    return lib_files, ":".join(lib_paths)
