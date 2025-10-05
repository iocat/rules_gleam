#!/bin/sh
set -eu

PACKAGE={PACKAGE}
BASE=$(echo "$0.runfiles/_main")

run() {
  exec erl \
    -pa "$BASE" \
    -eval "$PACKAGE@@main:run($PACKAGE)" \
    -noshell \
    -extra "$@"
}

run "$@"