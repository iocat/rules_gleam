load("//gleam:gleam_library.bzl", _gleam_library = "gleam_library")
load("//gleam:gleam_binary.bzl", _gleam_binary = "gleam_binary")
load("//gleam:gleam_erl_library.bzl", _gleam_erl_library = "gleam_erl_library")
load("//gleam:gleam_repository.bzl", _gleam_repository = "gleam_repository")

gleam_library = _gleam_library
gleam_binary = _gleam_binary
gleam_erl_library = _gleam_erl_library
gleam_repository = _gleam_repository