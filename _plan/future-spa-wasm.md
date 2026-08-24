# Future SPA and Wasm exploration

## Status

Exploration after the planned SSR, render, SSG, Wails, and DreeJS phases. This
file is not a release commitment.

## SPA question

SSR or SSG plus DreeJS may already cover most application interactions without
a client router or complete browser-owned state model. Explore a SPA runtime
only when real applications demonstrate a repeated limitation such as offline
navigation, application-wide client transitions, or high-frequency local state
that server-rendered component updates cannot serve well.

An experiment must measure:

- navigation latency and payloads against enhanced SSR;
- client bundle and memory cost;
- accessibility and focus behavior;
- data invalidation and mutation semantics;
- error recovery and offline behavior;
- duplicated Go and client business logic;
- operational complexity.

Do not label DreeJS a partial SPA or expand it incrementally without one
coherent navigation and state contract.

## Wasm question

Cloudflare Workers currently requires Go to run through WebAssembly rather than
as a native Go server. Wasm may also provide sandboxed compiler plugins or
browser-side Go, but those are separate use cases with different constraints.

Evaluate independently:

1. Cloudflare Workers server target.
2. Sandboxed processor execution through WASI.
3. Browser-side Go components.

Each experiment must account for binary size, startup, Go runtime compatibility,
filesystem and network APIs, debugging, deployment limits, dependency support,
and platform lock-in. Success in one use case does not justify the others.

## Promotion criteria

A future capability moves into a numbered plan only when:

- at least one real application cannot be served well by existing targets;
- a prototype demonstrates the user value;
- the implementation has a clear owner and maintenance cost;
- public contracts can remain separate from existing stable APIs;
- security, accessibility, performance, and deployment risks are understood.

## Explicit non-decision

Dreego does not currently commit to SPA, Go-to-JavaScript compilation,
browser-side Go, or a Cloudflare Workers target. The architecture should avoid
needlessly blocking them, but speculative abstractions must not delay the
planned targets.
