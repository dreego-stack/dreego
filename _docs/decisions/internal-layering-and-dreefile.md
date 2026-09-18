---
type: Decision
title: Internal layering and the dreefile package
description: Root internal/ is the shared implementation layer; internal/dreefile/ holds the compiler with ir, dreecode, gogen, and sections
tags: [architecture, compiler, internal, dreefile]
timestamp: 2026-09-18T00:00:00Z
---
# Internal layering and the dreefile package

**Date:** 2026-09-18
**Status:** Accepted

## Context

Two structural problems accumulated as the project grew:

1. The compiler lives in `internal/transpiler/`, a flat package with lexer,
   parser, per-section processors, per-language processors, code generation,
   discovery, formatting, configuration, accessibility checks, and i18n side by
   side. Related concerns (for example the duplicated expression-kind and
   source-position helpers in the root package and in `ir/`) drifted apart.
   A generic `transpiler` name also describes a mechanism rather than the
   artifact: the compiler consumes `.dreego` sources and is the `dreefile`
   compiler.
2. The target-neutral runtime implementation sits under `core/internal/`
   (`core/internal/render|context|i18n|server`) while shared implementation
   (`internal/session`, `internal/validate`, `internal/middleware`,
   `internal/gomod`) already lives at the repository root. The same conceptual
   layer is therefore split across two places, and `core/` carries
   implementation next to its public facade.

The root module and `core` are separate Go modules. Go's internal-visibility
rule is path-prefix based: `github.com/dreego-stack/dreego/internal/...` is
importable from any package whose path begins with
`github.com/dreego-stack/dreego`, which includes the `core` module. This already
works today for `internal/session`, `internal/validate`, and
`internal/middleware`, so moving more shared implementation to the root
introduces no new mechanism.

## Decision

### Root `internal/` is the shared implementation layer

Root `internal/` is the general implementation layer, roughly the role a
project-local standard library plays. `core/`, `adapter/ssr/`,
`adapter/wails/`, `cmd/dreego/`, and `dreegotest/` stay slim
facades, hosts, and binaries over it.

- `core/internal/render|context|i18n|server` move to `internal/render`,
  `internal/context`, `internal/i18n`, `internal/server`.
- `core/` keeps only the public facade and re-export layer.
- `cmd/dreego/internal/templates` moves to `internal/templates`.

No import cycle may be introduced: root `internal/` must never import `core` as
a Go package. Code generation may reference the public import path only as an
emitted string.

### The compiler becomes `internal/dreefile/`

`internal/transpiler/` becomes `internal/dreefile/`. Central orchestration stays
in `dreefile/generate.go` (`Run`, `RunCheck`, `buildPlan`, `buildRootPlan`,
`buildRootFile`) plus `plan.go` and `support.go`.

Shared packages sit directly on the dreefile level:

- `dreefile/ir/` — AST and types, the single source of truth.
- `dreefile/dreecode/` — the mini template language semantics: `{#if}`,
  `{#each}`, `{{ }}`, filter pipes, `{#slot}`, and expression scanning.
  `dreecode` is **not** the lexer/parser/token stack; see the front-end stage
  below.
- `dreefile/gogen/` — Go emission helpers: `GoLiteral`, `ToPascalCase`,
  escaping, and source positions.
- `dreefile/codegen/` — the shared generation state (`codegen.State`,
  `codegen.Layout`) threaded through the sections, discovery, and generation.
  It imports `ir`. It is **not** merged into `gogen`: `gogen` is stateless Go
  emission, while `codegen.State` is the mutable build state (component
  definitions, imports, message uses, Lua features) that nearly every section
  needs.
- `dreefile/jsoutput/` — the shared client emission artifact and `<script>`
  tag generation (former `js/output`: `Artifact`, `GenClient`, `GenClientTo`).
  It imports `ir`. It is shared because the body/html section and every client
  language emit the same script form.

### Front-end stage

The tokenizing and parsing stack keeps explicit homes at the dreefile level,
parallel to the shared packages:

- `dreefile/tokens/` — token types.
- `dreefile/lexer/` — the scanner; imports `tokens`.
- `dreefile/parser/` — the section parser; imports `tokens`, `lexer`, `ir`.

The front-end stage produces the `ir` tree; sections consume it. `dreecode` is
the mini-template-language semantics that sections use for control flow,
expressions, filters, and slots — not the token/lexer/parser stack.

### Sections

`dreefile/sections/` owns the section processors:

```text
dreefile/sections/
├── header/       DREEFILE, LAYOUT, COMPONENT, GOIMPORT
├── head/
├── server/
├── style/
├── body/html/
├── body/md/      Markdown processor home (see below)
├── client/       explicit client orchestrator (former js/process)
│   ├── js/
│   ├── ts/
│   └── lua/
```

Language subfolders exist **only** where real language variants exist: `body`
and `client`. A section that has one language stays flat.

The `client/` orchestrator (former `js/process`) is the single legal entry point
into the client languages. It imports its children `client/js`, `client/ts`, and
`client/lua`; the children never import the orchestrator or each other. The
orchestrator is a shared emission entry point rather than a section processor,
so `body/html` may import it without a section-to-section exception.

The `body/md/` package is the **Markdown processor home** (former `html/md`).
It owns both the `lang="md"` body transform and the `<md>`-tag region scanning
that currently lives in `parser/parser_md_tag.go`; that scanning moves out of
the parser into `body/md/`. The parser therefore imports no Markdown package.

Message-expression parsing (`ParseMessageExpression`) is
mini-template-language semantics and belongs to `dreecode`, not to the parser
stack; the sections that need it use `dreecode`.

### Discovery and code generation

`dreefile/route/`, `dreefile/component/`, `dreefile/layout/`,
`dreefile/assets/`, `dreefile/format/`, `dreefile/config/`, `dreefile/check/`,
and `dreefile/i18n/` own discovery and per-concern generation.

### Shared Markdown home is `internal/md`

`core/markdown.go` imports the Markdown implementation and `ir`. After the move
it must not import the compiler, so the shared Markdown implementation lives at
**`internal/md`** — *not* `internal/render/markdown`, because `render` is the
target-neutral render contract and must not grow a Markdown dependency.

`internal/md` is shared implementation and therefore imports **no compiler
package** (`ir` included). It exposes an HTML-string entry point and the
`Mode` values (`ModeSafe`, `ModeTrusted`). Exact names are a PR2 code decision;
the runtime's exported `core.MarkdownToHTML`/`MarkdownToHTMLTrusted` keep their
current public signatures.

The compiler keeps its `TemplateNode` path in `dreefile/sections/body/md/`
(`ToNodes`, `TransformNodes`): that adapter imports `ir` and may call
`internal/md`, but the dependency points compiler -> shared, never shared ->
compiler. The runtime only ever needs rendered text — `core/markdown.go` rejects
every non-`NodeText` node — so `core` calls `internal/md` directly on the string
path and never touches the node adapter or `ir`.

This replaces the earlier open item; it is decided, not deferred.

### Structural check is required

`_tests/sh/check-core-deps.sh` only checks **external** dependencies (it greps
out every `github.com/dreego-stack/dreego` and `golang.org/x/` import). It does
not catch layering or section-to-section violations. A dedicated structural
check for the rules above is therefore **required**, not optional, and lands as
its own slice before the structural moves can be trusted.

### Dependency rule

The rule must match the real code, not an idealized tree. The verified edges
are: every section imports `codegen.State`; the body/html output stage imports
the client orchestrator and the client script emitter; the client orchestrator
imports its three language children; and the front-end stage produces the `ir`
tree for the sections.

```text
sections/*            -> codegen, dreecode, gogen, ir, jsoutput
sections/body/html    -> sections/client            (client orchestrator, explicit entry)
sections/client       -> sections/client/{js,ts,lua} (parent -> child only)
sections/client/{js,ts,lua} -> codegen, jsoutput, ir
dreecode              -> ir
gogen                 -> ir
codegen               -> ir
jsoutput              -> ir
lexer                 -> tokens
parser                -> tokens, lexer, ir
dreefile root (generate.go, plan.go, support.go, codegen_*.go)
                      -> sections, front-end, codegen, dreecode, gogen, ir, jsoutput
```

The dreefile root orchestrates: it may import the sections, the front-end, and
the shared packages, but nothing imports the root back.

No section imports another **sibling** section family. The client orchestrator
(`sections/client/`) is a shared emission entry point, not a section processor:
it is the one legal target for `body/html -> sections/client`, and the language
children never import the orchestrator, `body/html`, or each other. A section
that needs shared behavior otherwise uses `codegen` (generation state),
`dreecode` (language semantics), `gogen` (Go emission), `jsoutput` (client
script emission), or `ir` (types).

`dreecode` owns the mini-template-language semantics, including message
expression parsing; the parser stack is not a section dependency.

### Layer summary

| Layer | Packages | May import | Must not import |
|---|---|---|---|
| Shared implementation | `internal/{render,context,i18n,server,middleware,session,validate,gomod,templates,md}` | standard library, `golang.org/x/*`, each other per direction | `core`, `adapter/*`, `cmd/*`, `dreegotest`, `internal/dreefile/...` as Go packages |
| Compiler shared | `internal/dreefile/{ir,dreecode,gogen,codegen,jsoutput}` | standard library, `ir` per the rule | other compiler packages outside the rule, shared implementation, `core`, `adapter/*` |
| Compiler front-end | `internal/dreefile/{tokens,lexer,parser}` | `tokens`, `lexer`, `ir` per the rule | sections, `core`, `adapter/*` |
| Compiler sections | `internal/dreefile/sections/...` | `codegen`, `dreecode`, `gogen`, `ir`, `jsoutput`, `sections/client` (orchestrator only) | sibling sections, `core`, `adapter/*` |
| Facades and hosts | `core`, `adapter/ssr`, `adapter/wails`, `cmd/dreego`, `dreegotest` | shared implementation, public facades | — |

The compiler dependency rule is a review invariant and is enforced by the
required structural check described above; generated import strings
(`github.com/dreego-stack/dreego/core`,
`github.com/dreego-stack/dreego/adapter/ssr`) are output, not imports of the
compiler package.

## Consequences

- The public API is unchanged. Generated code and applications keep importing
  only `github.com/dreego-stack/dreego/core` and, for HTTP hosting,
  `github.com/dreego-stack/dreego/adapter/ssr`.
- `internal/dreefile/` is importable only from inside the repository (CLI,
  `dreegotest`), exactly like `internal/transpiler/` was.
- The split packages in `core/internal/` keep their package names and contents;
  only their location changes. No re-export or forwarding package is added
  inside `core/`.
- The 300-line rule stays enforced per file inside the split.
- `_tests/sh/check-core-deps.sh` keeps the **external** dependency boundary
  checked after the moves. It does **not** check the layering rule, so the
  required structural check is a separate, necessary artifact.
- Structural invariants (no sibling-section imports; `dreecode`, `gogen`,
  `codegen`, and `jsoutput` import only `ir`; the client orchestrator is the
  only parent of the language children; root `internal/` never imports `core`)
  are asserted by the required dedicated check.
- The `internal/transpiler` name disappears; documentation, plans, and todos
  that cite `internal/transpiler/...` anchors must be re-pointed.
- `ir/mdtohtml.go` (`TranslateMdtohtml`) is codegen **string rewriting**, not
  AST/type code, so it is miscategorized in `ir`. Moving it to `gogen` or the
  server section is a small, non-blocking follow-up and is not part of the
  structural move.

## Supersedes

- The **location parts** of
  [Core runtime split into internal subpackages](core-internal-subpackages.md).
  Its runtime dependency direction (`context -> session, middleware`;
  `middleware -> session`; `server -> session, middleware, validate`) is carried
  over to the new root `internal/*` locations; the `core/internal/*` paths are
  superseded.
- The `core/internal/*` package map in
  [`_plan/phase-render-foundation.md`](../../_plan/phase-render-foundation.md)
  (`core/internal/app|render|compiler|assets|validate`). The capability
  conclusions of that completed phase remain valid; only the location map is
  superseded.
- The package **name and location** in
  [Transpiler as internal subpackage](../../cmd/dreego/_docs/decisions/transpiler-subpackage.md);
  its reasoning (a shared root `internal/` for multiple in-repo consumers) still
  holds, the compiler is now `internal/dreefile/`.

## See also

- [Phase: internal layering and dreefile restructure](../../_plan/phase-restructure.md)
- [Target-neutral application and first-party adapters](target-neutral-application-and-first-party-targets.md)
- [Dreego Architecture](../dreego-architecture.md)
