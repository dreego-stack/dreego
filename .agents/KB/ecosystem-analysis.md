# Relevanz-Analyse: React & Svelte Ecosystem 2025/2026

**Datum:** 23.07.2026
**Quelle:** State of JS 2025 Survey, Community-Analysen

## State of the Ecosystem

### React
- **Next.js dominiert (58.6%)** aber Satisfaction crasht: 68% → 55%. #1 Pain Point: "Excessive Complexity"
- **Astro** hat höchste Satisfaction (94%), stark wachsend
- **Remix v3** wird komplett neu gebaut, nicht mehr auf React, AI-first
- **shadcn/ui** dominiert UI-Libraries (Copy-Paste-Modell statt npm dependency)
- **RSC/Server Actions** sind kontrovers — Client/Server-Grenze verwirrt viele

### Svelte
- **SvelteKit: 88% Satisfaction** — File-based Routing, Form Actions, kein "use client"
- **Svelte 5 Runes:** `$state`, `$derived`, `$effect` — Signals als Sprach-Feature
- **shadcn-svelte** + Melt UI für Komponenten
- Built-in Stores reichen für State Management

### Generelle Trends
- **Signals** als Konzept setzt sich framework-übergreifend durch (#3 meistgewünschtes JS-Feature)
- **29% des Codes ist AI-generiert** — Framework muss AI-freundlich sein
- **Framework-Müdigkeit:** Entwickler nutzen nur 2.6 Frameworks in der Karriere
- **Kein Vendor Lock-in** ist Top-Pain-Point (Vercel-Kritik)

## Was Dreego adoptieren sollte

| Muster                     | Quelle        | Begründung                                                     |
|----------------------------|---------------|----------------------------------------------------------------|
| **File-based Routing**     | SvelteKit     | SvelteKit-Style (nicht Next.js App Router). Klar, einfach, bewährt |
| **SSR-First ohne "use client"** | SvelteKit | Explizite Scope-Sektionen statt unscharfer Grenzen              |
| **Compiler-Modell**        | Svelte        | Kein Runtime-Overhead. Svelte macht vor, wie gut das funktioniert |
| **Form Actions**           | SvelteKit/Next| Mappen natürlich auf Go-Handler. `g-submit` statt `r.ParseForm` |
| **Shadcn-Prinzip**         | shadcn/ui     | Copy-Paste-Komponenten statt npm-Dependency. Für dreego-ui übernehmen |
| **Tailwind**               | Universal     | Ist quasi Standard. Class Merging (`regeo.MergeClasses`) essenziell |
| **Single Binary**          | Go            | Kein Vercel-Lock-in. Überall deploybar. Top-Pain-Point gelöst   |
| **Signals-Konzept**        | Svelte/Solid  | Die `{#let}` + reaktive Updats via HTMX/Datastar abbilden       |

## Was Dreego vermeiden sollte

| Anti-Pattern               | Quelle        | Begründung                                                     |
|-----------------------------|---------------|----------------------------------------------------------------|
| **"use client" / RSC-Grenze** | Next.js     | Größter Pain Point. Verwirrt Developer. Dreego löst via Sektionen klarer |
| **Zu viele Rendering-Modi** | Next.js       | SSG, ISR, PPR, Edge... Overkill. SSR + HTMX partials reichen   |
| **Framework als Platform**  | Next/Vercel   | Kein Abhängigkeitsverhältnis zu Hosting-Provider                |
| **Breaking Changes**        | Next.js       | Next.js verliert Satisfaction durch ständige Brüche             |
| **Excessive Complexity**    | Next.js       | Kleines API-Surface, wenige Konzepte, Go-Idiome                 |
| **npm-Abhängigkeit**        | JS-Ökosystem  | Kein node_modules. Go-Module sind der Weg                       |

## Strategische Positionierung

### Zielgruppe
- Go-Entwickler, die SSR-Seiten bauen wollen ohne JS-Fatigue
- Teams, die von Next.js/SvelteKit kommen, aber Go im Backend wollen
- Entwickler, die ein Binary deployen wollen (Fly.io, VPS, Raspberry Pi)

### USP (Unique Selling Points)
1. **Single Binary Deployment** — das kann kein JS-Framework
2. **100% Compile-Time Safety** — kein Runtime-Template-Error
3. **Kein node_modules** — Go-Module lösen das Dependency-Problem
4. **SSR ohne Komplexität** — keine Client/Server-Grenze, klare Sektionen
5. **Kein Vendor Lock-in** — deploye überall

### Was Entwickler wirklich wollen (und Dreego liefert)
- **Static Typing** → Go nativ
- **Standard Library** → Go exzellent
- **SSR ohne Komplexität** → Sektionen-Modell
- **Kein Lock-in** → Single Binary, überall deploybar
- **Signals-Reaktivität** → HTMX + Datastar
