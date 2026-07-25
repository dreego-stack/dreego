
---
type: Reference
title: Rust Frontend Frameworks — Deep Dive & Dreego-Transfer
description: Deep dive into Rust frontend frameworks (Leptos, Dioxus, Yew) with transferable patterns for Dreego
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Rust Frontend Frameworks — Deep Dive & Dreego-Transfer

**Datum:** 23.07.2026
**Zweck:** Systematischer Vergleich der Rust-Frontend-Frameworks mit Fokus auf übertragbare Patterns für Dreego.

---

## Warum Rust-Frameworks für Dreego relevant sind

Rust und Go teilen entscheidende Eigenschaften:
- **Compiled, statically typed** → Compile-Time-Garantien möglich
- **Keine VM** → Direkte Code-Generierung, Single Binary
- **Std-Library-zentriert** → Wenig externe Abhängigkeiten
- **Kein JS-Ökosystem** → Eigenständige Lösungen für State, Routing, Templating
- **SSR als First-Class** → Alle großen Rust-Frameworks unterstützen SSR

Dreego-Spezifischer Kontext: Dreego ist SSR-First, braucht kein client-seitiges WASM-Framework. Wir lernen von Rust, wie Compile-Time-Safety, State-Management und Typed-Templating aussehen, aber im Kern bleibt Dreego ein Transpiler (.dreego → Go) mit SSR-Output.

---

## 1. Leptos (v0.8, 2026)

**Homepage:** https://leptos.dev | **Repo:** leptos-rs/leptos
**Ansatz:** Full-Stack, isomorphic, fine-grained reactivity
**Rendering-Modi:** CSR, SSR, Hydrate, Islands
**Inspiriert von:** Solid.js (Signals), Svelte (Compile-Time)

### 1.1 Reaktivität: Fine-Grained Signals

Leptos' Kerninnovation: **Fine-grained reactivity** ohne Virtual DOM. Das reaktive System ist ein gerichteter Graph aus drei Node-Typen:

```
Signals (src)  →  Memos (intermediate)  →  Effects (sink)
```

**Signals (Quellen):**
```rust
let (count, set_count) = signal(0);       // ReadSignal + WriteSignal
let count = RwSignal::new(0);             // Combined read/write
```

**Memos (abgeleitete Werte, gecached):**
```rust
let doubled = Memo::new(move |_| count.get() * 2);
```

**Effects (Side Effects, z.B. DOM-Updates):**
```rust
Effect::new(move |_| {
    log!("count changed: {}", count.get());
});
```

**Derived Signals (nicht gecached, leichtgewichtig):**
```rust
let doubled = move || count.get() * 2;  // re-evaluated at every read
```

**Entscheidendes Design-Detail:**
- Signals markieren sich immer als "Dirty" (keine `PartialEq`-Prüfung)
- Memos prüfen `PartialEq` bevor sie Subscriber benachrichtigen
- Effects werden maximal einmal pro Update-Zyklus ausgeführt (Diamond-Problem gelöst via "push-pull"-Algorithmus)
- Der komplette reaktive Graph arbeitet thread-safe (`Send + Sync`)

### 1.2 Server Functions — Die Killer-App

**Das Problem:** Client und Server teilen Code, aber API-Contract muss dreifach geschrieben werden (Client-Fetch, Server-Handler, Typ-Definition).

**Leptos' Lösung: `#[server]` Makro**

```rust
#[server]
pub async fn add_todo(title: String) -> Result<(), ServerFnError> {
    let mut conn = db().await?;
    sqlx::query("INSERT INTO todos (title) VALUES ($1)")
        .bind(title)
        .execute(&mut conn).await?;
    Ok(())
}

// Aufruf im Client — genau dieselbe Funktion!
#[component]
fn TodoForm() -> impl IntoView {
    view! {
        <button on:click=move |_| {
            spawn_local(async { add_todo("Buy milk".into()).await; });
        }>"Add"</button>
    }
}
```

**Was passiert unter der Haube:**
1. `#[server]` generiert *zwei Versionen* der Funktion via Conditional Compilation:
   - **Server-Build (ssr):** Echter Code, HTTP-Endpoint via POST (URL-encoded args, JSON return)
   - **Client-Build (hydrate/csr):** Stub der einen `fetch()`-Aufruf mit serialisierten Args macht
2. Argumente werden per `serde_qs` URL-encoded (kompatibel mit `<form>`)
3. Return-Values per `serde_json`
4. Custom Encodings möglich (JSON, Multipart, etc.)

**Warum ist das genial?**
- **Co-location:** DB-Query und Click-Handler in derselben Datei
- **Type-Safety über die Grenze:** Argument-Typen werden auf beiden Seiten geprüft
- **Progressive Enhancement:** `<ActionForm>` funktioniert auch ohne WASM/JS — Browser `<form>` reicht
- **Kein API-Contract:** Die Rust-Funktion *ist* der Contract

### 1.3 SSR + Hydration

Leptos kompiliert die App für **zwei Targets**:
1. **Server (ssr):** Rendert HTML, injiziert WASM-Loader
2. **Browser (hydrate):** Nimmt existierendes HTML, hängt Event-Listener dran

**Page-Load-Lifecycle:**
1. Browser `GET` → Server matched Route → rendert `<App/>` einmal → gibt HTML zurück
2. Browser lädt WASM parallel → WASM "wandert" das HTML ab ("hydration"), hängt Reaktivität an
3. Danach: Alle weiteren Navigationen sind **SPA-Navigationen** (kein Page Reload)

**Wichtig:** Der komplette Komponentenbaum muss auf Server UND Client laufen können. Das ist der Preis für Isomorphie.

### 1.4 Islands Architecture (Leptos 0.5+)

**Konzept:** Statt die ganze App zu hydratieren, nur kleine "Inseln" der Interaktivität aktivieren.

```rust
#[island]  // ← statt #[component]
fn Counter() -> impl IntoView {
    let (count, set_count) = signal(0);
    view! { <button on:click=move |_| set_count.update(|n| *n += 1)>"Count: {count}"</button> }
}
```

**Effekte:**
- **WASM-Größe sinkt drastisch:** 274 KB → 24 KB (Hello World), 400 KB → 200 KB (Tabs-Demo)
- **Server-Code in Components:** `#[component]` läuft NUR auf dem Server (kein WASM)
- **Children aus Server-Code:** Server-Side-Content kann als Children an Islands übergeben werden ("Donut-Architektur")
- **Context zwischen Islands:** `provide_context`/`expect_context` funktioniert island-übergreifend

### 1.5 Routing

- **Declaration-based:** `<Routes>` + `<Route path=... view=...>`
- **Nested Routes:** `<Route>` können `<Route>` children haben → Layouts
- **Params:** `/users/:id` → `let params = use_params();`
- **Queries:** `?search=...` → `let query = use_query();`
- **`<A/>` und `<Form/>`:** Router-aware Links/Forms, client-seitige Navigation
- **Fallback:** `<Routes fallback=|| view! { <NotFound/> }>`

### 1.6 Komponenten-Modell

```rust
#[component]
fn Button(
    #[prop(into)] class: String,     // Into<T> Konvertierung
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

### 1.7 Was Leptos einzigartig macht

1. **Fine-Grained Reactivity:** Nur exakt die DOM-Knoten werden aktualisiert, die ein Signal lesen — kein VDOM Diffing
2. **`#[server]`-Makro:** Die eleganteste Client-Server-Integration aller Rust-Frameworks
3. **Islands-Mode:** Reduziert WASM-Größe um 50%, ermöglicht serverseitige `std::fs` in Components
4. **Streaming SSR:** `<Suspense/>` + out-of-order Streaming
5. **`view!`-Makro:** RSX (React-like JSX in Rust) mit Compile-Time-Checks gegen HTML-Spec
6. **Reaktiver Graph ist `Send + Sync`:** Multi-Thread-SSR auf Tokio

---

## 2. Dioxus (v0.6/0.7, 2026)

**Homepage:** https://dioxuslabs.com | **Repo:** DioxusLabs/dioxus
**Ansatz:** Cross-Platform (Web, Desktop, Mobile) mit React-ähnlichem Modell
**Inspiriert von:** React, Flutter
**24.5k+ GitHub Stars**

### 2.1 Architektur: Virtual-DOM-basiert (anders als Leptos)

Im Gegensatz zu Leptos' fine-grained Reactivity verwendet Dioxus ein Virtual DOM:
- Komponenten re-rendern vollständig bei State-Änderung
- VDOM wird diffed, nur Änderungen landen im echten DOM
- Ähnlich wie React, aber mit Rust-Performance

### 2.2 Multi-Platform Rendering

```rust
// Derselbe Code, drei Targets:
#[cfg(target_platform = "web")]     // → WASM im Browser
#[cfg(target_platform = "desktop")] // → WGPU-basierter Renderer (Blitz)
#[cfg(target_platform = "mobile")]  // → iOS/Android via WGPU
```

**Desktop/Mobile:** Dioxus hat einen eigenen Renderer namens **Blitz** (basierend auf WGPU), der HTML/CSS nativ rendert — ohne WebView, ohne Electron. Das ist das "Flutter, but better"-Versprechen.

### 2.3 State Management: Signals (seit v0.5)

```rust
fn App() -> Element {
    let mut count = use_signal(|| 0);       // Signal Hook
    let half = use_memo(move || count() / 2); // Abgeleitet, gecached
    let data = use_resource(move || async {   // Async Data
        fetch_data(count()).await
    });

    use_effect(move || {
        log!("count changed to {count}");    // Side-Effekt
    });

    rsx! {
        button { onclick: move |_| count += 1, "Count: {count}" }
        div { "Half: {half}" }
    }
}
```

**State-Weitergabe:**
1. **Props:** `ReadOnlySignal<T>` — Standard
2. **Context:** `use_context_provider` / `use_context` — für Teilbäume
3. **Globals:** `GlobalSignal<T>` als `static` — app-weit

### 2.4 Fullstack Mode

```rust
// Server Function (identisch zu Leptos' Konzept)
#[server]
async fn save_user(name: String) -> Result<(), ServerFnError> {
    db::insert_user(&name).await?;
    Ok(())
}

// Im Client:
rsx! { button { onclick: move |_| { spawn(async { save_user("Alice".into()).await; }); } } }
```

**Hydration:** Server rendert HTML → serialisiert Server-Future-Daten als JSON → Client hydriert.

**Wichtiger Unterschied zu Leptos:**
- Dioxus benötigt `use_server_future` für non-deterministische Daten
- Hydration-Mismatches sind ein echtes Problem (Random-Werte, Timeouts)
- Dioxus' VDOM-Ansatz macht Hydration aufwändiger als Leptos' Fine-Grained-Approach

### 2.5 CLI & Dev Experience

```
dx new my-app        # Scaffolding
dx serve             # Dev-Server + Hot Reload
dx build             # Production Build
dx translate         # HTML → RSX Konvertierung
```

**Hot Reload:** Änderungen an RSX werden ohne State-Verlust im Browser aktualisiert.

### 2.6 Static Site Generation (SSG)

Dioxus unterstützt SSG via Fullstack-Mode: Der Server rendert alle Routen zu statischem HTML, das dann deployed wird.

### 2.7 Was Dioxus einzigartig macht

1. **True Cross-Platform:** Web, Desktop, Mobile — ein Codebase
2. **Blitz Renderer:** HTML/CSS nativ auf GPU gerendert, kein WebView
3. **React-Ähnlichkeit:** Entwickler von React/Flutter fühlen sich sofort zuhause
4. **`dx` CLI:** Erstklassige Developer Experience
5. **Ecosystem:** Größtes Rust-UI-Ökosystem (24.5k Stars)

### 2.8 Dioxus LiveView-Mode

Dioxus' SSR lässt sich als "LiveView" betreiben: Der Server hält den VirtualDom im Speicher, rendert bei Änderungen HTML-Diffs und schickt sie via WebSocket an den Client. Dies ist ähnlich zu Phoenix LiveView. Der Client muss kein WASM laden — alle Events werden per WebSocket zum Server geschickt.

---

## 3. Yew (v0.23, 2026)

**Homepage:** https://yew.rs | **Repo:** yewstack/yew
**Ansatz:** Elm-inspiriertes Komponenten-Framework (CSR-first, SSR optional)
**Ältestes Rust-Frontend-Framework** (seit 2017)

### 3.1 Architektur: Elm/React-Modell

Yew ist der Pionier — das erste Rust-Frontend-Framework. Es folgt dem Elm/React-Modell:

- **Virtual DOM** (wie React)
- **Komponenten** sind Rust-Funktionen mit `#[function_component]`
- **Message/Update-Architektur** (von Elm) via `use_reducer`
- **Props** sind reine Rust-Structs
- **SSR** ist additiv (nicht First-Class wie bei Leptos/Dioxus)

### 3.2 Komponenten-Modell

```rust
#[function_component]
fn HelloWorld() -> Html {
    html! { <p>{"Hello world"}</p> }
}
```

**Zwei Komponenten-Arten:**
1. **Function Components** (`#[function_component]`) — empfohlen
2. **Struct Components** (advanced, mehr Kontrolle)

**Lifecycle:** `use_state`, `use_effect`, `use_reducer`, `use_memo`, `use_callback`

### 3.3 State Management

Yew's State-Hooks sind React-ähnlich:

| Hook             | Verhalten                               | Scope    |
|-------------------|-----------------------------------------|----------|
| `use_state`       | State, re-rendert bei `set`            | Component|
| `use_state_eq`    | Wie `use_state`, aber mit `PartialEq`  | Component|
| `use_reducer`     | Elm-Architektur: Action → State        | Component|
| `use_reducer_eq`  | Reducer + `PartialEq`                  | Component|
| `use_memo`        | Memoized Computation                   | Component|
| `use_callback`    | Memoized Callback (Deps-Array)         | Component|
| `use_mut_ref`     | Mutable Ref (kein Re-Render)           | Component|

**Context:** `use_context` für Baum-weite State-Weitergabe.
**Agents:** Yew-spezifisches Konzept für Worker-Threads / Hintergrund-Tasks.

### 3.4 SSR

Yew's SSR ist funktional, aber nicht First-Class:
```rust
let renderer = yew::ServerRenderer::<App>::new();
let html = renderer.render().await;
```

Kein Hydration-Konzept, kein Server-Functions, kein File-Based Routing.

### 3.5 Yews Bedeutung heute

Historisch extrem wichtig (Go-To für Rust-WASM 2017–2021). Aber von Leptos und Dioxus in den meisten Aspekten überholt. Lehren:

- **Message/Reducer-Architektur** ist für komplexe State-Maschinen gut, aber für UI zu boilerplate-lastig
- **VDOM** ist konzeptionell einfach, aber fine-grained Reactivity (Leptos) performanter
- **Community-getriebene Entwicklung** — Yew wurde von der Community gebaut, nicht von einem Startup

---

## 4. Andere Rust-Frameworks

### 4.1 Sycamore (v0.9, 2024)

**Ansatz:** Fine-Grained Reactivity (wie Leptos, aber schlanker)

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
- Fine-grained reactivity (kein VDOM)
- SSR + Hydration
- Async/Suspense
- Built-in Routing (Client + Server)
- 3.3k Stars, aktiv gepflegt

**Relevanz:** Sycamore teilt viele Ideen mit Leptos, ist aber kleiner und weniger Full-Stack-orientiert. Proof dass Fine-Grained Reactivity in Rust gut funktioniert.

### 4.2 Perseus (veraltet)

Perseus war ein Full-Stack-Framework inspiriert von Next.js. Features:
- SSG (Static Site Generation) als First-Class-Feature
- i18n built-in
- State-Plugins (globaler State via Context)

Veraltet/verschmolzen. Die SSG-Ideen sind in Dioxus und Leptos aufgegangen.

### 4.3 MoonZoon (veraltet)

Minimalistisch, "no-JS, no-CSS, no-HTML" — komplett in Rust. War ein experimentelles Framework das Rust-APIs direkt in Browser-Elemente umsetzte.

---

## 5. Cross-Cutting Patterns

### 5.1 Was Rust-Frameworks können, das JS-Frameworks nicht können

| Fähigkeit                  | Ermöglicht durch              | Relevanz für Dreego |
|----------------------------|-------------------------------|---------------------|
| **Single Binary Deployment**| Rust compiliert zu WASM + Server-Binary | Go kann das nativ |
| **Compile-Time Template-Checks** | Makros + Typ-System | Transpiler macht das |
| **Typ-Sicherheit über Client/Server-Grenze** | `#[server]` Makro generiert beide Seiten | `g-submit`-Konzept |
| **Kein Runtime-Template-Error** | Compiler prüft alle Template-Referenzen | Transpiler-Checks |
| **Keine node_modules** | Cargo statt npm | Go-Module |
| **Thread-Safe Reactivity** | Rusts Ownership-System | Go-Goroutines |
| **Shared Code ohne API-Contract** | Server Functions | Go-Funktionen als Handler |

### 5.2 Compile-Time Safety — Wie Rust-Frameworks sie nutzen

1. **RSX-Makros (`view!`, `rsx!`, `html!`):** Alle Tags, Attribute, Event-Handler werden zur Compile-Zeit gegen HTML-Spec geprüft. Tippfehler in Attributen = Compiler-Fehler.

2. **Typ-geprüfte Props:** `#[component]`-Attribute können Rust-Typen sein. `Option<T>` für optionale Props, `Into<T>` für automatische Konvertierung.

3. **Route-Parameter als Typen:** `/:id` → `use_params()` gibt `Result<Params, _>` mit statisch abgeleiteten Typen.

4. **`#[server]`-Argumente:** Serialisierung wird via `serde` zur Compile-Zeit validiert.

### 5.3 WASM-Story

Alle Rust-Frameworks nutzen WASM für Client-Side-Interaktivität. Die WASM-Binary-Größe ist der größte Pain-Point:

- Basis-Leptos-App: ~274 KB (CSR+Hydrate)
- Islands-Mode: ~24 KB (nur interaktive Teile)
- `wasm-opt` + LTO: ~30-50% Reduktion

**Dreego lernt daraus:** WASM ist für Dreego nicht nötig (HTMX + Alpine.js ersetzen es), aber die **Islands-Architektur** ist das konzeptionelle Vorbild für Dreegos `<go>`-Block-Modell: Server-Code läuft OHNE Client-Bundle.

### 5.4 Emergente Patterns

| Pattern                     | Wer macht's          | Dreego-Übertragung                             |
|-----------------------------|----------------------|------------------------------------------------|
| **Server Functions**        | Leptos, Dioxus       | `g-submit` / `g-action` im `<go>`-Block       |
| **Fine-Grained Reactivity** | Leptos, Sycamore     | Alpine.js / Datastar (Client-Seite)            |
| **Signals als Primitive**   | Alle                 | `{#let}` + HTMX partials (Server-Äquivalent)   |
| **Islands Architecture**    | Leptos               | `<go>`-Block ist autark, keinerlei Client-Code |
| **Context API**             | Alle                 | `c.Session()`, `c.User()` im `<go>`            |
| **VDOM vs Fine-Grained**    | Dioxus/Leptos-Dualität| Dreego braucht beides nicht (SSR-only)          |
| **Progressive Enhancement** | Leptos `<ActionForm>` | Form-Actions, `<form>` ohne JS                 |
| **File-Based Routing**      | Leptos (manual), Dioxus (code) | File-Based via `.dreego`-Dateien       |
| **Suspense/Streaming**      | Leptos, Sycamore     | Go-Templates + HTMX partials                   |

---

## 6. Übertragbare Features für Dreego

### 6.1 Von Leptos übernehmen

| Feature                     | Dreego-Implementierung                                         | Priorität |
|-----------------------------|---------------------------------------------------------------|-----------|
| `#[server]` Makro           | `g-submit` / `g-action` → Go-Handler aus Form-Name generieren | **V1**    |
| Fine-Grained Reactivity     | Nicht nötig (SSR), aber Alpine.js für Client-Seite            | V1        |
| Islands-Architektur         | `<go>`-Block = Island, kein JS für Nicht-Interaktives         | V1 (implizit) |
| `<ActionForm>`              | `<form g-submit="login">` mit Struct-Validierung              | **V1**    |
| `provide_context`           | `c.Set("user", user)` im `<go>`-Block                         | V1        |
| `Memo`/`Derived`            | `{#let name = expr}` oder Go-Variablen im `<go>`              | V1        |
| Suspense/Streaming          | Go `http.Flusher` + HTMX SSE                                 | V2        |
| Nested Routing              | File-Based + `layout.dreego`                                   | V1        |
| `#[prop(into)]`             | Auto-Konvertierung via Go-Interfaces                          | V2        |
| `#[island]` Makro           | Dreego braucht's nicht — `<script>`-Block ist die Insel        | N/A       |

### 6.2 Von Dioxus übernehmen

| Feature                     | Dreego-Implementierung                                         | Priorität |
|-----------------------------|---------------------------------------------------------------|-----------|
| `dx` CLI                    | `dreego dev`, `dreego build`, `dreego new`                      | **V1**    |
| Hot Reload                  | File-Watcher + SSE Browser Reload                             | **V1**    |
| Static Site Generation      | `dreego build --static` → HTML-Dateien                         | V2        |
| `GlobalSignal`              | Go-Package-Level-Variablen (nicht empfohlen, aber möglich)    | V2        |
| `use_memo`                  | `<go>`-Block cached Berechnungen via `sync.Once`              | V2        |
| HTML → RSX Translator       | `dreego convert` → HTML nach `.dreego` (niedrige Prio)         | V3        |

### 6.3 Von Yew übernehmen

| Feature                     | Dreego-Implementierung                                         | Priorität |
|-----------------------------|---------------------------------------------------------------|-----------|
| Elm-Reducer                 | Nicht nötig, Go-Handler reichen                               | N/A       |
| Agents / Worker             | Go-Goroutines sind Worker                                     | V2        |

### 6.4 Was Dreego NICHT braucht (und warum)

| Rust-Feature                | Warum Dreego es nicht braucht                                  |
|-----------------------------|---------------------------------------------------------------|
| **WASM-Compilation**        | HTMX + Alpine.js ersetzen WASM-Interaktivität, 0 KB Overhead |
| **Client-Side Signals**     | Alpine.js `x-data` / Datastar Signals sind 15 KB              |
| **Virtual DOM**             | SSR + HTML-Fragment-Swaps via HTMX                            |
| **RSX/JSX-Makros**          | `.dreego`-Dateien sind HTML-nativ, kein Makro nötig            |
| **Hydration**               | HTMX lädt HTML nach, kein Hydration-Konzept nötig             |
| **Multi-Platform**          | Fokus auf Web, Single Binary reicht                           |

---

## 7. Architektonische Lehren

### 7.1 Das Server-Functions-Pattern ist dominant

Alle großen Rust-Frameworks haben Server Functions implementiert. Es ist DER Weg, Client-Server-Interaktion zu modellieren:

```
Leptos:   #[server] fn add_todo()  →  <ActionForm action=add_todo>
Dioxus:   #[server] fn save_user()  →  onclick: spawn(async { save_user().await })
Dreego:    <form g-submit="login">   →  func login(form LoginForm)
```

**Die Kernidee:** Die Rust/Go-Funktion IST die API. Kein separater REST-Endpoint, keine manuelle Serialisierung.

### 7.2 Compile-Time > Runtime

Rust-Frameworks maximieren Compile-Time-Checks. Dreego macht dasselbe via Transpiler:

| Rust                    | Dreego                                    |
|-------------------------|------------------------------------------|
| `view! { <h1>{x}</h1> }` (Makro prüft Tags) | `.dreego` → Transpiler prüft Template-Syntax |
| `#[server]` prüft Argument-Typen | `g-submit` prüft Struct-Tags          |
| Route-Params als Typen  | File-Based Routes werden zu Go-Handlern  |

### 7.3 Islands sind die Zukunft (aber Dreego ist schon da)

Rust-Frameworks entdecken gerade die Islands-Architektur (Leptos 0.5+). Dreego ist von Natur aus "Islands":

- `<go>`-Block = Server-Code, null Client
- `<script>`-Block = Client-Interaktivität, optional
- Template = reines HTML

Das ist die sauberste Trennung von Server/Client, die Rust-Frameworks erst über `#[island]`-Makros erreichen.

### 7.4 Das Diamond-Problem in reaktiven Systemen

Leptos' reaktiver Graph löst das Diamond-Problem elegant mit einem Push-Pull-Algorithmus. Für Dreego ist das akademisch interessant aber nicht direkt anwendbar — wir haben keinen client-seitigen reaktiven Graphen. Die Lehre: **Fine-grained Updates sind das Ziel**, egal ob über Signals (Leptos) oder HTML-Fragments (HTMX).

---

## 8. Konkrete Dreego-Feature-Ideen aus dieser Analyse

### 8.1 Server-Actions (inspiriert von Leptos `#[server]` + `<ActionForm>`)

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

### 8.2 Context-Kaskade (inspiriert von `provide_context`)

```html
<go>
    // In layout.dreego — für alle Kind-Seiten verfügbar
    user := c.Session.Get("user")
    c.Provide("user", user)
</go>
```

### 8.3 Derived Values (inspiriert von `Memo`/Derived Signals)

```html
<go>
    items := db.GetItems()
    totalPrice := lo.SumBy(items, func(i Item) float64 { return i.Price })
    itemCount := len(items)
</go>

<p>Total: {totalPrice} ({itemCount} items)</p>
```

### 8.4 CLI-First Dev-Experience (inspiriert von `dx` / `cargo-leptos`)

```
dreego new my-app              # Scaffold mit layout.dreego, app.go
dreego dev                     # Transpiler-Watch + Hot-Reload + Tailwind
dreego build                   # Single Binary mit embed
dreego routes                  # Liste aller generierten Routen
dreego add auth                # Addon installieren
dreego build --static          # SSG-Mode (V2)
```

### 8.5 Progressive Enhancement als First-Class (inspiriert von `<ActionForm>`)

Jedes `<form>` in Dreego funktioniert:
1. **Ohne JS:** Normales Browser-Formular, POST, Full-Page-Reload
2. **Mit HTMX:** `hx-post`, Fragment-Austausch, kein Reload
3. **Mit Alpine.js/Datastar:** Client-seitige Validierung, optimistisches UI

---

## 9. Fazit

Rust-Frontend-Frameworks haben in den letzten 3 Jahren eine enorme Entwicklung durchgemacht. Die zentralen Innovationen sind:

1. **Server Functions:** Funktionen statt API-Endpoints — das eleganteste Client-Server-Modell
2. **Fine-Grained Reactivity:** Nur das DOM-Element ändern, das sich wirklich geändert hat
3. **Islands Architecture:** Interaktivität als Opt-in, nicht als Default
4. **Compile-Time Safety:** Das Typ-System als Qualitätsgarant über die gesamte App

Dreego kann alle diese Konzepte aufgreifen — aber in einem grundlegend anderen Modell: Statt WASM+Signals setzt Dreego auf **SSR+HTMX+Alpine.js**. Das Ergebnis ist dasselbe (reaktive, interaktive Web-Apps), aber der Weg ist einfacher: kein WASM, kein Hydration-Konzept, kein Client-Server-Isomorphie-Zwang.

Die größte Lektion: **Das Server-Functions-Pattern ist zu gut, um es zu ignorieren.** Dreegos `g-submit`/`g-action` sollte das zentrale Interaktionsmodell werden.
