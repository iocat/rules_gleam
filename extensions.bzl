""" Extensions for bzlmod.

"""

load("//gleam_hex:repositories.bzl", "gleam_hex_repositories")
load("//gleam_tools:coded_versions.bzl", "VERSIONS")
load("//gleam_tools:repositories.bzl", "register_toolchains")
load("//internal:common.bzl", "extension_metadata")

def _compiler_extension(module_ctx):
    direct_deps = {}

    non_default_registrations = {}
    for mod in module_ctx.modules:
        if mod.name == "rules_gleam":
            continue
        for toolchain in mod.tags.toolchain:
            non_default_registrations[toolchain.version] = True

    default_toolchains = [mod.tags.toolchain for mod in module_ctx.modules if mod.name == "rules_gleam"]

    selected = None
    if len(non_default_registrations.keys()) >= 1:
        selected = sorted(non_default_registrations.keys(), reverse = True)[0]
        # buildifier: disable=print
        # print("NOTE: gleam toolchain has multiple redundancy, for versions {}, selected {}".format(", ".join(non_default_registrations.keys()), selected))

    else:
        selected = default_toolchains[0][0].version

    register_toolchains(
        version = selected,
    )
    direct_deps["gleam_toolchains"] = True

    gleam_toml = None
    for mod in module_ctx.modules:
        # if mod.name == "rules_gleam":
        #     continue
        for gleam_deps in mod.tags.deps:
            if gleam_toml != None:
                fail("There should be one gleam.toml defined, existing declaration at %s" % module_ctx.path(gleam_toml))
            gleam_toml = gleam_deps.gleam_toml

    hex_modules = []
    hex_modules = gleam_hex_repositories(
        module_ctx,
        gleam_toml = gleam_toml,
    )
    for hex_mod in hex_modules:
        direct_deps[hex_mod] = True
    
    first_module = module_ctx.modules[0]
    if first_module.is_root and first_module.name == "rules_gleam":
        direct_deps["gleam_hex_repositories_config"] = True

    return extension_metadata(
        module_ctx,
        root_module_direct_deps = direct_deps.keys(),
        root_module_direct_dev_deps = [],
        reproducible = True,
    )

gleam = module_extension(
    implementation = _compiler_extension,
    tag_classes = {
        "toolchain": tag_class(
            attrs = {
                "version": attr.string(doc = "Versions of gleam compiler", mandatory = True, values = VERSIONS.keys() + ["latest"]),
            },
        ),
        "deps": tag_class(
            attrs = {
                "gleam_toml": attr.label(
                    allow_single_file = [".toml"],
                    mandatory = True,
                    doc = "The gleam.toml file to be pulling deps from.",
                ),
            },
        ),
    },
)
