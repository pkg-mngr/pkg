# pkg

A simple cross-platform package manager for macOS and Linux.

## Installation

You can build it from source using Go:

```sh
go install github.com/noclaps/pkg@latest
```

or download a prebuilt binary in [Releases](https://github.com/noClaps/pkg/releases).

## Usage

```
USAGE: pkg [add | update | remove | info | list] [--init]

COMMANDS:
  add               Install packages.
  update            Update packages.
  remove            Remove packages.
  info              Get the info for a package.
  list              List installed packages

OPTIONS:
  --init            Initialise pkg
  -h, --help        Display this help and exit.
```

Initialise `pkg` with:

```sh
pkg --init
```

You can install packages by running:

```sh
pkg add go # or any other package
```

You can update installed packages with:

```sh
pkg update
```

You can also remove installed packages with:

```sh
pkg remove go
```

You can fetch the info for a package with:

```sh
pkg info go
```

You can list installed packages with:

```sh
pkg list
```

You can view the help by using `-h` or `--help`:

```sh
pkg -h
pkg --help
```
