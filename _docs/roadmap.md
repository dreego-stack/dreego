# Roadmap

Dreego is an SSR-first full-stack web framework for Go. It combines Go's type safety, simplicity, performance, and low resource usage with an intuitive component and file-based development experience inspired by Svelte and Astro.

Accessibility is a core design principle. Dreego aims to provide screen-reader-friendly tools, understandable diagnostics, semantic generated markup, keyboard-operable defaults, and honest guidance about what applications must still verify themselves.

This roadmap communicates direction, not release promises. Priorities may change as the core is tested in real applications.

## Planned for v0.0.x

The current priority is a small, dependable SSR core.

- Harden routing, rendering, components, layouts, forms, sessions, middleware, and the CLI.
- Replace package-global runtime state with an explicit `App` instance before the public API is stabilized.
- Remove speculative event bus, queue, key-value, and storage APIs from core; revisit shared contracts only after real plugin implementations.
- Improve diagnostics, documentation, accessibility, tests, fuzzing, race detection, and benchmarks.
- Establish accessibility checks for the CLI, documentation, generated blueprints, and official component defaults.
- Add a small set of reference applications that double as documentation and black-box tests of the complete CLI-to-HTTP workflow.
- Keep generated applications close to standard Go and `net/http`.
- Support progressive enhancement through HTMX, Alpine.js, and plain JavaScript.
- Identify the extension points needed by future plugins without promising a stable plugin contract before real implementations exist.
- Avoid an internal client runtime, SSG implementation, or expanded Wails target during core stabilization.

Concrete near-term work lives in [_todo/core](https://github.com/dreego-stack/dreego/tree/main/_todo/core), with one item per file.

## Planned between v0.1 and v1

The production focus remains SSR while the stable core is exercised by real applications and plugins.

- Maintain compatibility and document unavoidable breaking changes clearly.
- Build real plugins against the public interfaces and adjust the pre-v1 contracts when evidence requires it.
- Explore islands or hydration outside the stable core.
- Promote client-side capabilities only when they preserve Dreego's simplicity and do not require Go developers to adopt a complex JavaScript toolchain.
- Continue improving the SSR developer experience, production safety, and editor tooling.

## Planned for v1.x

SSR remains the dependable default. Work may expand to additional targets once the core contract is stable.

- Improve the existing Wails integration for desktop applications.
- Develop static site generation for targets such as Cloudflare Pages and GitHub Pages.
- Evaluate a supported islands or hydration model based on the pre-v1 experiments.
- Preserve shared `.dreego` components where target differences can remain explicit and understandable.

## Planned for v2 or later

Long-term work may make SSR, static sites, desktop applications, and future platforms first-class targets without weakening the Go-first model.

- Mature the target-agnostic rendering pipeline.
- Expand static content, collections, and deployment workflows.
- Explore mobile or other Go-powered application targets.
- Add transpiler extension points only after their contracts are proven by real use cases.

## Planned plugins

Plugins live in separate Go modules and are not part of the core roadmap. Individual idea files live in [_todo/plugins](https://github.com/dreego-stack/dreego/tree/main/_todo/plugins). This list is a readable overview, not a commitment or required sequence.

- Auth
- UI components and icons
- Admin tools
- Database integration and migrations
- Analytics
- Cache
- Charts
- Cluster and distributed state
- Feature flags and experiments
- Feedback
- Internationalization
- Background jobs and scheduling
- Mail
- Maps
- Markdown and content
- Notifications
- PDF generation
- Payments
- Progressive web app support
- Search
- SEO, structured data, and sitemaps
- File storage, uploads, and image processing
- DDoS protection and rate limiting
- Tailwind integration
- Developer tools, language server, and editor support
- Optional runtime scripting

Plugin ideas should remain outside core unless a small, provider-neutral interface is required by multiple proven implementations.
