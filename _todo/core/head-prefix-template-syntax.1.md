# Dedupe silently disabled for non-literal layout head prefixes

- Area: compiler
- Phase: head dedupe
- Goal: make the dedupe fallback for non-literal layout heads explicit instead of silently disabling dedupe.
- Gap: `ir.File.StaticBodyHead` (`internal/transpiler/ir/head_prefix.go:12-37`) returns false when a non-text node precedes `{#head}` or the prefix contains `{#`, `{{` or `[[`; `StaticBodyHeadTail` (`internal/transpiler/ir/head_prefix.go:45-66`) applies the same limit. A layout head with an i18n message or a component call before the placeholder therefore emits two `<title>` tags, contradicting the exactly-one guarantee in `core/_docs/head-section.md` and `core/_docs/layouts.md`.
- Acceptance: either support those shapes or emit a generate-time diagnostic; the limitation is documented.
- Depends on: nothing.
