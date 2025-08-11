# pkg

A simple package manager for macOS.

## Website build instructions

You'll need [Bun](https://bun.sh) and [Node](https://nodejs.org).

Clone the repository:

```sh
git clone https://github.com/noClaps/pkg.git
cd pkg
```

Copy the `packages/` directory into `web/public/`:

```sh
cp -r packages/ web/public/
```

Run the build script inside the `web/` directory:

```sh
cd web
bun run build
```

The built files will be in `.vitepress/dist/`.
