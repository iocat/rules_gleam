GleamToolsInfo = provider(
    doc = "Information about Gleam tools needed to complete compilation, specifically erl/js/gleam tools",
    fields = ["compiler", "js_prelude"],
)

GleamErlangToolsInfo = provider(
    doc = "Information about Erlang tools needed to compiler Gleam",
    fields = ["escript", "escript_label", "erl", "erlc", "otp"],
)

def _gleam_toolchain_impl(ctx):
    return platform_common.ToolchainInfo(
        gleamtools = GleamToolsInfo(
            # erl_command = ctx.attr.erl_command,
            compiler = ctx.file.compiler,
            js_prelude = ctx.file.js_prelude,
        ),
    )

gleam_toolchain = rule(
    implementation = _gleam_toolchain_impl,
    attrs = {
        "compiler": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "The gleam compiler tool to compile gleam with.",
            allow_single_file = True,
        ),
        "js_prelude": attr.label(
            allow_single_file = [".mjs"],
            mandatory = True,
            doc = "The js prelude provided by gleam compiler.",
        )
    },
    doc = "Defines the gleam tools to compile and execute Gleam.",
    provides = [platform_common.ToolchainInfo],
)

def _gleam_erlang_toolchain_impl(ctx):
    return platform_common.ToolchainInfo(
        gleamerlangtools = GleamErlangToolsInfo(
            escript = ctx.file.escript,
            erl = ctx.file.erl,
            erlc = ctx.file.erlc,
            escript_label = ctx.attr.escript,
            otp = ctx.files.otp,
        ),
    )

gleam_erlang_toolchain = rule(
    implementation = _gleam_erlang_toolchain_impl,
    attrs = {
        "escript": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "The escript tool to transpile Gleam and execute erlang binary.",
            allow_single_file = True,
        ),
         "erl": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "The erl tool to transpile Gleam and execute erlang binary.",
            allow_single_file = True,
        ),
         "erlc": attr.label(
            executable = True,
            mandatory = True,
            cfg = "exec",
            doc = "The erlc tool to transpile Gleam and execute erlang binary.",
            allow_single_file = True,
        ),
        "otp": attr.label(
            allow_files = True,
            doc = "Everything under OTP",
        )
    },
    doc = "Defines the Erlang tools to compile and execute Gleam.",
    provides = [platform_common.ToolchainInfo],
)