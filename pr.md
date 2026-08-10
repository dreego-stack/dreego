---
version: patch
---

- Feat: `dreego docs` reads local module docs via `go.mod` (module itself → `vendor/` → module cache) instead of an embedded copy or HTTP — no network, single source of truth per module
- Feat: `dreego docs -p <plugin>` reads a plugin's docs from `github.com/dreego-stack/<plugin>` (requires it in `go.mod`)
- Feat: `dreego docs --list` lists every core + plugin page from each module's `_docs/sitemap.json`
- Feat: per-module `_docs/sitemap.json` defines each module's doc pages (replaces the embedded mirror)
- Fix: remove `cli/dreego/embedded/`, `cli/dreego/embed.go`, `_scripts/sync-embedded-docs.sh` and the duplicated docs copy
- Feat: `find-binary` check in `make test` fails on any stray binary in the repo (NUL-byte heuristic; exceptions: `.DS_Store`, `.kilo/`, `.tmp/`, allowlisted image/font/archive types)
