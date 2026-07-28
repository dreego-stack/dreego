
---
type: Reference
title: Rust Frontend Frameworks — Deep Dive & Dreego Transfer
description: Deep dive into Rust frontend frameworks (Leptos, Dioxus, Yew) with transferable patterns for Dreego
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Rust Frontend Frameworks — Deep Dive & Dreego Transfer

**Date:** 2026-07-28
**Purpose:** Systematic comparison of Rust frontend frameworks focusing on transferable patterns for Dreego.

---

## Why Rust Frameworks Are Relevant for Dreego

Rust and Go share decisive properties:
- **Compiled, statically typed** → Compile-time guarantees possible
- **No VM** → Direct code generation, Single Binary
- **Std-Library-centric** → Few external dependencies
- **No JS ecosystem** → Standalone solutions for State, Routing, Templating
- **SSR as First-Class** → All major Rust frameworks support SSR

Dreego-specific context: Dreego is SSR-First, needs no client-side WASM framework. We learn from Rust how compile-time safety, state management, and typed templating look, but at its core Dreego remains a transpiler (.dreego → Go) with SSR output.

---

## 1. Leptos (v0.8, 2026)

**Homepage:** https://leptos.dev | **Repo:** leptos-rs/leptos
**Approach:** Full-Stack, isomorphic, fine-grained reactivity
**Rendering Modes:** CSR, SSR, Hydrate, Islands
**Inspired by:** Solid.js (Signals), Svelte (Compile-Time)

### 1.1 Reactivity: Fine-Grained Signals

Leptos' core innovation: **Fine-grained reactivity** without Virtual DOM. The reactive system is a directed graph of three node types:

```
Signals (src)  →  Memos (intermediate)  →  Effects (sink)
```

**Signals (sources):**
```rust
let (count, set_count) = signal(0);       // ReadSignal + WriteSignal
let count = RwSignal::new(0);             // Combined read/write
```

**Memos (derived values, cached):**
```rust
let doubled = Memo::new(move |_| count.get() * 2);
```

**Effects (side effects, e.g. DOM updates):**
```rust
Effect::new(move |_| {
    log!("count changed: {}", count.get());
});
```

**Derived Signals (not cached, lightweight):**
```rust
let doubled = move || count.get() * 2;  // re-evaluated at every read
```

**Critical Design Detail:**
- Signals always mark themselves as "Dirty" (no `PartialEq` check)
- Memos check `PartialEq` before notifying subscribers
- Effects are executed at most once per update cycle (Diamond problem solved via "push-pull" algorithm)
- The entire reactive graph is thread-safe (`Send + Sync`)

### 1.2 Server Functions — The Killer App

**The Problem:** Client and server share code, but API contract must be written three times (Client-Fetch, Server-Handler, Type Definition).

**Leptos' Solution: `#[server]` Macro**

```rust
#[server]
pub async fn add_todo(title: String) -> Result<(), ServerFnError> {
    let mut conn = db().await?;
    sqlx::query("INSERT INTO todos (title) VALUES ($1)")
        .bind(title)
        .execute(&mut conn).await?;
    Ok(())
}

// Call on client — exactly the same function!
#[component]
fn TodoForm() -> impl IntoView {
    view! {
        <button on:click=move |_| {
            spawn_local(async { add_todo("Buy milk".into()).await; });
        }>"Add"</button>
    }
}
```

**What happens under the hood:**
1. `#[server]` generates *two versions* of the function via conditional compilation:
   - **Server-Build (ssr):** Real code, HTTP endpoint via POST (URL-encoded args, JSON return)
   - **Client-Build (hydrate/csr):** Stub that makes a `fetch()` call with serialized args
2. Arguments are URL-encoded per `serde_qs` (compatible with `<form>`)
3. Return values per `serde_json`
4. Custom encodings possible (JSON, Multipart, etc.)

**Why is this brilliant?**
- **Co-location:** DB query and click handler in the same file
- **Type-Safety across the boundary:** Argument types checked on both sides
- **Progressive Enhancement:** `<ActionForm>` also works without WASM/JS — Browser `<form>` suffices
- **No API Contract:** The Rust function *is* the contract

### 1.3 SSR + Hydration

Leptos compiles the app for **two targets**:
1. **Server (ssr):** Renders HTML, injects WASM loader
2. **Browser (hydrate):** Takes existing HTML, attaches event listeners

**Page-Load-Lifecycle:**
1. Browser `GET` → Server matches route → renders `<App/>` once → returns HTML
2. Browser loads WASM in parallel → WASM "walks" the HTML ("hydration"), attaches reactivity
3. After: All further navigations are **SPA navigations** (no page reload)

**Important:** The entire component tree must run on server AND client. That's the price of isomorphism.

### 1.4 Islands Architecture (Leptos 0.5+)

**Concept:** Instead of hydrating the whole app, only activate small "islands" of interactivity.

```rust
#[island]  // ← instead of #[component]
fn Counter() -> impl IntoView {
    let (count, set_count) = signal(0);
    view! { <button on:click=move |_| set_count.update(|n| *n += 1)>"Count: {count}"</button> }
}
```

**Effects:**
- **WASM size drops drastically:** 274 KB → 24 KB (Hello World), 400 KB → 200 KB (Tabs demo)
- **Server code in components:** `#[component]` runs ONLY on server (no WASM)
- **Children from server code:** Server-side content can be passed as children to islands ("Donut architecture")
- **Context between islands:** `provide_context`/`expect_context` works across islands

### 1.5 Routing

- **Declaration-based:** `<Routes>` + `<Route path=... view=...>`
- **Nested Routes:** `<Route>` can have `<Route>` children → Layouts
- **Params:** `/users/:id` → `let params = use_params();`
- **Queries:** `?search=...` → `let query = use_query();`
- **`<A/>` and `<Form/>`:** Router-aware links/forms, client-side navigation
- **Fallback:** `<Routes fallback=|| view! { <NotFound/> }>`

### 1.6 Component Model

```rust
#[component]
fn Button(
    #[prop(into)] class: String,     // Into<T> conversion
    #[prop(optional)] disabled: bool, // Optional
    children: Children,               // Slots/Children
) -> impl IntoView {
    view! {
        <button class=class disabled=disabled>
            {children()}
        </button>
    }
}
```

### 1.7 What Makes Leptos Unique

1. **Fine-Grained Reactivity:** Only exactly the DOM nodes reading a signal are updated — no VDOM Diffing
2. **`#[server]` Macro:** The most elegant client-server integration of all Rust frameworks
3. **Islands Mode:** Reduces WASM size by 50%, enables server-side `std::fs` in components
4. **Streaming SSR:** `<Suspense/>` + out-of-order streaming
5. **`view!` Macro:** RSX (React-like JSX in Rust) with compile-time checks against HTML spec
6. **Reactive graph is `Send + Sync`:** Multi-thread SSR on Tokio

---

## 2. Dioxus (v0.6/0.7, 2026)

**Homepage:** https://dioxuslabs.com | **Repo:** DioxusLabs/dioxus
**Approach:** Cross-Platform (Web, Desktop, Mobile) with React-like model
**Inspired by:** React, Flutter
**24.5k+ GitHub Stars**

### 2.1 Architecture: Virtual-DOM-based (unlike Leptos)

Unlike Leptos' fine-grained reactivity, Dioxus uses a Virtual DOM:
- Components re-render completely on state change
- VDOM is diffed, only changes land in the real DOM
- Similar to React, but with Rust performance

### 2.2 Multi-Platform Rendering

```rust
// Same code, three targets:
#[cfg(target_platform = "web")]     // → WASM in browser
#[cfg(target_platform = "desktop")] // → WGPU-based renderer (Blitz)
#[cfg(target_platform = "mobile")]  // → iOS/Android via WGPU
```

**Desktop/Mobile:** Dioxus has its own renderer called **Blitz** (based on WGPU), which natively renders HTML/CSS — without WebView, without Electron. That's the "Flutter, but better" promise.

### 2.3 State Management: Signals (since v0.5)

```rust
fn App() -> Element {
    let mut count = use_signal(|| 0);         // Signal Hook
    let half = use_memo(move || count() / 2); // Derived, cached
    let data = use_resource(move || async {    // Async Data
        fetch_data(count()).await
    });

    use_effect(move || {
        log!("count changed to {count}");     // Side effect
    });

    rsx! {
        button { onclick: move |_| count += 1, "Count: {count}" }
        div { "Half: {half}" }
    }
}
```

**State Passing:**
1. **Props:** `ReadOnlySignal<T>` — Standard
2. **Context:** `use_context_provider` / `use_context` — for sub-trees
3. **Globals:** `GlobalSignal<T>` as `static` — app-wide

### 2.4 Fullstack Mode

```rust
// Server Function (identical to Leptos' concept)
#[server]
async fn save_user(name: String) -> Result<(), ServerFnError> {
    db::insert_user(&name).await?;
    Ok(())
}

// In client:
rsx! { button { onclick: move |_| { spawn(async { save_user("Alice".into()).await; }); } } }
```

**Hydration:** Server renders HTML → serializes server future data as JSON → Client hydrates.

**Important difference from Leptos:**
- Dioxus requires `use_server_future` for non-deterministic data
- Hydration mismatches are a real problem (random values, timeouts)
- Dioxus' VDOM approach makes hydration more complex than Leptos' fine-grained approach

### 2.5 CLI & Dev Experience

```
dx new my-app        # Scaffolding
dx serve             # Dev server + Hot Reload
dx build             # Production build
dx translate         # HTML → RSX conversion
```

**Hot Reload:** Changes to RSX are updated in the browser without state loss.

### 2.6 Static Site Generation (SSG)

Dioxus supports SSG via fullstack mode: The server renders all routes to static HTML, which is then deployed.

### 2.7 What Makes Dioxus Unique

1. **True Cross-Platform:** Web, Desktop, Mobile — one codebase
2. **Blitz Renderer:** HTML/CSS natively rendered on GPU, no WebView
3. **React-Similarity:** Developers from React/Flutter feel immediately at home
4. **`dx` CLI:** First-class developer experience
5. **Ecosystem:** Largest Rust UI ecosystem (24.5k Stars)

### 2.8 Dioxus LiveView Mode

Dioxus' SSR can be operated as "LiveView": The server keeps the VirtualDom in memory, renders HTML diffs on changes, and sends them via WebSocket to the client. This is similar to Phoenix LiveView. The client doesn't need to load WASM — all events are sent to the server via WebSocket.

---

## 3. Yew (v0.23, 2026)

**Homepage:** https://yew.rs | **Repo:** yewstack/yew
**Approach:** Elm-inspired component framework (CSR-first, SSR optional)
**Oldest Rust Frontend Framework** (since 2017)

### 3.1 Architecture: Elm/React Model

Yew is the pioneer — the first Rust frontend framework. It follows the Elm/React model:

- **Virtual DOM** (like React)
- **Components** are Rust functions with `#[function_component]`
- **Message/Update architecture** (from Elm) via `use_reducer`
- **Props** are pure Rust structs
- **SSR** is additive (not first-class like Leptos/Dioxus)

### 3.2 Component Model

```rust
#[function_component]
fn HelloWorld() -> Html {
    html! { <p>{"Hello world"}</p> }
}
```

**Two Component Types:**
1. **Function Components** (`#[function_component]`) — recommended
2. **Struct Components** (advanced, more control)

**Lifecycle:** `use_state`, `use_effect`, `use_reducer`, `use_memo`, `use_callback`

### 3.3 State Management

Yew's state hooks are React-like:

| Hook             | Behavior                                | Scope    |
|------------------|-----------------------------------------|----------|
| `use_state`       | State, re-renders on `set`            | Component|
| `use_state_eq`    | Like `use_state`, but with `PartialEq` | Component|
| `use_reducer`     | Elm architecture: Action → State       | Component|
| `use_reducer_eq`  | Reducer + `PartialEq`                  | Component|
| `use_memo`        | Memoized Computation                   | Component|
| `use_callback`    | Memoized Callback (Deps-Array)         | Component|
| `use_mut_ref`     | Mutable Ref (no re-render)            | Component|

**Context:** `use_context` for tree-wide state passing.
**Agents:** Yew-specific concept for worker threads / background tasks.

### 3.4 SSR

Yew's SSR is functional, but not first-class:
```rust
let renderer = yew::ServerRenderer::<App>::new();
let html = renderer.render().await;
```

No hydration concept, no server functions, no file-based routing.

### 3.5 Yew's Significance Today

Historically extremely important (Go-To for Rust-WASM 2017–2021). But overtaken by Leptos and Dioxus in most aspects. Lessons:

- **Message/Reducer architecture** is good for complex state machines, but too boilerplate-heavy for UI
- **VDOM** is conceptually simple, but fine-grained reactivity (Leptos) is more performant
- **Community-driven development** — Yew was built by the community, not by a startup

---

## 4. Other Rust Frameworks

### 4.1 Sycamore (v0.9, 2024)

**Approach:** Fine-Grained Reactivity (like Leptos, but leaner)

```rust
#[component]
fn Counter(initial: i32) -> View {
    let mut value = create_signal(initial);
    view! {
        button(on:click=move |_| value += 1) {
            "Count: " (value)
        }
    }
}
```

**Features:**
- Fine-grained reactivity (no VDOM)
- SSR + Hydration
- Async/Suspense
- Built-in Routing (Client + Server)
- 3.3k Stars, actively maintained

**Relevance:** Sycamore shares many ideas with Leptos, but is smaller and less full-stack oriented. Proof that fine-grained reactivity works well in Rust.

### 4.2 Perseus (deprecated)

Perseus was a full-stack framework inspired by Next.js. Features:
- SSG (Static Site Generation) as first-class feature
- i18n built-in
- State plugins (global state via context)

Deprecated/merged. The SSG ideas have been absorbed into Dioxus and Leptos.

### 4.3 MoonZoon (deprecated)

Minimalist, "no-JS, no-CSS, no-HTML" — completely in Rust. Was an experimental framework that mapped Rust APIs directly to browser elements.

---

## 5. Cross-Cutting Patterns

### 5.1 What Rust Frameworks Can Do That JS Frameworks Cannot

| Capability                   | Enabled By                     | Relevance for Dreego |
|------------------------------|--------------------------------|----------------------|
| **Single Binary Deployment** | Rust compiles to WASM + Server binary | Go can do this natively |
| **Compile-Time Template Checks** | Macros + Type system | Transpiler does this |
| **Type Safety across Client/Server boundary** | `#[server]` macro generates both sides | `g-submit` concept |
| **No Runtime Template Error** | Compiler checks all template references | Transpiler checks |
| **No node_modules** | Cargo instead of npm | Go modules |
| **Thread-Safe Reactivity** | Rust's ownership system | Go goroutines |
| **Shared code without API contract** | Server Functions | Go functions as handlers |

### 5.2 Compile-Time Safety — How Rust Frameworks Use It

1. **RSX Macros (`view!`, `rsx!`, `html!`):** All tags, attributes, event handlers are checked against HTML spec at compile time. Typos in attributes = compiler error.

2. **Type-checked Props:** `#[component]` attributes can be Rust types. `Option<T>` for optional props, `Into<T>` for automatic conversion.

3. **Route Parameters as Types:** `/:id` → `use_params()` returns `Result<Params, _>` with statically derived types.

4. **`#[server]` Arguments:** Serialization is validated at compile time via `serde`.

### 5.3 WASM Story

All Rust frameworks use WASM for client-side interactivity. The WASM binary size is the biggest pain point:

- Base Leptos app: ~274 KB (CSR+Hydrate)
- Islands mode: ~24 KB (only interactive parts)
- `wasm-opt` + LTO: ~30-50% reduction

**Dreego learns from this:** WASM is not needed for Dreego (HTMX + Alpine.js replace it), but the **Islands architecture** is the conceptual model for Dreego's `<go>` block model: Server code runs WITHOUT a client bundle.

### 5.4 Emergent Patterns

| Pattern                      | Who does it           | Dreego Transfer                                 |
|------------------------------|-----------------------|-------------------------------------------------|
| **Server Functions**         | Leptos, Dioxus        | `g-submit` / `g-action` in `<go>` block        |
| **Fine-Grained Reactivity**  | Leptos, Sycamore      | Alpine.js / Datastar (client-side)              |
| **Signals as Primitive**     | All                   | `{#let}` + HTMX partials (server equivalent)    |
| **Islands Architecture**     | Leptos                | `<go>` block is self-contained, zero client code |
| **Context API**              | All                   | `c.Session()`, `c.User()` in `<go>`             |
| **VDOM vs Fine-Grained**     | Dioxus/Leptos duality | Dreego needs neither (SSR-only)                 |
| **Progressive Enhancement**  | Leptos `<ActionForm>` | Form actions, `<form>` without JS               |
| **File-Based Routing**       | Leptos (manual), Dioxus (code) | File-Based via `.dreego` files        |
| **Suspense/Streaming**       | Leptos, Sycamore      | Go templates + HTMX partials                    |

---

## 6. Transferable Features for Dreego

### 6.1 Adopt from Leptos

| Feature                      | Dreego Implementation                                           | Priority |
|------------------------------|----------------------------------------------------------------|----------|
| `#[server]` Macro            | `g-submit` / `g-action` → Generate Go handler from form name   | **V1**   |
| Fine-Grained Reactivity      | Not needed (SSR), but Alpine.js for client-side                | V1       |
| Islands Architecture         | `<go>` block = Island, no JS for non-interactive               | V1 (implicit) |
| `<ActionForm>`               | `<form g-submit="login">` with struct validation               | **V1**   |
| `provide_context`            | `c.Set("user", user)` in `<go>` block                          | V1       |
| `Memo`/`Derived`             | `{#let name = expr}` or Go variables in `<go>`                 | V1       |
| Suspense/Streaming           | Go `http.Flusher` + HTMX SSE                                  | V2       |
| Nested Routing               | File-Based + `layout.dreego`                                    | V1       |
| `#[prop(into)]`              | Auto-conversion via Go interfaces                              | V2       |
| `#[island]` Macro            | Dreego doesn't need it — `<script>` block is the island         | N/A      |

### 6.2 Adopt from Dioxus

| Feature                      | Dreego Implementation                                           | Priority |
|------------------------------|----------------------------------------------------------------|----------|
| `dx` CLI                     | `dreego dev`, `dreego build`, `dreego new`                       | **V1**   |
| Hot Reload                   | File watcher + SSE Browser Reload                               | **V1**   |
| Static Site Generation       | `dreego build --static` → HTML files                            | V2       |
| `GlobalSignal`               | Go package-level variables (not recommended, but possible)     | V2       |
| `use_memo`                   | `<go>` block cached computations via `sync.Once`               | V2       |
| HTML → RSX Translator        | `dreego convert` → HTML to `.dreego` (low priority)             | V3       |

### 6.3 Adopt from Yew

| Feature                      | Dreego Implementation                                           | Priority |
|------------------------------|----------------------------------------------------------------|----------|
| Elm-Reducer                  | Not needed, Go handlers suffice                                 | N/A      |
| Agents / Worker              | Go goroutines are workers                                       | V2       |

### 6.4 What Dreego Does NOT Need (and Why)

| Rust Feature                 | Why Dreego Doesn't Need It                                     |
|------------------------------|----------------------------------------------------------------|
| **WASM Compilation**         | HTMX + Alpine.js replace WASM interactivity, 0 KB overhead    |
| **Client-Side Signals**      | Alpine.js `x-data` / Datastar signals are 15 KB               |
| **Virtual DOM**              | SSR + HTML fragment swaps via HTMX                             |
| **RSX/JSX Macros**           | `.dreego` files are HTML-native, no macro needed               |
| **Hydration**                | HTMX loads HTML afterward, no hydration concept needed         |
| **Multi-Platform**           | Focus on Web, Single Binary suffices                           |

---

## 7. Architectural Lessons

### 7.1 The Server Functions Pattern Is Dominant

All major Rust frameworks have implemented server functions. It is THE way to model client-server interaction:

```
Leptos:   #[server] fn add_todo()  →  <ActionForm action=add_todo>
Dioxus:   #[server] fn save_user()  →  onclick: spawn(async { save_user().await })
Dreego:    <form g-submit="login">   →  func login(form LoginForm)
```

**The core idea:** The Rust/Go function IS the API. No separate REST endpoint, no manual serialization.

### 7.2 Compile-Time > Runtime

Rust frameworks maximize compile-time checks. Dreego does the same via transpiler:

| Rust                     | Dreego                                     |
|--------------------------|-------------------------------------------|
| `view! { <h1>{x}</h1> }` (macro checks tags) | `.dreego` → Transpiler checks template syntax |
| `#[server]` checks argument types | `g-submit` checks struct tags          |
| Route params as types    | File-Based routes become Go handlers      |

### 7.3 Islands Are the Future (But Dreego Is Already There)

Rust frameworks are just discovering the Islands architecture (Leptos 0.5+). Dreego is naturally "Islands":

- `<go>` block = Server code, zero client
- `<script>` block = Client interactivity, optional
- Template = pure HTML

That is the cleanest separation of server/client, which Rust frameworks only achieve through `#[island]` macros.

### 7.4 The Diamond Problem in Reactive Systems

Leptos' reactive graph elegantly solves the Diamond problem with a push-pull algorithm. For Dreego, this is academically interesting but not directly applicable — we have no client-side reactive graph. The lesson: **Fine-grained updates are the goal**, whether via signals (Leptos) or HTML fragments (HTMX).

---

## 8. Concrete Dreego Feature Ideas from This Analysis

### 8.1 Server Actions (inspired by Leptos `#[server]` + `<ActionForm>`)

```html
<!-- user/login.dreego -->
<form g-action="Login">
    <input name="email" type="email" />
    <input name="password" type="password" />
    <button>Login</button>
</form>

<go>
    type LoginForm struct {
        Email    string `form:"email" validate:"required,email"`
        Password string `form:"password" validate:"min=8"`
    }

    func Login(c *dreego.Context, form LoginForm) error {
        user, err := db.Authenticate(form.Email, form.Password)
        if err != nil {
            return dreego.Redirect("/login?error=invalid")
        }
        c.Session.Set("user", user)
        return dreego.Redirect("/dashboard")
    }
</go>
```

### 8.2 Context Cascade (inspired by `provide_context`)

```html
<go>
    // In layout.dreego — available for all child pages
    user := c.Session.Get("user")
    c.Provide("user", user)
</go>
```

### 8.3 Derived Values (inspired by `Memo`/Derived Signals)

```html
<go>
    items := db.GetItems()
    totalPrice := lo.SumBy(items, func(i Item) float64 { return i.Price })
    itemCount := len(items)
</go>

<p>Total: {totalPrice} ({itemCount} items)</p>
```

### 8.4 CLI-First Dev Experience (inspired by `dx` / `cargo-leptos`)

```
dreego new my-app              # Scaffold with layout.dreego, app.go
dreego dev                     # Transpiler watch + Hot-Reload + Tailwind
dreego build                   # Single Binary with embed
dreego routes                  # List of all generated routes
dreego add auth                # Install addon
dreego build --static          # SSG mode (V2)
```

### 8.5 Progressive Enhancement as First-Class (inspired by `<ActionForm>`)

Every `<form>` in Dreego works:
1. **Without JS:** Normal browser form, POST, full page reload
2. **With HTMX:** `hx-post`, fragment swap, no reload
3. **With Alpine.js/Datastar:** Client-side validation, optimistic UI

---

## 9. Conclusion

Rust frontend frameworks have undergone enormous development in the last 3 years. The central innovations are:

1. **Server Functions:** Functions instead of API endpoints — the most elegant client-server model
2. **Fine-Grained Reactivity:** Only change the DOM element that actually changed
3. **Islands Architecture:** Interactivity as opt-in, not as default
4. **Compile-Time Safety:** The type system as quality guarantee across the entire app

Dreego can adopt all these concepts — but in a fundamentally different model: Instead of WASM+Signals, Dreego relies on **SSR+HTMX+Alpine.js**. The result is the same (reactive, interactive web apps), but the path is simpler: no WASM, no hydration concept, no client-server isomorphism requirement.

The biggest lesson: **The Server Functions pattern is too good to ignore.** Dreego's `g-submit`/`g-action` should become the central interaction model.
