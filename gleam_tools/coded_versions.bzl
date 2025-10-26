"""Gleam toolchain versions and download information."""

# From https://github.com/gleam-lang/gleam/releases
VERSIONS = {
    "v1.13.0": {
        "download_url_template": "https://github.com/gleam-lang/gleam/releases/download/{version}/gleam-{version}-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:858ec1297a3f5770324fb49ed0461b0c15c29bd3fe6e636e390f2e249d86a24a",
            "aarch64-apple-darwin": "sha256:2398d1a130b1bb406bdb4a5c2cea2dc9867c677cf48b3d6c45eb200a653dbb36",
            "aarch64-unknown-linux-musl": "sha256:7f45fd9b9cce8106851fb1c476a58193e5cca456cf06efd784b39680f839da1e",
            "x86_64-unknown-linux-musl": "sha256:8b372488e5ccaa54d8acc2feb9852c9e7916e480566049edd565caa1d8c74eec",
            "aarch64-pc-windows-msvc": "sha256:55956ac4358dd14874089f9726fc133d7765a53f73921e40fc64319655a0eb63",
            "x86_64-pc-windows-msvc": "sha256:a59358ebba1abd10d50593a1dff88c0c4c16dc97566d097dc724788458794128",
        },
    },
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

# From https://github.com/cocoa-xu/otp-build/releases
ERL_VERSIONS = {
    "v27.3.3": {
        "download_url_template": "https://github.com/cocoa-xu/otp-build/releases/download/{version}/otp-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:68ad3b38d8d414d12a50e3f370017a212cbf5b229826e71b354ddfb36220a27d",
            "arm64-apple-darwin": "sha256:ff5befc0fa7d3c5fb71ec76cb05fc80041fc816064fa6594fb35a54de7e81c62",
            "aarch64-linux-gnu": "sha256:28cfb4e07825b400a2b2674be02454e3def89d922ffa476b54bba08e8918bb37",
            "x86_64-linux-gnu": "sha256:fb05e9f406fc5b4f84ecc0ab79d58a5b5be788e600b36f7c1031f83f40de8a19",
        },
    },
    "v26.2.5.3": {
        "download_url_template": "https://github.com/cocoa-xu/otp-build/releases/download/{version}/otp-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:bf222a65543f655ebc86a8cdcf26020c7f5f82969ecf38a79e65ee925ff9307d",
            "arm64-apple-darwin": "sha256:3b2a92eba72dcd5cb78e483cb322a030a16ccb3cff8e95a454341e4c6796c853",
            "aarch64-linux-gnu": "sha256:79561eac9867928c585bb425ee2ca43cc7677d3dbc5d99fc6c100f5e9ba7ecf1",
            "x86_64-linux-gnu": "sha256:8da912108c8f63a9d8e0112e26419c947943ca01a5c60967564b30ff13f86fe7",
        },
    },
    "v25.3.2.14": {
        "download_url_template": "https://github.com/cocoa-xu/otp-build/releases/download/{version}/otp-{platform}.tar.gz",
        "platforms": {
            "x86_64-apple-darwin": "sha256:d77dd7869f67e517c0db45513f33c1fe641b9561af70366a5c82785b8dab51c0",
            "arm64-apple-darwin": "sha256:d8888102c59fe19d9399c4766fc787b3cd7130b66f643379afd402ad6c8cc091",
            "aarch64-linux-gnu": "sha256:9498724487b5d62473fe6bff009fa3c1051ed4016f3f464f1acfbbfee9256b04",
            "x86_64-linux-gnu": "sha256:4b2832966d38252ce2786e8f2b8e8324426f62d4ffc1a852aabc7057ce73d3fe",
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
    minor_parts = []
    if len(splitted) > 2:
        minor_parts = splitted[2].split("-")
    minor = int(0 if len(minor_parts) == 0 else minor_parts[0])
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

_stable_erl_semvers = sorted([key for key in ERL_VERSIONS.keys() if "rc" not in key], key = get_key, reverse = True)
ERL_VERSIONS["latest"] = {k: v for k, v in ERL_VERSIONS[_stable_erl_semvers[0]].items()}
ERL_VERSIONS["latest"]["download_url_template"] = ERL_VERSIONS["latest"]["download_url_template"].format(
    version = _stable_erl_semvers[0],
    # No partial template ;(
    platform = "{platform}",
)
