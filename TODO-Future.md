# TODO Future

Long-term ideas and larger architecture goals without a concrete near-term plan. Concrete, planned work lives in [TODO.md](TODO.md).

## Architecture (V2)

- **ssg.1** — Static Site Generation: `dreego build --static`, compile-time generation of all pages as static HTML, no runtime server. Frontmatter support, content collections. ADR: `ssg-wails-v2.md`
- **wails.1** — Wails Desktop Integration: `dreego build --wails`, WebView-based desktop apps, WailsContext (no `http.Request`), target-agnostic renderer, system API access. ADR: `ssg-wails-v2.md`
- **client-reactivity.1** — Client-side interactivity research (Alpine/islands/custom runtime)

## Ecosystem

### Naming decision

- **addon vs plugin** — decide whether the ecosystem is called "addons" or "plugins"; both currently describe the same thing (tracked in TODO.md as a decision entry)

### Plugin ideas

- **addon-ecosystem.1** — Addon ecosystem (auth, ui, admin, db)
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
