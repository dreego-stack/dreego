# Blockwebchain Index

Generated: 2026-07-26T04:12:41Z

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
- **components.1** — Component-System ({#use}, props)  (missing: plugin-interface.1)
- **csrf.1** — CSRF-Schutz (Core-Conditional)  (missing: session.1)
- **ddos-protection.1** — DDoS-Schutz (PoW + Rate-Limiting) — Plugin  (missing: plugin-interface.1)
- **devtools.1** — DevTools (LSP, VS Code, CLI-Niceties)  (missing: plugin-interface.1)
- **form-actions.1** — Form Actions (g-action / g-submit)  (missing: csrf.1)
- **ssg.1** — Static Site Generation (SSG)  (missing: plugin-interface.1)
- **wails.1** — Wails Desktop Integration  (missing: plugin-interface.1)

chain: 13 | web: 25 | next code: 14
