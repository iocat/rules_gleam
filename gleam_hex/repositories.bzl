# Accept gleam.toml, generate manifest.toml

load("//internal:common.bzl", "env_execute")

def _gleam_hex_repository(ctx):
    tar_package_url = _hex_tar_url(ctx.attr.module_name, ctx.attr.version)
    ctx.download_and_extract(
        url = tar_package_url,
        sha256 = ctx.attr.checksum,
        type = "tar",
        output = "hexpkg",
    )
    ctx.file("BUILD.bazel", """
package(
    default_visibility = ["//visibility:public"]
)

load("@gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    gazelle = "@rules_gleam//gazelle",
)

# gleam_library(
#     name = "gleam_stdlib",
#     srcs = [],
# )
    """)
    ctx.extract("hexpkg/contents.tar.gz", output = ".", strip_prefix = "src")
    ctx.delete("hexpkg")

    # ctx.rename("src", ".")
    # print(src)
    # if not src.exists:
    #     fail("There is no src directory")
    # result = env_execute(ctx, ["cp", "-R", src, "%s/.."])
    # print(result.stderr)
    print(env_execute(ctx, ["tree", "."]).stderr)
    # Invoke gazelle update-repos to create BUILDs target for
    # an external dep.
    pass

gleam_hex_repository = repository_rule(
    _gleam_hex_repository,
    doc = """
    Creates a repository with the Gleam hex repository.
    """,
    attrs = {
        "module_name": attr.string(doc = "The Hex package module name to download. This will be used as the root import path."),
        "checksum": attr.string(doc = "The checksum of the outerpackage we got from the manifest.toml file"),
        "version": attr.string(doc = "Semver version for the module"),
        "otp_app": attr.string(doc = "The otp_app from the module"),
    },
)

def _hex_tar_url(module, version):
    return "https://repo.hex.pm/tarballs/{module}-{version}.tar".format(
        module = module,
        version = version,
    )

# A macro (like a repository rule) to download hex repositories.
def gleam_hex_repositories(module_ctx, *, gleam_toml, _get_hex_repos = Label("@rules_gleam_internal_tools//:bin/get_hex_repos")):
    """Creates repositories with the Gleam hex repository.

    Args:
        module_ctx (module_ctx): The module context.
        *: Additional arguments.
        gleam_toml (Label): The path to the gleam.toml file to be included.
        _get_hex_repos (Label): The path to the get_hex_repos script to translate the manifest.toml to json
          that bazel can consume.

    Returns:
        A list of module names.
    """
    file = module_ctx.path(gleam_toml).dirname.get_child("manifest.toml")
    if not file.exists:
        fail("Could not find manifest.toml at the gleam.toml location %s" % gleam_toml.path)

    dependencies = module_ctx.execute([module_ctx.path(_get_hex_repos), "--manifest", file])
    if dependencies.return_code:
        fail("failed to read manifest.toml file: %s" % dependencies.stderr)

    dep_json = json.decode(dependencies.stdout)
    repos = dep_json.get("repos", default = [])

    for repo in repos:
        gleam_hex_repository(
            name = repo.get("module_name"),
            module_name = repo.get("module_name"),
            checksum = repo.get("checksum"),
            version = repo.get("version"),
            otp_app = repo.get("otp_app"),
        )

    return [repo.get('module_name') for repo in repos]
