# Context Methods

Dreego exposes three context types. They are not interchangeable: each one is
the boundary for one kind of code, and the method set is deliberately different.

| Type | Received by | Embeds |
|------|-------------|--------|
| `dreego.Context` | `<server>` sections and `g-action` handlers (interface) | `context.Context` |
| `*dreego.SSRContext` | the concrete value named `c` in route code (implements `dreego.Context`) | `context.Context` |
| `dreego.RenderContext` | generated and hand-written components, named `ctx` | `context.Context` |

Route `<server>` code and `g-action` handlers see `c`; components see `ctx`.
`dreego.Context` is the interface a `<server>` section is written against;
`*dreego.SSRContext` is the request-bound implementation. Since v0.10.9 the
interface carries the full request surface, so `c.Set(...)`,
`c.FormValue(...)`, `c.Query(...)`, and `c.DestroySession()` type-check directly
against `dreego.Context`.

## Context and SSRContext

`dreego.Context` is the interface; `*dreego.SSRContext` is the concrete
request-bound value. Every method below is on both.

| Method | Returns | Description |
|--------|---------|-------------|
| `c.Param("id")` | `string` | URL path parameter from a `[id]` segment |
| `c.Query("ref")` | `string` | URL query parameter `?ref=x` |
| `c.FormValue("name")` | `string` | POST form value; `""` on parse failure |
| `c.FormError()` | `error` | Non-nil when the last `FormValue` failed to parse the body |
| `c.Data("key")` | `any` | Request-local extension data |
| `c.Set("key", val)` | — | Store request-local data for nested calls |
| `c.Get("key")` | `string` | Retrieve string request-local data |
| `c.Delete("key")` | — | Remove request-local data |
| `c.Errors("field")` | `string` | Validation error message for a field |
| `c.Old("field")` | `string` | Previously submitted value after a validation failure |
| `c.SessionVal("key")` | `string` | Session value; `""` when no store is configured |
| `c.SetSessionVal("k","v")` | — | Write a session value with secure defaults |
| `c.DelSessionVal("key")` | — | Delete one session key |
| `c.DestroySession()` | — | Destroy the complete session (logout) |
| `c.CSRFToken()` | `string` | Current CSRF token from the session |
| `c.CSRFInput()` | `string` | Rendered hidden `csrf_token` input field |
| `c.Flash("key", "msg")` | — | Store a one-shot message in the session |
| `c.FlashGet("key")` | `string` | Read and consume a flash message |
| `c.FlashPeek("key")` | `string` | Read a flash message without consuming it |
| `c.RequestID()` | `string` | Request ID from `X-Request-ID` or the generated ID |
| `c.SessionError()` | `error` | Non-nil when the last session operation failed |
| `c.Redirect(url, code)` | `error` | Send a redirect and return `dreego.ErrRedirect` |

`*dreego.SSRContext` additionally exposes the raw request and writer plus the
response helpers used by typed `<server type="json">` and `<server type="xml">`
blocks. These are **not** part of the `dreego.Context` interface:

| Method / field | Description |
|--------|-------------|
| `c.R` | Raw `*http.Request` (use sparingly) |
| `c.W` | Raw `http.ResponseWriter` (use sparingly) |
| `c.JSON(status, data)` | Write a JSON response |
| `c.XML(status, data)` | Write an XML response |
| `c.Write(status, contentType, body)` | Write a raw response body |
| `c.Bind(target)` | Decode a JSON request body into `target` (1 MiB cap) |
| `c.Wants(mime)` | True when the `Accept` header matches `mime` |

## RenderContext

Components receive `dreego.RenderContext`, a smaller boundary. It has no
request, no response writer, no sessions, and no redirects. Convert
request-derived values to typed props in the route and pass them to the
component.

| Method | Returns | Description |
|--------|---------|-------------|
| `ctx.Data("key")` | `any` | Render-local data |
| `ctx.Set("key", val)` | — | Store render-local data |
| `ctx.Get("key")` | `string` | Retrieve string render-local data |
| `ctx.Delete("key")` | — | Remove render-local data |
| `ctx.Errors("field")` | `string` | Validation error message for a field |
| `ctx.Old("field")` | `string` | Previously submitted value |

Both context types also satisfy Go's `context.Context` (`Deadline`, `Done`,
`Err`, `Value`), so either value can be passed to functions that accept a
`context.Context`. This is how request-scoped services reach application code
(see [Calling App Code](app-code.md)).

## Choosing a Method

- Request data (`Param`, `Query`, `FormValue`) exists only on the route context.
- Session and flash state exists only on the route context.
- Response operations exist only on the route context.
- `Data`/`Set`/`Get`/`Delete` exist on both, but each keeps its own scope: a
  route's `Set` is not automatically a component's `Data`. Pass values to
  components as typed props.
- `Errors`/`Old` exist on both so a component that renders a form field can read
  the same validation state the route set.

## See Also

- [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md) — server, session, and config APIs
- [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md) — the component render boundary
- [Recipe: App Code](app-code.md) — service injection through the context
- [Forms](https://github.com/dreego-stack/dreego/blob/main/_docs/forms.md) — `Errors`, `Old`, and CSRF
