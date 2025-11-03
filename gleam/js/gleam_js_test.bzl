"""JS binary rule for Gleam."""

load("@aspect_rules_js//js:providers.bzl", "js_info")
load("@bazel_skylib//lib:paths.bzl", "paths")
load("//gleam:provider.bzl", "GLEAM_ARTEFACTS_DIR")
load("//gleam/js:build.bzl", "COMMON_ATTRS", "declare_inputs", "declare_libs", "declare_outputs", "declare_runtime_mjs", "get_gleam_compiler", "get_js_prelude")
load("//gleam/js:provider.bzl", "GleamJsPackageInfo")

def _gleam_js_test_impl(ctx):
    """Implementation for the gleam_js_test rule."""
    package_dir = ctx.label.package
    main_module = paths.join(package_dir, "gleam_test")
    test_main = ctx.actions.declare_file("gleam_test.gleam")
    dots_to_root = "./" if (
        paths.dirname(main_module) == "."
    ) else "".join(["../" for dot in paths.dirname(main_module).split("/")])
    ctx.actions.expand_template(
        template = ctx.file._main_test_module_tmpl,
        output = test_main,
        substitutions = {
            "{DOTS_TO_ROOT}": dots_to_root,
            "{MODULES_UNDER_TEST}": ", ".join(["\"" + paths.replace_extension(src.path, "") + "\"" for src in ctx.files.srcs]),
            "{TIMEOUT_SECS}": "300" if ctx.attr.timeout == "eternal" else "60" if ctx.attr.timeout == "long" else "30" if ctx.attr.timeout == "moderate" else "10",
        },
    )

    if not main_module and len(ctx.files.srcs) == 1:
        main_module = paths.replace_extension(ctx.files.srcs[0].path, "")

    if not main_module:
        fail("Main module is not provided. Please provide one via main_module attribute", "main_module")

    inputs = declare_inputs(ctx, ctx.files.srcs + [test_main])
    lib_inputs, lib_path = declare_libs(ctx, ctx.attr.deps)
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
                "$COMPILER compile-package --target javascript --package '.' --out '.' --lib %s --javascript-prelude %s" % (lib_path, get_js_prelude(ctx).path),
                # Move cache outside
                "mv ./%s/* ." % GLEAM_ARTEFACTS_DIR,
            ] + ([
                # Move mjs outside
                "mv ./%s/* ." % ctx.label.package,
            ] if ctx.label.package else [])),
            env = {
                "FORCE_COLOR": "true",
            },
        )

    # Accumulate runfiles.
    runfiles = ctx.runfiles(files = ctx.files.data)
    runfiles = runfiles.merge_all([dep[DefaultInfo].default_runfiles for dep in ctx.attr.deps])

    # Create the entry point script
    main_module_filename = main_module.split("@")[-1]
    output_entry_point = ctx.actions.declare_file(main_module_filename + "_main.mjs")
    main_js_path = "./" + main_module_filename + ".mjs"

    ctx.actions.expand_template(
        template = ctx.file._main_js,
        output = output_entry_point,
        is_executable = True,
        substitutions = {
            "{MAIN_JS}": main_js_path,
        },
    )

    # Executable
    output_sh_executable = ctx.actions.declare_file(main_module_filename + "_main.sh")
    ctx.actions.expand_template(
        template = ctx.file._main_js_sh,
        output = output_sh_executable,
        is_executable = True,
        substitutions = {
            "{MAIN_JS}": paths.join(ctx.label.package, output_entry_point.basename),
        },
    )

    (js_modules_depset, root_symlinks) = declare_runtime_mjs(ctx, outputs.js_files + [output_entry_point], ctx.attr.deps)

    # Accumulate runfiles.
    runfiles = ctx.runfiles(files = ctx.files.data + js_modules_depset.to_list(), root_symlinks = root_symlinks)
    transitive_runfiles = []
    for dep in ctx.attr.deps:
        if DefaultInfo in dep:
            transitive_runfiles.append(dep[DefaultInfo].default_runfiles)
    runfiles = runfiles.merge_all(transitive_runfiles)

    return [
        DefaultInfo(files = depset(outputs.js_files + [output_entry_point]), default_runfiles = runfiles, executable = output_sh_executable),
        GleamJsPackageInfo(
            module_names = outputs.module_names,
            js_module = js_modules_depset,
            gleam_cache = depset(direct = outputs.cache_files, transitive = [dep[GleamJsPackageInfo].gleam_cache for dep in ctx.attr.deps if GleamJsPackageInfo in dep]),
            strip_src_prefix = ctx.attr.strip_src_prefix,
            target_name = ctx.label.name,
        ),
        js_info(
            target = ctx.label,
            sources = js_modules_depset,
        ),
    ]

def gleam_js_test(size = "small", timeout = None, **kwargs):
    """Gleam test rule.

    Compiles Gleam source files into an Javascript Deno executable that can be run.

    Args:
        size: The size of the test. One of "small", "medium", "large", "enormous". Defaults to "small".
        timeout: The timeout for the test. One of "short", "moderate", "long", "eternal". Defaults to None.
        **kwargs: Other keyword arguments passed to the gleam_js_test rules.

    Returns:
        A Bazel target.
    """
    deps = kwargs.get("deps", [])

    deps.append("@rules_gleam//gleam/gleeunit:gleeunit_gleam_js")
    kwargs["deps"] = deps
    return _gleam_js_test(size = size, timeout = timeout, **kwargs)

_gleam_js_test = rule(
    implementation = _gleam_js_test_impl,
    test = True,
    attrs = dict(
        COMMON_ATTRS,
        srcs = attr.label_list(
            doc = "The list of gleam module files to compile under the current package.",
            mandatory = True,
            allow_files = [".gleam"],
        ),
        main_module = attr.string(doc = "The module name containing the main function. Must match the file name of one of the source. Default to the module at srcs[0]"),
        deps = attr.label_list(
            doc = "The list of dependent gleam modules.",
            providers = [GleamJsPackageInfo],
        ),
        _main_test_module_tmpl = attr.label(
            default = "//gleam/js/templates:gleam_test.gleam",
            allow_single_file = True,
            doc = "The gleam file template containing the test main function.",
        ),
        _main_js = attr.label(
            default = "//gleam/js/templates:main_js.mjs",
            allow_single_file = [".mjs"],
        ),
        _main_js_sh = attr.label(
            default = "//gleam/js/templates:main_js.sh",
            allow_single_file = [".sh"],
        ),
    ),
    toolchains = [
        "//gleam_tools:toolchain_type",
    ],
)
