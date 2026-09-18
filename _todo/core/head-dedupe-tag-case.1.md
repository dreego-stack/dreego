# Head dedupe tag detection is case-sensitive

- Area: compiler
- Phase: head dedupe
- Goal: make dedupe tag detection consistent with case-insensitive HTML.
- Gap: `HasHeadDedupeTag` (`internal/transpiler/ir/head_prefix.go:70-74`) matches only lowercase `<title`, `name="description"` and `name='description'`; `<TITLE>` or `NAME="description"` bypasses dedupe and produces two head tags.
- Acceptance: decide case-insensitive matching or document that lowercase is required; a test pins the choice.
- Depends on: nothing.
