# Accept gleam.toml, generate manifest.toml

load("//internal:common.bzl", "env_execute", "watch", "executable_extension")

TEMP_DIR = "__hexpkg"

def _gleam_hex_repository(ctx):
    tar_package_url = _hex_tar_url(ctx.attr.module_name, ctx.attr.version)
    ctx.report_progress("Downloading Repo")
    ctx.download_and_extract(
        url = tar_package_url,
        sha256 = ctx.attr.checksum,
        type = "tar",
        output = TEMP_DIR,
    )
    ctx.file("BUILD", """load("@gazelle//:def.bzl", "gazelle")
load("@rules_gleam//gleam:defs.bzl", "gleam_library")

package(
    default_visibility = ["//visibility:public"],
)

gazelle(
    name = "gazelle",
    gazelle = "@rules_gleam//gazelle",
)

""")
    ctx.extract("%s/contents.tar.gz" % TEMP_DIR, output = ".", strip_prefix = "src")
    ctx.delete(TEMP_DIR)

    _gazelle_label = Label("@rules_gleam_internal_tools//:bin/gazelle{}".format(executable_extension(ctx)))
    _gazelle = ctx.path(_gazelle_label)
    watch(ctx, _gazelle)

    cmd = [
        _gazelle,
        "-mode",
        "fix",
        "-gleam_external_repo",
        "-repo_root",
        ctx.path(""),
    ]
    cmd.append(ctx.path(""))
    ctx.report_progress("Runnning Gazelle")

    result = env_execute(ctx, cmd)
    if result.return_code:
        fail("failed to generate BUILD files for %s: %s. err: %s" % (
            ctx.attr.module_name,
            result.stdout,
            result.stderr,
        ))

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

    return [repo.get("module_name") for repo in repos]
