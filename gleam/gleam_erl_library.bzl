# Erl library for interoping with erlang via
# 
# @external(erlang, "erl", "function")
load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam:build.bzl", "declare_inputs", "declare_lib_files_for_dep", "declare_outputs", "get_gleam_compiler")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR", "GleamErlPackageInfo")

def _gleam_erl_library_impl(ctx):
    inputs = declare_inputs(ctx, ctx.files.srcs)
    outputs = declare_outputs(ctx, ctx.files.srcs, is_binary = False, main_module = "")
    _, lib_path = declare_lib_files_for_dep(ctx, []) # no deps, we need the path.

    working_root = paths.dirname(inputs.toml_file.path)
    gleam_compiler = get_gleam_compiler(ctx)
    if len(outputs.all_files):
        ctx.actions.run_shell(
            inputs = inputs.sources + [gleam_compiler],
            outputs = outputs.all_files,
            use_default_shell_env = True,
            mnemonic="GleamErlLibraryCompile",
            command = """
                COMPILER="$(pwd)/%s" &&
                cd %s &&
                $COMPILER compile-package --package '.' --target erlang --out '.' --lib %s &&
                mv ./%s/* ./ &&
                mv ./ebin/* ./
            """ % (gleam_compiler.path, working_root, lib_path, GLEAM_ARTEFACTS_DIR),
        )

    # Accumulate runfiles.
    runfiles = ctx.runfiles(files = ctx.files.data + outputs.beam_files)
    transitive_runfiles = []
    for runfiles_attr in (
        ctx.attr.data,
    ):
        for target in runfiles_attr:
            transitive_runfiles.append(target[DefaultInfo].default_runfiles)
    runfiles = runfiles.merge_all(transitive_runfiles)

    return [
        DefaultInfo(files = depset(outputs.erl_mods + outputs.beam_files + outputs.cache_files), runfiles = runfiles),
        GleamErlPackageInfo(
            module_names = outputs.module_names,
            erl_module = depset(direct = outputs.erl_mods),
            beam_module = depset(direct = outputs.beam_files),
            gleam_cache = depset(direct = outputs.cache_files),
        ),
    ]

# Provides GleamErlPackageInfo and DefaultInfo that includes targets that are .beam, .erl sources.
gleam_erl_library = rule(
    implementation = _gleam_erl_library_impl,
    attrs = {
        "srcs": attr.label_list(
            doc = "The list of gleam module files to compile under the current package.",
            mandatory = True,
            allow_files = [".erl"],
        ),
        "data": attr.label_list(doc = "The data available at runtime", allow_files = True),
    },
    toolchains = [
        "//gleam_tools:toolchain_type"
    ]
)
