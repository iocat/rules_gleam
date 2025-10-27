"""Defs for Gleam JS rules."""

load(":gleam_js_binary.bzl", _gleam_js_binary = "gleam_js_binary")
load(":gleam_js_library.bzl", _gleam_js_library = "gleam_js_library")
# load(":gleam_js_test.bzl", _gleam_js_test = "gleam_js_test")
load(":gleam_js_ffi_library.bzl", _gleam_js_ffi_library = "gleam_js_ffi_library")

gleam_js_binary = _gleam_js_binary
# gleam_js_test = _gleam_js_test
gleam_js_library = _gleam_js_library
gleam_js_ffi_library = _gleam_js_ffi_library
