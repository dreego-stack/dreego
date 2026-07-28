# Blockwebchain

Status via `python _todo/process.py`. Chain 1–25 done. Next code: **26**.

## v0.0.9 (current)

- **each-loop.1** — `{#each}` mit `$loop`-Variable ✅
- **verbatim.1** — `{#verbatim}` Block ✅
- **tag-prefix-fix.1** — scanTag: head vs header Fix ✅
- **template-filters.1** — `{var|raw}`, `{var|upper}` ✅
- **if-else.1** — `{#if}...{#else}...{/if}` ✅
- **each-else.1** — `{#each}...{#each else}...{/each}` ✅

## AVAILABLE NEXT (v0.0.10+)

- **static-assets.1** — Static Assets (static/ → embed.FS)
- **form-actions.1** — Form Actions (g-action / g-submit)
- **health-checks.1** — /health + /ready Endpoints
- **security-headers.1** — Security-Header
- **compression.1** — Gzip/Brotli Middleware

## Blocker

- **golden-tests.1** — Golden-File Tests (needs: dreegotest.1)
- **plugin-interface.1** — Plugin-Interface (blockiert alle Addons)

Siehe `_todo/index.md` für die vollständige Chain und den Dependency-Graph.
