---
version: patch
---

- Feat: `dreego generate --check` now verifies generated content byte-for-byte (no working-tree modification, no mtime-based false passes) with a path-level diff for missing, extra, and stale generated files.
- Feat: source discovery is restricted to the project's `dreego/` root; `vendor`, `node_modules`, nested `subapp/dreego`, and other out-of-root `routes`/`components`/`layouts` directories are ignored.
- Feat: layout lookup is route-local and cascades through parent route scopes; `default.dreego` and `layout.dreego` in the same `layouts` directory is an ambiguous-layout error.
- Feat: duplicate catch-all error-page patterns (e.g. `404.dreego` in `routes/` and `routes/(group)/`) now fail `dreego generate` with a diagnostic naming both source paths.
- Refactor: `Run` is split into `buildPlan` + `applyPlan` so normal and `--check` modes share one generation implementation (no drift).