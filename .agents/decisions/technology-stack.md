---
type: Decision
title: Technologie-Stack fur Dreego V1
description: Tech-Stack: Go, net/http, HTMX, Alpine.js
tags: [stack, v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---

# Technologie-Stack fur Dreego V1

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

Dreego soll auf bewährten Go-Bibliotheken aufbauen, statt alles neu zu erfinden. Jede Dependency wurde nach dem Kriterium "Selbst bauen vs. Dependency nutzen" evaluiert.

## Entscheidungen im Detail

### HTTP-Router: go-chi/chi
- **Warum nicht selbst bauen?** Zeitverschwendung. Chi ist ultra-schnell, 100% Go-Stdlib kompatibel, extrem flexibel und gut getestet.
- **Alternative:** `net/http` direkt (zu low-level für File-based Routing), gorilla/mux (nicht mehr maintained)

### Template Engine: Dreego Custom Transpiler
- **Warum selbst bauen?** Das ist das Herzstück von Dreego. Keine existierende Lösung bietet `.dreego` → Go-Code Transpilation.
- **Keine Alternative:** `a-h/templ` und `gomponents` sind cool, aber nicht das, was Dreego sein soll.

### Interaktivität: HTMX + Alpine.js + Datastar
- **Warum nicht selbst bauen?** Ein JS-Framework zu bauen ist ein Riesenprojekt und nicht das Ziel.
- **Datastar** nutzt SSE für Signale — ideal für Go's Concurrency.
- **HTMX** und **Alpine.js** sind extrem leicht und decken 95% aller Interaktivitätsfälle ab.

### CSS: Tailwind CLI
- **Warum nicht selbst bauen?** Einen CSS-Parser/Generator zu bauen ist unnötig komplex.
- Dreego ruft im Dev-Server die standalone Tailwind-Binary auf.

### Validierung: go-playground/validator
- **Warum nicht selbst bauen?** Zu viele Edge Cases. Der Validator ist ausgereift und deckt alle gängigen Fälle ab.

### Binary Packaging: embed (Go Stdlib)
- **Warum nicht selbst bauen?** Gibt es schon nativ in Go. Packt Tailwind, JS und Templates direkt ins Binary.

## Abgelehnte Dependencies

| Dependency      | Grund für Ablehnung                                  |
|-----------------|-----------------------------------------------------|
| Node.js / npm   | Zerstört Single-Binary-Versprechen                  |
| esbuild (V1)    | In V1 nur Vanilla JS — kein Bundler nötig           |
| TypeScript      | Komplexität, Scope Creep, erst in V2                |

## Konsequenzen

- `go.mod` in V1: chi, validator, goree/dreego (self)
- Tailwind wird als standalone Binary eingebunden (nicht als Go-Dependency)
- HTMX + Alpine.js werden als embedded Assets ausgeliefert
