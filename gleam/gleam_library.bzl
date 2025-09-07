load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam:build.bzl", "declare_inputs", "declare_lib_files_for_dep", "declare_outputs", "get_gleam_compiler")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR", "GleamErlPackageInfo")

def _gleam_library_impl(ctx):
    inputs = declare_inputs(ctx, ctx.files.srcs)
    lib_inputs, lib_path = declare_lib_files_for_dep(ctx, ctx.attr.deps)

    outputs = declare_outputs(ctx, ctx.files.srcs, is_binary = False, main_module = "")

    working_root = paths.dirname(inputs.toml_file.path)
    gleam_compiler = get_gleam_compiler(ctx)
    ctx.actions.run_shell(
        inputs = inputs.sources + lib_inputs + [gleam_compiler],
        outputs = outputs.all_files,
        use_default_shell_env = True,
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
        ctx.attr.srcs,
        ctx.attr.data,
        ctx.attr.deps,
    ):
        for target in runfiles_attr:
            transitive_runfiles.append(target[DefaultInfo].default_runfiles)
    runfiles = runfiles.merge_all(transitive_runfiles)

    return [
        DefaultInfo(files = depset(outputs.erl_mods + outputs.beam_files + outputs.cache_files), runfiles = runfiles),
        GleamErlPackageInfo(
            module_names = outputs.module_names,
            erl_module = depset(direct = outputs.erl_mods, transitive = [dep[GleamErlPackageInfo].erl_module for dep in ctx.attr.deps]),
            beam_module = depset(direct = outputs.beam_files, transitive = [dep[GleamErlPackageInfo].beam_module for dep in ctx.attr.deps]),
            gleam_cache = depset(direct = outputs.cache_files, transitive = [dep[GleamErlPackageInfo].gleam_cache for dep in ctx.attr.deps]),
        ),
    ]

# Provides GleamErlPackageInfo and DefaultInfo that includes targets that are .beam, .erl sources.
gleam_library = rule(
    implementation = _gleam_library_impl,
    attrs = {
        "srcs": attr.label_list(
            doc = "The list of gleam module files to compile under the current package.",
            mandatory = True,
            allow_files = [".gleam"],
        ),
        "data": attr.label_list(doc = "The data available at runtime", allow_files = True),
        "deps": attr.label_list(
            doc = "The list of dependent gleam modules.",
            providers = [GleamErlPackageInfo],
        ),
    },
    toolchains = [
        "//gleam_tools:toolchain_type"
    ]
)
