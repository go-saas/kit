// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

import { themes as prismThemes } from "prism-react-renderer";

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "GO-SAAS-KIT",
  tagline: "Open source SAAS starter kit built with golang and 💙",
  favicon: "img/favicon.ico",

  // Set the production url of your site here
  url: "https://go-saas.github.io",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/kit/",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "go-saas", // Usually your GitHub org/user name.
  projectName: "kit", // Usually your repo name.

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["zh-Hans", "en"],
    localeConfigs: {
      "zh-Hans": {
        label: "中文",
      },
      en: {
        label: "English",
      },
    },
  },

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editLocalizedFiles: true,
          sidebarCollapsible: true,
          sidebarCollapsed: false,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/go-saas/kit/tree/main/docs",
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/go-saas/kit/tree/main/blog",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: "img/desc.png",
      navbar: {
        title: "GO-SAAS-KIT",
        logo: {
          alt: "GO-SAAS-KIT Logo",
          src: "img/logo.png",
        },
        items: [
          {
            type: "doc",
            docId: "intro/overview",
            position: "left",
            label: "Tutorial",
          },
          { to: "/blog", label: "Blog", position: "left" },
          {
            type: "localeDropdown",
            position: "right",
          },
          {
            href: "https://github.com/go-saas/kit",
            position: "right",
            className: "header-github-link",
            "aria-label": "GitHub repository",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: "Docs",
            items: [
              {
                label: "Tutorial",
                to: "/docs/intro/overview",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "Blog",
                to: "/blog",
              },
              {
                label: "GitHub",
                href: "https://github.com/go-saas",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} GO-SAAS-KIT.`,
      },
      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
      },
      colorMode: {
        defaultMode: "dark",
        // disableSwitch: true,
        respectPrefersColorScheme: false,
      },
    }),
  markdown: {
    mermaid: true,
  },
  themes: ["@docusaurus/theme-mermaid"],
};

module.exports = config;
