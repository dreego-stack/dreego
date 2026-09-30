---
version: minor
---

- Breaking: a website root is now one directory with `dreego.config.json` that hosts multiple apps; each app is a subdirectory with its own `routes/` tree and its own Go package
- Breaking: `dreego generate` emits `<app>/dree.go` with `var App dreego.Registrar` instead of a root-level `Register(app)`; generated route packages keep `func Register(app *dreego.App) error`
- Breaking: a root-level `routes/` tree is rejected with a migration hint; create an app subdirectory (for example `dreego/www/routes/`) instead
- Feat: `dreego.New(registrars ...Registrar)` builds an App directly from generated app registrars; it panics on a registration error, which `dreego generate` already prevents
- Feat: shared `layouts/` and `components/` at the website root are generated once and used by every app; an app may override them with its own `layouts/`/`components/`
- Feat: an app may ship `<app>/dreego.config.json`; it overrides the root defaults field by field (no deep merge)
- Docs: `dreego.config.json` is the website root marker; hosts and ports live in `main.go`, so one binary can serve several apps on several ports
