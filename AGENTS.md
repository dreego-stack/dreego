# Agent Instructions for Dreego

## Aktuelle Phase: v0.0.1 → v0.0.2

Transpiler, File-based Routing, Layout, Middleware, CLI funktionieren.
Version v0.0.1 ist getaggt. Jetzt: Context-Refactoring aus TODO.md.

## Datei-Struktur

```
repo-root/
├── TODO.md                 ← NACHSTE Code-Anderungen (kurz, priorisiert)
├── ROADMAP.md              ← Release-Pipeline (high-level)
├── CHANGELOG.md            ← Was in welcher Version kam
├── README.md               ← Projekt-Ubersicht
├── LICENSE                 ← MPL-2.0
├── _docs/                  ← Offentliche Dokumentation
│
.agents/                    ← Knowledge Base (OKF-Format)
├── index.md                 ← Start hier (OKF TOC)
├── log.md                   ← Anderungshistorie
├── PlanTODO.md              ← Vollstandiger Plan aller Features
├── tips.md                  ← 50 Tipps + Beachtungsliste
├── KB/                      ← Referenz-Material
├── decisions/               ← Architektur-Entscheidungen (ADR)
├── concepts/                ← Ausgearbeitete Konzepte
└── guides/                  ← Coding-Standards, Skills, OKF-Konventionen
```

## Skills

- [Knowledge Base](.agents/guides/knowledge-base.md) — Wie die Knowledge Base gepflegt wird
- [Changelog](.agents/guides/changelog.md) — Wie CHANGELOG.md und Versionierung funktioniert
- [Open Knowledge Format](.agents/guides/open-knowledge-format.md) — OKF-Konventionen (YAML-Frontmatter, Typen, Links)

## Commit-Konvention

Siehe [Changelog-Guide](.agents/guides/changelog.md) fur den vollstandigen Workflow.

## Projekt: dreego

- **Name:** dreego | **Package:** `dreego` | **Modul:** `codeberg.org/dreego/dreego`
- **Org:** codeberg.org/dreego | **Mirror:** github.com/LukasLow/dreego
- **Ansatz:** Compile-Time Transpiler (.dreego → Go-Code) + net/http + HTMX/Alpine.js

## Coding Rules

- Max 120 Zeilen pro Datei, eine logische Sache pro Datei
- Keine Kommentare im Code
- Go 1.22+, Standard Library bevorzugen
- Core-Code liegt in `dreego-core/` (single package), Plugins in `dreego-plugin/`
- Build via `dreego` CLI, nicht direkt `go build`
- Generierte `dree.go` nicht committed

## Type Safety

1. Build-Zeit (`dreego generate`, `go build`)
2. Start-Zeit (`server.Listen()`)
3. Laufzeit (pro Request, so lokal wie moglich)

Kein `map[string]string`, kein `interface{}`-Cast, kein String-Key im Core.

---

## Architektur-Garantien fur V2

### 1. Target-Agnostische Transpiler-Pipeline
V1: `TargetSSR`, V2: `TargetSSG`, `TargetWails`. → [decisions/ssg-wails-v2](.agents/decisions/ssg-wails-v2.md)

### 2. `<go>`-Block: Kein hartes `*http.Request`
Losung: `dreego.Context` Interface. → [decisions/context-design](.agents/decisions/context-design.md)

### 3. Transpiler-Pipeline mit Erweiterungspunkten
→ [decisions/typescript-v2](.agents/decisions/typescript-v2.md)

### 4. Plugin-Interface: First Release = Final Contract
→ [concepts/addon-ecosystem](.agents/concepts/addon-ecosystem.md)

### 5. File-based Routing: Crawlable fur SSG
→ [decisions/routing-and-components](.agents/decisions/routing-and-components.md)

### 6. Asset-System: Dual-Mode (Embedded + Disk)

### 7. Template-Rendering ohne HTTP-Server

### 8. CLI-Interface: Reservierte Flags
`dreego build --static | --wails | --mobile`

vergesse die CLI.md nicht und dann kannst du beim demo server starten statt einfach die binary direkt starten und ausschalten wann auch immer und vergessen einfach einen timer setzen