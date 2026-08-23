# Build Hooks

Dreego plugins can run build steps during `dreego build`, `dreego run`, and
`dreego dev`. The steps run after `.dreego` → Go codegen and before the Go
compiler — so plugins can generate or transform assets (CSS, JS, images) that
the Go binary then embeds.

## How It Works

1. The CLI reads `go.mod` and finds all `require` lines matching
   `github.com/dreego-stack/*`.
2. For each plugin, the CLI resolves its on-disk directory (vendor/ or module
   cache).
3. If the plugin root contains a `dreego-plugin.json` file, the CLI reads its
   `build.steps`.
4. Steps with `"when": "pre-build"` run in the project root directory, in
   alphabetical order by plugin name, after `dreego generate` and before
   `go build`.

## dreego-plugin.json

Place this file at the root of the plugin module:

```json
{
  "build": {
    "steps": [
      {
        "cmd": "npx tailwindcss -i input.css -o www/static/tailwind.css --minify",
        "when": "pre-build"
      }
    ]
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `build.steps` | array | Build steps to run |
| `steps[].cmd` | string | Shell command (run via `sh -c`) |
| `steps[].when` | string | When to run. Only `"pre-build"` is supported. |

The `cmd` runs in the **project root** (where `go.mod` lives), not in the
plugin directory. This lets plugins reference project-relative paths like
`www/static/tailwind.css`.

## When Steps Run

- `dreego build` — yes (after generate, before go build)
- `dreego run` — yes (calls build)
- `dreego dev` — yes (calls build on each change)
- `dreego generate` — no (generate only transpiles, no build)

## Rules

1. A plugin without `dreego-plugin.json` is silently skipped.
2. A failing build step stops the build with an error.
3. Plugins run in alphabetical order by module path.
4. Only `github.com/dreego-stack/*` modules are discovered. Other dependencies
   are ignored.
5. The command runs with the project root as the working directory.

## Example: Tailwind

Add the plugin to your project:

```sh
go get github.com/dreego-stack/plugin-tailwind
```

The plugin's `dreego-plugin.json` runs `npx tailwindcss` during `dreego build`
and writes the compiled CSS to `www/static/tailwind.css`. The Go binary embeds
this file via the Dreego static-asset system. No Node.js needed at runtime.