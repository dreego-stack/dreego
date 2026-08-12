# Getting Started

## Installation

```bash
go install github.com/dreego-stack/dreego/cli/dreego@latest
```

## New Project

```bash
mkdir myapp && cd myapp
go mod init myapp
mkdir -p dreego/routes
```

## First Page

```html
<!-- dreego/routes/get.dreego -->
<head>
    <title>My App</title>
</head>

<go>
    message := "Hello from Dreego!"
</go>

<div>
    <h1>{message}</h1>
</div>
```

## Build & Run

```bash
dreego generate    # transpiles .dreego → Go code
dreego build       # generate + go build
dreego run         # build + start server
dreego run -d      # with debug logging (JSONL)
```

## main.go

```go
package main

import (
    _ "myapp/dreego/gen"
    dreego "github.com/dreego-stack/dreego/core"
)

func main() {
    dreego.Listen(":8080")
}
```

## Adding a Layout

Create `dreego/layouts/default.dreego` — wraps all pages:

```html
<head><title>My App</title></head>

<div>
    <nav><a href="/">Home</a> | <a href="/about">About</a></nav>
    <main>{#slot}</main>
</div>

<style>
    nav { padding: 1rem; background: #1e293b; }
    nav a { color: #e2e8f0; margin-right: 1rem; }
</style>
```

## Creating a Component

Create `dreego/components/Card.dreego`:

```
Component Card (title string)

<div>
    <article class="card">
        <h2>{title}</h2>
        <div>{#slot}</div>
    </article>
</div>

<style>
.card { border: 1px solid #e2e8f0; padding: 1rem; border-radius: 8px; }
</style>
```

Use it in any route or layout:

```html
<@Card title="Welcome">
    <p>This is the card body.</p>
</@Card>
```

Components are auto-discovered — no import needed.

## Dynamic Routes

Create `dreego/routes/users/[id]/get.dreego`:

```html
<head><title>User {c.Param("id")}</title></head>

<go>
    userID := c.Param("id")
</go>

<div>
    <h1>User: {userID}</h1>
</div>
```

Visiting `/users/42` shows "User: 42".

## See Also

- [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md) — full component docs
- [Routing](https://github.com/dreego-stack/dreego/blob/main/_docs/routing.md) — dynamic segments, groups, methods
- [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md) — SSRContext, sessions, config
- [CLI Reference](https://github.com/dreego-stack/dreego/blob/main/_docs/cli.md)
- [Docs Index](https://github.com/dreego-stack/dreego/blob/main/_docs/index.md)
