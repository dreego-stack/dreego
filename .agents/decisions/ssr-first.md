
---
type: Decision
title: SSR-First Architektur
description: Dreego rendert HTML auf dem Server mit HTMX und Alpine.js für Interaktivität
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# SSR-First Architektur

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

Dreego soll ein Webframework sein. Die Frage: Client-Side Rendering (CSR/SPA) oder Server-Side Rendering (SSR)?

## Entscheidung

**SSR-First.** Dreego rendert alles auf dem Server. Der Client bekommt fertiges HTML.

Interaktivität kommt nicht durch ein Client-Side Framework, sondern durch:
- **HTMX** für Server-Interaktionen (Partial Page Updates ohne Reload)
- **Alpine.js** für lokale UI-Interaktionen (Dropdowns, Tabs, Modals)
- **Datastar** (optional) für SSE-basierte Echtzeit-Updates

## Begründung

1. **Keine JS-Build-Hölle:** Kein node_modules, kein Webpack, kein Vite
2. **0 MB State-Synchronisierung:** State existiert nur auf dem Go-Server
3. **Perfektes SEO:** Alles ist statisches HTML beim First Load
4. **Schneller FCP (First Contentful Paint):** Auch auf schwachen Mobilgeräten
5. **Einfachere Architektur:** Kein API-Layer, keine JSON-Serialisierung, keine Client-State-Stores
6. **Direkter DB-Zugriff:** `<go>`-Block kann direkt auf die Datenbank zugreifen

## Vergleich

| Aspekt              | CSR (React/Svelte)        | SSR (Dreego)                  |
|---------------------|---------------------------|------------------------------|
| Initial Load        | JS muss laden + hydrieren | Sofort fertiges HTML         |
| SEO                 | Schwierig (SSR nötig)     | Perfekt                      |
| State Management    | Client + Server sync      | Nur Server                   |
| Bundle Size         | 100+ KB JS                | ~10 KB (HTMX + Alpine)       |
| Deployment          | Node.js + Static Files    | Single Go Binary             |

## Gegenargumente

- **"SSR fühlt sich langsam an bei Interaktionen"** → HTMX tauscht nur HTML-Fragmente aus, kein Full-Page-Reload
- **"Kein SPA-Feeling"** → Alpine.js + View Transitions API für flüssige Übergänge
- **"Weniger interaktiv"** → Datastar streamt DOM-Updates via SSE (wie Phoenix LiveView)

## Konsequenzen

- Kein Client-Side Router nötig
- Keine API-Schicht zwischen Template und Datenbank
- Tailwind + HTMX + Alpine.js sind feste Core-Dependencies
- V2: SSG (Static Site Generation) für rein statische Seiten
