# Dreego Architecture

Dreego is an SSR-first Go framework with a compile-time `.dreego` transpiler.
Applications keep explicit Go ownership: generated code registers routes and
assets on an application, while the HTTP host is selected separately.

## Build flow

```text
.dreego sources
    -> lexer and parser
    -> semantic sections
    -> HTML and JavaScript language processors
    -> shared intermediate representation
    -> generated dree.go files
    -> normal Go compiler
    -> application binary
```

The HTML pipeline accepts HTML and Markdown. The JavaScript pipeline accepts
JavaScript, pinned TypeScript, and Dreego Browser Lua. See
[Anatomy of a `.dreego` File](file-anatomy.md).

## Runtime flow

```text
HTTP request
    -> adapter/ssr host
    -> application middleware
    -> generated route handler
    -> route server code
    -> head and body rendering
    -> optional layout
    -> HTTP response
```

`core` owns the public application, context, component, validation, and session
contracts. Implementation packages below `core/internal` separate server,
middleware, rendering, session, context, and validation responsibilities.
`adapter/ssr` is the first-party HTTP host.

## Repository boundaries

```text
core/                  public runtime facade
core/internal/         runtime implementation
adapter/ssr/           first-party HTTP host
internal/transpiler/   .dreego compiler and code generation
cmd/dreego/            command-line application
dreegotest/            public test helpers
_tests/go/             integration and regression tests
_docs/                 versioned public documentation
_docs/decisions/       architecture decisions
_todo/                 one file per open work item
```

The repository has one root Go module and one release tag. Generated application
`dree.go` files are build artifacts and are not committed by Dreego itself.

## Core and plugins

Core contains the capabilities required by a normal SSR application and keeps
its dependency boundary small. Optional providers and capabilities live in
separate plugin repositories with their own Go modules, releases, tests, and
CI. SSE, WebSockets, Tailwind integration, storage providers, and similar
features are plugins rather than core packages.

A provider-neutral core interface is introduced only after real implementations
demonstrate the same small contract. Plugins register behavior through the
owning application and must not weaken unrelated defaults.

## Client code

Dreego does not require a general browser runtime. Plain JavaScript and
TypeScript compile without one. Browser Lua links only the semantic helpers used
by the application into one generated asset. HTMX, Alpine.js, and normal browser
JavaScript remain valid progressive-enhancement choices.

## Generated ownership

Route, component, layout, and asset discovery happens during generation.
Generated registration is explicit so multiple applications can exist without
process-global route state. Errors should fail during generation or build when
possible, with `.dreego` source locations.

## Product boundary

SSR remains the production baseline. The experimental Wails v3 Phase 1 host is
the next target slice, followed by incremental DreeJS capabilities. Broader
Wails v3 Phase 2 support waits for a stable upstream release and evidence from
the Phase 1 reference application. Future extension points are kept only when
they are inexpensive and proven by real applications. See the
[Roadmap](roadmap.md) and architecture decisions for current direction.
