---
area: product-architecture
phase: v0.1-blocker
priority: 0
---
# Resolve the remaining pre-v0.1 product decisions

## Goal
Record the remaining decisions that require the project owner's judgment. No
implementation should guess these answers. Resolve this item before the App,
routing, component, frontmatter, or plugin contracts are stabilized.

## How to use this item
For each question, choose one option or write a different decision. Update the
linked architecture decision and affected implementation todos when the choice
is made. The recommendations are starting points, not decisions already made.

## Decision 1: Optional route segments

Example source:

```text
routes/docs/[[lang]]/get.dreego
```

### Option A: Generate both paths

```text
GET /docs
GET /docs/de
```

Impact:
- Concise for pages with one optional segment.
- One source file owns two route patterns.
- Conflict detection and URL generation become more complex.
- Multiple optional segments create a rapidly growing number of paths.

### Option B: Remove optional segments

```text
routes/docs/get.dreego        -> GET /docs
routes/docs/[lang]/get.dreego -> GET /docs/{lang}
```

Impact:
- Every source maps to one explicit route pattern.
- Routing, conflict detection, testing, and generated registration stay simple.
- Shared rendering may need a component or Go helper.

Recommendation: Option B. Explicit routes fit Dreego's simplicity goal.

Decision: Pending.

## Decision 2: `g-action` and HTTP method ownership

Example source:

```text
routes/login/post.dreego
```

### Option A: The method file is authoritative

`post.dreego` registers only `POST /login`. A separate
`routes/login/get.dreego` renders the initial form. Validation failure may
render a response from the POST handler without creating a GET route.

Impact:
- Preserves the rule that directories define URLs and files define methods.
- GET and POST implementations remain small and explicit.
- Shared form markup should live in a component or helper.

### Option B: A `g-action` file registers GET and POST

Impact:
- A complete form flow can live in one file.
- The filename no longer describes all routes created by that file.
- Route ownership and conflict diagnostics become less predictable.

Recommendation: Option A.

### Successful action return

Choose whether `return nil` means:

- Option 1: success with an automatic `303` redirect to the GET route; concise,
  but implicit.
- Option 2: successful handlers must explicitly redirect or write a response;
  more verbose, but fully visible in Go code.

Recommendation: Option 2 because explicit behavior is easier to understand and
debug with a screen reader.

Decision: Pending.

## Decision 3: Component calls, props, nesting, and slots

Example:

```dreego
<@Card title="Hello" count={3} />
<@Card count={3} title="Hello" />
```

### Option A: Props are named and order-independent

Both calls are identical. Unknown, duplicate, and missing required props fail
during generation. Components may render nested components and scoped slots.

Impact:
- Matches the documented syntax and common developer expectations.
- Requires a component symbol table and lexical slot scoping in the generator.
- Produces much safer diagnostics than positional binding.

### Option B: Props are positional

Impact:
- The generator is simpler.
- Attribute names become misleading and reordering attributes can silently
  change behavior.
- Component composition remains much less useful.

Recommendation: Option A. Named props are central to the intended authoring
experience. Component nesting and slots should work predictably rather than be
advertised as partial features.

Decision: Pending.

## Decision 4: Frontmatter

Example:

```dreego
---
title: About
draft: false
---
<h1>About</h1>
```

### Option A: Frontmatter is part of the `.dreego` language

Impact:
- Useful for page metadata, documentation, and future SSG work.
- The generator must parse it, type its values, expose it deliberately, and
  report source-aware errors.
- Adds language surface that SSR applications may not need before v1.

### Option B: Remove it from the public language for now

Impact:
- Keeps v0.1 focused on SSR behavior already used by real applications.
- The parser helper may remain internal or be removed.
- Frontmatter can return when a real metadata consumer proves its contract.

Recommendation: Option B until an SSR use case needs it.

Decision: Pending.

## Decision 5: Meaning of the type-safety rule

The current absolute rule forbids string keys in Core, while sessions and
request-local data naturally use keys.

### Option A: Ban dynamic keyed data everywhere

Impact:
- Strongest compile-time guarantees.
- Session values and request-local extension data require generated or generic
  typed wrappers.
- Plugin interoperability becomes harder.

### Option B: Apply the rule to primary application data and generated APIs

String-keyed boundary data remains allowed where the domain is inherently
dynamic, such as HTTP headers, form fields, session stores, frontmatter, and
request-local extension state. These APIs must make type conversion and missing
values explicit.

Impact:
- Keeps strong types where Dreego controls the schema.
- Acknowledges unavoidable dynamic web boundaries honestly.
- Public claims must say "typed generated application code" rather than
  "no string keys anywhere".

Recommendation: Option B.

Decision: Pending.

## Decision 6: Provisional plugin contract for v0.1

The current all-in-one Plugin interface includes routes, middleware, assets,
startup, and shutdown even though parts are not wired into the server.

### Option A: Keep an App-bound all-in-one interface

Impact:
- Smallest conceptual change from the current API.
- Every plugin must implement capabilities it may not use.
- The unproven contract remains public during v0.1.

### Option B: Use ordinary registration functions for v0.1

Example:

```go
func Register(app *dreego.App) error
```

Assets and lifecycle are added only by plugins that need them through explicit
App methods or later capability interfaces.

Impact:
- Very small provisional contract.
- Fits Go conventions and explicit App ownership.
- Lifecycle semantics can be added after real plugins prove the need.

### Option C: Define small optional capability interfaces now

Impact:
- Supports lifecycle and assets before v0.1.
- Risks designing abstractions from hypothetical plugins again.

Recommendation: Option B. Build two or three real plugins before standardizing
capability interfaces.

Decision: Pending.

## Acceptance criteria
- Every decision above records the selected option and rationale.
- Normative architecture documents contain one current answer per decision.
- Implementation todos are updated to depend on the selected contracts.
- Public documentation stops advertising rejected or unimplemented behavior.
