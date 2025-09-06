GleamErlPackageInfo = provider(
    "Provide compiled Gleam package.",
    fields = {
        "module_names": "the names of the modules this package provided, .e.g hello@world",
        "erl_module": "depset of Erlang compilation output files.",
        "beam_module": "depset of Beam module compilation output files.",
        "gleam_cache": "depset of Gleam cache and cache_meta files for this module.",
    },
)

GleamErlBinaryPackageInfo = provider(
    "Provide compiled Gleam binary package ready for execute.",
    fields = {
        "module_names": "the names of the modules this package provided, .e.g hello@world",
        "erl_module": "depset of Erlang compilation output files.",
        "beam_module": "depset of Beam module compilation output files.",
        "executable": "depset of a Beam executable.",
        "beam_manifest": "depset of a Beam app manifest file",
    },
)

GLEAM_ARTEFACTS_DIR = "_gleam_artefacts"
GLEAM_EBIN_DIR = "ebin"
