---
version: patch
---

- Feat: ship project templates from `cmd/dreego/internal/templates/` with a shared `_common` layer, per-template `template.json` metadata, and the `web-minimal` template
- Feat: add `-t`/`--template <name>` selection and `-l`/`--list` to `dreego new` and `dreego init`, both defaulting to `web-minimal`
- Feat: add `dreego task [args...]` to forward to the external `task` binary (runs `task --list` with no arguments and forwards the exit code)
- Breaking: remove `cmd/dreego/blueprints/` and the `default`/`landing` blueprint names; scaffold paths move to `cmd/dreego/internal/templates/` with `web-minimal` as the replacement
- Docs: document the template overlay, metadata format, substitution rules, and CLI surface in `_docs/decisions/template-scaffolds.md`
