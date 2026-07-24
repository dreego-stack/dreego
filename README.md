# dreego — Go-Webframework

Ein SSR-First Webframework fur Go. Dateiendung `.dreego`, Transpiler-Ansatz, Single Binary.

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

- **Compile-Time Transpiler** — `.dreego` → Go-Code, null Runtime-Overhead
- **File-based Routing** — `dreego/routes/get.dreego` → `/`
- **5 Sektionen** — `<head>`, `<go>`, `<div>`, `<script>`, `<style>`
- **Template-Logik** — `{#if}`, `{#each}`, `{var}`
- **Layout-System** — `dreego/layouts/default.dreego` mit `{#slot}`
- **CSS-Scoping** — `data-scope` via Source-Hash
- **Middleware** — RequestLogging (JSONL), Redirects, Rewrites
- **Single Binary** — `dreego build` → eine Datei deployen
- **CLI** — `dreego generate`, `dreego build`, `dreego run -d -t 60`

## Quick Start

```bash
go install codeberg.org/dreego/dreego/cmd/dreego@v0.0.1
mkdir myapp && cd myapp
go mod init myapp
mkdir -p dreego/routes
# ... schreibe dreego/routes/get.dreego ...
dreego generate
dreego run
```

## Dokumentation

Siehe [`_docs/`](_docs/index.md) und [`.agents/`](.agents/_index.md).

## License

## License

MPL-2.0
