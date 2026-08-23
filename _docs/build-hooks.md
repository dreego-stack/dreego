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
6. Every plugin build step requires explicit approval before it runs (see
   [Security: Build-Step Approval](#security-build-step-approval)).

## Security: Build-Step Approval

Plugin build steps run shell commands in your project directory. To prevent
a malicious plugin from running arbitrary commands unnoticed, `dreego build`
prompts for approval the first time it encounters a plugin build step.

### Interactive prompt

When `dreego build` finds a `pre-build` step that has not been approved yet and
stdin is a terminal, it prints:

```
plugin "github.com/dreego-stack/plugin-tailwind" wants to run "npx tailwindcss -i input.css -o www/static/tailwind.css --minify" in this repo. Approve? y/N:
```

Type `y` (or `yes`, case-insensitive) to approve and run the command. Type `N`,
any other input, or press Enter to reject — the build then fails with
`plugin ... build step not approved by user`.

The prompt shows the plugin module path and the **exact** command string, so
you can verify what will run before approving.

### Saved approvals

Approved commands are recorded in `dreego-build.json` at the project root:

```json
{
  "approvedHooks": {
    "github.com/dreego-stack/plugin-tailwind:npx tailwindcss -i input.css -o www/static/tailwind.css --minify": true
  }
}
```

The approval key is `pluginModulePath:commandString`. Subsequent builds reuse
the saved approval and skip the prompt. If a plugin changes its command, the
key changes and `dreego build` prompts again.

`dreego-build.json` can be committed to share approvals with your team. This
ensures everyone runs the same approved commands without re-prompting.

### Non-interactive environments (CI)

In CI, stdin is not a terminal. Without a prior approval and without `--yes`,
the build fails with a helpful message:

```
plugin github.com/dreego-stack/plugin-tailwind build step not approved: npx tailwindcss ...
Run 'dreego build' interactively to approve, or use 'dreego build --yes', or pre-approve in dreego-build.json
```

Use one of these options for CI:

- `dreego build --yes` — auto-approve all plugin build steps for this run
  (approvals are saved to `dreego-build.json`).
- Commit a `dreego-build.json` with the approved keys so CI reuses them.

### `--yes` flag

`dreego build --yes` auto-approves every plugin build step without prompting.
It is intended for CI or other trusted, non-interactive environments. The
approvals are written to `dreego-build.json` so subsequent builds (interactive
or not) reuse them.

## Example: Tailwind

Add the plugin to your project:

```sh
go get github.com/dreego-stack/plugin-tailwind
```

The plugin's `dreego-plugin.json` runs `npx tailwindcss` during `dreego build`
and writes the compiled CSS to `www/static/tailwind.css`. The Go binary embeds
this file via the Dreego static-asset system. No Node.js needed at runtime.