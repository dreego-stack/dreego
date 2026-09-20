# Wails desktop project template for `dreego new`

- Area: cli/templates
- Phase: Wails v3
- Goal: ship a Wails desktop project template for `dreego new`.
- Gap: `_docs/decisions/template-scaffolds.md` currently states that further templates, such as a Wails desktop template, are deliberately deferred.
- Acceptance: the template reuses the embedded template system in `cmd/dreego/internal/templates/` with its `_common` overlay and `template.json` metadata, declares the additional Dreego module through the `adapter`/`extraRequires` extension point, and produces a project that uses `adapter/wails` without a `package.json` or a Dreego HTTP listener.
- Depends on: the planned breaking-change release; it changes the templates and the generated application surface, so the template must wait until that pull request has landed.
