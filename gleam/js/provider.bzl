"""Gleam JS provider."""

# Information about a compiled Gleam JS package.
GleamJsPackageInfo = provider(
    "GleamJsPackageInfo",
    fields = {
        "module_names": "depset of strings of module names",
        "js_module": "depset of .js files",
        "gleam_cache": "depset of gleam cache files",
        "strip_src_prefix": "string of prefix to strip from source files",
    },
)
