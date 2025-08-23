import fs from "node:fs";
import packageTemplate from "./package.tmpl.md" with { type: "text" };
import indexTemplate from "./packages-index.tmpl.md" with { type: "text" };

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

function $(cmd: string) {
  const [arg0, ...args] = cmd.split(" ");
  return new Deno.Command(arg0!, { args: args }).outputSync();
}

$("rm -rf packages public");
$("cp -r ../packages public");
$("cp pkg.png public");
$("mkdir packages");

const manifests = fs.globSync("./public/*.json")
  .map(
    (manifest) =>
      JSON.parse(fs.readFileSync(manifest, { encoding: "utf8" })) as Manifest,
  ).sort((a, b) => a.name.localeCompare(b.name));

const index = indexTemplate.replace(
  "{{ manifests }}",
  manifests
    .map(
      (m) =>
        `- [${m.name}](./${m.name}) — ${m.description}{data-name="${m.name}" data-desc="${m.description}"}`,
    )
    .join("\n"),
);
Deno.writeTextFileSync("./packages/index.md", index);

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

  let completionsScripts = "";
  if (pkg.scripts.completions) {
    completionsScripts = `### Completions

::: code-group

`;
    for (const platform in pkg.scripts.completions) {
      completionsScripts += `
\`\`\`sh [${platform}]
${formatData(pkg.scripts.completions[platform]!.join("\n"), pkg)}
\`\`\`
`;
    }
    completionsScripts += ":::";
  }

  const sha256 = Object.entries(pkg.sha256)
    .map(([platform, sha256]) => `| ${platform} | \`${sha256}\` |`)
    .join("\n");

  const dependencies = pkg.dependencies
    ? `Dependencies:
${pkg.dependencies.map((dep) => `- [${dep}](./${dep}.md)`).join("\n")}
`
    : "";

  const caveats = pkg.caveats
    ? `
::: warning CAVEATS
${pkg.caveats}
:::`
    : "";

  const page = packageTemplate
    .replaceAll("{{ name }}", pkg.name)
    .replaceAll("{{ description }}", pkg.description)
    .replaceAll("{{ version }}", pkg.version)
    .replaceAll("{{ homepage }}", pkg.homepage)
    .replaceAll("{{ sha256 }}", sha256)
    .replaceAll("{{ dependencies }}", dependencies)
    .replaceAll("{{ caveats }}", caveats)
    .replaceAll("{{ scripts.install }}", installScripts.join("\n"))
    .replaceAll("{{ scripts.latest }}", latestScript)
    .replaceAll("{{ completions }}", completionsScripts);

  Deno.writeTextFileSync(`packages/${pkg.name}.md`, page);
}
