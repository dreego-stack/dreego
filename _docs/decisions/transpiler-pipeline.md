
---
type: Decision
title: Transpiler Pipeline (Lexer → Parser → AST → CodeGen)
description: Compile-time transpiler pipeline with lexer, parser, IR, language processors, and code generation
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Transpiler Pipeline (Lexer → Parser → AST → CodeGen)

**Date:** 2026-07-28
**Status:** Accepted — pipeline core is current; see note below

> **Current direction:** The compile-time pipeline remains accepted. SSR is the
> production implementation. The historical universal
> `Target` interface and context examples below are not accepted APIs. Planned
> work uses a typed render foundation and adds Wails through capabilities
> proven by implementations. Static site generation is not planned. See
> [target-neutral-application-and-first-party-targets](target-neutral-application-and-first-party-targets.md).
**Review:** GLM-5.2 Expert Review (.tmp/output3.md)

## Current package boundaries

The implementation follows the pipeline with packages named after their
responsibility:

```text
internal/transpiler/
├── tokens/       token definitions
├── lexer/        source text to tokens
├── parser/       tokens to IR
├── ir/           intermediate representation and source metadata
├── codegen/      mutable generation state and layout metadata
├── html/         processors that produce HTML IR plus shared output generation
│   ├── html/      HTML input
│   ├── md/        Markdown input
│   ├── head/      head input
│   ├── css/       CSS input
│   └── output/    HTML IR to generated Go
└── js/           processors that produce JavaScript output
    └── js/        JavaScript input
```

The first directory in the processor matrix names the normalized output
language and the second names the source language. Future `js/ts` and `js/lua`
processors may therefore share JavaScript output generation without implying a
Lua-to-Go processor.

The root transpiler package owns project discovery, generation planning, and
the narrow entry points used by the CLI and `dreegotest`. It may provide small
compatibility helpers for its package-level tests, but it does not add facade
packages that mirror every lower-level function. The IR package contains the
data exchanged between pipeline stages; mutable code-generation state belongs
to `codegen`.

## Context

Dreego is a compile-time transpiler. `.dreego` files are normalized through
HTML, JavaScript, and Go output pipelines before Dreego emits ordinary Go and
browser assets for SSR or Wails hosts.

## Decision

**Phase 0: Single-Pass Scanner.** No lexer/parser/AST — only state machine.

```
scan.go      — Hand scanner: recognizes sections, {#tags}, {{ expressions }}
codegen.go   — Pattern → Go source string
target_ssr.go — Wraps render(ctx) as http.HandlerFunc
```

~150 lines, 0 dependencies. Output: `func render(ctx dreego.Context) string` — target-agnostic.

**Phase 1+: Formal Pipeline.** As input languages grow, use formal separation
from lexer to parser, IR, processor output, and code generation.

## 0.0.1 Architecture (Minimal)

```go
// scan.go — State machine
func scan(src []byte) (*File, error) {
    // 1. Section split: <server>, <style>, Rest=Template
    // 2. Template scan: {#if} → if block, {#each} → for block, {{ value }} → interpolation
    // 3. Stack-based: []string for nested tags
}

// codegen.go — Generate Go code
func codegen(f *File) string {
    // <server> block → copy-paste
    // {{ value }} → context-aware escaped output
    // {#if cond} → if cond {
    // {#each xs as x} → for _, x := range xs {
    // <style> → extract + scope hash
}
```

## Pitfalls

1. `{#` vs `{` — `#` is the discriminator. Read `{`, peek next character
2. Nested tags — Stack (depth counter), no AST
3. Output: render functions are not directly `http.HandlerFunc`; tests and Wails need non-HTTP rendering

## 1. Pipeline Interfaces

```go
type Source struct {
    Path    string
    Content []byte
}

type File struct {
    Path     string
    Head     *HeadSection
    Server   []ServerSection
    Body     *BodySection
    Bodies   []BodySection
    Client   *ClientSection
    Style    *StyleSection
}

type Lexer interface {
    Lex(src Source) ([]Token, error)
}

type Parser interface {
    Parse(tokens []Token, path string) (*File, error)
}

type AST = *File

type CodeGen interface {
    Generate(file *File, target Target) (*GeneratedFile, error)
}

type Target interface {
    Name() string
    ContextType() string
    Post(in *GeneratedFile) error
}
```

## 2. Section Separation

A **Pre-Lexer (Section Splitter)** splits the file at block tags:

```
<head>...</head>     → HeadSection
<server>...</server> → ServerSection
<body>...</body>     → BodySection
<style>...</style>   → StyleSection
<client>...</client> → ClientSection
```

Each section gets a specialized sub-parser:
- **Head:** HTML parser → `HeadSection{Metas, Scripts, Links}`
- **Server:** Go is passed to the Go compiler for validation
- **Body:** Dreego template lexer → `{#if}`, `{#each}`, expressions
- **Client:** Raw JavaScript string; optional languages require processors
- **Style:** V1 = Raw string + scope hash, V2 = PostCSS pipe

## 3. V2 Extension Points (Strategy Interfaces)

Belong in the core, not in plugins:

```go
type ScriptProcessor interface {
    Process(src []byte, scope string) ([]byte, error)
    // V1: Identity (pass through 1:1)
    // V2: esbuild TS compiler
}

type StyleProcessor interface {
    Process(src []byte, scope string) ([]byte, error)
    // V1: Apply scope hash
    // V2: PostCSS + Tailwind
}
```

## 4. AST Format

**Own AST** — not Go's `ast` package. Rationale: Go's `ast` is for source code transformation, not for mixed HTML/Go sections. The `<server>` block is parsed via `go/parser` and embedded as an identity pass — so user Go remains validatable.

Generator output always goes through `go/format.Source` — never hand-written indentation.

## 5. Error Handling

Fail-loud, compile-time. No best-effort, no runtime panics.

- Lexer/Parser collect `Errors []Diag{Pos, Level, Msg}` up to cap (20), then abort
- Critical errors → `dreego generate` exit ≠ 0, no output
- `<server>` block: `go/parser` errors are mapped to `.dreego` lines via source map
- Template syntax errors abort immediately

## 6. Host output

SSR and Wails share render contracts while retaining different host lifecycles.

**SSR:**
```go
func indexSSR(w http.ResponseWriter, r *http.Request) {
    ctx := dreego.NewSSRContext(r, w)
    user, err := loadUser(ctx)
    _ = user
    _ = err
}
```

**Wails:**
```go
func IndexWails(ctx *dreego.WailsContext) (string, error) {
    user, _ := loadUser(ctx)
    return EvalTemplate(ctx, user), nil
}
```

## Consequences

- Generated files: `pages/index_dreego.go` (not committed)
- `dreego generate` must run before `go build`
- The transpiler currently emits SSR handlers and target-neutral render functions.
- Wails is the planned second host after TypeScript-to-JavaScript and
  Lua-to-JavaScript complete the multi-language phase.
