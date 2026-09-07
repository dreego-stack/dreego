---
area: transpiler
phase: multi-language-dreego
---
# Shared JavaScript output stage

## Goal

Create one normalized JavaScript artifact stage shared by JavaScript,
TypeScript, and Lua client inputs before implementing either compiler-backed
language processor.

## Acceptance criteria

- Existing `<client lang="js">` behavior is preserved through the new stage.
- The stage owns deterministic JavaScript assets, diagnostics, and optional
  source-map metadata without knowing the input language.
- Input processors cannot write generated application files directly.
- Integration tests prove identical JavaScript output before and after the
  refactor.
