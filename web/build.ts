type Manifest = {
  name: string;
  description: string;
  homepage: string;
  version: string;
  sha256: string | Record<string, string>;
  url: string | Record<string, string>;
  dependencies?: string[];
  caveats?: string;
  platforms?: Record<string, {
    url: string;
    sha256: string;
    scripts?: {
      install?: string[];
      latest?: string[];
      completions?: string[];
    };
  }>;
  scripts: {
    install: string[] | Record<string, string[]>;
    latest: string[] | Record<string, string[]>;
    completions?: string[] | Record<string, string[]>;
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
  // Helper function to process script replacements
  const processScript = (script: string | undefined) => {
    if (!script) return undefined;
    return script
      .replaceAll("{{ version }}", pkg.version)
      .replaceAll("{{ pkg.bin_dir }}", "$PKG_HOME/bin")
      .replaceAll("{{ pkg.opt_dir }}", "$PKG_HOME/opt")
      .replaceAll("{{ pkg.tmp_dir }}", "$PKG_HOME/tmp")
      .replaceAll(
        "{{ pkg.completions.zsh }}",
        "$PKG_HOME/share/zsh/site-functions",
      );
  };

  // Determine if this is old or new format and get platforms
  const isNewFormat = typeof pkg.sha256 === 'object' && typeof pkg.url === 'object';
  const platforms = isNewFormat ? Object.keys(pkg.sha256 as Record<string, string>) : ['legacy'];
  
  // Process install scripts
  let installScript: string;
  if (Array.isArray(pkg.scripts.install)) {
    // Global install script
    installScript = processScript(pkg.scripts.install.join("\n")) || "";
  } else {
    // Platform-specific install scripts - show the first platform as example
    const firstPlatform = platforms[0];
    const platformScript = pkg.scripts.install[firstPlatform];
    installScript = processScript(platformScript?.join("\n")) || "";
  }

  // Process latest scripts
  let latestScript: string;
  if (Array.isArray(pkg.scripts.latest)) {
    // Global latest script
    latestScript = processScript(pkg.scripts.latest.join("\n")) || "";
  } else {
    // Platform-specific latest scripts - show the first platform as example
    const firstPlatform = platforms[0];
    const platformScript = pkg.scripts.latest[firstPlatform];
    latestScript = processScript(platformScript?.join("\n")) || "";
  }

  // Process completions scripts
  let completionsScript: string | undefined;
  if (pkg.scripts.completions) {
    if (Array.isArray(pkg.scripts.completions)) {
      // Global completions script
      completionsScript = processScript(pkg.scripts.completions.join("\n"));
    } else {
      // Platform-specific completions scripts - show the first platform as example
      const firstPlatform = platforms[0];
      const platformScript = pkg.scripts.completions[firstPlatform];
      completionsScript = processScript(platformScript?.join("\n"));
    }
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

${isNewFormat ? `## Supported Platforms

${platforms.map(platform => {
  const urls = pkg.url as Record<string, string>;
  const sha256s = pkg.sha256 as Record<string, string>;
  const url = urls[platform];
  const sha256 = sha256s[platform];
  return `- **${platform}**: [Download](${url?.replace("{{ version }}", pkg.version) || "#"}) (SHA256: \`${sha256}\`)`;
}).join("\n")}` : `## Download

- **URL**: [Download](${(pkg.url as string)?.replace("{{ version }}", pkg.version) || "#"})
- **SHA256**: \`${pkg.sha256}\``}

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

${Array.isArray(pkg.scripts.install) ? 
  `### Install ${isNewFormat ? "(Global)" : ""}

\`\`\`sh
${installScript}
\`\`\`` :
  `### Install Scripts by Platform

${platforms.map(platform => {
    const platformScript = (pkg.scripts.install as Record<string, string[]>)[platform];
    const processed = processScript(platformScript?.join("\n"));
    return `#### ${platform}

\`\`\`sh
${processed || "No install script"}
\`\`\``;
  }).join("\n\n")}`
}

${Array.isArray(pkg.scripts.latest) ? 
  `### Latest Version Check ${isNewFormat ? "(Global)" : ""}

\`\`\`sh
${latestScript}
\`\`\`` :
  `### Latest Version Check Scripts by Platform

${platforms.map(platform => {
    const platformScript = (pkg.scripts.latest as Record<string, string[]>)[platform];
    const processed = processScript(platformScript?.join("\n"));
    return `#### ${platform}

\`\`\`sh
${processed || "No latest script"}
\`\`\``;
  }).join("\n\n")}`
}

${
  pkg.scripts.completions
    ? Array.isArray(pkg.scripts.completions) ?
      `### Completions ${isNewFormat ? "(Global)" : ""}

\`\`\`sh
${completionsScript}
\`\`\`` :
      `### Completions Scripts by Platform

${platforms.map(platform => {
    const platformScript = (pkg.scripts.completions as Record<string, string[]>)[platform];
    const processed = processScript(platformScript?.join("\n"));
    return processed ? `#### ${platform}

\`\`\`sh
${processed}
\`\`\`` : "";
  }).filter(Boolean).join("\n\n")}`
    : ""
}
`;

  Bun.write(`packages/${pkg.name}.md`, page);
}
