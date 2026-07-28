# Getting Started

## Installation

```bash
go install codeberg.org/dreego/dreego/cmd/dreego@v0.0.1
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
dreego run -t 30   # auto-stop after 30 seconds
```

## main.go

```go
package main

import (
    _ "myapp/dreego/gen"
    core "codeberg.org/dreego/dreego/dreego-core"
)

func main() {
    core.Listen(":8080")
}
```

## See Also

- [CLI Reference](/_docs/cli.md) — all CLI commands
- [Routing](/_docs/routing.md) — file-based routing
- [Docs Index](/_docs/index.md)
