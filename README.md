
```markdown
# rules_gleam

This repository contains Bazel rules for working with the [Gleam](https://gleam.run) programming language.

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Rules](#rules)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

To use the rules, add the following to your `WORKSPACE` file:

```python
# WORKSPACE
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "rules_gleam",
    urls = ["https://github.com/your-username/rules_gleam/archive/main.zip"],
    strip_prefix = "rules_gleam-main",
)

load("@rules_gleam//:deps.bzl", "rules_gleam_dependencies", "rules_gleam_toolchains")

rules_gleam_dependencies()

rules_gleam_toolchains()
```

Then, set up the Gleam toolchain in your `BUILD` files:

```python
# BUILD
load("@rules_gleam//:gleam.bzl", "gleam_binary", "gleam_library")

gleam_library(
    name = "example_lib",
    srcs = ["src/example.gleam"],
    deps = [
        "//path/to/dependency",
    ],
)

gleam_binary(
    name = "example_bin",
    srcs = ["src/main.gleam"],
    main_module = "main",
    deps = [
        ":example_lib",
    ],
)
```

## Usage

### Building a Gleam Library

To build a Gleam library, use the `gleam_library` rule:

```python
gleam_library(
    name = "my_lib",
    srcs = ["src/lib.gleam"],
    deps = [
        "//path/to/dependency",
    ],
)
```

### Building a Gleam Binary

To build a Gleam binary, use the `gleam_binary` rule:

```python
gleam_binary(
    name = "my_bin",
    srcs = ["src/main.gleam"],
    main_module = "main",
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
- `srcs` (mandatory): List of `.gleam` source files.
- `deps`: List of dependencies.

### `gleam_binary`

Builds a Gleam binary.

**Attributes:**

- `name` (mandatory): A unique name for this target.
- `srcs` (mandatory): List of `.gleam` source files.
- `main_module` (mandatory): The main module containing the `main` function.
- `deps`: List of dependencies.

## Examples

You can find example usage of these rules in the [`examples`](examples) directory.

### Basic Example

```python
# examples/basic/BUILD
load("@rules_gleam//:gleam.bzl", "gleam_binary", "gleam_library")

gleam_library(
    name = "example_lib",
    srcs = ["src/example.gleam"],
)

gleam_binary(
    name = "example_bin",
    srcs = ["src/main.gleam"],
    main_module = "main",
    deps = [
        ":example_lib",
    ],
)
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```

Feel free to make any additional customizations or adjustments as needed!