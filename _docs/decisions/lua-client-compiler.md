---
type: Decision
title: Browser Lua compiler and feature-linked runtime
status: Accepted
---

# Browser Lua compiler and feature-linked runtime

## Context

Dreego supports JavaScript and TypeScript client sections, but both retain
JavaScript semantics. Lua is intended as a maintainable authoring language for
browser behavior, not as a plugin engine or server language.

## Decision

`js/lua` is a first-party compiler inside the Dreego transpiler. It owns a
lexer, parser, AST, semantic lowering, JavaScript emitter, and a registry of
small runtime helpers. It does not embed a Lua VM, WebAssembly runtime, Node,
or an external Lua compiler.

Compilation returns JavaScript and the semantic features it requires. The
shared JavaScript output stage resolves transitive helper dependencies and
emits only the deterministic runtime subset used by the compiled source.

Semantically equivalent operations become direct JavaScript. Operations whose
Lua behavior differs use a helper. Unsupported syntax or behavior fails during
generation with a `.dreego` source location.

## MVP language

The first release supports:

- comments and semicolon or newline separated statements;
- local declarations, assignments, and expression statements;
- strings, numbers, booleans, and `nil`;
- arithmetic, concatenation, comparison, and Lua `and`, `or`, and `not`;
- `if`, `elseif`, and `else` blocks;
- local and anonymous functions with lexical closures;
- bare returns and single return values;
- function calls, dotted browser object access, and colon-spelled browser methods;
- `print`, mapped through a linked runtime helper.

Lua truthiness and logical operators retain Lua behavior. In particular, zero
and empty strings are true, and `and` and `or` return operands.

## Deliberately excluded from the MVP

- tables, metatables, and metamethods;
- named global functions, varargs, and multiple return values;
- loops, iterators, and coroutines;
- modules, `require`, dynamic `load`, bytecode, and the debug library;
- filesystem, process, socket, and native-library APIs;
- a claim of complete Lua 5.x compatibility.

These are added through tested patch releases only when their browser contract
and runtime cost are explicit. Browser-inapplicable standard-library features
may remain permanently unsupported.

Colon-spelled calls target native browser methods. They preserve JavaScript's
receiver binding and do not inject a second Lua-style `self` argument.

## Output contract

The internal compiler result contains code and a sorted feature set. Runtime
helpers are namespaced under `dreegoLua`; application code receives no other
compiler-owned globals. Identical source and compiler versions produce
identical output.

## Consequences

- Dreego users can author browser behavior without maintaining JavaScript.
- Simple Lua stays close to readable JavaScript output.
- The runtime grows only with used semantic features at generation time.
- Full Lua compatibility is not implied and no VM is shipped.
