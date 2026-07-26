# Blockwebchain Index

Generated: 2026-07-26T03:58:06Z

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
- **session.1** — Session-Interface (Cookie Store im Core)
- **static-assets.1** — Static Assets (static/ → embed.FS)

## BLOCKED
- **csrf.1** — CSRF-Schutz (Core-Conditional)  (missing: session.1)
- **form-actions.1** — Form Actions (g-action / g-submit)  (missing: csrf.1)

chain: 13 | web: 4 | next code: 14
