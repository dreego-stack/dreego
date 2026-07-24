# Agent Instructions for Dreego

## Aktuelle Phase: v0.0.1 — Erster Prototyp

Transpiler, File-based Routing, Layout, Middleware, CLI funktionieren.
Version v0.0.1 ist getaggt. Jetzt: Feature-Arbeit aus ROADMAP.md und thinking-list.md.

## Knowledge Base

Use `.agents/` as your Obsidian-style knowledge base.
Read `.agents/_index.md` first for orientation.
Write freely in `.agents/` — it's your brain, not source code.
Keep notes linked with `[[wiki-style links]]` where useful.

### Struktur

```
.agents/
├── _index.md              ← Start hier
├── thinking-list.md        ← Detaillierte Feature-Liste & offene Fragen
├── KB/                    ← Referenz-Material
├── decisions/             ← Architektur-Entscheidungen (ADR-Format)
├── concepts/              ← Ausgearbeitete Konzepte
└── guides/                ← Coding-Standards, Architektur-Guide
```

## Projekt: dreego

- **Name:** dreego
- **Dateiendung:** `.dreego`
- **Package:** `dreego`
- **Modul:** `codeberg.org/dreego/dreego`
- **Org:** codeberg.org/dreego
- **Mirror:** github.com/LukasLow/dreego
- **Ziel:** Ein SSR-First Go-Webframework auf Augenhohe mit Phoenix, Next.js, SvelteKit
- **Ansatz:** Compile-Time Transpiler (.dreego → Go-Code) + net/http + HTMX/Alpine.js
- **V1:** MVP mit Transpiler, File-based Routing, `{#if}`, `{#each}`
- **V2:** TypeScript, SSG, Wails, erweiterte Template-Logik

## Architektur-Garantien fur V2 (MUSSEN in V1 beachtet werden)

Siehe unveranderte Sektion 1-8 unten.

## Coding Rules

- Max 120 Zeilen pro Datei
- Eine logische Sache pro Datei
- Keine Kommentare im Code
- Package names kurz, sauber, ohne Hyphen
- Go 1.22+, Standard Library bevorzugen
- Build & Start via `dreego` CLI, nicht direkt `go build`
- Generierte `dree.go` Dateien werden nicht committed (`.gitignore`)

## Versioning & Tags

- Version: SemVer (`v0.0.1`)
- Tag: `git tag -a v0.0.1 -m "v0.0.1: first prototype"` nach Commit
- Neue Version wenn:
  - Neues Feature (MINOR bump in 0.x: `v0.0.1` → `v0.0.2`)
  - Bug-Fix (PATCH bump: `v0.0.1` → `v0.0.2`)
  - Doc-/Tool-Verbesserung, die das Repo besser macht

## Docs & Changelog

- `_docs/` — offentliche Dokumentation (immer aktuell halten)
- `README.md` — Projekt-Ubersicht (bei Major-Changes updaten)
- `CHANGELOG.md` — bei jedem Tag aktualisieren
- `TODO.md` — offene Docs/Features

Agent MUSS nach bedeutenden Changes `_docs/`, `README.md`, `CHANGELOG.md` aktualisieren.
Nicht bei jedem Commit, aber wenn ein Feature fertig ist oder ein Tag erstellt wird.

## Type Safety

Das Framework muss typensicher sein. Fehler sollen so fruh wie moglich auftreten:
1. Build-Zeit (`dreego generate`, `go build`)
2. Start-Zeit (`server.Listen()`)
3. Laufzeit (pro Request, so lokal wie moglich)

---

## Architektur-Garantien fur V2 (unverandert)

### 1. Target-Agnostische Transpiler-Pipeline

**V1:** `TargetSSR` (generiert `http.HandlerFunc`)
**V2:** `TargetSSG`, `TargetWails`

→ [[decisions/ssg-wails-v2]]

### 2. `<go>`-Block: Kein hartes `*http.Request`

Losung: `dreego.Context` Interface.

→ [[decisions/context-design]]

### 3. Transpiler-Pipeline mit Erweiterungspunkten

Parse → AST → CodeGen. Jede Phase austauschbar.

→ [[decisions/typescript-v2]]

### 4. Plugin-Interface: First Release = Final Contract

→ [[concepts/addon-ecosystem]]

### 5. File-based Routing: Crawlable fur SSG

Zentrale `dreego/gen/dree.go` listet alle Routen.

→ [[decisions/routing-and-components]]

### 6. Asset-System: Dual-Mode (Embedded + Disk)

### 7. Template-Rendering ohne HTTP-Server

### 8. CLI-Interface: Reservierte Flags

```
dreego build --static    # SSG (V2)
dreego build --wails     # Wails (V2)
dreego build --mobile    # Mobile (Zukunft)
```
