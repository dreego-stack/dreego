# Roadmap

Dreego's roadmap is directional, not a release promise. It frames Dreego as
SSR-first, accessibility-first, Go-native, Svelte/Astro-inspired.

The goal is to stay in v0.x as long as possible — everything that would
traditionally be "v1" or "v2" lands in v0.x patches when it can be added
without breaking stability. v1 is not "the big new thing" but "now everything
is stable and promised."

---

## v0.1 — Release Candidate (current)

### Completed
- Explicit App-Ownership of all runtime state
- Removed speculative APIs (EventBus, Queue, KVStore, Storage)
- SSR hardening: timeouts, graceful shutdown, panic recovery, CSRF, sessions
- Accessibility quality gate (CLI, docs, blueprints, components)
- Reference apps as blackbox tests (hello, forms, components, plugin)
- Plugin build-hook mechanism (dreego-plugin.json + interactive approval)
- 4 plugins on GitHub (example, sse, websocket, tailwind)
- Demo app (example-with-plugins) with 8 components

### What remains
- Tag v0.1

---

## v0.2 — Non-HTTP Render Mode

The single biggest unblock: `app.Render(page, data)` without an HTTP server.

### What
- Decouple `SSRContext` from `http.ResponseWriter` — use `io.Writer`
- Add `App.Render(page string, data map[string]any) (string, error)`
- Session/cookie methods become optional (only when real `*http.Request` exists)

### Unblocks
- Wails: WebView calls `app.Render()` instead of HTTP
- SSG: loop over routes, render, save as `.html` files
- Testing: `app.Render()` directly instead of httptest.Server
- CLI preview: `dreego preview /about` without a server

### Effort
Medium. SSRContext already abstracts `*http.Request` via the Context interface.
Generated handlers call `NewSSR(w, r)` — needs to become `NewSSR(io.Writer, *http.Request)`
with nil-safe `r`. No transpiler change for the base.

---

## v0.3 — Transpiler Extension Points

Opens the transpiler for plugins. The bottleneck after v0.2.

### What
```go
type SectionProcessor interface {
    Tag() string
    Process(node ASTNode) (string, error)
}
```
Plugins register new section types via `dreego-plugin.json`. The transpiler
discovers processors from `go.mod` (same mechanism as build-hooks).

### Unblocks
- `<script lang="lua">` — plugin-lua via gopher-lua (pure Go, ~100KB)
- `<markdown>` — plugin-markdown
- `<svg>` — typed SVG components
- Custom directives — `{#chart}`, `{#mermaid}` as plugins

### Effort
Large. Lexer has hardcoded `sectionTags`. Parser has hardcoded tag switches.
Codegen has hardcoded functions per section. Must be formalized.
ADR `transpiler-extensions.1.md` says v2+ — pulled forward to v0.3.

---

## v0.4 — First Transpiler Plugins

Proves the extension points work with at least 2 real processors.

### What
1. **plugin-markdown** — `<markdown>` section, renders Markdown to HTML
2. **plugin-lua** — `<script lang="lua">` section, runs Lua server-side via gopher-lua

### Why two
The ADR requires "at least 2 real processors validate the contract." A single
processor could accidentally satisfy the interface. Two prove it is generic.

### Effort
Medium — once extension points exist, processors are small.

---

## v0.5 — Plugin Lifecycle Hooks + Streaming SSR

### Plugin Lifecycle
```go
type LifecyclePlugin interface {
    OnStart(app *dreego.App) error
    OnShutdown(app *dreego.App) error
}
```
Today `App.Shutdown()` only stops the HTTP server. Plugins with background
work (Redis connections, cron jobs, WebSocket hubs) have no clean shutdown.
The `plugin-contract.1.md` ADR explicitly requires this.

### Streaming SSR
Today `renderXxx` builds the complete HTML string, then writes once. No streaming.
v0.5: codegen writes incrementally to `io.Writer` instead of `strings.Builder`.
Flush points after `</section>`, `</main>`. Browser can render the header + hero
immediately while the rest is still being built.

### Effort
Medium. Codegen change (StringBuilder to io.Writer) + compress middleware coordination.

---

## v0.6 — Client-side Codegen Experiment

The Svelte-Runes model: `.dreego` compiles to Go (server) AND JS (client).

### What
```dreego
<div>
    <button @click={increment}>Count: {{ count }}</button>
</div>

<go server>
    count := 0
</go>

<go client>
    function increment() { count++ }
</go>
```

The transpiler generates:
1. Go code for SSR (as today)
2. JS code for client reactivity (new)

### Status
Experiment — not for stable core. Lives outside core per `client-reactivity.1.md`.
Promotion requires evidence from real apps.

### Effort
Very large. Second codegen target, state serialization, event binding, security
model. Realistic v0.6+.

---

## v0.7 — SSG (Static Site Generation)

Now possible because v0.2 Non-HTTP Rendering exists.

### What
```sh
dreego build --static  # generates static HTML files
```
- Loop over all routes, `app.Render()`, write `dist/index.html` etc.
- Dynamic routes: `/users/[id]` → generate for each known ID
- Content collections: blog posts from Markdown frontmatter
- Deployment: Cloudflare Pages, GitHub Pages

### Effort
Medium. The hard work was v0.2. SSG is "call Render for each route, write files."

---

## v0.8 — Wails Support

Now possible because v0.2 Non-HTTP Rendering exists.

### What
- `app.Render(page, data)` in a Wails WebView
- Desktop APIs as a plugin (not in core)
- Reference desktop app showing a Dreego app in a window
- No local HTTP server needed

### Effort
Medium. Wails v3 embeds a WebView. Dreego calls `app.Render()` instead of HTTP.
Desktop APIs (file dialog, tray, clipboard) as a plugin.

---

## v0.9 — Stabilization + TypeScript

### TypeScript
Per ADR `typescript-v2.md`: `lang="ts"` is reserved. In v0.9 it is activated —
esbuild compiles TS in `<script lang="ts">` blocks to JS. Build-hook like
plugin-tailwind.

### Stabilization
- Convert ~60 compile-only integration tests to behavioral assertions
- Plugin-contract validation: Auth + UI + Infrastructure plugins test different capabilities
- Performance benchmarks
- Fuzzing
- Race detection in CI for all plugins

---

## v1.0 — Stable Promise

v1 is not "the big new thing." v1 is "now everything is stable and promised."

### What v1 means
1. **Plugin contract stable** — `Register(app, typedOptions) error` is the promised interface
2. **Transpiler extension points stable** — Processor interface is promised
3. **Non-HTTP rendering stable** — `app.Render()` is a promised API
4. **Client reactivity optional** — If v0.6 experiment succeeded, it is an optional plugin
5. **SSG + Wails stable** — Both build on v0.2 and are production-ready
6. **Accessibility guaranteed** — CLI, docs, blueprints, components all accessible

### What v1 is NOT
- No new language
- No new core API
- No architecture revolution
- Everything built in v0.x is now promised

---

## v1.x — Production Hardening

### What comes
- **Mobile** — Explore mobile targets
- **Content Collections** — Blog posts, tag pages, pagination for SSG
- **Multi-site** — One Dreego instance, multiple websites
- **i18n** — Internationalization as core feature or plugin
- **DevTools/LSP** — Language server for `.dreego` files
- **Observability** — Prometheus + OpenTelemetry as plugins
- **Auth** — Login/logout, OAuth, session management as plugin

### What v1.x is NOT
- No breaking changes to core
- No new core packages (everything via plugins)
- No transpiler architecture change (extension points suffice)

---

## v2.x — Ecosystem Maturity

### What comes
- **Target-agnostic rendering** — One `.dreego` file → Go (SSR) + JS (Client) + Wails (Desktop) + SSG (Static)
- **Transpiler plugins mature** — Lua, Markdown, SVG, Charts, Mermaid as established plugins
- **Mobile** — Go-powered mobile targets
- **Enterprise** — Multi-tenant, RBAC, audit logs as plugins
- **Plugin Marketplace** — Discovery, installation, versioning

### What v2.x is NOT
- No rewrite. Dreego v2 builds on v1, does not replace it.
- No new template language. `.dreego` stays `.dreego`.

---

## What Cannot Be Done Via Plugins

| Feature | Needs Core/Transpiler? | v0.x Version |
|---------|------------------------|-------------|
| Non-HTTP Rendering | Core (SSRContext + App) | v0.2 |
| Transpiler Extension Points | Transpiler (Lexer/Parser/Codegen) | v0.3 |
| Lua/Markdown Sections | Extension Points (plugin uses them) | v0.4 |
| Plugin Lifecycle Hooks | Core (App) | v0.5 |
| Streaming SSR | Transpiler (Codegen) | v0.5 |
| Client-side Codegen | Transpiler (second target) | v0.6 |
| SSG | Core (Non-HTTP Rendering) | v0.7 |
| Wails | Core (Non-HTTP Rendering) | v0.8 |
| TypeScript | Build-Hook (plugin possible) | v0.9 |
| Stabilization | Tests, docs | v0.9 |

Everything else is a plugin. Auth, SSE, WebSocket, Tailwind, observability,
i18n, search, SEO, charts, maps, mail, payments, PWA, cache, jobs — all
plugins, no core code.

---

## Non-Binding

This roadmap is directional and non-binding. Priorities may shift based on
real-world usage, plugin ecosystem feedback, and community input.