# Dreego implementation plans

This directory contains durable implementation direction for maintainers and
implementation agents. Plans use named capability phases rather than semantic
versions. Releases may complete part of a phase, and one phase may span several
releases.

## Reading order

1. [Product architecture](00-product-architecture.md)
2. [Completed SSR foundation](phase-ssr-foundation.md)
3. [Completed render foundation](phase-render-foundation.md)
4. [Completed multi-language Dreego](phase-multilanguage-dreego.md)
5. [Wails v3 Phase 1 and deferred Phase 2](phase-wails.md)
6. [DreeJS foundation and extensions](phase-dreejs.md)
7. [Live updates](phase-live.md)
8. [Stabilization](phase-stabilization.md)
9. [Future SPA and Wasm exploration](future-spa-wasm.md)

The render foundation's original behavioral inventory remains in
[phase-render-inventory.md](phase-render-inventory.md).

## Phase order

```text
SSR and render foundation (complete)
              |
              v
multi-language Dreego: md -> html, ts -> js, lua -> js
              |
              v
Wails v3 Phase 1 (planned v0.8)
              |
              v
DreeJS foundation (planned v0.9)
              |
              v
DreeJS data islands and extensions (planned from v0.10)
              |
              v
polling and live transports
              |
              v
stabilization and the v1 contract review
```

Wails v3 Phase 2 is a separate deferred phase. It begins only after an upstream
stable Wails v3 release and practical evidence from Dreego's Phase 1 reference
application. It is not assigned to v0.9 or v0.10.

Static site generation is not a planned phase. Dynamic SSR with explicit cache
policy and invalidation is the supported web deployment direction.

## Rules for implementation agents

- Read `AGENTS.md`, this index, the current phase, and every named dependency
  before changing code.
- Treat examples as contract sketches until a todo or accepted ADR promotes an
  exact API.
- Implement one `_todo/` item per pull request unless the item explicitly says
  otherwise.
- Start behavioral changes with an integration test under `_tests/go/`.
- Keep generated application APIs typed. Do not replace page or component data
  with `map[string]any` for convenience.
- Do not introduce a universal host, processor, cache, or live interface before
  real implementations prove the same small contract.
- Keep optional tools and provider dependencies outside the dependency-free
  core boundary.
- Update the relevant plan when implementation evidence changes an assumption.

## Relationship to `_todo/`

Plans explain a complete phase. Todos are the executable work queue. Promote a
slice into one concrete `_todo/<area>/<id>.md` file when its dependencies are
met and one pull request can complete it.

## Phase gate

A phase is complete only when:

- public behavior has black-box integration coverage;
- generated code compiles through the supported Dreego CLI workflow;
- documentation and migration guidance match released behavior;
- accessibility, security, race, and dependency checks pass;
- a reference application exercises the capability;
- unresolved design questions are answered or explicitly deferred.
