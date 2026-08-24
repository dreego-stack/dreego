# Product architecture

## Objective

Dreego is a Go-native application platform for building websites and desktop
applications with Svelte-inspired `.dreego` components. Developers write Go,
HTML, CSS, and optional client code without managing a JavaScript build system.
Dreego generates ordinary Go and browser assets through explicit, reproducible
build steps.

The target audience ranges from solo developers and small businesses to teams
that require predictable operations, strong typing, accessibility, and a path
to horizontally scalable applications. Breadth alone is not the quality bar:
each supported target must be deep enough to deploy and operate in production.

## Product promise

The same target-neutral application and component model can be hosted by
different first-party targets:

```text
.dreego source
      |
      v
parse -> typed model -> validate -> render plan
                                  |         |
                                  |         +-> head, styles, assets, client metadata
                                  v
                            typed HTML renderer
                                  |
                    +-------------+-------------+
                    |             |             |
                    v             v             v
                  SSR            SSG          Wails
```

SSR, SSG, and Wails are first-party target packages in the monorepo because
they depend on the same compiler, render contracts, component metadata, asset
rules, diagnostics, and compatibility policy. Optional provider integrations
remain external plugins.

## Public package direction

The intended package shape is:

```text
github.com/dreego-stack/dreego
github.com/dreego-stack/dreego/target/ssr
github.com/dreego-stack/dreego/target/ssg
github.com/dreego-stack/dreego/target/wails
```

The root package owns target-neutral application declarations, typed render
contracts, routes, components, and shared context capabilities. A host is
selected explicitly:

```go
app := dreego.New()
if err := www.Register(app); err != nil {
    log.Fatal(err)
}
if err := ssr.Run(app, ssr.Options{Address: ":8080"}); err != nil {
    log.Fatal(err)
}
```

The exact API is decided by the render-foundation implementation. This example
must not be copied into released documentation before it compiles.

## Target model

A target is a first-party host or build pipeline, not a generic feature flag.

- SSR binds a prepared application to `net/http` and request capabilities.
- SSG enumerates build-time route inputs and writes HTML and assets.
- Wails binds rendered HTML, assets, navigation, and a typed host bridge to a
  desktop WebView without requiring a local HTTP server.
- DreeJS is optional browser output shared by targets. It is not a target.
- Islands are not a separate product concept. DreeJS supplies narrowly scoped
  dynamic components through local code, fetch, polling, streams, or live
  connections.

One application may be used by more than one target. Components should not
branch on target names. They may require explicit capabilities, and a build
must fail when the selected target cannot provide them.

## Capability checks

Targets and processors exchange small capability declarations instead of a
single broad `Target` interface. Candidate capabilities include:

- static output;
- request and response access;
- server routes;
- client JavaScript assets;
- persistent server connections;
- desktop host bridge;
- build-time route enumeration.

Capabilities are added only when an implementation needs them. Missing
capabilities are build errors with the source location, requesting component or
plugin, selected target, and a concrete remediation.

## Section model

The planned `.dreego` root sections describe purpose, while `lang` describes
source language:

```html
<server lang="go"></server>
<head></head>
<body lang="html"></body>
<style lang="css"></style>
<client lang="js"></client>
```

Defaults are Go, HTML, CSS, and JavaScript. The semantic section names are part
of the pre-v0.1 implementation; optional languages remain future processors.

Only one body section exists per component. Dreego parses template constructs
such as `<@Component>`, `{#if}`, `{#each}`, `{ expression }`, and escaped output
outside the body-language processor. A Markdown body processor therefore
cannot consume, reinterpret, or erase Dreego control flow or component calls.

An HTML `<script>` nested inside `<body lang="html">` remains an ordinary HTML
element. Only the root `<client>` section represents client source processed by
Dreego.

## Processor and plugin boundary

JavaScript, Go, HTML, and CSS defaults are owned by the monorepo. Optional
languages and integrations live in external plugin repositories with their own
modules, dependencies, releases, tests, and CI.

Examples:

- TypeScript processor: client TypeScript to checked JavaScript;
- Markdown processor: body Markdown to HTML-compatible body nodes;
- Lua processor: client Lua to JavaScript;
- Tailwind plugin: managed external tool plus generated CSS;
- Stripe or MapLibre plugin: typed Go registration, components, assets, and
  optional processor or build capabilities.

Compiler plugins run through a versioned process protocol rather than Go's
native `plugin` package or an embedded Lua extension VM. The protocol must
carry structured inputs, outputs, diagnostics, source positions, assets, and
required target capabilities. It is validated with at least two real language
processors before a stability promise.

## Managed tools

Normal projects do not require developers to operate npm or Node directly.
Official plugins may install a pinned external tool after an explicit warning
and approval. Installation records version, checksum, permissions, and cache
location. CI supports a non-interactive allowlist and fails instead of silently
downloading an unapproved tool.

TypeScript is expected to be an official external processor. Raw JavaScript
remains the dependency-free core path. TypeScript checking must use a real
TypeScript type checker; syntax stripping alone is not called type safety.

## DreeJS direction

DreeJS is the umbrella name for generated browser support. Its implementation
is modular:

```text
dreejs/
├── core/
├── fetch/
├── poll/
├── stream/
└── live/
```

Generation includes only used modules. A static component emits no runtime.
Local presentation state may run in the browser, but authoritative business
state remains in Go or another explicit backend. A countdown may update and
switch a prepared view locally; payment completion, permissions, inventory,
and similar facts require server confirmation.

## Enterprise quality direction

Production suitability is measured through stable contracts, security,
accessibility, observability, performance, deployment documentation, upgrade
paths, and operational failure behavior. It is not inferred from the number of
targets or plugins.

Distributed sessions, pub/sub, caches, jobs, locks, and presence belong to
provider plugins. Dreego does not promise transparent synchronization of
arbitrary Go process memory. Shared core interfaces are introduced only after
multiple real implementations prove a common contract.

## Non-goals for the planned v0.x line

- A complete single-page-application runtime.
- Go-to-JavaScript compilation.
- Cloudflare Workers through Wasm.
- An embedded Lua VM used as Dreego's plugin engine.
- Provider-specific databases, queues, caches, auth, billing, or maps in core.
- One universal target interface that hides incompatible runtime models.

SPA and Wasm remain explicit future investigations after SSR, SSG, Wails, and
DreeJS demonstrate where those capabilities provide additional user value.
