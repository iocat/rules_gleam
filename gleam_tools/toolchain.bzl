GleamToolsInfo = provider(
    doc = "Information about Gleam tools needed to complete compilation, specifically erl/js/gleam tools",
    fields = ["compiler", "erl_command"],
)

def _gleam_toolchain_impl(ctx):
    return platform_common.ToolchainInfo(
        gleamtools = GleamToolsInfo(
            # erl_command = ctx.attr.erl_command,
            compiler = ctx.file.compiler,
        ),
    )

gleam_toolchain = rule(
    implementation = _gleam_toolchain_impl,
    attrs = {
        # "erl_command": attr.label(
        #     executable = True,
        #     mandatory = True,
        #     cfg = "target",
        #     doc = "The erl tool to execute erlang binary.",
        # ),
        "compiler": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "The gleam compiler tool to compile gleam with.",
            allow_single_file = True,
        ),
    },
    doc = "Defines the gleam tools to compile and execute Gleam.",
    provides = [platform_common.ToolchainInfo],
)
