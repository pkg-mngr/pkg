[← See all packages](./index.md)

# {{ name }}

Install command:

```sh
pkg add {{ name }}
```

{{ description }}

Version: `{{ version }}`

Homepage: {{ homepage }}

Manifest: [{{ name }}.json](/{{ name }}.json)

| Platform | SHA256 Checksum |
| -------- | --------------- |
{{ sha256 }}

{{ dependencies }}

{{ caveats }}

## Scripts

### Install

::: code-group

{{ scripts.install }}

:::

### Latest

```sh
{{ scripts.latest }}
```

{{ completions }}
