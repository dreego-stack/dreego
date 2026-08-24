# Dreego implementation plans

This directory contains the durable implementation direction for Dreego. It is
written for maintainers and implementation agents that need more detail than
the public, non-binding roadmap provides.

The files describe intended architecture, dependencies, boundaries, acceptance
criteria, and verification gates. They are not claims about currently released
behavior. Current behavior remains documented under `_docs/`.

## Reading order

1. [Product architecture](00-product-architecture.md)
2. [v0.1 SSR foundation](v0.1-ssr-foundation.md)
3. [v0.2 render foundation](v0.2-render-foundation.md)
4. [v0.3 language processors](v0.3-language-processors.md)
5. [v0.4 SSG target](v0.4-ssg-target.md)
6. [v0.5 Wails target](v0.5-wails-target.md)
7. [v0.6 DreeJS foundation](v0.6-dreejs-foundation.md)
8. [v0.7 DreeJS data and live updates](v0.7-dreejs-data-live.md)
9. [v0.8 stabilization](v0.8-stabilization.md)
10. [Future SPA and Wasm exploration](future-spa-wasm.md)

Version labels describe sequencing, not deadlines. Dreego is expected to stay
in v0.x for a long time. A phase may be split across several releases when its
contracts need validation.

## Rules for implementation agents

- Read `AGENTS.md`, this index, the current phase, and every dependency named by
  that phase before changing code.
- Treat examples as contract sketches until a todo or accepted ADR promotes an
  exact API.
- Implement one `_todo/` item per pull request unless the item explicitly says
  otherwise.
- Start behavioral changes with an integration test under `_tests/go/`.
- Keep generated application APIs typed. Do not replace page or component data
  with `map[string]any` for convenience.
- Do not introduce a universal target or plugin interface before at least two
  implementations prove the same small contract.
- Keep the root module dependency-free where the current dependency checks
  require it. Optional tools and their dependencies belong to external plugin
  repositories or managed tool processes.
- Update the relevant plan first when implementation evidence changes an
  architectural assumption.

## Relationship to `_todo/`

Plans explain a complete phase. Todos are the executable work queue. A plan can
contain many candidate slices without making all of them active work. Promote a
slice into one concrete `_todo/<area>/<id>.md` file when its dependencies are
met and the next pull request can complete it.

## Phase gate

A phase is complete only when:

- its public behavior is covered by black-box integration tests;
- generated code compiles through the supported Dreego CLI workflow;
- documentation and migration guidance match actual behavior;
- accessibility, security, race, and dependency checks pass where applicable;
- at least one reference application exercises the new capability;
- unresolved design questions are either answered or explicitly deferred.
