#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

# Argument provided by reusable workflow caller, see
# https://github.com/bazel-contrib/.github/blob/d197a6427c5435ac22e56e33340dff912bc9334e/.github/workflows/release_ruleset.yaml#L72
TAG=$1
# The prefix is chosen to match what GitHub generates for source archives
PREFIX="rules_gleam-${TAG:1}"
ARCHIVE="rules_gleam-$TAG.tar.gz"
ARCHIVE_TMP=$(mktemp)

# NB: configuration for 'git archive' is in /.gitattributes
git archive --format=tar --prefix=${PREFIX}/ ${TAG} >$ARCHIVE_TMP

############
# Patch up the archive to have integrity hashes for built binaries that we downloaded in the GHA workflow.
# Now that we've run `git archive` we are free to pollute the working directory.

# Delete the placeholder file
tar --file $ARCHIVE_TMP --delete ${PREFIX}/tools/integrity.bzl

mkdir -p ${PREFIX}/tools
cat >${PREFIX}/tools/integrity.bzl <<EOF
"Generated during release by release_prep.sh"

RELEASED_BINARY_INTEGRITY = $(
  jq \
    --from-file .github/workflows/integrity.jq \
    --slurp \
    --raw-input go-binaries/*.sha256
)
EOF

# Append that generated file back into the archive
tar --file $ARCHIVE_TMP --append ${PREFIX}/tools/integrity.bzl

# END patch up the archive
############

gzip <$ARCHIVE_TMP >$ARCHIVE
SHA=$(shasum -a 256 $ARCHIVE | awk '{print $1}')

cat << EOF
## Using Bzlmod with Bazel 6

1. Enable with \`common --enable_bzlmod\` in \`.bazelrc\`.
2. Add to your \`MODULE.bazel\` file:

\`\`\`starlark
bazel_dep(name = "rules_gleam", version = "${TAG:1}")

# Next, follow the instructions at
# https://github.com/iocat/rules_gleam/blob/${TAG}/README.md
\`\`\`

EOF

echo "\`\`\`" 