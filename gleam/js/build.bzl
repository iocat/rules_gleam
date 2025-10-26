"""Helper functions for Gleam JS rules."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam_tools:toolchain.bzl", "gleam_toolchain_info")
load("//gleam:provider.bzl", "GleamJsPackageInfo")

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

GLEAM_ARTEFACTS_DIR = "_gleam_artefacts"

def get_gleam_compiler(ctx):
    toolchain = ctx.toolchains["//gleam_tools:toolchain_type"]
    return toolchain[gleam_toolchain_info].gleam_compiler
def declare_inputs(ctx, srcs):
    """Prepares and stages source files for the Gleam compiler.

def get_js_prelude(ctx):
    toolchain = ctx.toolchains["//gleam_tools:toolchain_type"]
    return toolchain[gleam_toolchain_info].js_prelude
    The resulting directory structure passed to the compiler will be:

def declare_inputs(ctx, srcs, is_binary = False, main_module = "", main_template = None):
    toml_file = None
    for f in ctx.files.srcs:
        if f.basename == "gleam.toml":
            toml_file = f
            break
    if not toml_file:
        fail("gleam.toml not found in srcs")
    .
    ├── gleam.toml
    └── src/
        ├── some_source.gleam          (symlink)
        └── another_source.gleam       (symlink)

    Args:
        ctx: The Bazel rule context object.
        srcs (list): A list of source `File` objects to be prepared.

    Returns:
        A `struct` containing the prepared input files
    """
    toml_file = ctx.actions.declare_file("gleam.toml")
    ctx.actions.write(
        toml_file,
        """
name = "placeholder"
version = "0.0.0"
target = "javascript"
        """,
    )

    # Build inputs.
    sources = []
    for src in srcs:
        input_src = ctx.actions.declare_file(paths.join("src", strip_src_prefix(ctx, src.path)))
        ctx.actions.symlink(output = input_src, target_file = src)
        sources.append(input_src)

    sources.append(toml_file)
    return struct(
        sources = srcs,
        toml_file = toml_file,
        sources = sources,
    )

def declare_outputs(ctx, srcs, is_binary = False, main_module = ""):
def declare_outputs(ctx, srcs):
    """Calculates module names and declares compilation output files.

    For each source file provided, this function generates a corresponding
    module name by replacing path separators with '@'. It then declares
    the expected output files for compilation, including .mjs and .cache artifacts.

    Args:
        ctx: The Bazel rule context object, used for declaring files.
        srcs (list): A list of source `File` objects.
    Returns:
        A struct containing lists of the generated names and files.
    """
    js_files = []
    cache_files = []
    module_names = []

    for src in srcs:
        module_name = paths.replace_extension(src.path, "").replace("/", "@")
        src_path = strip_src_prefix(ctx, src.path)
        module_name = paths.replace_extension(src_path, "").replace("/", "@")
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
def declare_libs(ctx, deps):
    """Stages dependency artifacts into a standard library directory.

    This function assembles a directory structure that can be passed to the
    compiler via a library flag (`--lib`). It populates the directory with the
    necessary compilation artifacts from transitive dependencies.

    The created directory structure looks like this:

        lib/javascript
        └── placeholder
            └── _gleam_artefacts
                ├── *.mjs
                └── *.cache

    Args:
        ctx: The Bazel rule context object.
        deps (list): A list of targets providing `GleamJsPackageInfo`.

    Returns:
        tuple: A tuple containing:
            - depset: A depset of all files staged in the lib directory.
            - str: The path to the root of the created `lib` directory.
    """
    lib_inputs = []
    lib_path = paths.join("lib", "javascript")

    for dep in deps:
        if GleamJsPackageInfo in dep:
            for js_file in dep[GleamJsPackageInfo].js_module.to_list():
                lib_files.append(js_file)
            for cache_file in dep[GleamJsPackageInfo].gleam_cache.to_list():
                lib_files.append(cache_file)
            lib_paths.append(paths.dirname(dep[GleamJsPackageInfo].js_module.to_list()[0].path))
    return lib_files, ":".join(lib_paths)
            target = dep[GleamJsPackageInfo]
            for file in target.js_module.to_list() + target.gleam_cache.to_list():
                link = ctx.actions.declare_file(paths.join(lib_path, "placeholder", GLEAM_ARTEFACTS_DIR, paths.basename(file.path)))
                ctx.actions.symlink(
                    output = link,
                    target_file = file,
                )
                lib_inputs.append(link)

    return (depset(lib_inputs), paths.join(lib_path, "placeholder") if lib_inputs else "")

def strip_src_prefix(ctx, src):
    """
    Strip the prefix from the source path.

    Args:
        ctx (object): The Bazel context.
        src (str): The source path.

    Returns:
        str: The source path with the prefix stripped off.
    """
    prefix = ctx.attr.strip_src_prefix
    if not prefix.endswith("/"):
        prefix = prefix + "/"
    if prefix:
        return src.removeprefix(prefix)
    return src
