# Blockwebchain Index

Generated: 2026-07-26T04:22:37Z

Chain: 1–13 | Next status code: **14**

## CHAIN (History)
- `01` **transpiler.1** — Transpiler-Pipeline (Lexer → Parser → AST → CodeGen)
- `02` **routing.1** — File-based Routing
- `03` **layout.1** — Layout-System ({#slot} + {#head})
- `04` **middleware.1** — Middleware-System (RequestLogging + Redirect/Rewrite)
- `05` **cli.1** — dreego CLI (generate, build, run)
- `06` **config.1** — dreego/config.json (Redirects, Rewrites, Logging)
- `07` **context-refactoring.1** — Context Interface + SSRContext (map[string]string → any)
- `08` **recovery.1** — Recovery-Middleware (Panic → 500)
- `09` **xss.1** — XSS Auto-Escaping ({variable} → html.EscapeString)
- `10` **error-pages.1** — Custom Error-Pages (404.dreego + 500.dreego)
- `11` **bracket-routes.1** — [id] Brackets für dynamische Segmente
- `12` **route-groups.1** — (group)/ Route Groups (unsichtbar in URL)
- `13` **flat-gen.1** — Flat Gen-Package (gen/routes.go statt per-dir dree.go)

## AVAILABLE NEXT
- **api-json.1** — API-Routen + JSON Responses
- **compression.1** — Gzip/Brotli Compression Middleware
- **deployment.1** — Deployment-Strategie (Docker, Single-Binary, Graceful)
- **documentation.1** — docs.dreego.dev + Tutorial + Examples
- **dreegotest.1** — dreegotest — Testing-Package
- **each-loop.1** — {#each} mit $loop-Variable
- **health-checks.1** — /health + /ready Endpoints
- **hot-reload.1** — Hot Reload (Dev-Server + SSE)
- **observability.1** — Observability (Prometheus, OpenTelemetry, Request-ID)
- **plugin-interface.1** — Plugin-Interface (Frozen for v1)
- **scaffolding.1** — dreego new + Generatoren
- **security-headers.1** — Security-Header (CSP, HSTS, X-Frame, X-Content-Type)
- **session.1** — Session-Interface (Cookie Store im Core)
- **static-assets.1** — Static Assets (static/ → embed.FS)
- **template-filters.1** — Template-Filter ({var|raw}, {var|upper})
- **verbatim.1** — {#verbatim} Block (Raw-Output)

## BLOCKED
- **addon-ecosystem.1** — Addon-Ökosystem (auth, ui, admin, db)  (missing: plugin-interface.1, components.1, session.1)
- **api-swagger.1** — Swagger/OpenAPI Auto-Generation  (missing: api-json.1)
- **cache-interface.1** — Caching Interface (Memory, Redis)  (missing: plugin-interface.1)
- **components.1** — Component-System ({#use}, props)  (missing: plugin-interface.1)
- **csrf.1** — CSRF-Schutz (Core-Conditional)  (missing: session.1)
- **ddos-protection.1** — DDoS-Schutz (PoW + Rate-Limiting) — Plugin  (missing: plugin-interface.1, middleware-hooks.1)
- **devtools.1** — DevTools (LSP, VS Code, CLI-Niceties)  (missing: plugin-interface.1)
- **dreego-analytics.1** — dreego-analytics (Privacy-friendly, Server-Side)  (missing: plugin-interface.1, middleware-hooks.1)
- **dreego-cache.1** — dreego-cache (Caching: Memory, Redis)  (missing: plugin-interface.1, cache-interface.1)
- **dreego-charts.1** — dreego-charts (Chart.js/Canvas Components)  (missing: plugin-interface.1, components.1)
- **dreego-features.1** — dreego-features (Feature-Flags, A/B-Testing)  (missing: plugin-interface.1, middleware-hooks.1)
- **dreego-i18n.1** — dreego-i18n (Internationalisierung)  (missing: plugin-interface.1, middleware-hooks.1)
- **dreego-icons.1** — dreego-icons (Lucide/Heroicons Components)  (missing: plugin-interface.1, components.1)
- **dreego-jobs.1** — dreego-jobs (Background-Jobs, Cron, Queue)  (missing: plugin-interface.1, queue-interface.1)
- **dreego-mail.1** — dreego-mail (E-Mail SMTP/Resend/Postmark)  (missing: plugin-interface.1, email-interface.1)
- **dreego-map.1** — dreego-map (MapLibre/Leaflet Components)  (missing: plugin-interface.1, components.1)
- **dreego-markdown.1** — dreego-markdown (Markdown-Rendering, Frontmatter)  (missing: plugin-interface.1)
- **dreego-notify.1** — dreego-notify (Multi-Channel Notifications)  (missing: plugin-interface.1, email-interface.1, event-bus.1)
- **dreego-pdf.1** — dreego-pdf (PDF-Generierung aus HTML)  (missing: plugin-interface.1)
- **dreego-polar.1** — dreego-polar (Payments via Polar.sh)  (missing: plugin-interface.1)
- **dreego-pwa.1** — dreego-pwa (Service Worker, Offline-Caching)  (missing: plugin-interface.1)
- **dreego-search.1** — dreego-search (Volltextsuche)  (missing: plugin-interface.1)
- **dreego-seo.1** — dreego-seo (Meta-Tags, OG, JSON-LD, Sitemap)  (missing: plugin-interface.1, middleware-hooks.1)
- **dreego-storage.1** — dreego-storage (File-Uploads, Progress, Resize)  (missing: plugin-interface.1, storage-interface.1)
- **email-interface.1** — Email-Sending Interface (SMTP, Resend, Postmark)  (missing: plugin-interface.1)
- **event-bus.1** — Pub/Sub Event-Bus (Core-Interface)  (missing: plugin-interface.1)
- **form-actions.1** — Form Actions (g-action / g-submit)  (missing: csrf.1)
- **middleware-hooks.1** — Plugin-Middleware-Hooks (app.Use FIFO)  (missing: plugin-interface.1)
- **queue-interface.1** — Background-Job-Queue Interface  (missing: plugin-interface.1)
- **route-hooks.1** — Plugin-Route-Registration  (missing: plugin-interface.1)
- **ssg.1** — Static Site Generation (SSG)  (missing: plugin-interface.1)
- **storage-interface.1** — File-Storage Interface (S3, R2, Local)  (missing: plugin-interface.1)
- **wails.1** — Wails Desktop Integration  (missing: plugin-interface.1)

chain: 13 | web: 49 | next code: 14
