type Manifest = {
  name: string;
  description: string;
  homepage: string;
  version: string;
  sha256: string;
  url: string;
  dependencies: string[];
  caveats?: string;
  scripts: {
    install: string[];
    latest: string[];
    completions?: string[];
  };
};

await Bun.$`rm -rf packages`;

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

for (const pkg of manifests) {
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

SHA256 Checksum: \`${pkg.sha256}\`

Dependencies:
${pkg.dependencies.map((dep) => `- [${dep}](./${dep}.md)`).join("\n")}

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

\`\`\`sh
${pkg.scripts.install.join("\n")}
\`\`\`

### Latest

\`\`\`sh
${pkg.scripts.latest.join("\n")}
\`\`\`

${
  pkg.scripts.completions
    ? `
### Completions

\`\`\`sh
${pkg.scripts.completions.join("\n")}
\`\`\`
`
    : ""
}
`;

  Bun.write(`packages/${pkg.name}.md`, page);
}
