```markdown
# rules_gleam

This repository contains Bazel rules for working with the [Gleam](https://gleam.run) programming language.

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Rules](#rules)
  - [`gleam_library`](#gleam_library)
  - [`gleam_binary`](#gleam_binary)
  - [`gleam_test`](#gleam_test)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

**Note:** These rules are designed for use with Bazel Modules and do not support the legacy `WORKSPACE` setup.

To use the rules, add the following to your `MODULE.bazel` file:

```starlark
# MODULE.bazel

bazel_dep(name = "rules_gleam", version = "0.0.1")
```

Then, in your `BUILD.bazel` files, load the rules you need:

```starlark
# BUILD.bazel
load("@rules_gleam//gleam:defs.bzl", "gleam_binary", "gleam_library", "gleam_test")
```

### Managing Hex Dependencies

To use external dependencies from [Hex](https://hex.pm/), you'll need a `gleam.toml` and a `manifest.toml` file.

1.  **Define Dependencies in `gleam.toml`**: Add your Hex dependencies to the `[dependencies]` or `[dev-dependencies]` sections of your `gleam.toml` file.

    ```toml
    # gleam.toml
    [dependencies]
    gleam_stdlib = ">= 0.51.0 and < 2.0.0"
    gleam_json = ">= 3.0.2 and < 4.0.0"

    [dev-dependencies]
    gleeunit = ">= 1.0.0 and < 2.0.0"
    ```

2.  **Fetch Dependencies**: Run `gleam deps download` to fetch the dependencies and generate a `manifest.toml` file.

Commit and save both "gleam.toml", and "manifest.toml". The rest can be discarded.

3.  **Configure Bazel Module**: In your `MODULE.bazel` file, use the `gleam.deps` extension to declare your dependencies.

    ```starlark
    # MODULE.bazel
    gleam = use_extension("@rules_gleam//:extensions.bzl", "gleam")
    gleam.deps(gleam_toml = "//:gleam.toml")
    use_repo(gleam, "hex_gleam_json", "hex_gleam_stdlib", "hex_gleeunit")
    ```

4.  **Use Dependencies in `BUILD.bazel`**: You can now reference the Hex packages in your `BUILD.bazel` file.

    ```starlark
    # BUILD.bazel
    gleam_binary(
        name = "example",
        srcs = ["example.gleam"],
        deps = [
            "@hex_gleam_json//gleam",
            "@hex_gleam_stdlib//gleam",
        ],
    )
    ```

    Gazelle will automatically update your repository with these dependencies.

## Usage

### Building a Gleam Library

To build a Gleam library, use the `gleam_library` rule. This rule compiles your Gleam source files into Erlang modules.

```starlark
gleam_library(
    name = "my_lib",
    srcs = ["my_lib.gleam"],
    deps = [
        "//path/to/another:lib",
    ],
)
```

### Building a Gleam Binary

To build a Gleam binary, use the `gleam_binary` rule. This creates an executable script that runs your Gleam application.

```starlark
gleam_binary(
    name = "my_app",
    srcs = ["main.gleam"],
    main_module = "main",
    deps = [
        ":my_lib",
    ],
)
```

### Testing a Gleam Module

To test a Gleam module, use the `gleam_test` rule. This will compile your test files and run them using the `gleeunit` test runner.

```starlark
gleam_test(
    name = "my_lib_test",
    srcs = ["my_lib_test.gleam"],
    deps = [
        ":my_lib",
    ],
)
```

## Rules

### `gleam_library`

Builds a Gleam library.

**Attributes:**

- `name` (mandatory): A unique name for this target.
- `srcs` (mandatory): A list of `.gleam` source files to be compiled.
- `deps`: A list of other `gleam_library` targets that this library depends on.
- `data`: A list of data files needed by the library at runtime.
- `strip_src_prefix`: A string to strip from the beginning of the source file paths.

### `gleam_binary`

Builds a Gleam binary.

**Attributes:**

- `name` (mandatory): A unique name for this target.
- `srcs` (mandatory): A list of `.gleam` source files to be compiled.
- `main_module` (mandatory): The name of the module containing the `main` function.
- `deps`: A list of `gleam_library` targets that this binary depends on.
- `data`: A list of data files needed by the binary at runtime.
- `strip_src_prefix`: A string to strip from the beginning of the source file paths.

### `gleam_test`

Builds and runs a Gleam test.

**Attributes:**

- `name` (mandatory): A unique name for this target.
- `srcs` (mandatory): A list of `.gleam` test files. Test file names must end with `_test.gleam` or `_tests.gleam`.
- `deps`: A list of `gleam_library` targets that the test depends on.
- `size`: The size of the test. Can be `small`, `medium`, `large`, or `enormous`.
- `timeout`: The timeout for the test. Can be `short`, `moderate`, `long`, or `eternal`.
- `data`: A list of data files needed by the test at runtime.
- `strip_src_prefix`: A string to strip from the beginning of the source file paths.

## Gazelle Integration

This repository provides a Gazelle extension that can automatically generate `BUILD.bazel` files for your Gleam projects.

### Setup

To use the Gazelle extension, first, set it up in your root `BUILD.bazel` file:

```starlark
# BUILD.bazel
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    gazelle = "@rules_gleam//gazelle",
)
```

### Usage

Once configured, you can run Gazelle from the command line:

```sh
bazel run //:gazelle
```

Gazelle will scan your project and generate `gleam_library`, `gleam_binary`, and `gleam_test` rules automatically.

### Directives

The Gleam Gazelle extension supports the following directive:

- `gleam_visibility`: Specifies the visibility of the generated targets. You can add this as a comment in your `BUILD.bazel` file.

  ```starlark
  # gazelle:gleam_visibility //my/project:__subpackages__
  ```

## Examples

You can find example usage of these rules in the [`examples`](examples) directory.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```

Feel free to make any additional customizations or adjustments as needed!