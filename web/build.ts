type Manifest = {
  name: string;
  description: string;
  homepage: string;
  version: string;
  sha256: Record<string, string>;
  url: Record<string, string>;
  dependencies?: string[];
  caveats?: string;
  scripts: {
    install: Record<string, string[]>;
    latest: string[];
    completions?: Record<string, string[]>;
  };
};

await Bun.$`rm -rf packages`;

// Copy manifests from ../packages to ./public
await Bun.$`cp ../packages/*.json ./public/`;

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
  // Helper function to process script replacements
  const processScript = (scripts: string[]) => {
    return scripts
      .join("\n")
      .replaceAll("{{ version }}", pkg.version)
      .replaceAll("{{ pkg.bin_dir }}", "$PKG_HOME/bin")
      .replaceAll("{{ pkg.opt_dir }}", "$PKG_HOME/opt")
      .replaceAll("{{ pkg.tmp_dir }}", "$PKG_HOME/tmp")
      .replaceAll(
        "{{ pkg.completions.zsh }}",
        "$PKG_HOME/share/zsh/site-functions",
      );
  };

  // Get platforms from sha256 keys
  const platforms = Object.keys(pkg.sha256);
  
  // Create code groups for install scripts
  const createInstallCodeGroup = () => {
    const groups = platforms.map(platform => {
      const script = processScript(pkg.scripts.install[platform] || []);
      return `:::code-group

\`\`\`sh [${platform}]
${script}
\`\`\`

:::`;
    }).join('\n\n');
    
    return `:::code-group

${platforms.map(platform => {
  const script = processScript(pkg.scripts.install[platform] || []);
  return `\`\`\`sh [${platform}]
${script}
\`\`\``;
}).join('\n\n')}

:::`;
  };

  // Create code groups for completions scripts (if they exist)
  const createCompletionsCodeGroup = () => {
    if (!pkg.scripts.completions) return '';
    
    return `:::code-group

${platforms.map(platform => {
  const script = pkg.scripts.completions?.[platform];
  if (!script) return '';
  const processed = processScript(script);
  return `\`\`\`sh [${platform}]
${processed}
\`\`\``;
}).filter(Boolean).join('\n\n')}

:::`;
  };

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

## Supported Platforms

${platforms.map(platform => {
  const url = pkg.url[platform];
  const sha256 = pkg.sha256[platform];
  return `- **${platform}**:
  - URL: \`${url?.replace("{{ version }}", pkg.version) || "#"}\`
  - SHA256: \`${sha256}\`
  `;
}).join("\n")}

${
  pkg.dependencies && pkg.dependencies.length > 0
    ? `## Dependencies

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

${createInstallCodeGroup()}

### Latest Version Check

\`\`\`sh
${processScript(pkg.scripts.latest)}
\`\`\`

${pkg.scripts.completions ? `### Completions

${createCompletionsCodeGroup()}` : ''}
`;

  Bun.write(`packages/${pkg.name}.md`, page);
}
