# Agent Instructions for Dreego

## Aktuelle Phase: Planung / Konzeption

Dreego befindet sich in der Konzeptionsphase. Wir planen das Framework, schreiben noch keinen Code.
Architektur-Entscheidungen werden dokumentiert, Konzepte ausgearbeitet.

## Knowledge Base

Use `.agents/` as your Obsidian-style knowledge base.
Read `.agents/_index.md` first for orientation.
Write freely in `.agents/` — it's your brain, not source code.
Keep notes linked with `[[wiki-style links]]` where useful.

### Struktur

```
.agents/
├── _index.md              ← Start hier
├── KB/                    ← Referenz-Material (PDF-Konvertierungen, externe Quellen)
│   └── dreego-concept.md   ← Vollständiges Gemini-Chat-Konzept
├── decisions/             ← Architektur-Entscheidungen (ADR-Format)
│   ├── name-dreego.md
│   ├── technology-stack.md
│   ├── transpiler-vs-runtime.md
│   ├── typescript-v2.md
│   ├── sections-in-dreego.md
│   ├── ssr-first.md
│   ├── no-catch-tag.md
│   └── file-based-routing.md
├── concepts/              ← Ausgearbeitete Konzepte
│   ├── dreego-architecture.md
│   ├── dreego-sections.md
│   ├── template-logic.md
│   ├── addon-ecosystem.md
│   └── roadmap.md
└── guides/                ← Arbeits-Anleitungen
    ├── architecture.md
    └── coding-standards.md
```

## Projekt: dreego

- **Name:** dreego
- **Dateiendung:** `.dreego`
- **Package:** `dreego`
- **Ziel:** Ein SSR-First Go-Webframework auf Augenhöhe mit Phoenix, Next.js, SvelteKit
- **Ansatz:** Compile-Time Transpiler (.dreego → Go-Code) + Chi-Router + HTMX/Alpine.js
- **V1:** MVP mit Transpiler, File-based Routing, `{#if}`, `{#each}`
- **V2:** TypeScript, SSG, Wails, erweiterte Template-Logik

## Architektur-Garantien für V2 (MÜSSEN in V1 beachtet werden)

Diese Punkte werden in V1 noch nicht umgesetzt, aber die Architektur muss von Anfang an
darauf vorbereitet sein. Wenn einer dieser Punkte in V1 falsch designed wird, ist die
V2-Implementierung extrem aufwändig oder unmöglich.

### 1. Target-Agnostische Transpiler-Pipeline
V1 ist SSR-only, aber die Code-Generation muss ein abstraktes Target-Interface haben.
Nicht den HTTP-Handler direkt im Code-Generator hart verdrahten.

**V1:** `TargetSSR` (generiert `http.HandlerFunc`)
**V2:** `TargetSSG` (generiert statische HTML-Dateien), `TargetWails` (generiert Wails-kompatiblen Code)

→ Entscheidung: [[decisions/ssg-wails-v2]]

### 2. `<go>`-Block: Kein hartes `*http.Request`
Der `<go>`-Block bekommt in V1 den HTTP-Request. Aber für SSG gibt es keinen Request
(Build-Zeit), und für Wails gibt es System-APIs statt HTTP.

**Lösung:** Ein `dreego.Context` Interface, das:
- In SSR: `*http.Request` wrappt
- In SSG: Build-Time-Daten + Dateisystem-Zugriff bietet
- In Wails: System-APIs + Fenster-Context bietet

NICHT `r *http.Request` direkt in generierten Code schreiben. Immer über `dreego.Context` gehen.

→ Entscheidung: [[decisions/context-design]]

### 3. Transpiler-Pipeline mit Erweiterungspunkten
Parse → AST → CodeGen. Jede Phase muss austauschbar/erweiterbar sein.

- `<script>`-Verarbeitung: V1 = No-Op, V2 = esbuild-TS-Compiler
- `<style>`-Verarbeitung: V1 = Sammeln, V2 = PostCSS/Tailwind-Integration
- CodeGen: V1 = Go-HTTP-Handler, V2 = SSG/Wails

→ Entscheidung: [[decisions/typescript-v2]]

### 4. Plugin-Interface: First Release = Final Contract
Das `dreego.Plugin` Interface aus [[concepts/addon-ecosystem]] muss vor dem ersten
Release final sein. Spätere Änderungen sind Breaking Changes für alle Addons.

**Muss abdecken:**
- Middleware-Injection (SSR)
- Route-Registrierung (SSR + SSG)
- Asset-Bereitstellung (`//go:embed`)
- Transpiler-Hooks (Custom-Tags wie `<dreego:map />`)
- Context-Erweiterung (`c.User()`, `c.Session()`)

### 5. File-based Routing: Crawlable für SSG
Für SSG muss der Transpiler alle Routen kennen, ohne den Server zu starten.
Route-Registrierung muss deklarativ genug sein, dass ein Static Analyzer sie findet.

**NICHT:** `init()`-Funktionen mit `mux.HandleFunc()` (imperativ, nicht crawlable)
**SONDERN:** Generierte `dreego_router.go` die alle Routen zentral listet

→ Entscheidung: [[decisions/file-based-routing]]

### 6. Asset-System: Dual-Mode (Embedded + Disk)
`//go:embed` packt Assets ins Binary (SSR). Für SSG müssen Assets als Dateien auf
Disk geschrieben werden.

Beide Pfade müssen von Anfang an coexistieren können. Kein Code der annimmt, dass
Assets NUR embedded oder NUR auf Disk sind.

### 7. Template-Rendering ohne HTTP-Server
Die Render-Funktion einer `.dreego`-Datei muss ohne laufenden HTTP-Server aufrufbar sein
(für SSG und Tests).

```go
// SSR
page.RenderSSR(w, r)

// SSG (V2) — kein http.ResponseWriter, kein *http.Request
html := page.RenderStatic(props)
```

### 8. CLI-Interface: Reservierte Flags
Diese Flags in V1 bereits definieren (zeigen "Coming in V2"):

```
dreego build --static    # SSG: Statische HTML-Dateien generieren
dreego build --wails     # Wails: Desktop-App-Code generieren
dreego build --mobile    # Mobile: App-Code generieren (Zukunft)
```

## Coding Rules (für später, wenn Code geschrieben wird)

- Max 120 Zeilen pro Datei
- Eine logische Sache pro Datei
- Keine Kommentare im Code
- Package names kurz, sauber, ohne Hyphen
- Go 1.26+, Standard Library bevorzugen
- Build & Start via `make` oder `dreego` CLI, nicht direkt `go build`
- Generierte `*_dreego.go` Dateien werden nicht committed

Vor jedem Commit stelle sicher alle was im code geändert wurde oder seit dem gechatet wurde du nun mit den *.md in sync hälst  


Das framework sollte schon typesage sein somit erros sollten im generate oder build prozess auftreten und spätendens beim start des wegserver füher immer besser und wenn es die bei der laufzeit gibt dann sollten die so lokal wie moeglich sein