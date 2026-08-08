# Blockwebchain Index

Generated: 2026-08-08T12:53:50Z

Chain: 1–49 | Next status code: **50**

## CHAIN (History)
- `01` **transpiler.1** — Transpiler Pipeline (Lexer → Parser → AST → CodeGen)
- `02` **routing.1** — File-based Routing
- `03` **layout.1** — Layout System ({#slot} + {#head})
- `04` **middleware.1** — Middleware System (RequestLogging + Redirect/Rewrite)
- `05` **cli.1** — dreego CLI (generate, build, run)
- `06` **config.1** — dreego/config.json (Redirects, Rewrites, Logging)
- `07` **context-refactoring.1** — Context Interface + SSRContext (map[string]string → any)
- `08` **recovery.1** — Recovery Middleware (Panic → 500)
- `09` **xss.1** — XSS Auto-Escaping ({variable} → html.EscapeString)
- `10` **error-pages.1** — Custom Error Pages (404.dreego + 500.dreego)
- `11` **bracket-routes.1** — [id] Brackets for Dynamic Segments
- `12` **route-groups.1** — (group)/ Route Groups (invisible in URL)
- `13` **flat-gen.1** — Flat Gen Package (gen/routes.go instead of per-dir dree.go)
- `14` **session.1** — Session Interface (Cookie Store in Core)
- `15` **csrf.1** — CSRF Protection (Core-Conditional)
- `16` **ci-check.1** — dreego generate --check (CI Mode)
- `17` **components.1** — Component System ({#use}, props)
- `18` **component-handler.1** — ComponentHandler (Buffered Mode + Functional Options)
- `19` **named-slots.1** — Named Slots ({#slot header}...{/slot})
- `20` **each-loop.1** — {#each} with $loop variable
- `21` **verbatim.1** — {#verbatim} Block (Raw-Output)
- `22` **tag-prefix-fix.1** — scanTag: Tag Prefix Matching Fix (head vs header)
- `23` **template-filters.1** — Template Filters ({var|raw}, {var|upper})
- `24` **if-else.1** — {#else} in {#if}-Block
- `25` **each-else.1** — {#each else} — Empty List Fallback
- `26` **static-assets.1** — Static Assets (dreego/static/ → inline Handler)
- `27` **dreego-fmt.1** — dreego fmt (Formatter)
- `28` **scaffolding.1** — dreego new + Generators
- `29` **health-checks.1** — /health + /ready Endpoints
- `30` **security-headers.1** — Security Headers (nosniff, frame, referrer, permissions)
- `31` **compression.1** — Gzip Compression Middleware
- `32` **api-json.1** — Content-Type Routing (JSON, XML, Custom)
- `33` **form-actions.1** — Form Actions (g-action / g-submit)
- `34` **request-id.1** — Request-ID Middleware (X-Request-ID)
- `35` **security-cookie.1** — Harden Session and CSRF Cookie Flags
- `36` **security-csp.1** — Add Content-Security-Policy Header
- `37` **deployment.1** — Deployment Strategy (Docker, Single-Binary, Graceful)
- `38` **servemux-cache.1** — Cache Built Middleware/Router Stack
- `39` **codegen-errors.1** — Replace Silent CodeGen Failures with Errors
- `40` **security-session.1** — Document or Encrypt Session Payload
- `41` **typed-forms.1** — Typed Form Binding and Validation
- `42` **dreegotest.1** — dreegotest — Testing Package
- `43` **golden-tests-core.1** — Golden Code Tests for Generator Output
- `44` **plugin-interface.1** — Plugin Interface (Frozen for v1)
- `45` **middleware-hooks.1** — Plugin Middleware Hooks (app.Use FIFO)
- `46` **route-hooks.1** — Plugin Route Registration
- `47` **docs-extensibility.1** — Extensible dreego docs Command
- `48` **frontmatter.1** — Frontmatter Support in .dreego
- `49` **dev-server.1** — Dev Server with Hot Reload

## AVAILABLE NEXT
- **addon-ecosystem.1** — Addon Ecosystem (auth, ui, admin, db)
- **api-swagger.1** — Swagger/OpenAPI Auto-Generation
- **cache-interface.1** — Caching Interface (Memory, Redis)
- **client-reactivity.1** — Client-Side Reactivity for .dreego
- **ddos-protection.1** — DDoS Protection (PoW + Rate-Limiting) — Plugin
- **devtools.1** — DevTools (LSP, VS Code, CLI-Niceties)
- **documentation.1** — docs.dreego.dev + Tutorial + Examples
- **dreego-analytics.1** — dreego-analytics (Privacy-friendly, Server-Side)
- **dreego-charts.1** — dreego-charts (Chart.js/Canvas Components)
- **dreego-features.1** — dreego-features (Feature-Flags, A/B-Testing)
- **dreego-feedback.1** — dreego feedback (POST endpoint)
- **dreego-i18n.1** — dreego-i18n (Internationalization)
- **dreego-icons.1** — dreego-icons (Lucide/Heroicons Components)
- **dreego-map.1** — dreego-map (MapLibre/Leaflet Components)
- **dreego-markdown.1** — dreego-markdown (Markdown Rendering, Frontmatter)
- **dreego-pdf.1** — dreego-pdf (PDF Generation from HTML)
- **dreego-polar.1** — dreego-polar (Payments via Polar.sh)
- **dreego-pwa.1** — dreego-pwa (Service Worker, Offline-Caching)
- **dreego-search.1** — dreego-search (Full-Text Search)
- **dreego-seo.1** — dreego-seo (Meta-Tags, OG, JSON-LD, Sitemap)
- **email-interface.1** — Email Sending Interface (SMTP, Resend, Postmark)
- **event-bus.1** — Pub/Sub Event Bus (Core Interface)
- **golden-tests.1** — Golden File Tests for Generator
- **observability.1** — Observability (Request-ID, Metrics, Tracing)
- **queue-interface.1** — Background Job Queue Interface
- **ssg.1** — Static Site Generation (SSG)
- **storage-interface.1** — File Storage Interface (S3, R2, Local)
- **tailwind-plugin.1** — Tailwind CSS Build Plugin
- **wails.1** — Wails Desktop Integration

## BLOCKED
- **dreego-cache.1** — dreego-cache (Caching: Memory, Redis)  (missing: cache-interface.1)
- **dreego-cluster.1** — dreego-cluster (Multi-Node, Distributed State)  (missing: cache-interface.1, event-bus.1)
- **dreego-jobs.1** — dreego-jobs (Background Jobs, Cron, Queue)  (missing: queue-interface.1)
- **dreego-mail.1** — dreego-mail (Email SMTP/Resend/Postmark)  (missing: email-interface.1)
- **dreego-notify.1** — dreego-notify (Multi-Channel Notifications)  (missing: email-interface.1, event-bus.1)
- **dreego-storage.1** — dreego-storage (File Uploads, Progress, Resize)  (missing: storage-interface.1)

## REJECTED
- **hot-reload.1** — Hot Reload (Dev Server + SSE)
- **live-reload.1** — Live Reload Proxy (SSE + Script Injection)
- **smart-recompile.1** — Smart Recompile (Text vs Go Detection)

chain: 49 | web: 38 | next code: 50
