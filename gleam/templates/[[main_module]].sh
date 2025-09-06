#!/bin/sh
set -eu

PACKAGE={PACKAGE}
BASE=$(dirname "$0")

run() {
  exec erl \
    -pa "$BASE" \
    -eval "$PACKAGE@@main:run($PACKAGE)" \
    -noshell \
    -extra "$@"
}

run "$@"