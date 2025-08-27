import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "pkg",
  description: "A simple cross-platform package manager",
  srcExclude: ["README.md", "**.tmpl.md"],
  head: [["link", { rel: "icon", href: "/pkg.png" }]],
  themeConfig: {
    logo: "/pkg.png",
    sidebar: [
      { text: "CLI", link: "/" },
      { text: "Packages", link: "/packages" },
      { text: "Self-hosting", link: "/self-hosting" },
    ],
    outline: [2, 6],
    socialLinks: [{ icon: "github", link: "https://github.com/pkg-mngr/pkg" }],
  },
});
