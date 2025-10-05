
"""Defines the gleam_repository rule."""

def _gleam_repository_impl(_ctx):
    return [DefaultInfo()]

# Provides dummy repository that carry metadata about a Gleam external dependency.
# You generally do not need to use this target as it's an implementation details
# of rules_gleam.
gleam_repository = rule(
    implementation = _gleam_repository_impl,
    attrs = dict(
        module_name = attr.string(
            doc = "The name of the Bazel module carries this repository.",
            mandatory = True,
        ),
        gleam_modules = attr.string_list(
            doc = "The list of import paths this repository carries.",
        ),
        repo_file = attr.label(
            allow_single_file = ["REPO"],
            doc = "The REPO file that declares this repository. Needed to pull transitive module.",
        ),
    )
)
