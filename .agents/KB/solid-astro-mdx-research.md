
---
type: Reference
title: Solid.js, Astro & MDX — Deep-Dive Research
description: Research on Solid.js Signals, Astro Islands architecture, and MDX for Dreego feature adoption
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Solid.js, Astro & MDX — Deep-Dive Research

**Date:** 2026-07-28
**Sources:** docs.solidjs.com, docs.astro.build, mdxjs.com

---

## 1. Solid.js

### What Makes Solid Different from React/Svelte?

| Aspect | React | Svelte | Solid |
|--------|-------|--------|-------|
| Rendering | Virtual DOM → Diff → DOM | Compiler → direct DOM update | Fine-grained Reactivity → no VDOM, no compiler |
| Component Model | Components re-render completely | Compiler knows what changes | Components run ONCE, only reactive bindings update |
| State | `useState` (Hook, immutable) | `$state()` (Rune, mutable via Compiler) | `createSignal()` (getter/setter, mutable) |
| Philosophy | "UI = f(state)" | "Compiler makes it efficient" | "Every cell updates itself" |

**Core:** Solid components are **not render functions** — they are executed once. The return value is real DOM. Only signal-driven attributes/textNodes update on changes. React components, in contrast, run completely anew on each state change.

### How Signals Work in Solid (Fine-Grained Reactivity)

**Mechanism (Observer Pattern):**

1. `createSignal(initialValue)` → `[getter, setter]`
2. `createEffect(fn)` registers `fn` as subscriber
3. On `getter()` call: adds current subscriber to signal's subscriber list
4. `setter(newValue)`: notifies all subscribers (only on value change)
5. **No dirty checking, no VDOM diff** — signal knows exactly which DOM nodes to update

**Important Properties:**
- **Synchronous:** Tracking runs synchronously. `setTimeout` in effect → no tracking
- **Memos:** `createMemo()` caches derived values (like `$derived` in Svelte)
- **Resources:** `createResource(fetcher)` converts async to sync (like Svelte's `{#await}`)
- **Stores:** `createStore()` creates proxy-based signals for nested objects
- **No Stale Closure Problem:** Because components only run once and signals directly update the DOM

**This is relevant for Dreego because:**
- Dreego's Architecture (SSR + HTMX partials) is essentially the same as signals, just at the server level
- State changes → only the affected HTML fragment is re-rendered and swapped via HTMX
- No Virtual DOM needed — direct DOM update through HTML swap

### SolidStart Features

SolidStart is the meta-framework (like Next.js for React, SvelteKit for Svelte):

- **File-based Routing** — `routes/` directory, nested layouts, dynamic routes `[id].tsx`, catch-all `[...slug].tsx`, route groups `(groupName)/`, escaping nested routes `users(details)/`
- **Multiple Rendering Modes:** CSR, SSR (Sync, Async, Streaming), SSG (Pre-rendering via `crawlLinks: true`)
- **Server Functions** — `"use server"` Directive, database access without API endpoint. Integrated with `query()` (Solid Router)
- **API Routes** — GET/POST/etc. handlers in `routes/`
- **Middleware** — Request/Response interception
- **Sessions, Auth, WebSocket Endpoints**
- **Vinxi/Nitro** — agnostic bundler + server runtime (no vendor lock-in)
- **Deployment Presets:** Vercel, Netlify, Cloudflare, Node
- **No own router** — uses `@solidjs/router` (separates meta from router)

### Solid's Killer Features

1. **True Reactivity without Compiler or VDOM** — minimal runtime code
2. **Components run only once** — no re-rendering, only targeted DOM updates
3. **No Stale Closure Problem** — because components don't re-render
4. **No Dependency Arrays** — no `useEffect([dep])`, no `useMemo`, no `useCallback`
5. **SSR + Hydration without double rendering** — Solid only hydrates data bindings, not entire components
6. **Isomorphic Server Functions** — same function runs on server and client (with `"use server"`)
7. **Extremely small bundle size** — <7 KB gzipped for Solid Core
8. **JSX without React** — JSX compiles to real DOM operations (not to React.createElement)

---

## 2. Astro

### Islands Architecture

Astro popularized the **Islands Architecture** (originally by Etsy/Katie Sylor-Miller, 2019, then Jason Miller/Preact, 2020):

**Concept:**
- The page is a **sea of static HTML**
- Interactive components are **islands** of JavaScript in the static sea
- Each island is **isolated** loaded and hydrated
- Islands can load **in parallel** — a slow image carousel doesn't block the header

**Two Types of Islands:**
1. **Client Islands** — interactive UI components (React/Svelte/Vue/Solid)
2. **Server Islands** (`server:defer`) — dynamic server content that streams parallel to the main content (e.g. user avatar, product ratings)

### Partial Hydration — Client Directives

Astro's core mechanism: **JavaScript is removed by default.** Only explicitly marked components get JS. The directives:

| Directive | Behavior |
|-----------|----------|
| `client:load` | Load immediately (highest priority) |
| `client:idle` | When browser is idle (requestIdleCallback) |
| `client:visible` | When component enters viewport (IntersectionObserver) |
| `client:media="(max-width: 50em)"` | Only at matching media query |
| `client:only="react"` | Client-side only, no SSR |

### Content Collections

Astro's content management system:

- **`src/content.config.ts`** — Central configuration of all collections
- **Zod Schema** per collection → TypeScript types automatically generated, editor intellisense
- **Built-in Loaders:**
  - `glob()` — Directory of Markdown/MDX/Markdoc/JSON/YAML/TOML files
  - `file()` — Single file with array of entries
- **Custom Loader API** — CMS, database, API → everything integrable
- **Reference System** — Collections can reference each other (`reference('authors')`)
- **`getCollection()`, `getEntry()`** — Typed query API
- **`render()`** — Renders Markdown/MDX to HTML + `<Content />` Component
- **Filter API** — `getCollection('blog', ({data}) => data.draft !== true)`
- **Live Collections** — For real-time data (live updates without rebuild)
- **Route Generation** — Automatically generate pages from collection entries

### What Makes Astro Unique (94% Satisfaction)?

1. **Zero-JS Default** — It's impossible to accidentally ship JavaScript
2. **MPA instead of SPA** — Multi-Page architecture instead of Single-Page-App. Page transitions load new HTML documents (faster, simpler, SEO)
3. **Content-Driven Design** — Optimized for content sites (Blog, Marketing, Docs, E-Commerce), not for web apps
4. **Server-First** — Rendering happens on the server, not in the browser (like PHP/Laravel/Rails)
5. **"Opt in to Complexity"** — Start with HTML+CSS, add frameworks/JS as needed
6. **UI-Agnostic** — React, Preact, Svelte, Vue, Solid, HTMX, Web Components — all usable in parallel
7. **`.astro` Syntax** — Superset of HTML: every valid HTML is a valid Astro template
8. **View Transitions Router** — SPA-like animations between MPAs
9. **Dev Toolbar** — Integrated dev tools in the browser
10. **Adapter System** — Separation of framework and deployment target

---

## 3. MDX (and Astro's MDX Integration)

### How Does Markdown-to-Component Rendering Work?

MDX = Markdown + JSX:

```
# Hello, world!                           ← Markdown
<div className="note">                    ← JSX (Component!)
  > Some notable things in a block quote! ← Markdown IN JSX
</div>
```

**Processing Chain:**
1. **Parse** MDX text → MDAST (Markdown AST) + JSX nodes
2. **Transform** Remark plugins modify the MDAST
3. **Compile** → JavaScript (JSX becomes `createElement`-Calls, Markdown becomes HTML strings)
4. **Evaluate** → Execute in JS runtime (React/Preact/Vue)

**Important Features:**
- **Import/Export** — `import` and `export` statements in Markdown
- **Expressions** — `{Math.PI * 2}` in Markdown
- **Custom Components** — Replace HTML elements with own components: `export const components = {blockquote: CustomBlockquote}`
- **Frontmatter** — YAML/TOML at the beginning of the file
- **ESM Support** — Full JavaScript module syntax

### Astro's MDX Integration

Astro extends MDX:

- **Content Collections + MDX** — `.mdx` as collection entries with Zod schema and typing
- **Astro Components in MDX** — `.astro` Components directly imported and used in `.mdx`
- **Custom Components Mapping** — `<Content components={{h1: Heading}} />`
- **Own MDX Compiler** — `@astrojs/mdx` with `recmaPlugins`, `optimize` option
- **Frontmatter as First-Class** — `{frontmatter.title}` directly usable in MDX
- **Separate Processor** — MDX can use a different Markdown processor than `.md` files
- **Hybrid Mode** — Static pages + MDX content: essentially a headless CMS in Git

---

## For Dreego: Adopt / Don't Adopt

### Solid.js

**ADOPT:**
- [x] **Signals as conceptual model** — Dreego's SSR + HTMX is functionally equivalent: State change → only affected HTML fragment update. The mental model is the same.
- [x] **No Virtual DOM** — Solid proves that VDOM is not necessary. Dreego does the same at the server level (direct HTML rendering).
- [x] **Server Functions model** — `"use server"` = Dreego's `<go>` block. The `<go>` block runs exclusively server-side and interacts with DB/APIs — exactly the same concept.
- [x] **Isomorphic Code** — SolidStart's design principle: code runs on server and client. Dreego could use `dreego generate` to generate Go code that contains both server and client logic (via GopherJS/WASM in V2).
- [x] **No Dependency Arrays** — Dreego's template engine doesn't need `useEffect` equivalents. The `<go>` block model is simpler.
- [x] **Resource Pattern** — `createResource` = Dreego's `{#await}` tag. Make async data treatable as synchronous.

**DO NOT ADOPT:**
- [ ] **Client-side Signal System** — Dreego relies on HTMX + Alpine.js. A custom JS signal system would be redundancy and unnecessary JS.
- [ ] **JSX Syntax** — Dreego uses HTML template syntax (like Svelte), not JSX. JSX is too strongly tied to JavaScript/React.
- [ ] **Vinxi/Nitro as Build Tool** — Dreego has Go's `go build`. Doesn't need a JS build tool.
- [ ] **SolidStart as Architecture Model** — SolidStart is less established than SvelteKit/Next.js/Astro. For file-based routing, better to use SvelteKit/Astro as model.

### Astro

**ADOPT:**
- [x] **Islands Architecture** — This is THE decisive insight for Dreego. `.dreego` pages are by default static (Zero JS). Interactive "islands" are declared via HTMX/Alpine.js. No framework JS on the page except what was explicitly marked as interactive.
- [x] **"Zero JS by Default"** — This principle is perfect for Dreego. No framework JS code in the output, only what the developer explicitly wants.
- [x] **Partial Hydration Concept → "Partial HTML Swap"** — Astro's `client:visible` etc. can be transferred to Dreego: HTMX partials can be loaded lazy, on-visible, on-interaction.
- [x] **Content Collections** — Dreego urgently needs a similar concept for `.dreego` pages. A `dreego.collections.toml` with schema definition, automatic type generation, query API. That would be a killer feature for content sites.
- [x] **MPA Approach (Multi-Page-App)** — Dreego is by definition MPA. Astro proves that MPA is the right approach for content sites and SPAs are overkill. This validates Dreego's basic architecture.
- [x] **Adapter System** — Astro's separation of framework core and deployment adapter is architecturally clean. Dreego could do the same: Core Framework + Deployment Presets (Single Binary, Docker, Fly.io, VPS).
- [x] **".astro Syntax as HTML Superset" → Dreego's Template Syntax** — `.dreego` Template should lean on HTML (as it already does). Every valid HTML is a valid `.dreego` template — this massively lowers the entry barrier.
- [x] **View Transitions** — Astro's MPA+SPA animations. Dreego could achieve similar with HTMX's `hx-swap` + CSS View Transitions API.
- [x] **"Opt in to Complexity"** — Adopt design phrase for Dreego. Start simple, add complexity only when needed.
- [x] **Dev Toolbar** — A CLI dev server with built-in debug toolbar (shows current route, `<go>` block data, HTMX requests) would be a differentiating feature.
- [x] **Server Islands (`server:defer`)** — For Dreego this would mean: Parts of a page can stream asynchronously. `<div dreego:defer>` → is separately rendered and loaded via SSE/HTMX. Ideal for slow DB queries.
- [x] **UI-Agnosticism** — Astro supports React, Vue, Svelte in parallel. Dreego could allow multiple client frameworks in the `<script>` block (Alpine, Datastar, Petite-Vue). For V2 even WASM components.
- [x] **CSR Fallback (`client:only`)** — Some pages need full client interactivity. Dreego could have `regeo:client="true"` as a page-level directive to make a page fully client-side.

**DO NOT ADOPT:**
- [ ] **Server-First Rendering in JS** — Astro renders that on Node/Deno. Dreego does it in Go. The mechanism is fundamentally different (Compiler vs. Runtime-Renderer).
- [ ] **Vite-based Dev Server** — Dreego doesn't need a JS bundler. The dev server is Go-native.
- [ ] **npm/Node Ecosystem** — Dreego has Go modules. Astro's strength is integration with the JS ecosystem. Dreego can't replicate that (and doesn't want to).
- [ ] **Astro Components on Client and Server** — Dreego's `<go>` block and Template are strictly separated (Server vs. Client). Astro's `.astro` Components are a hybrid form.
- [ ] **Integration API in JS style** — Dreego needs a plugin system, but in Go idioms (Interfaces, not JS functions).

### MDX

**ADOPT:**
- [x] **Markdown + Components** — Dreego could give `.dreego` files a `{#md}` block concept in which Markdown is mixed with template components.
- [x] **Custom Component Mapping** — Like MDX's `export const components = {h1: MyHeading}` → Dreego could define templates that convert HTML tags in `.dreego` Components.
- [x] **Frontmatter as First-Class** — Every `.dreego` file could optionally have frontmatter (YAML/TOML) that is available in the `<go>` block as a `Meta` variable.
- [x] **Content Collections + Schema** — Like MDX in Astro Collections: `.dreego` pages in a collection with Zod-like schema (Go struct tags?), automatic type generation.
- [x] **Markdown in Template Blocks** — Template code (`{#if}`, `{#each}`) within Markdown content. "Literate Programming" for websites.
- [x] **ESM/Import Pattern → Dreego Addon System** — MDX's `import` Statement → Dreego's `{#use addon}` in template. Import components from addons.
- [x] **Remark/Rehype Plugin System → Dreego Markdown Pipeline** — A pipeline of Markdown transformers (Syntax Highlighting, Table of Contents, Link Rewriting). Implementable as Go pipeline.

**DO NOT ADOPT:**
- [ ] **JSX in Markdown** — Only works with JS runtime. Dreego would need to invent its own component syntax for Markdown.
- [ ] **Evaluate Step** — MDX evaluates JS code at runtime. Dreego compiles everything beforehand — no runtime evaluate. Safer, but less dynamic.
- [ ] **npm Plugin Ecosystem** — MDX's strength is the Remark/Rehype ecosystem. Dreego can't replicate this 1:1. A custom Markdown plugin system in Go would be needed.
- [ ] **MDX as File Format** — Dreego has `.dreego` as format. MDX would be a second, competing format. Better: Integrate MDX-like features into `.dreego`.

---

## Concrete Suggestions for Dreego

### 1. Islands Architecture for `.dreego`
```
<!-- Page is static by default -->
<div>
  <h1>{title}</h1>

  <!-- Interactive island via HTMX -->
  <div hx-get="/api/counter" hx-trigger="load">
    <!-- loads afterward -->
  </div>

  <!-- Client-side island via Alpine -->
  <div x-data="{ open: false }">
    <button @click="open = !open">Toggle</button>
  </div>
</div>
```

### 2. Content Collections via `dreego.collections.toml`
```toml
[collections.blog]
source = "content/blog"
pattern = "**/*.dreego"
[schema]
title = "string"
date = "time"
draft = "bool"
```

### 3. Server Islands via `dreego:defer`
```html
<!-- Parallel server rendering, streamed via SSE -->
<UserAvatar dreego:defer user="{userId}" />
<ProductReviews dreego:defer product="{productId}" />
```

### 4. Frontmatter in `.dreego`
```html
---
title: My Blog Post
date: 2026-07-28
draft: false
---

<go>
    // Meta.title, Meta.date, Meta.draft available
</go>

<h1>{Meta.title}</h1>
```

### 5. Markdown in Templates (MDX-like)
```html
<go>
    posts := query.Posts()
</go>

{#md}
# Welcome to My Blog

Here are the latest posts:

{#each posts as post}
- [{post.Title}](/post/{post.Slug})
{#end}
{/md}
```
