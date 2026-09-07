---
area: transpiler
phase: multi-language-dreego
---
# TypeScript client processor

## Goal

Implement first-party `<client lang="ts">` through a pinned, explicitly
approved TypeScript toolchain and the shared JavaScript output stage.

## Acceptance criteria

- Real TypeScript type checking fails generation with `.dreego` source ranges.
- Builds are reproducible and never install an unpinned latest version.
- Developers do not maintain a project-owned npm pipeline for the standard
  workflow.
- Generated JavaScript, assets, and source maps use the common output stage.
