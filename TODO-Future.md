# TODO Future

Long-term ideas and larger architecture goals without a concrete near-term plan. Concrete, planned work lives in [TODO.md](TODO.md).

## Architecture (V2)

- **ssg.1** — Static Site Generation: `dreego build --static`, compile-time generation of all pages as static HTML, no runtime server. Frontmatter support, content collections. ADR: `ssg-wails-v2.md`
- **wails.1** — Wails Desktop Integration: `dreego build --wails`, WebView-based desktop apps, WailsContext (no `http.Request`), target-agnostic renderer, system API access. ADR: `ssg-wails-v2.md`
- **client-reactivity.1** — Client-side interactivity research (Alpine/islands/custom runtime)
- **transpiler-extensions.1** — Plugin-registered section tags: extend the 5 built-in tags (`go`, `div`, `head`, `script`, `style`) with plugin-defined ones (e.g. `<markdown>`, `<svg>`, `<chart>`) via a `SectionProcessor` interface. Requires the formal Lexer→Parser→AST→CodeGen pipeline first (see `transpiler-pipeline.md`). Plugin discovery reuses the `dreego docs` pattern (go.mod Requires → `findModDir` → module cache). V2.

## Ecosystem

### Plugin ideas

- **plugin-ecosystem.1** — Plugin ecosystem (auth, ui, admin, db)
- **dreego-analytics.1** — Privacy-friendly, server-side analytics
- **dreego-cache.1** — Caching (Memory, Redis)
- **dreego-charts.1** — Chart.js/Canvas components
- **dreego-cluster.1** — Multi-node, distributed state
- **dreego-features.1** — Feature flags, A/B testing
- **dreego-feedback.1** — Feedback POST endpoint
- **dreego-i18n.1** — Internationalization
- **dreego-icons.1** — Lucide/Heroicons components
- **dreego-jobs.1** — Background jobs, cron, queue
- **dreego-mail.1** — Email sending (SMTP/Resend/Postmark)
- **dreego-map.1** — MapLibre/Leaflet components
- **dreego-markdown.1** — Markdown rendering + frontmatter
- **dreego-notify.1** — Multi-channel notifications
- **dreego-pdf.1** — PDF generation from HTML
- **dreego-polar.1** — Payments via Polar.sh
- **dreego-pwa.1** — Service worker, offline caching
- **dreego-search.1** — Full-text search
- **dreego-seo.1** — Meta tags, OG, JSON-LD, sitemap
- **dreego-storage.1** — File uploads, progress, resize
- **cache-interface.1** — Caching interface (Memory, Redis) — prerequisite for `dreego-cache`
- **email-interface.1** — Email sending interface — prerequisite for `dreego-mail`
- **ddos-protection.1** — DDoS protection (PoW + rate-limiting)
- **tailwind-plugin.1** — Tailwind CSS build plugin
- **devtools.1** — DevTools (LSP, VS Code, CLI niceties)
- **dreego-scripting.1** — Optional runtime scripting via `flow` (github.com/LukasLow/flow, MPL-2.0): embed flow as a plugin to let users run typed scripts without recompiling. Separate concern from the compile-time transpiler — flow is a runtime language, so it does NOT extend `dreego generate` (that stays Go-only). Only relevant if runtime-scriptable logic is ever wanted; keep out of core.
