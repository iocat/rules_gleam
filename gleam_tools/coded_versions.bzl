"""Gleam toolchain versions and download information."""

VERSIONS = {
    "v1.13.0-rc1": {
        "download_url_template": "https://github.com/gleam-lang/gleam/releases/download/{version}/gleam-{version}-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:2237e801b0d7a01df3ffd67c1de06d04045dc496ecb8d6d33c727504398f2972",
            "aarch64-apple-darwin": "sha256:a8f095d9f38413b2ad66a82b32465eaea6cf39c34927ea5a07b3aacaf0ae9d49",
            "aarch64-unknown-linux-musl": "sha256:8725b6457fe0f2fbfe044046a2941b200d8e05a11bea961da9fd17047fc3d50d",
            "x86_64-unknown-linux-musl": "sha256:3fdc4becaf5e3f6f53ed1e89527e415241a4eab7d81316aeac9a73b55da62922",
            "aarch64-pc-windows-msvc": "sha256:d2f2d92dbb5ef08a08296aa289733d5dbde43b2d7f5a73baa702b530cc58ac18",
            "x86_64-pc-windows-msvc": "sha256:29e2794b58db2aab45e074e4d4df0aa8f837a7decc137e778a4a7cbb4c84e714",
        },
    },
    "v1.12.0": {
        "download_url_template": "https://github.com/gleam-lang/gleam/releases/download/{version}/gleam-{version}-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:fca2a8cf5424cfa096710c191bbf8c02b635ed0e6c0e2ed38b95eba5e8df0302",
            "aarch64-apple-darwin": "sha256:885396e845fbbc014253dd95da493439b785641022fb92c45657b2b936cf317f",
            "aarch64-unknown-linux-musl": "sha256:eb9707e42e452d5bbfb48b12d81571a098238a5afb4b7d16c1d07df7089e2bde",
            "x86_64-unknown-linux-musl": "sha256:039a87bd7294d3cfd2425f56e8ffef508b170ecec42e760806833fb1e0319d49",
            "aarch64-pc-windows-msvc": "sha256:44e820854af78a4a0daaa3282524e87f12c418bd76f1999869ec10fb0b6aec76",
            "x86_64-pc-windows-msvc": "sha256:24e89e4f8bdb5a80ed214f6d57d9118a43ba3a339a8d1831b2975808a4988648",
        },
    },
    "v1.11.1": {
        "download_url_template": "https://github.com/gleam-lang/gleam/releases/download/{version}/gleam-{version}-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:f250284e4998f4e7d274d60f8d4f9d5d5cf1519a86bd21065ba77dbff6065054",
            "aarch64-apple-darwin": "sha256:a2d492592d86539b1f51ae81f344dff71870b23479bf4d00d78512c6f9eaf1e5",
            "aarch64-unknown-linux-musl": "sha256:22d4c9299f57f712210df7f48d1c84f99f13abfd9e067a6b290abfdab384161d",
            "x86_64-unknown-linux-musl": "sha256:31649fab05f982c51d7553cc4a2d3e615ee49348f7162fe75d48f71331618033",
            "aarch64-pc-windows-msvc": "sha256:a409e9502dc20a427c118bb9ce64afeb8b78ab5652bf6cce3759ab7327ba3c93",
            "x86_64-pc-windows-msvc": "sha256:1cb54473bed54d74e213f993a7f269ed27afa761d693be681b9bf1a3e618d4d5",
        },
    },
}

# MAKE THE LATEST VERSION.

def get_key(semver):
    """
    Get the version, major, minor, and release candidate count from the given semantic version.

    Args:
        semver: The semantic version.

    Returns:
        version (int): The version.
        major (int): The major.
        minor (int): The minor.
        rc (None): The release candidate count, if any.
    """
    splitted = semver[1:].split(".")
    version = int(splitted[0])
    major = int(splitted[1])
    minor_parts = splitted[2].split("-")
    minor = int(minor_parts[0])
    rc = None
    if len(minor_parts) > 1:
        rc = minor_parts[1]

    return (version, major, minor, rc)

_stable_semvers = sorted([key for key in VERSIONS.keys() if "rc" not in key], key = get_key, reverse = True)

VERSIONS["latest"] = {k: v for k, v in VERSIONS[_stable_semvers[0]].items()}
VERSIONS["latest"]["download_url_template"] = VERSIONS["latest"]["download_url_template"].format(
    version = _stable_semvers[0],
    # No partial template ;(
    platform = "{platform}",
)
