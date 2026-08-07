---
id: tailwind-plugin.1
title: Tailwind CSS Build Plugin
status: planned
phase: v0.0.26
requires:
  - plugin-interface.1
created: 2026-08-07
changed: 2026-08-07
---

# Tailwind CSS Build Plugin

Compile Tailwind classes from .dreego templates into a static CSS file
during `dreego build`, replacing the Tailwind CDN with a self-hosted,
offline-friendly build step.

## Research summary (2026-08-07, web research with sources)

- Tailwind v4 is MIT licensed (tailwindcss, @tailwindcss/cli,
  @tailwindcss/oxide all MIT) — embedding in an MPL-2.0 framework is
  fine; only a MIT notice (THIRD_PARTY_NOTICES) is required.
- No pure-Go Tailwind compiler exists. All Go integrations wrap the
  official binary/CLI (e.g. kikihakiem/tailwindcss-go downloads the
  standalone binary, scriptogre/tailwindcss-go-tool uses the Go 1.24
  `go tool` feature).
- Tailwind v4 uses automatic content detection (text tokenization) —
  .dreego files can be scanned directly; classes live in class="..."
  attributes. No need to scan generated Go code.
- Standalone binary (no Node, ~80-110 MB) exists; Hugo removed
  standalone support in v0.161 (maintenance concern) — npm CLI is the
  more maintained path.
- Hugo is the best reference: it writes a class-statistics file
  (hugo_stats.json) and feeds it to the Tailwind CLI via @source.

## Plan

1. Official plugin `plugins/tailwind/` (own go.mod, stdlib only).
2. Build step in `dreego build` (or a `dreego build --tailwind` flag):
   - generate input.css with `@import "tailwindcss"` + `@source` on the
     dreego/ directory (or automatic content detection)
   - ensure the standalone binary (download + cache, pattern from
     kikihakiem/tailwindcss-go) or npm CLI
   - run `tailwindcss -i input.css -o static/css/output.css --minify`
   - embed the result via //go:embed
3. Fallback without a build step: document Open Props or Pico CSS
   (both MIT, zero-build) as alternatives.
4. Tests: plugin interface satisfaction, build step produces CSS with
   classes found in a .dreego file, offline behavior.

## Sources

- https://tailwindcss.com/blog/tailwindcss-v4
- https://tailwindcss.com/docs/detecting-classes-in-source-files
- https://tailwindcss.com/docs/installation/tailwind-cli
- https://github.com/tailwindlabs/tailwindcss/releases
- https://gohugo.io/functions/css/tailwindcss
- https://github.com/kikihakiem/tailwindcss-go
- https://github.com/scriptogre/tailwindcss-go-tool
- https://registry.npmjs.org/tailwindcss/latest (MIT)
- https://www.mozilla.org/en-US/MPL/2.0/FAQ (MPL+MIT compatibility)
