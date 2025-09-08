# Accept gleam.toml, generate manifest.toml

def _gleam_hex_repository(repository_ctx):
    # Invoke gazelle update-repos to create BUILDs target for
    # an external dep.
    pass

gleam_hex_repository = repository_rule(
    _gleam_hex_repository,
    doc = """
    Creates a repository with the Gleam hex repository.
    """,
    attrs = {}
)

def _gleam_hex_repositories(repository_ctx):
    pass

gleam_hex_repositories = repository_rule (
    _gleam_hex_repositories,
    doc = """
    From the manifest.toml file, declare deps for the repositories.
    """
)