# Dreego Linear Plan

Loose, reorderable timeline. Versioning stays at `v0.x.y` for the foreseeable future. `v0.x` marks larger milestones and may include breaking changes; `y` is continuous. `v1.0.0` is reserved for a stable, trustworthy release and is not a near-term target.

Plugins live in separate repos and must not be required by core. Core must never import a plugin package.

## v0.0.20

- `security-cookie.1` — harden session and CSRF cookie flags
- `security-csp.1` — add Content-Security-Policy header

## v0.0.21

- `servemux-cache.1` — cache built middleware/router stack
- `codegen-errors.1` — replace silent CodeGen failures with errors
- `security-session.1` — document or encrypt session payload

## v0.0.22

- `dreegotest.1` — testing package for routes and components
- `golden-tests-core.1` — golden tests for generated Go code
- `typed-forms.1` — int/bool/slice binding, custom validators, improve `email` validator

## v0.0.23

- `frontmatter.1` — parse frontmatter and expose typed metadata
- `dev-server.1` — `dreego dev` with file watcher and auto-regenerate
- `docs-extensibility.1` — design how `dreego docs` can read plugin docs

## v0.0.24

- `plugin-interface.1` — frozen plugin contract
- `middleware-hooks.1` — plugin middleware hooks via app.Use
- `route-hooks.1` — plugin route registration

## unlock: plugin ecosystem

After plugin-interface.1 + middleware-hooks.1 + route-hooks.1 the following plugin blocks unblock. They are tracked here for visibility but implemented in separate repos under `codeberg.org/dreego/<name>`:

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

## v0.0.25

- `documentation.1` — docs site, tutorial, examples
- `api-swagger.1` — Swagger/OpenAPI auto-generation
- `observability.1` — core request-id done; plugin metrics/tracing later

## v0.0.30 area

- `ssg.1` — `dreego build --static`
- `wails.1` — `dreego build --wails`
- `client-reactivity.1` — client-side reactivity research and prototype

## Notes

- Core must never import a plugin package.
- `dreego docs` extensibility is an open design question: plugin docs could be embedded, read from local paths, or fetched from known URLs. Not decided yet.
- This plan is intentionally linear and can be shifted as priorities change.
