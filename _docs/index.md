# Dreego Documentation

Agent-friendly CLI docs. Use `dreego docs [path]` to read any file. Default: this page.

## Getting Started

- [Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md) — Quick Start guide
- [Semantic Sections Migration](https://github.com/dreego-stack/dreego/blob/main/_docs/semantic-sections-migration.md) — Move legacy root sections to `server`, `body`, and `client`
- [CLI](https://github.com/dreego-stack/dreego/blob/main/_docs/cli.md) — CLI Reference
- [Config](https://github.com/dreego-stack/dreego/blob/main/_docs/config.md) — dreego.config.json
- [Routing](https://github.com/dreego-stack/dreego/blob/main/_docs/routing.md) — File-based Routing
- [Layouts](https://github.com/dreego-stack/dreego/blob/main/_docs/layouts.md) — `{#slot}` / `{#head}` and route head behavior
- [Middleware](https://github.com/dreego-stack/dreego/blob/main/_docs/middleware.md) — Middleware System
- [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md) — SSRContext, Listen, sessions, config
- [Session Encryption](https://github.com/dreego-stack/dreego/blob/main/_docs/session-encryption.md) — AES-256-GCM encrypted session cookies
- [Progressive Enhancement](https://github.com/dreego-stack/dreego/blob/main/_docs/progressive-enhancement.md) — HTMX, Alpine.js, and plain JavaScript without an internal client runtime
- [Output Safety](https://github.com/dreego-stack/dreego/blob/main/_docs/security.md) — Context-aware escaping and URL scheme validation
- [Accessibility](https://github.com/dreego-stack/dreego/blob/main/_docs/accessibility.md) — Framework accessibility guarantees, diagnostics, and blueprint defaults

## Framework

- [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md) — Component System (`<@Name>`, slots, scoped CSS)
- [Template Logic](https://github.com/dreego-stack/dreego/blob/main/_docs/template-logic.md) — `{#if}`, `{#each}`, `$loop`, `{#verbatim}`, filters
- [Markdown Body](https://github.com/dreego-stack/dreego/blob/main/_docs/markdown.md) — Write `<body lang="md">` in Markdown, rendered to HTML at build time
- [Plugin Interfaces](https://github.com/dreego-stack/dreego/blob/main/_docs/plugin-interfaces.md) — Plugin System (planned)
- [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md) — Plugin model, middleware + route hooks

## Planned architecture

- [Roadmap](https://github.com/dreego-stack/dreego/blob/main/_docs/roadmap.md) — Public v0.x product direction
- [Implementation plans](https://github.com/dreego-stack/dreego/tree/main/_plan) — Detailed architecture, phase dependencies, acceptance criteria, and worker guidance
- [Target-neutral application decision](https://github.com/dreego-stack/dreego/blob/main/_docs/decisions/target-neutral-application-and-first-party-targets.md) — Root App plus explicit SSR, SSG, and Wails hosts
- [Semantic sections decision](https://github.com/dreego-stack/dreego/blob/main/_docs/decisions/semantic-sections-and-language-processors.md) — Implemented `server`, `head`, `body`, `style`, and `client` model

## Development

- [Testing](https://github.com/dreego-stack/dreego/blob/main/_docs/testing.md) — Integration Test Strategy
- [Reference Applications](https://github.com/dreego-stack/dreego/blob/main/_docs/reference-apps.md) — End-to-end example apps under `_tests/fixtures/`
- [Benchmarks](https://github.com/dreego-stack/dreego/blob/main/_docs/benchmarks.md) — Code generation and request benchmarks
- [Deployment](https://github.com/dreego-stack/dreego/blob/main/_docs/deployment.md) — Build, Cross-Compile, Containers
- [Architecture](https://github.com/dreego-stack/dreego/blob/main/_docs/dreego-architecture.md) — Architecture Overview
- [Dev Server](https://github.com/dreego-stack/dreego/blob/main/_docs/dev-server.md) — `dreego dev` watcher + auto-reload

## Decisions

- [Architecture Decisions](https://github.com/dreego-stack/dreego/tree/main/_docs/decisions) — ADRs (context design, routing, middleware, forms, session, transpiler, ...)

## Meta

- [Compatibility](https://github.com/dreego-stack/dreego/blob/main/_docs/compatibility.md) — Breaking-change policy and the v0.1 stability promise
- [README](https://github.com/dreego-stack/dreego/blob/main/README.md)
- [CHANGELOG](https://github.com/dreego-stack/dreego/blob/main/CHANGELOG.md)
- [Open work](https://github.com/dreego-stack/dreego/tree/main/_todo) — One item per file
