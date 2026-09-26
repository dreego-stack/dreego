# Recipe: Calling App Code from `<server>` Sections

A `<server>` section is ordinary Go. It can call functions from your own
application packages, use third-party libraries, and read request-scoped
services — not just the standard library.

This recipe covers three patterns:

1. importing and calling a function directly (v0.10.9 and later),
2. a hand-written **glue** file next to the route,
3. **service injection** through middleware context.

## Importing Go packages

A `.dreego` file declares Go imports with the `GOIMPORT` header directive:

```text
GOIMPORT { path1, alias "path2" }
```

`GOIMPORT` accepts any package that resolves in the application:

- standard-library packages (`os`, `encoding/json`, `sync`),
- packages in your own module (`myapp/internal/store`),
- dependencies listed in your application's `go.mod`
  (`github.com/example/otel`).

An entry may carry an alias: `alias "path"`. The base name is used when no
alias is given. Two imports with the same base name require aliases; a collision
without one fails at `dreego generate` and names both paths.

A package that is not in `go.mod` fails at `dreego generate` with a diagnostic
naming the path, the words `not in go.mod`, and the matching `go get` command.
Dreego never runs `go get` for you.

`GOIMPORT` applies to the generated package for routes, layouts, and components.
`strings`, `net/http`, and `fmt` are also detected automatically from the code.

> **Version note:** before v0.10.9, `GOIMPORT` accepted only a fixed
> standard-library allowlist. v0.10.9 removes the allowlist and adds the alias
> form, so application and third-party packages can be imported directly.
> `_docs/file-anatomy.md` and
> `_docs/decisions/explicit-dreefile-header-grammar.md` now document the
> free+alias form as well.

## Pattern 1: direct import and call

```dreego
GOIMPORT { myapp/internal/greeting }

<server>
message := greeting.For(c.Param("name"))
</server>

<body><p>{{ message }}</p></body>
```

`greeting` is the base name of the import path. For a package whose base name is
not a valid identifier or collides with another import, add an alias:

```text
GOIMPORT { appauth "myapp/internal/auth", oauthauth "example.com/auth" }
```

The alias is the name used inside `<server>` (`appauth.Register(...)`).

## Pattern 2: the glue pattern

A **glue file** is a hand-written `.go` file placed in the same directory as the
route. `dreego generate` only writes and manages `dree.go` (guarded by the
generated-file marker), so hand-written files in the route folder are never
touched.

Since v0.10.9 each route directory is its own Go package, named after the
directory (sanitized). The generated `dree.go` and the glue file therefore share
one package, and `<server>` can call the glue function without an import:

```go
// www/routes/register/glue.go
package register

import (
	"myapp/internal/auth"

	dreego "github.com/dreego-stack/dreego/core"
)

func CreateUser(c dreego.Context, email, password string) error {
	return auth.CreateUser(c, email, password)
}
```

```dreego
<!-- www/routes/register/+page.dreego -->
<server>
	type RegisterForm struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required,min=8"`
	}

	func Register(c dreego.Context, form RegisterForm) error {
		if err := CreateUser(c, form.Email, form.Password); err != nil {
			return err
		}
		c.Flash("notice", "account created")
		return c.Redirect("/", 303)
	}
</server>

<body>
	<form g-action="Register" method="post">
		{{ c.CSRFInput()|raw }}
		<label for="email">Email</label>
		<input id="email" name="email" type="email" value="{{ c.Old("email") }}">
		<label for="password">Password</label>
		<input id="password" name="password" type="password">
		<button type="submit">Create account</button>
	</form>
</body>
```

In releases before v0.10.9 all route files compiled into one `routes` package,
so the glue file was `www/routes/glue.go`. With per-folder packages the glue
file lives in the route folder itself.

Use glue when you need behavior that is awkward to inline: wrapping a function
with a different signature, adapting a typed error, or keeping helper code out
of the template. When a plain import is enough, prefer Pattern 1.

## Pattern 3: service injection via middleware context

Request-scoped dependencies (a database handle, a configured client, a logger)
should not be package-level globals. Inject them through the request context in
an `app.Use` middleware, then read them in `<server>`.

```go
// internal/services/services.go
package services

import "context"

type dbKey struct{}

func WithDB(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

func DBFrom(ctx context.Context) (*DB, bool) {
	db, ok := ctx.Value(dbKey{}).(*DB)
	return db, ok
}
```

```go
func InjectDB(db *DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(services.WithDB(r.Context(), db)))
		})
	}
}
```

```go
// main.go
if err := app.Use(services.InjectDB(db)); err != nil {
	log.Fatal(err)
}
```

```dreego
GOIMPORT { myapp/internal/services }

<server>
	db, ok := services.DBFrom(c)
	if !ok {
		return "", fmt.Errorf("database service missing")
	}
	items, err := db.List(c, c.Param("list"))
	if err != nil {
		return "", err
	}
</server>

<body>
	<ul>
		{#each items as item}<li>{{ item.Title }}</li>{/each}
	</ul>
</body>
```

`*dreego.SSRContext` and `dreego.RenderContext` both embed Go's
`context.Context`, so `c` can be passed directly to any function that accepts a
`context.Context`. This is the supported way to hand request-scoped state to
application code without exporting a global.

Keep the injected value narrow: inject a small interface or a concrete service,
and return a typed error from the lookup so a missing service fails loudly
instead of panicking.

## Declarations and scope

The leading declaration block of a `<server>` section (`type`, `func`, `var`,
`const`, plus any top-level `func`) is emitted at package level; statements stay
inside the render function. With one package per route folder, a package-level
declaration is visible to the other route files **in that folder** and must be
unique within it.

Keep request-local values after a statement so they stay inside the render
function. Keep shared mutable state package-level only when it truly belongs to
the application, and synchronize it with a mutex.

## See Also

- [Server Section](https://github.com/dreego-stack/dreego/blob/main/_docs/server-section.md) — section scope and methods
- [File Anatomy](https://github.com/dreego-stack/dreego/blob/main/_docs/file-anatomy.md) — header directives
- [Context Methods](context-methods.md) — which methods exist on which context
- [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md) — registration from `main.go`
