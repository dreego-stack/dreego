---
area: transpiler
phase: pre-v0.1
---
# Rename root sections by semantic purpose

## Goal

Atomically replace the language-coupled root section names with semantic names
before the `.dreego` language receives its v0.1 promise:

```text
<go>     -> <server>
<div>    -> <body>
<script> -> <client>
```

`<head>` and `<style>` remain unchanged. Default languages are Go, HTML, CSS,
and JavaScript. The first implementation may accept only the defaults while
preserving the `lang` field in the parsed model for later processors.

## Required semantics

- Exactly one body section is allowed per component or renderable route method.
- A root `<client>` section is client source. A `<script>` nested inside
  `<body lang="html">` remains ordinary HTML and is emitted normally.
- `<@Component>`, template expressions, escaped output, slots, `{#if}`, and
  `{#each}` remain owned by Dreego rather than a future body-language processor.
- Unknown section/language pairs fail with path, line, column, and remediation.
- No compatibility aliases remain after the migration. The generated code,
  scaffolds, fixtures, and docs move together.

## Acceptance criteria

- Add permanent integration tests that fail against the old parser and prove
  the new five-section model.
- Add a regression test for a nested HTML `<script>` inside `<body>`.
- Add parser and generator tests for omitted and explicit default `lang`
  attributes.
- Reject legacy root `<go>`, `<div>`, and `<script>` with a migration diagnostic.
- Update lexer/parser/AST/codegen, formatter ordering, scaffolds, fixtures,
  examples, syntax highlighting contract, and every current user document.
- Add a migration section mapping all old syntax to new syntax.
- Keep handwritten files below 300 lines and preserve dependency checks.
- `make test` and race checks pass through `smd`.

## Dependencies

- The release and tag history must be coherent before the breaking migration
  is released.
- The accepted design is recorded in
  `_docs/decisions/semantic-sections-and-language-processors.md`.

## Out of scope

- TypeScript, Markdown, or Lua processor implementation.
- Multiple body sections.
- DreeJS behavior.
- The v0.2 root-package and render-foundation migration.
