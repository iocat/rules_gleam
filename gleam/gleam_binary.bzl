load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam:build.bzl", "declare_inputs", "declare_lib_files_for_dep", "declare_outputs", "get_gleam_compiler")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR", "GleamErlPackageInfo")

def _gleam_binary_impl(ctx):
    main_module = ctx.attr.main_module
    if main_module == "" and len(ctx.files.srcs) == 1:
        main_module = paths.replace_extension(ctx.files.srcs[0].path, "").replace("/", "@")
    if main_module == "":
        fail("Main module is not provided. Please provide one via main_module attribute", "main_module")

    inputs = declare_inputs(ctx, ctx.files.srcs, is_binary = True, main_module = main_module, main_template = ctx.file._main_erl)
    lib_inputs, lib_path = declare_lib_files_for_dep(ctx, ctx.attr.deps)

    outputs = declare_outputs(ctx, ctx.files.srcs, is_binary = True, main_module = main_module)

    working_root = paths.dirname(inputs.toml_file.path)
    gleam_compiler = get_gleam_compiler(ctx)
    ctx.actions.run_shell(
        inputs = inputs.sources + lib_inputs + [gleam_compiler],
        outputs = outputs.all_files_include_binary,
        use_default_shell_env = True,
        command = """
            COMPILER="$(pwd)/%s" &&
            cd %s &&
            $COMPILER compile-package --package '.' --target erlang --out '.' --lib %s &&
            mv ./%s/* ./ &&
            mv ./ebin/* ./
        """ % (gleam_compiler.path, working_root, lib_path, GLEAM_ARTEFACTS_DIR),
    )

    erl_mod_depset = depset(direct = outputs.erl_mods + [inputs.binary_erl_mod], transitive = [dep[GleamErlPackageInfo].erl_module for dep in ctx.attr.deps])

    # Manifest
    ctx.actions.expand_template(
        template = ctx.file._app_manifest,
        output = outputs.beam_app_manifest,
        substitutions = {
            "{PACKAGE}": main_module,
            "{APPS_COMMA_SEP}": "",
            "{MODS_COMMA_SEP}": ",\n\t\t\t\t".join([paths.replace_extension(paths.basename(erl_mod.path), "") for erl_mod in erl_mod_depset.to_list()]),
        },
    )

    # Binary
    # Entry point .sh
    ctx.actions.expand_template(
        template = ctx.file._sh_entrypoint,
        output = outputs.output_entry_point,
        is_executable = True,
        substitutions = {
            "{PACKAGE}": main_module,
        },
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

    # Symlinks needed .beam dependencies for target executable to run..
    dep_beam_symlinks = []
    for dep_beam_modules in [dep[GleamErlPackageInfo].beam_module for dep in ctx.attr.deps]:
        for dep_beam_module in dep_beam_modules.to_list():
            dep_beam = ctx.actions.declare_file(paths.basename(dep_beam_module.path))
            ctx.actions.symlink(output = dep_beam, target_file = dep_beam_module)
            dep_beam_symlinks.append(dep_beam)
    runfiles = runfiles.merge(ctx.runfiles(files = dep_beam_symlinks))

    return [
        DefaultInfo(files = depset(outputs.beam_files + [outputs.output_entry_point, outputs.beam_app_manifest]), default_runfiles = runfiles, executable = outputs.output_entry_point),
        GleamErlPackageInfo(
            module_names = outputs.module_names,
            erl_module = erl_mod_depset,
            beam_module = depset(direct = outputs.beam_files, transitive = [dep[GleamErlPackageInfo].beam_module for dep in ctx.attr.deps]),
            gleam_cache = depset(direct = outputs.cache_files, transitive = [dep[GleamErlPackageInfo].gleam_cache for dep in ctx.attr.deps]),
        ),
    ]

# Provides GleamErlPackageInfo and DefaultInfo that includes targets that are .beam, .erl sources.
gleam_binary = rule(
    implementation = _gleam_binary_impl,
    executable = True,
    attrs = {
        "srcs": attr.label_list(
            doc = "The list of gleam module files to compile under the current package.",
            mandatory = True,
            allow_files = [".gleam"],
        ),
        "main_module": attr.string(doc = "The module name containing the main function. Must match the file name of one of the source."),
        "data": attr.label_list(doc = "The data available at runtime"),
        "deps": attr.label_list(
            doc = "The list of dependent gleam modules.",
            providers = [GleamErlPackageInfo],
        ),
        "_main_erl": attr.label(
            default = "//gleam/templates:[[main_module]]@@main.erl",
            allow_single_file = True,
        ),
        "_app_manifest": attr.label(
            default = "//gleam/templates:[[main_module]].app",
            allow_single_file = True,
        ),
        "_sh_entrypoint": attr.label(
            default = "//gleam/templates:[[main_module]].sh",
            allow_single_file = True,
        ),
    },
    toolchains = [
        "//gleam_tools:toolchain_type"
    ]
)
