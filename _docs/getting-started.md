# Getting Started

## Installation

```bash
go install codeberg.org/dreego/dreego/cmd/dreego@v0.0.1
```

## Neues Projekt

```bash
mkdir myapp && cd myapp
go mod init myapp
mkdir -p dreego/routes
```

## Erste Seite

```html
<!-- dreego/routes/get.dreego -->
<head>
    <title>Meine App</title>
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
dreego generate    # transpiliert .dreego → Go-Code
dreego build       # generate + go build
dreego run         # build + server starten
dreego run -d      # mit Debug-Logging (JSONL)
dreego run -t 30   # auto-stop nach 30 Sekunden
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
