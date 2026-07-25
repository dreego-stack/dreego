
---
type: Reference
title: Dreego Ecosystem Research 2025-2026
description: Research on State of JS 2025 and current web framework trends relevant to Dreego
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Dreego Ecosystem Research 2025-2026

> Research zu State of JS 2025 und aktuellen Web-Framework-Trends.

## React Ecosystem

### Meta-Frameworks
- **Next.js** dominiert Nutzung (58.6%), aber Satisfaction crasht: 68% → 55% in einem Jahr
  - #1 Pain Point: "excessive complexity"
  - Meistkommentiert, meistgehasst
- **Astro** hat höchste Satisfaction (94%), stark wachsend, Content-Sites
- **Remix** bei 68% Zufriedenheit — v3 wird komplett neu gebaut, **nicht mehr auf React**, AI-first Design
- **Gatsby** tot (24% Satisfaction)
- **TanStack Start** steigt stark ein (235 Write-ins)

### UI Libraries
- **shadcn/ui** (auf Radix) dominiert — Copy-paste-Modell statt npm dependency
- Tailwind quasi Standard für Styling

### State Management
- Weg von Redux/MobX → React built-ins + TanStack Query
- **Signals** als Konzept setzt sich framework-übergreifend durch (Solid, Svelte 5 Runes, Angular, Preact)
- Ist #3 der meistgewünschten JS-Features

### Forms
- React Hook Form, TanStack Form
- Trend zu Server Actions

### Auth
- Auth.js v5, Clerk, Lucia

### ORMs
- Drizzle und Kysely steigen (typed SQL)
- Prisma verliert (Performance-Kritik)

### Neue Trends
- RSC/Server Actions kontrovers — Client/Server-Grenze verwirrt viele
- `useOptimistic`, Transitions — Async als First-Class
- **Isomorphic-first Bewegung** (TanStack Start, SvelteKit, SolidStart): SSR ohne Architektur-Wechsel

## Svelte Ecosystem

- **SvelteKit:** 88% Satisfaction (A-tier). File-based Routing, Form Actions
- Kein "use client" — alles SSR, client explizit
- Wird oft als "was Next.js sein sollte" bezeichnet
- **Svelte 5 Runes:** `$state`, `$derived`, `$effect` — Signals als Sprach-Feature
- UI: shadcn-svelte, Melt UI (headless), Skeleton
- State: Built-in Stores reichen meist

## Generelle Trends & Pain Points

### Top 5 Pain Points (Meta-Frameworks)
1. Excessive complexity
2. Next.js-spezifisch
3. SSR complexity
4. Performance
5. Vendor lock-in (Vercel)

### Top 5 Pain Points (State Management)
1. Code Architecture
2. State Management
3. Managing Dependencies
4. Date Management
5. Performance

### Fehlende JS-Features (Top-Wünsche)
1. Static Typing (5595 votes)
2. Standard Library (4952)
3. **Signals** (2823)
4. Pipe Operator
5. Pattern Matching

### Nutzungsmuster
- SPA: 89% | SSR: 59% | SSG: 46% | MPA: 40% | Streaming SSR: 10%

### KI-Trend
- 29% des Codes AI-generiert
- "React is the Last Framework"-These
- Remix 3 designed AI-friendly durch DSL-Reduktion
- Cursor #2 Editor

### Framework-Müdigkeit
- Entwickler nutzen durchschnittlich nur 2.6 Frameworks in der Karriere
- Niemand springt mehr zwischen Frameworks

## Relevanz für Dreego

Siehe [KB/ecosystem-analysis](../KB/ecosystem-analysis.md) für die detaillierte Relevanz-Analyse.
