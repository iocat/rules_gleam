
load("@bazel_features//:features.bzl", "bazel_features")

def env_execute(ctx, arguments, environment = {}, **kwargs):
    """Executes a command in for a repository rule.
    
    It prepends "env -i" to "arguments" before calling "ctx.execute".
    Variables that aren't explicitly mentioned in "environment"
    are removed from the environment. This should be preferred to "ctx.execute"
    in most situations.

    Args:
      ctx: The Bazel context.
      arguments: The command to execute.
      environment: A dictionary of environment variables to set.
      **kwargs: Additional keyword arguments to pass to ctx.execute.

    Returns:
        The result of the execution.
    """
    if ctx.os.name.startswith("windows"):
        return ctx.execute(arguments, environment = environment, **kwargs)
    env_args = ["env", "-i"]
    environment = dict(environment)
    for var in ["TMP", "TMPDIR"]:
        if var in ctx.os.environ and not var in environment:
            environment[var] = ctx.os.environ[var]
    for k, v in environment.items():
        env_args.append("%s=%s" % (k, v))
    arguments = env_args + arguments
    return ctx.execute(arguments, **kwargs)

def executable_extension(ctx):
    extension = ""
    if ctx.os.name.startswith("windows"):
        extension = ".exe"
    return extension

def watch(ctx, path):
    # Versions of Bazel that have ctx.watch may no longer explicitly watch
    # labels on which ctx.path is called and/or labels in attributes. Do so
    # explicitly here, duplicate watches are no-ops.
    if hasattr(ctx, "watch"):
        ctx.watch(path)

def extension_metadata(
        module_ctx,
        *,
        root_module_direct_deps = None,
        root_module_direct_dev_deps = None,
        reproducible = False):
    """Returns the extension metadata for the given module context.

    Args:
        module_ctx: The module context.
        root_module_direct_deps: The direct dependencies of the root module.
        root_module_direct_dev_deps: The direct development dependencies of the root module.
        reproducible: Whether the metadata should be reproducible.

    Returns:
        The extension metadata, or None if the module context does not have an extension metadata attribute.
    """
    if not hasattr(module_ctx, "extension_metadata"):
        return None
    metadata_kwargs = {}
    if bazel_features.external_deps.extension_metadata_has_reproducible:
        metadata_kwargs["reproducible"] = reproducible
    return module_ctx.extension_metadata(
        root_module_direct_deps = root_module_direct_deps,
        root_module_direct_dev_deps = root_module_direct_dev_deps,
        **metadata_kwargs
    )