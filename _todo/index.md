# Blockwebchain Index

Generated: 2026-07-26T04:01:30Z

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
- **deployment.1** — Deployment-Strategie (Docker, Single-Binary, Graceful)
- **dreegotest.1** — dreegotest — Testing-Package
- **each-loop.1** — {#each} mit $loop-Variable
- **hot-reload.1** — Hot Reload (Dev-Server + SSE)
- **plugin-interface.1** — Plugin-Interface (Frozen for v1)
- **scaffolding.1** — dreego new + Generatoren
- **session.1** — Session-Interface (Cookie Store im Core)
- **static-assets.1** — Static Assets (static/ → embed.FS)
- **template-filters.1** — Template-Filter ({var|raw}, {var|upper})
- **verbatim.1** — {#verbatim} Block (Raw-Output)

## BLOCKED
- **addon-ecosystem.1** — Addon-Ökosystem (auth, ui, admin, db)  (missing: plugin-interface.1, components.1, session.1)
- **components.1** — Component-System ({#use}, props)  (missing: plugin-interface.1)
- **csrf.1** — CSRF-Schutz (Core-Conditional)  (missing: session.1)
- **devtools.1** — DevTools (LSP, VS Code, CLI-Niceties)  (missing: plugin-interface.1)
- **form-actions.1** — Form Actions (g-action / g-submit)  (missing: csrf.1)
- **ssg.1** — Static Site Generation (SSG)  (missing: plugin-interface.1)
- **wails.1** — Wails Desktop Integration  (missing: plugin-interface.1)

chain: 13 | web: 17 | next code: 14
