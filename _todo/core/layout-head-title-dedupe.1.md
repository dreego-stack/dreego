# Layout head merge dedupes title only at root-level `<head>`

- Area: compiler
- Phase: layout and head composition
- Goal: make the head merge produce exactly one `<title>` (and one meta description) per page regardless of whether the layout declares `<head>` at root level or body level.
- Gap: the runtime title/meta dedupe in `internal/transpiler/html/output/templ.go:94-106` is guarded by `layout.File.Head != nil`, and `File.Head` is populated only for a root-level `<head>` section (`internal/transpiler/parser/parser.go:107`). Layouts written in the canonical body-level shape (`<body><html><head>…{#head}…</head>…`) therefore get `headPrefix == ""`, the dedupe never runs, and a layout `<title>` plus a route `<title>` emit two `<title>` elements.
- Gap: `web-minimal`, `demo-ssr/blog`, and the new `web-app` template all use the body-level shape. `web-app` worked around the bug by dropping its layout title; `web-minimal` and both demos still rely on a static layout title and would double-title as soon as a route adds one.
- Acceptance: a route `<title>` wins over the layout title in BOTH layout shapes; `web-minimal` and the demos keep working; a regression test covers the body-level shape; `_docs/getting-started.md`'s claim about layout/head composition stays true.
- Depends on: nothing (found while building the `web-app` template).
