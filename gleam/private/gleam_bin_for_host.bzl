load("@platforms//host:constraints.bzl", "HOST_CONSTRAINTS")
load("//gleam:build.bzl", "get_gleam_compiler")

def _ensure_target_cfg(ctx):
    if "-exec" in ctx.bin_dir.path or "/host/" in ctx.bin_dir.path:
        fail("//gleam is only meant to be used with 'bazel run', not as a tool. " +
             "If you need to use it as a tool (e.g. in a genrule), please " +
             "open an issue for your use case.")

def _gleam_bin_for_host_impl(ctx):
    """Exposes the go binary of the current Go toolchain for the host."""
    _ensure_target_cfg(ctx)

    compiler = get_gleam_compiler(ctx)
    bin_link = ctx.actions.declare_file("gleam")
    ctx.actions.symlink(output = bin_link, target_file = compiler, is_executable = True)
    return [
        DefaultInfo(
            files = depset([compiler]),
            executable = bin_link,
            runfiles = ctx.runfiles(files = [bin_link, compiler])
        ),
    ]

gleam_bin_for_host = rule(
    implementation = _gleam_bin_for_host_impl,
    toolchains = ["//gleam_tools:toolchain_type"],
    # Resolve a toolchain that runs on the host platform.
    exec_compatible_with = HOST_CONSTRAINTS,
    executable = True,
)