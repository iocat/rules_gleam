""" Extensions for bzlmod.

"""

load("//gleam_tools:coded_versions.bzl", "VERSIONS")
load("//gleam_tools:repositories.bzl", "register_toolchains")

def _compiler_extension(module_ctx):
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

gleam = module_extension(
    implementation = _compiler_extension,
    tag_classes = {
        "toolchain": tag_class(
            attrs = {
                "version": attr.string(doc = "Versions of gleam compiler", mandatory = True, values = VERSIONS.keys()),
            },
        ),
    },
)
