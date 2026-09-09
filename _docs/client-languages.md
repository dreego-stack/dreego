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

## Language scope

JavaScript, TypeScript, and Browser Lua are the complete supported client
language set for the foreseeable future. Dreego prioritizes compatibility,
diagnostics, and predictable generated output for these three processors over
adding more source languages.
