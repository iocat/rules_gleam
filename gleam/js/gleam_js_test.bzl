"""JS test rule for Gleam."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam/js:build.bzl", "COMMON_ATTRS", "declare_inputs", "declare_lib_files_for_dep", "declare_outputs", "get_gleam_compiler")
load("//gleam/js:provider.bzl", "GleamJsPackageInfo")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR")

def _gleam_js_test_impl(ctx):
    """Implementation for the gleam_js_test rule."""
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
            mnemonic = "GleamJsTestCompile",
            command = " && ".join([
                "COMPILER=$(pwd)/%s" % gleam_compiler.path,
                "cd %s" % working_root,
                "$COMPILER test --target javascript --out . --lib %s" % lib_path,
                "mv ./%s/* ." % GLEAM_ARTEFACTS_DIR,
            ]),
            env = {
                "FORCE_COLOR": "true",
            },
        )

    # Create the executable shell script
    executable = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.write(
        output = executable,
        # The test runner is output as 'test.js' by `gleam test`
        content = "#!/bin/sh\nexec node test.js \"$@\"\n",
        is_executable = True,
    )

    # Gather runfiles
    runfiles = ctx.runfiles(files = ctx.files.data + outputs.js_files)
    runfiles = runfiles.merge_all([dep[DefaultInfo].default_runfiles for dep in ctx.attr.deps])

    return [
        DefaultInfo(files = depset([executable]), runfiles = runfiles, executable = executable),
    ]

gleam_js_test = rule(
    implementation = _gleam_js_test_impl,
    test = True,
    attrs = dict(
        COMMON_ATTRS,
        srcs = attr.label_list(
            mandatory = True,
            allow_files = [".gleam"],
        ),
        deps = attr.label_list(
            providers = [GleamJsPackageInfo],
            # Automatically include gleeunit for tests.
            default = ["@hex_gleeunit//gleam/js"],
        ),
    ),
    toolchains = ["//gleam_tools:toolchain_type"],
)