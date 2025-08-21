type Manifest = {
  name: string;
  description: string;
  homepage: string;
  version: string;
  sha256: Record<string, string>;
  url: Record<string, string>;
  dependencies: string[];
  caveats?: string;
  scripts: {
    install: Record<string, string[]>;
    latest: string[];
    completions?: Record<string, string[]>;
  };
};

await Bun.$`rm -rf packages public`;
await Bun.$`cp -r ../packages public`;

const manifests = await Promise.all(
  Array.from(new Bun.Glob("*.json").scanSync("./public")).map(
    (manifest) => Bun.file(`./public/${manifest}`).json() as Promise<Manifest>,
  ),
).then((manifests) => manifests.sort((a, b) => a.name.localeCompare(b.name)));

const index = `
<script setup>
import Search from "../.vitepress/components/Search.vue";
</script>

# Packages

<Search />

${manifests.map((m) => `- [${m.name}](./${m.name}) — ${m.description}{data-name="${m.name}" data-desc="${m.description}"}`).join("\n")}
`;
Bun.write("./packages/index.md", index);

function formatData(data: string, pkg: Manifest): string {
  return data
    .replaceAll("{{ version }}", pkg.version)
    .replaceAll("{{ pkg.bin_dir }}", "$PKG_HOME/bin")
    .replaceAll("{{ pkg.opt_dir }}", "$PKG_HOME/opt")
    .replaceAll("{{ pkg.tmp_dir }}", "$PKG_HOME/tmp")
    .replaceAll(
      "{{ pkg.completions.zsh }}",
      "$PKG_HOME/share/zsh/site-functions",
    );
}

for (const pkg of manifests) {
  const installScripts: string[] = [];
  for (const platform in pkg.scripts.install) {
    installScripts.push(`
\`\`\`sh [${platform}]
${formatData(pkg.scripts.install[platform]!.join("\n"), pkg)}
\`\`\`
  `);
  }

  const latestScript = formatData(pkg.scripts.latest.join("\n"), pkg);

  const completionsScripts: string[] = [];
  for (const platform in pkg.scripts.completions) {
    completionsScripts.push(`
\`\`\`sh [${platform}]
${formatData(pkg.scripts.completions[platform]!.join("\n"), pkg)}
\`\`\`
  `);
  }

  const page = `
[← See all packages](./index.md)

# ${pkg.name}

Install command:

\`\`\`sh
pkg add ${pkg.name}
\`\`\`

${pkg.description}

Version: \`${pkg.version}\`

Homepage: ${pkg.homepage}

Manifest: [${pkg.name}.json](/${pkg.name}.json)

| Platform | SHA256 Checksum |
| -------- | --------------- |
${Object.entries(pkg.sha256)
  .map(([platform, sha256]) => `| ${platform} | \`${sha256}\` |`)
  .join("\n")}

${
  pkg.dependencies
    ? `Dependencies:
${pkg.dependencies.map((dep) => `- [${dep}](./${dep}.md)`).join("\n")}
`
    : ""
}

${
  pkg.caveats
    ? `
::: warning CAVEATS
${pkg.caveats}
:::`
    : ""
}

## Scripts

### Install

::: code-group

${installScripts.join("\n")}

:::

### Latest

\`\`\`sh
${latestScript}
\`\`\`

${
  completionsScripts.length > 0
    ? `
### Completions

::: code-group

${completionsScripts.join("\n")}

:::
`
    : ""
}
`;

  Bun.write(`packages/${pkg.name}.md`, page);
}
