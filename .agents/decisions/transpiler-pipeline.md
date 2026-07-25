
---
type: Decision
title: Transpiler-Pipeline (Lexer → Parser → AST → CodeGen)
description: Compile-Time Transpiler Pipeline mit Lexer, Parser, AST und CodeGen für drei Targets
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Transpiler-Pipeline (Lexer → Parser → AST → CodeGen)

**Datum:** 23.07.2026
**Status:** Akzeptiert
**Review:** GLM-5.2 Expert Review (.tmp/output3.md)

## Kontext

Dreego ist ein Compile-Time Transpiler. `.dreego`-Dateien müssen in Go-Code umgewandelt werden — für 3 Targets (SSR, SSG, Wails).

## Entscheidung

**Phase 0: Single-Pass Scanner.** Kein Lexer/Parser/AST — nur State-Machine.

```
scan.go      — Hand-Scanner: erkennt Sections, {#tags}, {var}
codegen.go   — Pattern → Go-Source-String
target_ssr.go — Wrappt render(ctx) als http.HandlerFunc
```

~150 Zeilen, 0 Dependencies. Output: `func render(ctx dreego.Context) string` — target-agnostisch.

**Phase 1+: Formale Pipeline.** Wenn Features wachsen (<script> TS, SSG-Codegen) → formale Trennung Lexer→Parser→AST→CodeGen. Der Scanner wächst mit.

## 0.0.1 Architektur (Minimal)

```go
// scan.go — State-Machine
func scan(src []byte) (*File, error) {
    // 1. Section-Split: <go>, <style>, Rest=Template
    // 2. Template-Scan: {#if} → if-Block, {#each} → for-Block, {var} → Interpolation
    // 3. Stack-basiert: []string für verschachtelte Tags
}

// codegen.go — Go-Code erzeugen
func codegen(f *File) string {
    // <go>-Block → copy-paste
    // {var} → fmt.Sprintf("%s", var)
    // {#if cond} → if cond {
    // {#each xs as x} → for _, x := range xs {
    // <style> → extrahieren + scope-hash
}
```

## Fallstricke

1. `{#` vs `{` — `#` ist Discriminator. `{` lesen, peek nächstes Zeichen
2. Verschachtelte Tags — Stack (Depth-Counter), kein AST
3. Output: `render(ctx dreego.Context) string` — nicht direkt `http.HandlerFunc`. SSG/Wails brauchen kein HTTP

## 1. Pipeline-Interfaces

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

## 2. Sektions-Trennung

Ein **Pre-Lexer (Section-Splitter)** zerlegt die Datei an Block-Tags:

```
<head>...</head>    → HeadSection
<go>...</go>        → GoSection
<script>...</script> → ScriptSection
<style>...</style>  → StyleSection
Alles dazwischen     → TemplateSection
```

Jede Section bekommt spezialisierten Sub-Parser:
- **Head:** HTML-Parser → `HeadSection{Metas, Scripts, Links}`
- **Go:** Wird direkt an `go/parser` übergeben (Compiler validiert User-Go)
- **Template:** Dreego-Template-Lexer → `{#if}`, `{#each}`, `{var}`
- **Script:** V1 = Raw-String-Blob, V2 = TypeScript-Parser
- **Style:** V1 = Raw-String + Scope-Hash, V2 = PostCSS-Pipe

## 3. V2-Erweiterungspunkte (Strategy-Interfaces)

Gehören in den Core, nicht in Plugins:

```go
type ScriptProcessor interface {
    Process(src []byte, scope string) ([]byte, error)
    // V1: Identity (1:1 durchreichen)
    // V2: esbuild-TS-Compiler
}

type StyleProcessor interface {
    Process(src []byte, scope string) ([]byte, error)
    // V1: Scope-Hash anwenden
    // V2: PostCSS + Tailwind
}
```

## 4. AST-Format

**Eigenes AST** — nicht Go's `ast` package. Begründung: Go's `ast` ist für Source-Code-Transformation, nicht für gemischte HTML/Go-Sektionen. Der `<go>`-Block wird via `go/parser` geparst und als Identity-Pass eingebettet — so bleibt User-Go validierbar.

Generator-Output geht immer durch `go/format.Source` — nie handgeschriebene Einrückungen.

## 5. Error-Handling

Fail-loud, Compile-Time. Kein Best-Effort, keine Runtime-Panics.

- Lexer/Parser sammeln `Errors []Diag{Pos, Level, Msg}` bis Cap (20), dann Abbruch
- Critical-Fehler → `dreego generate` exit ≠ 0, kein Output
- `<go>`-Block: `go/parser`-Fehler werden via Source-Map auf `.dreego`-Zeilen gemappt
- Template-Syntaxfehler brechen sofort ab

## 6. CodeGen-Output pro Target

Alle drei Targets teilen denselben `EvalTemplate`-Kern — nur Context-Factory und Dispatcher unterscheiden sich.

**SSR:**
```go
func indexSSR(w http.ResponseWriter, r *http.Request) {
    ctx := dreego.NewSSRContext(r, w)
    user, err := loadUser(ctx)
    ctx.Render("index", map[string]any{"user": user})
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

## Konsequenzen

- Generierte Dateien: `pages/index_dreego.go` (nicht committed)
- `dreego generate` muss vor `go build` laufen
- Transpiler ist target-agnostisch — Target wird via CLI-Flag gewählt
