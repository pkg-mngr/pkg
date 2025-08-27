# Self-hosting

This guide will help you host your own remote manifest host.

## Setting up `pkg`

In order to point your installed `pkg` at a specific manifest host, you can create an environment variable called `PKG_MANIFEST_HOST` and point it to the root directory of your host. For example, if you're serving packages from `https://pkg.example.com`, you would set it to:

```sh
PKG_MANIFEST_HOST="https://pkg.example.com"
```

You can override this globally by setting it in your shell config file (`.bashrc`, `.zshrc`, etc.) or for a specific installation. Once you set a remote for a package, it will be saved in the lockfile, and all future updates will be fetched from that remote, even if the `PKG_MANIFEST_HOST` environment variable is not set, or set to a different value.

If the `PKG_MANIFEST_HOST` environment variable is unset, the default manifest host is `https://pkg.zerolimits.dev`.

## Setting up your remote

At the minimum your remote must serve your package manifests **from the root**. For example, the Go package JSON file would need to be served from `https://pkg.example.com/go.json`.

If your manifests are in a subdirectory, be sure to include that subdirectory in the `PKG_MANIFEST_HOST` environment variable (e.g. if you're serving them from `https://example.com/packages/go.json`, you would set `PKG_MANIFEST_HOST` to be `https://example.com/packages`).

Each package manifest must follow the schema in https://github.com/pkg-mngr/pkg/blob/main/package.schema.json. If they do not, the `pkg` CLI may fail to install them, or result in unexpected behavior for certain commands.

### CLI Search Command

In order to support the `pkg search` command, you need to have an `index.json` served from the same path as the other package manifests. The `index.json` must have all of your packages listed in the following format:

```json
{
  "name": {
    "version": "...",
    "description": "..."
  }
  // ...
}
```

For example:

```json
// https://pkg.example.com/index.json
{
  "pkg": {
    "version": "0.3.0",
    "description": "A simple cross-platform package manager"
  },
  "go": {
    "version": "1.25.0",
    "description": "The Go programming language"
  }
}
```

Now, running `pkg search go` will display all results with "go" in the name or description. Note that `pkg search` uses the `PKG_MANIFEST_HOST` environment variable at runtime, so if you're getting your packages from multiple remotes, you may need to set the `PKG_MANIFEST_HOST` variable before running `pkg search`.
