PLATFORMS = {
    "x86_64-apple-darwin": struct(
        compatible_with = [
            "@platforms//os:macos",
            "@platforms//cpu:x86_64",
        ],
    ),
    "aarch64-apple-darwin": struct(
        compatible_with = [
            "@platforms//os:macos",
            "@platforms//cpu:aarch64",
        ],
    ),
    # Technically unsupported, since cli builds were in linux.
    "aarch64-pc-windows-msvc": struct(
        compatible_with = [
            "@platforms//os:windows",
            "@platforms//cpu:aarch64",
        ],
    ),
    "x86_64-pc-windows-msvc": struct(
        compatible_with = [
            "@platforms//os:windows",
            "@platforms//cpu:x86_64",
        ],
    ),
    "aarch64-unknown-linux-musl": struct(
        compatible_with = [
            "@platforms//os:linux",
            "@platforms//cpu:aarch64",
        ],
    ),
    "x86_64-unknown-linux-musl": struct(
        compatible_with = [
            "@platforms//os:linux",
            "@platforms//cpu:x86_64",
        ],
    ),
}

ERL_PLATFORMS = {
    # No windows support yet, relying on https://github.com/cocoa-xu/otp-build/releases
    # without a windows build.
    "x86_64-apple-darwin": struct(
        compatible_with = [
            "@platforms//os:macos",
            "@platforms//cpu:x86_64",
        ],
    ),
    "arm64-apple-darwin": struct(
        compatible_with = [
            "@platforms//os:macos",
            "@platforms//cpu:aarch64",
        ],
    ),
    "aarch64-linux-gnu": struct(
        compatible_with = [
            "@platforms//os:linux",
            "@platforms//cpu:aarch64",
        ],
    ),
    "x86_64-linux-gnu": struct(
        compatible_with = [
            "@platforms//os:linux",
            "@platforms//cpu:x86_64",
        ],
    ),
}