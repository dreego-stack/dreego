---
type: concept
status: draft
related: _docs/components.md, _docs/testing.md, _todo/core/component-correctness.1.md
---
# Component Call Contract

## Problem

Component calls currently work for the happy path, but the contract for named
props, typed expressions, default/named slots, and error reporting is implicit.
Without a written contract, edge cases such as duplicate props, unknown props,
wrong-type expressions, and slot leakage across siblings are handled
inconsistently or only fail late at `go build`. This document defines the
selected contract that the component-correctness work must implement.

## Goals

1. Calls use named props only; positional binding is not supported.
2. Props are order-independent within a single call.
3. Errors are source-aware and fail at `dreego generate` whenever possible.
4. Default and named slots are scoped to one invocation and cannot leak.
5. Self-closing calls reject children; calls with children provide default slot
   content unless named slots are used.
6. Every increment ships with HTTP-level assertions on rendered output.

## Acceptance Criteria

### Named props

- `Component Card (title string, count int)`
  Call `<@Card title="Hi" count={3}/>` renders.
  Call `<@Card count={3} title="Hi"/>` renders identically.

### Error contract

Errors are reported as:

```
<file>:<line>:<col>: <component> <prop>: <message>
```

Examples:

```
components/Card.dreego:7:12: Card title: missing required prop "title"
routes/index.dreego:4:3: Card count: unknown prop "count"
routes/index.dreego:4:18: Card title: duplicate prop "title"
routes/index.dreego:5:9: Card title: expected string, got int
```

- Missing required props produce one error per missing prop.
- Unknown props produce one error per unknown prop.
- Duplicate props produce one error per duplicated name.
- Wrong-type expression props produce one error per mismatch.
- Prop type checking uses simple literal inference only: it distinguishes string, int, and other literal kinds directly visible in the call. It does not infer types from arbitrary Go expressions or calling scope via AST analysis.

### Self-closing with children

A call written as `<@Card ... />` with following sibling content or nested tags
inside the same parent is invalid only when the parser treats the self-closing
tag as an open tag that later finds `</@Card>`. The exact diagnostic is:

```
<file>:<line>:<col>: Card: self-closing call must not contain children
```

### Slot fallback

When a component declares `{#slot}` and the call has no children, the slot
renders as empty string. No error is produced unless the component author marks
the slot as required in a future version.

### Nested component inside slot

A call may place another component call inside a slot:

```dreego
<@Card title="Outer">
    <@Icon name="star"/>
</@Card>
```

The inner component is rendered first and the resulting string is inserted into
the outer component's slot buffer.

### Named slot syntax

In a component definition, named slots use `{#slot <name>}{/slot}`.
In a call, named slot content uses `{#slot <name>}content{/slot}`.
The default slot is `{#slot}` and never has a closing tag.
Named slot names must match a declared slot exactly; an unknown named slot is
an error and produces:

```
<file>:<line>:<col>: Card: unknown slot "footer"
```

### Source location accuracy

Every diagnostic points to the line of the component call in the calling file,
not the component definition. Prop-level diagnostics point to the line where
`<@Name` appears.

### Sibling slot leakage

After rendering a component call with children, the slot buffer is reset so the
next sibling call does not receive leftover content. A test with two adjacent
`<@Card>` calls where only the first has children must show the second card with
an empty slot.

## Boundaries / Always Ask First Never

- Do not introduce positional prop binding.
- Do not make props required by default; a missing required prop is an error, a
  missing optional prop uses the default if present, otherwise zero value.
- Do not emit runtime diagnostics for cases that can be detected at generate
  time.
- Do not allow slot content to leak outside the component call that owns it.
- Do not add external dependencies to `core/`.
- Do not exceed 300 lines per new file.
- Do not keep German comments or messages in generated code or tests.

## Success Criteria

- `make test` passes.
- `_tests/go/component_call_*.go` covers all sections above.
- No existing component test regresses.
- Core remains free of external dependencies.
- Diagnostics match the message format documented above.

