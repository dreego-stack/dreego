---
version: patch
---

- Breaking: generated output is now per-directory — `dreego/gen/` is gone; every directory with `.dreego` sources gets its own `dree.go` (`www/routes/dree.go`, `www/components/dree.go`, `www/layouts/dree.go`, plus `www/dree.go` with `Register(app)`)
- Breaking: the website root is any directory with `dreego.config.json` (renamed from `config.json`); website directories are freely named, multiple websites per repo are supported
- Feat: user code now imports only the website package (`"myapp/www"`) and calls `www.Register(app)` — no generated `gen` import
