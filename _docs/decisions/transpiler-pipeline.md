
---
type: Decision
title: Transpiler Pipeline (Lexer → Parser → AST → CodeGen)
description: Compile-time transpiler pipeline with lexer, parser, AST, and CodeGen for three targets
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Transpiler Pipeline (Lexer → Parser → AST → CodeGen)

**Date:** 2026-07-28
**Status:** Accepted — pipeline core is current; see note below

> **Current direction:** The compile-time pipeline remains accepted. SSR is the
> current implementation and v0.1 foundation. The historical universal
> `Target` interface and context examples below are not accepted APIs. Planned
> v0.x work extracts a typed render foundation and adds explicit first-party
> target packages using capabilities proven by implementations. See
> [target-neutral-application-and-first-party-targets](target-neutral-application-and-first-party-targets.md).
**Review:** GLM-5.2 Expert Review (.tmp/output3.md)

## Context

Dreego is a compile-time transpiler. `.dreego` files must be converted to Go code — for 3 targets (SSR, SSG, Wails).

## Decision

**Phase 0: Single-Pass Scanner.** No lexer/parser/AST — only state machine.

```
scan.go      — Hand scanner: recognizes sections, {#tags}, {var}
codegen.go   — Pattern → Go source string
target_ssr.go — Wraps render(ctx) as http.HandlerFunc
```

~150 lines, 0 dependencies. Output: `func render(ctx dreego.Context) string` — target-agnostic.

**Phase 1+: Formal Pipeline.** As features grow (<script> TS, SSG codegen) → formal separation Lexer→Parser→AST→CodeGen. The scanner grows with it.

## 0.0.1 Architecture (Minimal)

```go
// scan.go — State machine
func scan(src []byte) (*File, error) {
    // 1. Section split: <go>, <style>, Rest=Template
    // 2. Template scan: {#if} → if block, {#each} → for block, {var} → Interpolation
    // 3. Stack-based: []string for nested tags
}

// codegen.go — Generate Go code
func codegen(f *File) string {
    // <go> block → copy-paste
    // {var} → fmt.Sprintf("%s", var)
    // {#if cond} → if cond {
    // {#each xs as x} → for _, x := range xs {
    // <style> → extract + scope hash
}
```

## Pitfalls

1. `{#` vs `{` — `#` is the discriminator. Read `{`, peek next character
2. Nested tags — Stack (depth counter), no AST
3. Output: `render(ctx dreego.Context) string` — not directly `http.HandlerFunc`. SSG/Wails don't need HTTP

## 1. Pipeline Interfaces

```go
type Source struct {
    Path    string
    Content []byte
}

type File struct {
    Path     string
    Head     *HeadSection
    Go       *GoSection
    Template *TemplateSection
    Script   *ScriptSection
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
<go>...</go>         → GoSection
<script>...</script> → ScriptSection
<style>...</style>   → StyleSection
Everything in between → TemplateSection
```

Each section gets a specialized sub-parser:
- **Head:** HTML parser → `HeadSection{Metas, Scripts, Links}`
- **Go:** Passed directly to `go/parser` (compiler validates user Go)
- **Template:** Dreego template lexer → `{#if}`, `{#each}`, `{var}`
- **Script:** V1 = Raw string blob, V2 = TypeScript parser
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

**Own AST** — not Go's `ast` package. Rationale: Go's `ast` is for source code transformation, not for mixed HTML/Go sections. The `<go>` block is parsed via `go/parser` and embedded as an identity pass — so user Go remains validatable.

Generator output always goes through `go/format.Source` — never hand-written indentation.

## 5. Error Handling

Fail-loud, compile-time. No best-effort, no runtime panics.

- Lexer/Parser collect `Errors []Diag{Pos, Level, Msg}` up to cap (20), then abort
- Critical errors → `dreego generate` exit ≠ 0, no output
- `<go>` block: `go/parser` errors are mapped to `.dreego` lines via source map
- Template syntax errors abort immediately

## 6. CodeGen Output per Target

All three targets share the same `EvalTemplate` core — only context factory and dispatcher differ.

**SSR:**
```go
func indexSSR(w http.ResponseWriter, r *http.Request) {
    ctx := dreego.NewSSRContext(r, w)
    user, err := loadUser(ctx)
    _ = user
    _ = err
}
```

**SSG:**
```go
func IndexStatic(ctx *dreego.SSGContext, w io.Writer) error {
    posts, _ := loadPosts(ctx)
    fmt.Fprint(w, EvalTemplate(ctx, posts))
    return nil
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
- Transpiler currently emits SSR handlers; planned v0.x work first extracts the
  typed render foundation and then adds explicit SSG and Wails target packages
