import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "pkg",
  description: "A simple cross-platform package manager for macOS and Linux",
  srcExclude: ["README.md"],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      {
        text: "Packages",
        link: "/packages",
      },
    ],
    aside: false,
    socialLinks: [{ icon: "github", link: "https://github.com/noClaps/pkg" }],
  },
});
