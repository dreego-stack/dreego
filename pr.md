---
version: patch
---

- Feat: `dreego docs` reads local module docs via `go.mod` (module itself → `vendor/` → module cache) instead of an embedded copy or HTTP — no network, single source of truth per module
- Feat: `dreego docs -p <plugin>` reads a plugin's docs from `github.com/dreego-stack/<plugin>` (requires it in `go.mod`)
- Feat: `dreego docs --list` lists every core + plugin page from each module's `_docs/sitemap.json`
- Feat: per-module `_docs/sitemap.json` defines each module's doc pages (replaces the embedded mirror)
- Fix: remove `cli/dreego/embedded/`, `cli/dreego/embed.go`, `_scripts/sync-embedded-docs.sh` and the duplicated docs copy
- Fix: `run-timer-sigterm` flake — timer logic extracted into `scheduleStop` with deterministic unit tests (`TestScheduleStopSendsSIGTERM`, `TestScheduleStopFallsBackToKill`); integration test uses prebuilt CLI + 1s timer instead of double `go run` + 3s wait
