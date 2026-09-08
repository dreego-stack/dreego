# Client Languages

Dreego normalizes client source to JavaScript. Each language has an explicit
compatibility target and stability level. Support for a language in `<client>`
does not imply support in `<server>`, `<body>`, or other sections.

## Primary language contracts

| Input | Compatibility target | Status | Dreego behavior |
|---|---|---|---|
| JavaScript | ECMAScript 2022 browser baseline | Primary | Passed through without type checking, polyfills, or syntax lowering |
| TypeScript | TypeScript 7.0.2, strict mode, ES2022 output | Primary | Type checked and transpiled by the pinned native compiler |
| Lua | Dreego Browser Lua 0.5, informed by Lua 5.4 | Stabilizing | Compiled by Dreego with an explicitly documented browser subset and feature-linked helpers |

JavaScript, TypeScript, and Lua are the intended production language set. Lua
joins that promise only after its v0.5 hardening and feature series completes.

JavaScript newer than ES2022 may pass through unchanged, but Dreego does not
claim compatibility or provide polyfills for it. Applications remain
responsible for selecting browsers that understand their raw JavaScript.

TypeScript uses `--strict`, `--target ES2022`, and `--module preserve`. The
compiler version is part of the reproducible build contract. Generated Go
models are exposed to the type checker where Dreego owns their schema.

Browser Lua is not advertised as complete Lua 5.4. Its exact supported syntax,
semantics, browser boundary, and exclusions are listed in [Browser Lua](lua.md).
This avoids silently accepting code with different number, string, table,
metatable, coroutine, module, or standard-library behavior.

## Experimental languages

Client Go and Starlark are planned as isolated experiments after the stable Lua
patch series:

- client Go explores a browser-safe, statically typed Go subset that can share
  Dreego-generated models and Wails bindings;
- Starlark explores concise Python-like client authoring with a deliberately
  small and deterministic language surface.

Experimental means no compatibility promise. Syntax, processor identifiers,
generated output, diagnostics, and the experiment itself may change or be
removed before v1. An experiment must live in its own processor package, add no
runtime behavior to applications that do not use it, and must not weaken the
stable JavaScript, TypeScript, or Lua paths.

Neither experiment may claim compatibility with full Go or Python. Before a
processor is accepted, its proposal must define:

1. the precise source-language subset and JavaScript target;
2. type and runtime semantics at the DOM and Wails boundaries;
3. generated asset and bundle budgets;
4. diagnostics, source mapping, accessibility, and security behavior;
5. removal criteria and evidence required for promotion.

The final `lang` spelling and opt-in mechanism remain intentionally unreserved
until each experiment has a tested vertical slice. This keeps removal cheap and
prevents an exploratory name from becoming an accidental public API.
