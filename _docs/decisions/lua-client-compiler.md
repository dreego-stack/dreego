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

The supported language includes:

- comments and semicolon or newline separated statements;
- local declarations, assignments, and expression statements;
- strings, numbers, booleans, and `nil`;
- arithmetic, concatenation, comparison, and Lua `and`, `or`, and `not`;
- `if`, `elseif`, and `else` blocks;
- local and anonymous functions with lexical closures;
- bare returns and single return values;
- while loops, numeric for loops, and break;
- mutable tables with deterministic one-based sequence length;
- function calls, dotted browser object access, and colon-spelled browser methods;
- `print`, mapped through a linked runtime helper.

Lua truthiness and logical operators retain Lua behavior. In particular, zero
and empty strings are true, and `and` and `or` return operands.

## Deliberately excluded

- metatables and metamethods;
- named global functions, varargs, and multiple return values;
- generic iterators and coroutines;
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

Built-in names follow lexical scope. Local declarations, local functions, and
parameters may shadow a browser builtin without changing unrelated scopes.

Every compiled source block is isolated in its own JavaScript function scope.
The shared JavaScript output stage escapes HTML script-end sequences for Lua,
TypeScript, and raw JavaScript before writing an inline script element.
Lua output carries a sanitized `dreego:///` source URL and starting line so
browser runtime stacks retain the original `.dreego` source identity.

Browser Lua tables use private `Map` storage. Sequence fields begin at one,
assigning nil removes a key, and length is the contiguous prefix ending before
the first missing positive integer. This deterministic rule selects the
sequence border rather than exposing JavaScript array or object semantics.

## Consequences

- Dreego users can author browser behavior without maintaining JavaScript.
- Simple Lua stays close to readable JavaScript output.
- The runtime grows only with used semantic features at generation time.
- Full Lua compatibility is not implied and no VM is shipped.
