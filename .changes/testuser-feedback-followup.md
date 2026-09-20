---
version: patch
---

- Breaking: remove the `dreego init` command; `dreego new` is the single scaffold entry point (it already creates the module, `init` did not)
- Bug: `dreego new` now always writes the module-qualified import (the removed `init` could emit the broken relative import `"./www"`)
- Feat: the scaffolded `main.go` declares the listening port as a `port` constant, with `DREEGO_PORT` as the container override; the `Dockerfile`/`docker-compose.yml` read that variable instead of hardcoding 8080
- Fix: the `web-app` starter no longer writes handwritten `notes_store.go` into the generated `routes` package; the notes store lives in the route's `<server>` section
- Fix: `dreego docs` works without a `go.mod` (resolves first-party docs from the module cache); `--list` no longer exits when no `go.mod` exists
- Docs: update `getting-started`, `cli`, `deployment`, and the template decision for the removed `init`, the port constant, and the `web-app` route store
- Docs: add `_todo/plugins/websocket-hub-broadcast.1.md` for the plugin's missing hub registration and README drift
