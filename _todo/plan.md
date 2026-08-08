# Dreego Linear Plan

Loose, reorderable timeline. Versioning stays at `v0.x.y` for the foreseeable future. `v0.x` marks larger milestones and may include breaking changes; `y` is continuous. `v1.0.0` is reserved for a stable, trustworthy release and is not a near-term target.

Official plugins live in separate repos under `github.com/dreego-stack/`. Each plugin has its own `go.mod` and requires `github.com/dreego-stack/dreego`. Core must never import a plugin package — plugins depend on Core, never the other way around.

## Done

- **v0.0.20** — `security-cookie.1`, `security-csp.1`
- **v0.0.21** — monorepo plugin layout (`plugins/`, `go.work`, test moves)
- **v0.0.22** — `servemux-cache.1`, `codegen-errors.1`, `security-session.1`
- **v0.0.23** — `dreegotest.1` (planned), `golden-tests-core.1` (planned), `typed-forms.1` (planned)
- **v0.0.24** — `scaffold-fix.1`, `layout-head.1`, `scoped-css.2`, `component-attr-props.1`, `typed-forms.1`, `dreegotest.1`, `golden-tests-core.1` + deterministic port-schema runner
- **v0.0.25** — `plugin-interface.1`, `middleware-hooks.1`, `route-hooks.1`, `docs-extensibility.1`, `docs-embed.1`, `frontmatter.1`, `dev-server.1` ✅
- **v0.0.26** — post-v0.0.25 bugfix batch + SSE example plugin + `--version`/`-v` flags ✅

## unlock: plugin ecosystem

After plugin-interface.1 + middleware-hooks.1 + route-hooks.1 the following plugin blocks unblock. They are tracked here for visibility and implemented as separate repos under `github.com/dreego-stack/` (each with its own `go.mod`):

- dreego-auth
- dreego-db
- dreego-cache
- dreego-storage
- dreego-i18n
- dreego-seo
- dreego-mail
- dreego-jobs
- dreego-analytics
- dreego-features
- dreego-pwa
- dreego-markdown
- dreego-icons
- dreego-charts
- dreego-map
- dreego-search
- dreego-pdf
- dreego-polar
- dreego-admin
- dreego-devtools

## v0.0.26

- `documentation.1` — docs site, tutorial, examples
- `api-swagger.1` — Swagger/OpenAPI auto-generation
- `observability.1` — core request-id done; plugin metrics/tracing later

## v0.0.30 area

- `ssg.1` — `dreego build --static`
- `wails.1` — `dreego build --wails`
- `client-reactivity.1` — client-side reactivity research and prototype

## Notes

- Core must never import a plugin package.
- Official plugins live in separate repos under `github.com/dreego-stack/` with their own `go.mod`.
- `dreego docs` resolves plugin docs from local `plugins/<name>/_docs/` (docs-extensibility.1, done in v0.0.25); external-repo docs for community plugins remain an open design question.
- This plan is intentionally linear and can be shifted as priorities change.