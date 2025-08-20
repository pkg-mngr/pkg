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
  const installScript = pkg.scripts.install
    .join("\n")
    .replaceAll("{{ version }}", pkg.version)
    .replaceAll("{{ pkg.bin_dir }}", "$PKG_HOME/bin")
    .replaceAll("{{ pkg.opt_dir }}", "$PKG_HOME/opt")
    .replaceAll("{{ pkg.tmp_dir }}", "$PKG_HOME/tmp")
    .replaceAll(
      "{{ pkg.completions.zsh }}",
      "$PKG_HOME/share/zsh/site-functions",
    );
  const latestScript = pkg.scripts.latest
    .join("\n")
    .replaceAll("{{ version }}", pkg.version)
    .replaceAll("{{ pkg.bin_dir }}", "$PKG_HOME/bin")
    .replaceAll("{{ pkg.opt_dir }}", "$PKG_HOME/opt")
    .replaceAll("{{ pkg.tmp_dir }}", "$PKG_HOME/tmp")
    .replaceAll(
      "{{ pkg.completions.zsh }}",
      "$PKG_HOME/share/zsh/site-functions",
    );
  const completionsScript = pkg.scripts.completions
    ?.join("\n")
    .replaceAll("{{ version }}", pkg.version)
    .replaceAll("{{ pkg.bin_dir }}", "$PKG_HOME/bin")
    .replaceAll("{{ pkg.opt_dir }}", "$PKG_HOME/opt")
    .replaceAll("{{ pkg.tmp_dir }}", "$PKG_HOME/tmp")
    .replaceAll(
      "{{ pkg.completions.zsh }}",
      "$PKG_HOME/share/zsh/site-functions",
    );

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

\`\`\`sh
${installScript}
\`\`\`

### Latest

\`\`\`sh
${latestScript}
\`\`\`

${
  completionsScript
    ? `
### Completions

\`\`\`sh
${completionsScript}
\`\`\`
`
    : ""
}
`;

  Bun.write(`packages/${pkg.name}.md`, page);
}
