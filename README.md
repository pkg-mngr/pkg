# pkg

A simple cross-platform package manager.

## Installation

You can build it from source using Go:

```sh
go install github.com/noclaps/pkg@latest
```

or download a prebuilt binary in [Releases](https://github.com/noClaps/pkg/releases).

## Usage

```
USAGE: pkg [add | update | remove | info | list | platform] [--init]

COMMANDS:
  add               Install packages.
  update            Update packages.
  remove            Remove packages.
  info              Get the info for a package.
  list              List installed packages
  platform          Show current platform information

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

## Platform Support

`pkg` now supports platform-specific package configurations, allowing packages to provide different URLs, checksums, and installation scripts for different operating systems and architectures.

### Supported Platforms

- **Linux**: `linux-x86_64`, `linux-arm64`
- **macOS**: `macos-x86_64`, `macos-arm64`

### Platform Detection

You can check your current platform with:

```sh
pkg platform
```

This will show your current platform.

### Package Manifest Format

Packages can now specify platform-specific configurations:

```json
{
  "name": "example-package",
  "platforms": {
    "linux-x86_64": {
      "url": "https://example.com/linux-x86_64.tar.gz",
      "sha256": "abc123...",
      "scripts": {
        "install": ["./install-linux.sh"],
        "latest": ["curl -s https://api.github.com/repos/owner/repo/releases/latest"]
      }
    },
    "macos-x86_64": {
      "url": "https://example.com/macos-x86_64.tar.gz",
      "sha256": "def456...",
      "scripts": {
        "install": ["./install-mac.sh"]
      }
    }
  },
  "scripts": {
    "install": ["./install-fallback.sh"],
    "latest": ["./check-latest.sh"]
  }
}
```

### Fallback Behavior

If no platform-specific configuration is found, `pkg` will fall back to the default `url`, `sha256`, and `scripts` fields. This ensures backward compatibility with existing packages.
