
---
type: Reference
title: 50 Tipps + Beachtungsliste
description: 50 development tips and checklist covering DX, architecture, templating, plugins, and performance
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# 50 Tipps + Beachtungsliste

**Quelle:** Gemini-Chat, 25.07.2026

---

## 1. Developer Experience (DX) & Ergonomie (1–10)

| # | Tipp | Status |
|---|------|--------|
| 1 | Fehler mit `.dreego`-Zeilennummern, nie aus `dree.go` | geplant |
| 2 | Farbiges CLI-Output (grune Haken, rote Fehler) | geplant |
| 3 | `dreego init my-app` — funktionierendes Minimalprojekt | geplant |
| 4 | Flache Ordnerstruktur erlauben | ✅ `dreego/routes/get.dreego` reicht |
| 5 | Farblich hervorgehobene HTTP-Logs im Dev-Mode | geplant (aktuell: JSONL) |
| 6 | `static/`-Ordner mit `embed.FS` ausliefern | geplant |
| 7 | Kompaktes CLI-Output, nur Anderungen zeigen | geplant |
| 8 | Hot-Reload in <1s | geplant (V2) |
| 9 | `dreego build` → Single Binary | ✅ |
| 10 | `/health`-Endpoint per Config aktivierbar | geplant |

## 2. Architekturentwurf & Routing (11–20)

| # | Tipp | Status |
|---|------|--------|
| 11 | Explizite HTTP-Methoden-Suffixe | ✅ `get.dreego`, `post.dreego` |
| 12 | Typensichere Pfad-Parameter (int/string) | geplant |
| 13 | Wildcard-Routen `[...all].dreego` | geplant (Entscheidung steht) |
| 14 | `_middleware.go` pro Ordner | geplant (Entscheidung steht) |
| 15 | Kein globaler State, Race-Conditions vermeiden | refactoring (runtime globals) |
| 16 | `context.Context` durchreichen | ✅ `c.R.Context()` |
| 17 | Custom `404.dreego` und `500.dreego` | geplant |
| 18 | Form-Data Parser (`x-www-form-urlencoded`) | geplant |
| 19 | SSE / Streaming aus `<go>`-Block | geplant |
| 20 | `dreego.Redirect(ctx, "/login")` | geplant |

## 3. Templating & Component-System (21–30)

| # | Tipp | Status |
|---|------|--------|
| 21 | Klare Trennung: `<go>`, `<script>`, `<style>`, HTML | ✅ |
| 22 | Typensichere Props fur Komponenten | geplant (Entscheidung steht) |
| 23 | Automatisches XSS-Escaping | geplant |
| 24 | Raw-HTML-Escape-Hatch | geplant |
| 25 | CSS-Scope-Hash fur `<style>`-Block | ✅ `data-scope="hash"` |
| 26 | JS-Inlining aus `<script>` | ✅ |
| 27 | Slot-System fur Layouts | ✅ `{#slot}` |
| 28 | Conditional Rendering `{#if}`, `{#each}` | ✅ |
| 29 | Zero-JS-Mode (keine `<script>`-Sektion) | ✅ implizit (nichts = kein JS) |
| 30 | Head-Management (`<title>`, `<meta>`) | ✅ `<head>`-Block |

## 4. Plugins, Typensicherheit & Okosystem (31–40)

| # | Tipp | Status |
|---|------|--------|
| 31 | Compile-Time Plugin-Muster (Go-Pakete) | ✅ `dreego.Plugin` Interface |
| 32 | Interface-Driven Design | ✅ Capability-basiert |
| 33 | Typensichere Context-Values | geplant (Entscheidung steht) |
| 34 | Plugin-Konfiguration validieren | geplant |
| 35 | Offizielles Auth-Plugin | geplant (dreego-auth) |
| 36 | i18n-Plugin mit Typsicherheit | geplant (V2) |
| 37 | Zero CGO fur Cross-Compiling | ✅ (keine CGO-Abhangigkeiten) |
| 38 | CLI-Plugin-Hooks | geplant |
| 39 | Plugin-Lebenszyklus (`OnStart`, `OnShutdown`) | geplant (Entscheidung steht) |
| 40 | Plugin-Dokumentations-Vorlage | geplant |

## 5. Performance, Sicherheit & Barrierefreiheit (41–50)

| # | Tipp | Status |
|---|------|--------|
| 41 | Sichere HTTP-Header (CSP, X-Frame-Options) | geplant |
| 42 | Gzip/Brotli-Komprimierung | geplant |
| 43 | Zero-Allocation im Hot-Path | geplant (V2) |
| 44 | CSRF-Schutz | geplant (Entscheidung steht) |
| 45 | Graceful Shutdown | geplant |
| 46 | VS Code Syntax Highlighting Extension | geplant |
| 47 | ARIA-Warnungen im CLI | geplant (V3) |
| 48 | Go-Doc Kommentare fur Core | geplant |
| 49 | Dockerfile-Vorlage (Multi-Stage auf Scratch) | geplant |
| 50 | Showcase-Projekte (Todo, Blog) | geplant |

## Zusammenfassung

- ✅ Bereits umgesetzt: 11 von 50
- Geplant (Entscheidung steht): 22 von 50
- Geplant (V2/V3): 17 von 50
