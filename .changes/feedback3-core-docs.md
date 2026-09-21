---
version: patch
---

- Bug: `c.Write` appends `; charset=utf-8` only when the caller-supplied content type has no charset, so a caller-provided charset is no longer duplicated
- Feat: `SafeURL` allows the `webcal` and `caldav` schemes for calendar subscriptions
- Docs: `dreego.config.json` is documented in the website root (`www/` by default), not the project root
- Docs: the CLI install path `github.com/dreego-stack/dreego/cmd/dreego@…` is documented and the module-root pitfall explained
