# Render result contract

Status: accepted for v0.2

## Context

Non-HTTP targets need to invoke generated rendering without constructing an
HTTP request or response. The first draft exposed `HTML`, `Head`, and `Assets`
before the generator could compose child metadata without loss. Exporting those
fields would create a contract that nested components and pages did not honor.

## Decision

`dreego.Result` contains only `HTML []byte` in v0.2. The bytes are the complete
renderer output and remain byte-identical to the SSR response body. Head markup
and scoped inline styles remain in their existing positions inside those bytes.

Generated components and pure GET pages both implement the same typed entry
contract through `dreego.Component`:

```go
type Component interface {
    Render(RenderContext) (Result, error)
}
```

Generated pure GET pages expose `Page<Name>() dreego.Component`. Pages with
request-dependent server sections remain SSR-only and do not expose a
render-neutral page constructor.

Rendering buffers output before returning. If rendering fails, callers receive
a zero `Result` and the error; partial output is not part of the public
contract. SSR may add streaming later as a host-specific behavior without
changing the target-neutral result.

`Result` intentionally does not expose separate head, asset, client-runtime, or
document-versus-fragment metadata in v0.2. Additive metadata fields require a
real producer and lossless composition tests before they become public.

## Consequences

- Components and pure pages render without `net/http` capabilities.
- SSR and non-HTTP rendering share the same generated renderer and exact bytes.
- Targets can distinguish page and component entry points by the generated API,
  not by different result shapes.
- A future metadata extension must define ordering, deduplication, nesting, and
  error behavior before adding fields to `Result`.
