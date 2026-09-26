---
type: Decision
title: Explicit Dreefile header grammar
description: Uppercase header directives replace the legacy Component, import, and from forms; the legacy forms are rejected
tags: [transpiler, grammar, dreefile, v0.x]
timestamp: 2026-09-18T00:00:00Z
---
# Explicit Dreefile header grammar

**Date:** 2026-09-18
**Status:** Accepted; grammar parsing and formatting and the `GOIMPORT` codegen consumer are implemented, the `LAYOUT` codegen consumer follows in a later slice

## Context

The original `.dreego` header used `Component Name (props)` for a component
declaration and Python-style `import` / `from` lines for component and Go
imports. That surface mixed three unrelated concerns, made the component name
depend on a declaration instead of the file, and left the file kind implicit in
the directory name. It also collided with the need for a single explicit Go
import channel, because a top-level `import "…"` line was ambiguous between a
component source and a Go package.

## Decision

A `.dreego` file declares its header with uppercase keywords and lowercase
values. One rule governs the surface: **uppercase is a keyword, lowercase is a
value.**

```text
DREEFILE component (title string)     component kind; props on the DREEFILE line
DREEFILE layout                       layout kind
DREEFILE page                         page kind (a missing DREEFILE is also a page)
LAYOUT "www/layouts/admin.dreego"     explicit layout path
COMPONENT "www/components" IMPORT { Card, Card as ProductCard }
GOIMPORT { sync, encoding/json }      Go import channel
```

- A component's name comes from its **filename**: `Card.dreego` becomes
  `<@Card>`.
- `COMPONENT ... IMPORT` imports components from a path; `as` creates a
  generator-global alias, consistent with the global component registry.
- `LAYOUT` is parsed and reserved; its codegen consumer arrives with the
  layout-chaining slice.
- `GOIMPORT` emits Go imports into the generated route, component, and layout
  package. Any package that resolves in the application is accepted: the
  standard library, packages in the application's own module, and dependencies
  listed in `go.mod`. An entry may carry an alias (`alias "path"`). An
  unresolvable package fails at `dreego generate` with a diagnostic naming the
  path, `not in go.mod`, and the matching `go get` command; a base-name or alias
  collision fails with a diagnostic naming both paths. Type safety is preserved
  by resolving against `go.mod` rather than executing arbitrary packages
  dynamically.

The legacy `Component X (..)`, top-level `import "…"`, and
`from "…" import {..}` forms are rejected at `dreego generate` with a
`file:line:col` diagnostic naming the replacement. This is a deliberate
breaking change before v0.1; the repository sources and documentation are
migrated in the same series.

## Consequences

- The file kind is explicit in the file instead of inferred from the directory.
- Component call names are derived from filenames, so a file rename changes the
  callable name.
- `GOIMPORT` becomes the single channel for generated Go imports, which moves
  Go imports into `<server>` and re-scopes the server-import work to an explicit
  channel resolved against the application's `go.mod`. v0.10.9 removes the
  earlier standard-library allowlist and adds the alias form; the grammar
  keywords above are unchanged.
- The migration cost is deliberately paid before v0.1; no compatibility aliases
  are provided.
- `dreego fmt` preserves legacy header lines verbatim rather than silently
  rewriting them, so a rejected file is not corrupted before its author fixes
  it.
- `dreego fmt` is a formatter, not a validator: `dreego generate` remains the
  single gate that rejects a legacy header. `dreego fmt --check` therefore
  exits 0 on a file that `dreego generate` rejects, and this divergence is
  deliberate. A formatter that refused legacy files would block formatting an
  entire repository because of one unmigrated file, while still not telling the
  author how to migrate it; `generate` already fails early with a
  `file:line:col` diagnostic naming the replacement.
- `dreego fmt` canonicalizes the header to at most one blank line between
  directives. Duplicate blank lines in the header are collapsed to a single
  blank line so the header is stable under repeated formatting.

## See also

- [Phase: explicit Dreefile composition](../../_plan/phase-dreefile.md)
- [Components](../components.md)
- [File Anatomy](../file-anatomy.md)
