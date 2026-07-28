
---
type: Reference
title: Dreego Ecosystem Research 2025-2026
description: Research on State of JS 2025 and current web framework trends relevant to Dreego
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Dreego Ecosystem Research 2025-2026

> Research on State of JS 2025 and current web framework trends.

## React Ecosystem

### Meta-Frameworks
- **Next.js** dominates usage (58.6%), but satisfaction crashes: 68% → 55% in one year
  - #1 Pain Point: "excessive complexity"
  - Most commented, most hated
- **Astro** has highest satisfaction (94%), strongly growing, content sites
- **Remix** at 68% satisfaction — v3 is being completely rebuilt, **no longer on React**, AI-first design
- **Gatsby** dead (24% satisfaction)
- **TanStack Start** entering strongly (235 write-ins)

### UI Libraries
- **shadcn/ui** (on Radix) dominates — Copy-paste model instead of npm dependency
- Tailwind basically standard for styling

### State Management
- Moving away from Redux/MobX → React built-ins + TanStack Query
- **Signals** as a concept establishing itself framework-agnostically (Solid, Svelte 5 Runes, Angular, Preact)
- Is #3 of the most wanted JS features

### Forms
- React Hook Form, TanStack Form
- Trend towards Server Actions

### Auth
- Auth.js v5, Clerk, Lucia

### ORMs
- Drizzle and Kysely rising (typed SQL)
- Prisma losing (performance criticism)

### New Trends
- RSC/Server Actions controversial — Client/Server boundary confuses many
- `useOptimistic`, Transitions — Async as first-class
- **Isomorphic-first movement** (TanStack Start, SvelteKit, SolidStart): SSR without architecture switch

## Svelte Ecosystem

- **SvelteKit:** 88% Satisfaction (A-tier). File-based Routing, Form Actions
- No "use client" — everything SSR, client explicit
- Often called "what Next.js should be"
- **Svelte 5 Runes:** `$state`, `$derived`, `$effect` — Signals as language feature
- UI: shadcn-svelte, Melt UI (headless), Skeleton
- State: Built-in Stores usually suffice

## General Trends & Pain Points

### Top 5 Pain Points (Meta-Frameworks)
1. Excessive complexity
2. Next.js-specific
3. SSR complexity
4. Performance
5. Vendor lock-in (Vercel)

### Top 5 Pain Points (State Management)
1. Code Architecture
2. State Management
3. Managing Dependencies
4. Date Management
5. Performance

### Missing JS Features (Top Wishes)
1. Static Typing (5595 votes)
2. Standard Library (4952)
3. **Signals** (2823)
4. Pipe Operator
5. Pattern Matching

### Usage Patterns
- SPA: 89% | SSR: 59% | SSG: 46% | MPA: 40% | Streaming SSR: 10%

### AI Trend
- 29% of code AI-generated
- "React is the Last Framework" thesis
- Remix 3 designed AI-friendly through DSL reduction
- Cursor #2 Editor

### Framework Fatigue
- Developers use on average only 2.6 frameworks in their career
- Nobody jumps between frameworks anymore

## Relevance for Dreego

See [KB/ecosystem-analysis](../KB/ecosystem-analysis.md) for the detailed relevance analysis.
