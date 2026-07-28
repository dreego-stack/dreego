# dreego — Go Web Framework

SSR-First web framework for Go. File extension `.dreego`, transpiler approach, single binary.

```html
<!-- dreego/routes/get.dreego -->
<head><title>Dreego</title></head>

<go>
    message := "Hello from Dreego!"
</go>

<div>
    <h1>{message}</h1>
</div>
```

```bash
dreego generate && dreego run
```

## Features

- **Compile-Time Transpiler** — `.dreego` → Go code, zero runtime overhead
- **File-based Routing** — `dreego/routes/get.dreego` → `/`
- **5 Sections** — `<head>`, `<go>`, `<div>`, `<script>`, `<style>`
- **Template Logic** — `{#if}`, `{#each}`, `{var}` + `{#else}`, `$loop`, `{#verbatim}`, `{var|raw|upper}`
- **Component System** — `dreego/components/`, `<@Name>`, named slots, scoped CSS
- **Layout System** — `dreego/layouts/default.dreego` with `{#slot}` + `{#head}`
- **Static Assets** — `dreego/static/` → inline handlers with MIME detection
- **CSS Scoping** — `data-scope` via source hash
- **Middleware** — RequestLogging (JSONL), Redirects, Rewrites
- **Single Binary** — `go build` → deploy one file
- **CLI** — `dreego init`, `dreego generate [--force] [--check]`

## Quick Start

```bash
go install codeberg.org/dreego/dreego/cmd/dreego@latest
dreego init myapp
cd myapp
go mod init myapp
dreego generate
go run .
```

## Documentation

See [`_docs/`](_docs/) and [`.agents/`](.agents/index.md).

## License

MPL-2.0
