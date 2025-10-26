"""JS binary rule for Gleam."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam/js:build.bzl", "COMMON_ATTRS", "declare_inputs", "declare_lib_files_for_dep", "declare_outputs", "get_gleam_compiler")
load("//gleam/js:provider.bzl", "GleamJsPackageInfo")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR")

def _gleam_js_binary_impl(ctx):
    """Implementation for the gleam_js_binary rule."""
    main_module = ctx.attr.main_module
    if not main_module:
        if len(ctx.files.srcs) == 1:
            # Use short_path to get path relative to workspace root
            main_module = paths.replace_extension(ctx.files.srcs[0].short_path, "").replace("/", "@")
        else:
            fail("main_module is required for gleam_js_binary when there are multiple sources.")

    inputs = declare_inputs(ctx, ctx.files.srcs)
    lib_inputs, lib_path = declare_lib_files_for_dep(ctx, ctx.attr.deps)
    outputs = declare_outputs(ctx, ctx.files.srcs)

    working_root = paths.dirname(inputs.toml_file.path)
    gleam_compiler = get_gleam_compiler(ctx)

    if len(outputs.all_files):
        ctx.actions.run_shell(
            inputs = inputs.sources + lib_inputs + [gleam_compiler],
            outputs = outputs.all_files,
            use_default_shell_env = True,
            mnemonic = "GleamJsBinaryCompile",
            command = " && ".join([
                "COMPILER=$(pwd)/%s" % gleam_compiler.path,
                "cd %s" % working_root,
                "$COMPILER compile-package --target javascript --out . --lib %s" % lib_path,
                "mv ./%s/* ." % GLEAM_ARTEFACTS_DIR,
            ]),
            env = {
                "FORCE_COLOR": "true",
            },
        )

    # Create the JS entry point
    entry_point_js = ctx.actions.declare_file(ctx.label.name + "_entry.js")
    ctx.actions.expand_template(
        template = ctx.file._main_js_tmpl,
        output = entry_point_js,
        substitutions = {
            "{main_module_path}": "./" + main_module + ".js",
        },
    )

    # Create the executable shell script
    executable = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.write(
        output = executable,
        content = "#!/bin/sh\nexec node $0.js \"$@\"\n",
        is_executable = True,
    )

    # Gather runfiles
    runfiles = ctx.runfiles(files = ctx.files.data + outputs.js_files + [entry_point_js])
    runfiles = runfiles.merge_all([dep[DefaultInfo].default_runfiles for dep in ctx.attr.deps])

    return [
        DefaultInfo(files = depset([executable]), runfiles = runfiles, executable = executable),
    ]

gleam_js_binary = rule(
    implementation = _gleam_js_binary_impl,
    executable = True,
    attrs = dict(
        COMMON_ATTRS,
        srcs = attr.label_list(mandatory = True, allow_files = [".gleam"]),
        main_module = attr.string(doc = "The Gleam module containing the main function."),
        deps = attr.label_list(providers = [GleamJsPackageInfo]),
        _main_js_tmpl = attr.label(
            default = Label("//gleam/js/templates:main.js.tmpl"),
            allow_single_file = True,
        ),
    ),
    toolchains = ["//gleam_tools:toolchain_type"],
)