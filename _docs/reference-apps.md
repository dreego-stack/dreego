# Reference Applications

Four small, documented applications under `_tests/fixtures/` verify the
complete public CLI-to-HTTP workflow: `dreego generate` → `go build` → start →
HTTP request. Each app uses only public Dreego APIs — the same code a user
writes in a real project. The integration tests in
`_tests/go/reference_apps_test.go` copy a fixture into a temp dir, run the CLI,
build the binary, start it, and assert on HTTP responses.

## hello — Minimal usage

`_tests/fixtures/hello/`

Teaches the smallest possible Dreego app:

- `www/routes/get.dreego` — one route file per URL, one method per file
- `<head>` with `<title>` and meta tags
- `<go>` block with a local variable rendered via `{{ message }}`
- `www/routes/about/get.dreego` — nested directory route
- `www/routes/users/[id]/get.dreego` — dynamic segment with `c.Param("id")`
- `www/routes/404.dreego` — custom not-found page
- `main.go` — `dreego.New()` + `www.Register(app)` + `app.Listen(addr)`

Run it:

```bash
cd _tests/fixtures/hello
dreego generate
go run .
# GET /            → 200 "Hello from Dreego"
# GET /about       → 200 "About this app"
# GET /users/42    → 200 "User 42"
# GET /missing     → 404 custom page
```

## forms — Forms and sessions

`_tests/fixtures/forms/`

Teaches declarative form handling and session state:

- `g-action="AddEntry"` on a `<form>` — Dreego generates the POST handler from
  the `EntryForm` struct and the `AddEntry(c dreego.Context, form EntryForm) error`
  function
- `form:"name"` and `validate:"required"` struct tags
- `c.Old("name")` / `c.Errors("name")` re-render on validation failure
- hidden `csrf_token` field with `c.CSRFToken()` (CSRF is on by default)
- `c.Redirect("/entries", 303)` Post-Redirect-Get
- `{#if}` / `{#each}` template logic on the entries page
- `c.SessionVal` / `c.SetSessionVal` in `www/routes/counter/` — a plain
  POST route (no `g-action`) that increments a session counter
- `main.go` — `dreego.NewCookieStore(secret)` + `app.SetSessionStore(store)`

The guestbook stores entries in a package-level `entries` slice. This is demo
state only — it is not a production pattern (no persistence, no locking).

Run it:

```bash
cd _tests/fixtures/forms
dreego generate
go run .
# GET /            → 200 guestbook form
# POST /           → 303 redirect to /entries (with CSRF token)
# GET /entries     → 200 list with the posted entry
# GET /counter     → 200 "Count: 0"
# POST /counter    → 200, then GET /counter → "Count: 1"
```

## components — Components

`_tests/fixtures/components/`

Teaches the component system:

- `Component ProductCard (name string, price string, inStock bool)` — typed
  props, `{#if inStock}` branch, scoped `<style>` (emits `data-scope=`)
- `Component PageShell (title string)` — default slot `{#slot}` wrapping child
  content
- `<@ProductCard name={product.Name} .../>` — expression props from a Go struct
- `{#each products as product}` — loop over a slice
- `www/routes/products/[id]/get.dreego` — dynamic route reusing the same
  components

Run it:

```bash
cd _tests/fixtures/components
dreego generate
go run .
# GET /            → 200 shop grid with "In stock" and "Sold out" badges
# GET /products/1  → 200 "Dreego Mug"
# GET /products/2  → 200 "Dreego Tee"
```

## plugin — A plugin

`_tests/fixtures/plugin/`

Teaches the plugin model: an ordinary Go package that registers routes against
the owning App before the generated routes.

- `plugin/plugin.go` — `Register(app *dreego.App, options Options) error` using
  `app.Register(http.MethodGet, ...)` with static and `{id}` patterns
- `main.go` — `plugin.Register(app, plugin.Options{Prefix: "/plugin"})` before
  `www.Register(app)`
- `www/routes/get.dreego` — the app's own home page

Run it:

```bash
cd _tests/fixtures/plugin
dreego generate
go run .
# GET /                → 200 "Plugin demo"
# GET /plugin/hello    → 200 "Hello from the plugin"
# GET /plugin/hello/42 → 200 "Hello 42"
# GET /plugin/health   → 200 "plugin ok"
```

## How the tests work

`dreegotest.ServeFixture(t, name)` (in `dreegotest/fixture.go`):

1. Copies `_tests/fixtures/<name>` into a temp dir
2. Rewrites the `replace` directive in `go.mod` to point at the repo root
3. Runs `dreego generate` via the cached CLI binary
4. Runs `go build`
5. Starts the binary on a free port (the fixture reads `PORT` from the
   environment)
6. Returns a `dreegotest.Client` with a cookie jar for session/CSRF cookies

The tests run in CI as part of `go test ./_tests/go/...` (see
`_tests/test.sh`), so every PR verifies that the public CLI-to-HTTP workflow
still works end to end.

## See Also

- [Testing](https://github.com/dreego-stack/dreego/blob/main/_docs/testing.md)
- [Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md)
- [Forms](https://github.com/dreego-stack/dreego/blob/main/_docs/forms.md)
- [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md)
- [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md)
