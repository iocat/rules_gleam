load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR", "GLEAM_EBIN_DIR", "GleamErlPackageInfo")

def declare_out_file_with_ext(ctx, src, ext, fully_qual_path = True):
    """Declare a file with the given extension

    Args:
        ctx: The Bazel context
        src: The source file
        ext: The extension to append to the file name
        fully_qual_path: True if the declared file should have a fully qualified path, false otherwise

    Returns:
        A file object representing the declared file
    """
    file_path = ""
    if fully_qual_path:
        file_path = paths.replace_extension(strip_src_prefix(ctx,src.path), ext).replace("/", "@")
    else:
        file_path = paths.replace_extension(paths.basename(strip_src_prefix(ctx,src.path)), ext)

    file = ctx.actions.declare_file(file_path)
    return file

def _get_gleam_package_name(ctx):
    return ctx.label.package.replace("/", "_")

# Declare outputs by the compiler.
def declare_outputs(ctx, srcs, *, is_binary, main_module):
    """Calculates module names and declares compilation output files.

    For each source file provided, this function generates a corresponding
    module name by replacing path separators with '@'. It then declares
    the expected output files for compilation, including .erl, .beam,
    and cache artifacts.

     .
    ├── main@module@@main.app
    ├── main@module@@main.beam
    ├── main@module@@main.sh
    ├── module.app
    ├── module.beam
    ├── module.cache
    ├── module.cache_meta
    ├── module.erl

    Args:
        ctx: The Bazel rule context object, used for declaring files.
        srcs (list): A list of source `File` objects.
        is_binary: whether this compiles binary.
        main_module: the name of the main module.

    Returns:
        A dictionary containing lists of the generated names and files.
    """
    all_files = []
    module_names = []
    output_erl_mods = []
    output_beam_files = []
    output_cache_files = []
    output_binary_beam = None
    output_entry_point = None
    beam_app_manifest = None

    for src in srcs:
        src_path = strip_src_prefix(ctx, src.path)
        module_name = paths.replace_extension(src_path, "").replace("/", "@")
        module_names.append(module_name)
        ext = paths.split_extension(src_path)[1]
        if ext == ".gleam":
            # Gleam produces other artefacts.
            output_erl_mods.append(declare_out_file_with_ext(ctx, src, ".erl"))
            output_beam_files.append(declare_out_file_with_ext(ctx, src, ".beam"))
            output_cache_files.append(declare_out_file_with_ext(ctx, src, ".cache"))
            output_cache_files.append(declare_out_file_with_ext(ctx, src, ".cache_meta"))
        if ext == ".erl":
            # Erl produces only beam, since there's only one bytecode compilation pass.
            output_beam_files.append(declare_out_file_with_ext(ctx, src, ".beam", fully_qual_path = False))
        all_files.extend(output_erl_mods)
        all_files.extend(output_beam_files)
        all_files.extend(output_cache_files)

    if is_binary:
        output_binary_beam = ctx.actions.declare_file(main_module + "@@main.beam")
        output_entry_point = ctx.actions.declare_file(main_module + ".sh")
        beam_app_manifest = ctx.actions.declare_file(main_module + ".app")
        output_beam_files.append(output_binary_beam)

    return struct(
        # All, except binary which has a separate compilation/creation process
        all_files = all_files,
        module_names = module_names,
        erl_mods = output_erl_mods,
        beam_files = output_beam_files,
        cache_files = output_cache_files,
        all_files_include_binary = all_files + [output_binary_beam] if is_binary else [],
        # output_binary_erl_mod = output_binary_erl_mod,
        output_binary_beam = output_binary_beam,
        output_entry_point = output_entry_point,
        beam_app_manifest = beam_app_manifest,
    )

def declare_lib_files_for_dep(ctx, deps):
    """Stages dependency artifacts into a standard library directory.

    This function assembles a directory structure that can be passed to the
    compiler via a library flag (e.g., --lib lib/erlang). It populates the
    directory with the necessary compilation artifacts from transitive
    dependencies.

    The created directory structure looks like this:

        lib/erlang
        └── placeholder -- this is a random name to pack everything
            ├── _gleam_artefacts
            │   ├── *.cache
            │   ├── *.cache_meta
            │   └── *.erl
            └── ebin
                └── *.beam

    Args:
        ctx: The Bazel rule context object.
        deps (list): A list of GleamErlPackageInfo targets
            providing the necessary artifact files.

    Returns:
        tuple: A tuple containing:
            - depset: A depset of all_files files staged in the lib directory.
            - str: The path to the root of the created `lib` directory.
    """
    lib_gleam_art_files = []

    # Beam modules are typically namespaced, unless it's an external module or
    # modules for FFI, then there might be conflicts.
    lib_beam_files = []
    for dep in deps:
        target = dep[GleamErlPackageInfo]
        for file in target.erl_module.to_list() + target.gleam_cache.to_list():
            lib_gleam_art_files.append(file)
        for file in target.beam_module.to_list():
            lib_beam_files.append(file)

    lib_inputs = []
    lib_path = paths.join("lib", "erlang")
    for lib_gleam_art_file in lib_gleam_art_files:
        link = ctx.actions.declare_file(paths.join(lib_path, "placeholder", GLEAM_ARTEFACTS_DIR, paths.basename(lib_gleam_art_file.path)))
        ctx.actions.symlink(
            output = link,
            target_file = lib_gleam_art_file,
        )
        lib_inputs.append(link)

    seen_beam_module = {}
    for lib_gleam_beam_file in lib_beam_files:
        beam_dir = paths.dirname(lib_gleam_beam_file.path)
        beam_base = paths.basename(lib_gleam_beam_file.path)
        if beam_base in seen_beam_module:
            if seen_beam_module.get(beam_base) != beam_dir:
                fail(
                    """Beam module {MODULE} is causing duplicates, existed 
                    as {EXISTED}, probably because of Gleam Erlang FFI has 
                    conflicting name. Note that gleam_erl_library() does not create
                    namespaces like a Gleam module. Make sure your FFI module is unique!""".format(
                        MODULE = lib_gleam_beam_file.path,
                        EXISTED = seen_beam_module.get(beam_base),
                    ),
                )
        else:
            seen_beam_module.update([(beam_base, beam_dir)])
            link = ctx.actions.declare_file(paths.join(lib_path, "placeholder", GLEAM_EBIN_DIR, paths.basename(lib_gleam_beam_file.path)))
            ctx.actions.symlink(
                output = link,
                target_file = lib_gleam_beam_file,
            )
            lib_inputs.append(link)
    return (lib_inputs, paths.join("lib", "erlang") if len(lib_inputs) > 0 else ".")

def declare_inputs(ctx, srcs, *, is_binary = False, main_module = "", main_template = ""):
    """Prepares and stages source files for the Gleam compiler.

    The resulting directory structure passed to the compiler will be:

    .
    ├── gleam.toml
    └── src/
        ├── some_source.gleam          (symlink)
        └── another_source.gleam       (symlink)

    Args:
        ctx: The Bazel rule context object.
        srcs (list): A list of source `File` objects to be prepared.
        is_binary: Whether this is a binary
        main_module: the main module.
        main_template (File): the template used for the .erl main file.

    Returns:
        A `struct` containing the prepared input files
    """
    toml_file = ctx.actions.declare_file("gleam.toml")
    ctx.actions.write(
        toml_file,
        """
name = "placeholder"
version = "0.0.0"
        """.format(package_name = _get_gleam_package_name(ctx)),
    )

    # Build inputs.
    sources = []
    for src in srcs:
        input_src = ctx.actions.declare_file(paths.join("src", strip_src_prefix(ctx, src.path)))
        ctx.actions.symlink(output = input_src, target_file = src)
        sources.append(input_src)

    binary_erl_mod = None
    if is_binary:
        binary_erl_mod = ctx.actions.declare_file(paths.join("src", main_module + "@@main.erl"))
        ctx.actions.expand_template(
            template = main_template,
            output = binary_erl_mod,
            substitutions = {
                "{PACKAGE}": main_module,
            },
        )
        sources.append(binary_erl_mod)

    sources.append(toml_file)
    return struct(
        toml_file = toml_file,
        binary_erl_mod = binary_erl_mod,
        sources = sources,
    )

def get_gleam_compiler(ctx):
    return ctx.toolchains["//gleam_tools:toolchain_type"].gleamtools.compiler


def strip_src_prefix(ctx, src):    
    """
    Strip the prefix from the source path.

    Args:
        ctx (object): The Bazel context.
        src (str): The source path.

    Returns:
        str: The source path with the prefix stripped off.
    """

    if src.startswith(ctx.bin_dir.path):
        src = src.removeprefix(ctx.bin_dir.path + "/")
    prefix = ctx.attr.strip_src_prefix
    if not prefix.endswith("/"):
        prefix = prefix + "/"
    if prefix:
        return src.removeprefix(prefix)
    return src

COMMON_ATTRS = {
    "strip_src_prefix": attr.string(
        doc = """Strip the prefix from the source path.
        Note this might break compilation because the import paths are not patched to follow.

        This is typically for external repository for moving external dependency inline with 
        the current repository.""",
    ),
    "data": attr.label_list(doc = "The data available at runtime", allow_files = True),
}
