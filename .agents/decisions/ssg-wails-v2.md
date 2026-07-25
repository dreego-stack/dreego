
---
type: Decision
title: SSG & Wails Integration in V2
description: SSG und Wails Integration als gleichwertige Output-Modi neben SSR in V2
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# SSG & Wails Integration in V2

**Datum:** 23.07.2026
**Status:** Akzeptiert (geplant für V2)

## Kontext

Dreego V1 ist SSR-Only (Server-Side Rendering). Für viele Use Cases ist das ausreichend. Aber zwei wichtige Szenarien brauchen statischen Output:

1. **SSG (Static Site Generation):** `.dreego`-Seiten zu statischem HTML/JS/CSS kompilieren
2. **Wails:** `.dreego`-Komponenten in Go-Desktop-Apps verwenden
3. **Mobile (zukünftig):** Gleiche Komponenten via Wails Mobile / Gomobile

## Entscheidung

**V1: SSR-Only** (kein SSG). Der Fokus liegt auf dem Transpiler + Router.
**V2: SSG + Wails** als gleichwertige Output-Modi neben SSR.

## SSG Use Cases

| Use Case               | Beispiel                                  |
|------------------------|-------------------------------------------|
| Cloudflare Pages       | Statische HTML-Dateien auf Edge deployen  |
| GitHub Pages           | Projekt-Doku, Landing Pages               |
| S3/Cloudflare R2       | Pure Static Sites, 0 Server-Kosten        |
| Blog                   | Markdown → .dreego → statisches HTML       |
| Dokumentation          | docs.dreego.dev selbst mit SSG gebaut      |

## Wails Use Cases

| Use Case               | Beispiel                                  |
|------------------------|-------------------------------------------|
| Desktop-App            | Go-Backend + Dreego-Frontend = Native App  |
| Tray-App               | Menüleisten-App mit Dreego-UI              |
| Mobile (zukünftig)     | Gleiche Codebase für iOS/Android          |

## Vorteil: Code-Reuse

```
.dreego Komponenten
       │
       ├── SSR (Web)          — Chi-Server, HTML on-the-fly
       ├── SSG (Static)       — Statische HTML-Dateien
       ├── Wails (Desktop)    — Native Fenster, System-APIs
       └── Mobile (später)    — iOS/Android via Wails Mobile
```

Dieselbe `.dreego`-Datei rendert in vier verschiedenen Kontexten.
Kein JS-Framework kann das — weil sie alle eine JS-Runtime brauchen.

## Architektur-Vorbereitung in V1

Auch wenn SSG/Wails erst in V2 kommen, muss die Architektur in V1 vorbereitet sein:

1. **Transpiler-Pipeline mit Output-Modi:** Der Code-Generator hat ein `Target`-Interface
   - `TargetSSR` — Go-HTTP-Handler (V1)
   - `TargetSSG` — Statische HTML-Dateien (V2)
   - `TargetWails` — Wails-kompatible Go-Funktionen (V2)
2. **`dreego build --static`** — CLI-Flag bereits in V1 reservieren (tut nichts, zeigt "Coming in V2")
3. **Keine SSR-spezifischen Annahmen im Template:** `<go>`-Block kann in V1 nur Server-Code, aber Template ist target-agnostisch

## Inspiration aus der Rust-Welt

- **Dioxus:** Gleiche Komponenten für Web, Desktop (Blitz), Mobile
- **Leptos:** SSR + Hydration + Islands — zeigt, dass multi-target Architektur funktioniert
- **Yew:** War CSR-only, hat später SSR addiert — schwieriger als von Anfang an multi-target zu designen

## Konsequenzen

- V1: `dreego build` erzeugt SSR-Binary (HTTP-Server)
- V2: `dreego build --static` erzeugt `dist/` mit HTML-Dateien
- V2: `dreego build --wails` erzeugt Wails-kompatiblen Code
- `dreego.config.json` bekommt ein `target`-Feld: `"ssr" | "ssg" | "wails"`
- Gleiche `.dreego`-Komponenten in allen Targets verwendbar
