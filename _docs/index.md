# Dreego Documentation

This reference is organized like a handbook while remaining available through
the existing documentation command. Run `dreego docs` for this table of
contents or `dreego docs <path>` for a chapter.

## Getting Started

- [Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md) — Quick Start guide
- [File Anatomy](https://github.com/dreego-stack/dreego/blob/main/_docs/file-anatomy.md) — How one `.dreego` file is divided into semantic sections
- [Semantic Sections Migration](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/semantic-sections-migration.md) — Move legacy root sections to `server`, `body`, and `client`
- [CLI](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/cli.md) — CLI Reference
- [Config](https://github.com/dreego-stack/dreego/blob/main/core/_docs/config.md) — dreego.config.json
- [Internationalization](https://github.com/dreego-stack/dreego/blob/main/core/_docs/i18n.md) — Compile-time catalogs, locale resolution, and formatting
- [Routing](https://github.com/dreego-stack/dreego/blob/main/core/_docs/routing.md) — File-based Routing
- [Layouts](https://github.com/dreego-stack/dreego/blob/main/core/_docs/layouts.md) — `{#slot}` / `{#head}` and route head behavior
- [Middleware](https://github.com/dreego-stack/dreego/blob/main/adapter/ssr/_docs/middleware.md) — Middleware System
- [Runtime API](https://github.com/dreego-stack/dreego/blob/main/core/_docs/runtime.md) — SSRContext, Listen, sessions, config
- [Session Encryption](https://github.com/dreego-stack/dreego/blob/main/adapter/ssr/_docs/session-encryption.md) — AES-256-GCM encrypted session cookies
- [Progressive Enhancement](https://github.com/dreego-stack/dreego/blob/main/_docs/progressive-enhancement.md) — HTMX, Alpine.js, and plain JavaScript without an internal client runtime
- [Output Safety](https://github.com/dreego-stack/dreego/blob/main/_docs/security.md) — Context-aware escaping and URL scheme validation
- [Accessibility](https://github.com/dreego-stack/dreego/blob/main/_docs/accessibility.md) — Framework accessibility guarantees, diagnostics, and blueprint defaults

## Sections and client languages

- [Server Section](https://github.com/dreego-stack/dreego/blob/main/core/_docs/server-section.md) — Go declarations and route handlers
- [Head Section](https://github.com/dreego-stack/dreego/blob/main/core/_docs/head-section.md) — Metadata and layout head composition
- [HTML Body](https://github.com/dreego-stack/dreego/blob/main/core/_docs/body-html.md) — HTML templates, expressions, and embedded processors
- [Markdown Body](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/markdown.md) — `<body lang="md">` and the `<md>` custom tag
- [Style Section](https://github.com/dreego-stack/dreego/blob/main/core/_docs/style-section.md) — Route CSS and component scoping
- [Client JavaScript](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/client-javascript.md) — Built-in dependency-free browser code
- [Client TypeScript](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/client-typescript.md) — TypeScript compiler setup and generated JavaScript
- [Browser Lua](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/lua.md) — Supported Lua syntax, semantics, browser boundary, and exclusions
- [Client Languages](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/client-languages.md) — Language support and compatibility targets

## Templates and composition

- [Components](https://github.com/dreego-stack/dreego/blob/main/core/_docs/components.md) — Component System (`<@Name>`, slots, scoped CSS)
- [Template Logic](https://github.com/dreego-stack/dreego/blob/main/core/_docs/template-logic.md) — `{#if}`, `{#each}`, `$loop`, `{#verbatim}`, filters
- [Plugin Interfaces](https://github.com/dreego-stack/dreego/blob/main/_docs/plugin-interfaces.md) — Plugin System (planned)
- [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md) — Plugin model, middleware + route hooks

## Planned architecture

- [Roadmap](https://github.com/dreego-stack/dreego/blob/main/_docs/roadmap.md) — Public v0.x product direction
- [Implementation plans](https://github.com/dreego-stack/dreego/tree/main/_plan) — Detailed architecture, phase dependencies, acceptance criteria, and worker guidance
- [Target-neutral application decision](https://github.com/dreego-stack/dreego/blob/main/_docs/decisions/target-neutral-application-and-first-party-targets.md) — Root App plus explicit SSR and Wails hosts
- [SSR over SSG decision](https://github.com/dreego-stack/dreego/blob/main/_docs/decisions/ssr-over-ssg.md) — Dynamic SSR, cache-aware islands, and no planned SSG target
- [Semantic sections decision](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/decisions/semantic-sections-and-language-processors.md) — Implemented `server`, `head`, `body`, `style`, and `client` model

## Development

- [Testing](https://github.com/dreego-stack/dreego/blob/main/dreegotest/_docs/testing.md) — Integration Test Strategy
- [Go 1.23 knowledge patch](https://github.com/dreego-stack/dreego/blob/main/_docs/handbook/go-1.23.md) — Changes an agent with Go 1.22 knowledge must know
- [Go 1.24 knowledge patch](https://github.com/dreego-stack/dreego/blob/main/_docs/handbook/go-1.24.md) — Changes an agent with Go 1.23 knowledge must know
- [Go 1.25 knowledge patch](https://github.com/dreego-stack/dreego/blob/main/_docs/handbook/go-1.25.md) — Changes an agent with Go 1.24 knowledge must know
- [Go 1.26 knowledge patch](https://github.com/dreego-stack/dreego/blob/main/_docs/handbook/go-1.26.md) — Changes an agent with Go 1.25 knowledge must know
- [Go 1.27 knowledge patch](https://github.com/dreego-stack/dreego/blob/main/_docs/handbook/go-1.27.md) — Changes an agent with Go 1.26 knowledge must know
- [Reference Applications](https://github.com/dreego-stack/dreego/blob/main/dreegotest/_docs/reference-apps.md) — End-to-end example apps under `_tests/fixtures/`
- [Benchmarks](https://github.com/dreego-stack/dreego/blob/main/dreegotest/_docs/benchmarks.md) — Code generation and request benchmarks
- [Deployment](https://github.com/dreego-stack/dreego/blob/main/adapter/ssr/_docs/deployment.md) — Build, Cross-Compile, Containers
- [Architecture](https://github.com/dreego-stack/dreego/blob/main/_docs/dreego-architecture.md) — Architecture Overview
- [Dev Server](https://github.com/dreego-stack/dreego/blob/main/cmd/dreego/_docs/dev-server.md) — `dreego dev` watcher + auto-reload

## Decisions

- [Architecture Decisions](https://github.com/dreego-stack/dreego/tree/main/_docs/decisions) — ADRs (context design, routing, middleware, forms, session, transpiler, ...)

## Meta

- [Compatibility](https://github.com/dreego-stack/dreego/blob/main/_docs/compatibility.md) — Breaking-change policy and the v0.1 stability promise
- [v0.1 to v0.2 Migration](https://github.com/dreego-stack/dreego/blob/main/_docs/migration-v0.1-to-v0.2.md) — Historical migration to explicit application and render ownership
- [README](https://github.com/dreego-stack/dreego/blob/main/README.md)
- [CHANGELOG](https://github.com/dreego-stack/dreego/blob/main/CHANGELOG.md)
- [Open work](https://github.com/dreego-stack/dreego/tree/main/_todo) — One item per file
