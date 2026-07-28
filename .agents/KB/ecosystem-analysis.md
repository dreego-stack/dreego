
---
type: Reference
title: Relevance Analysis: React & Svelte Ecosystem 2025/2026
description: Relevance analysis of React & Svelte ecosystem patterns for Dreego's strategic positioning
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

**Date:** 2026-07-28
**Source:** State of JS 2025 Survey, Community Analyses

## State of the Ecosystem

### React
- **Next.js dominates (58.6%)** but satisfaction crashes: 68% → 55%. #1 Pain Point: "Excessive Complexity"
- **Astro** has highest satisfaction (94%), growing strongly
- **Remix v3** is being completely rebuilt, no longer on React, AI-first
- **shadcn/ui** dominates UI libraries (Copy-Paste model instead of npm dependency)
- **RSC/Server Actions** are controversial — Client/Server boundary confuses many

### Svelte
- **SvelteKit: 88% Satisfaction** — File-based Routing, Form Actions, no "use client"
- **Svelte 5 Runes:** `$state`, `$derived`, `$effect` — Signals as language feature
- **shadcn-svelte** + Melt UI for components
- Built-in Stores suffice for state management

### General Trends
- **Signals** as a concept is establishing itself framework-agnostically (#3 most wanted JS feature)
- **29% of code is AI-generated** — Framework must be AI-friendly
- **Framework Fatigue:** Developers use only 2.6 frameworks in their career
- **No Vendor Lock-in** is top pain point (Vercel criticism)

## What Dreego Should Adopt

| Pattern                     | Source        | Rationale                                                      |
|-----------------------------|---------------|----------------------------------------------------------------|
| **File-based Routing**      | SvelteKit     | SvelteKit-Style (not Next.js App Router). Clear, simple, proven |
| **SSR-First without "use client"** | SvelteKit | Explicit scope sections instead of fuzzy boundaries              |
| **Compiler Model**          | Svelte        | No runtime overhead. Svelte demonstrates how well this works    |
| **Form Actions**            | SvelteKit/Next| Map naturally to Go handlers. `g-submit` instead of `r.ParseForm` |
| **Shadcn Principle**        | shadcn/ui     | Copy-paste components instead of npm dependency. Adopt for dreego-ui |
| **Tailwind**                | Universal     | Basically standard. Class merging (`regeo.MergeClasses`) essential |
| **Single Binary**           | Go            | No Vercel lock-in. Deployable everywhere. Top pain point solved   |
| **Signals Concept**         | Svelte/Solid  | Map `{#let}` + reactive updates via HTMX/Datastar                |

## What Dreego Should Avoid

| Anti-Pattern               | Source        | Rationale                                                      |
|-----------------------------|---------------|----------------------------------------------------------------|
| **"use client" / RSC boundary** | Next.js     | Biggest pain point. Confuses developers. Dreego solves via sections more clearly |
| **Too many rendering modes** | Next.js      | SSG, ISR, PPR, Edge... Overkill. SSR + HTMX partials suffice   |
| **Framework as Platform**   | Next/Vercel   | No dependency on hosting provider                              |
| **Breaking Changes**        | Next.js       | Next.js loses satisfaction through constant breaks             |
| **Excessive Complexity**    | Next.js       | Small API surface, few concepts, Go idioms                     |
| **npm Dependency**          | JS Ecosystem  | No node_modules. Go modules are the way                        |

## Strategic Positioning

### Target Audience
- Go developers who want to build SSR pages without JS fatigue
- Teams coming from Next.js/SvelteKit who want Go in the backend
- Developers who want to deploy a single binary (Fly.io, VPS, Raspberry Pi)

### USP (Unique Selling Points)
1. **Single Binary Deployment** — no JS framework can do this
2. **100% Compile-Time Safety** — no runtime template errors
3. **No node_modules** — Go modules solve the dependency problem
4. **SSR without complexity** — no client/server boundary, clear sections
5. **No Vendor Lock-in** — deploy anywhere

### What Developers Really Want (and Dreego Delivers)
- **Static Typing** → Go natively
- **Standard Library** → Go excellent
- **SSR without complexity** → Section model
- **No Lock-in** → Single Binary, deployable everywhere
- **Signals Reactivity** → HTMX + Datastar
