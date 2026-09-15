---
type: Decision
title: Project templates with a shared overlay
description: cmd/dreego/internal/templates/ holds one _common layer plus per-template directories with template.json metadata
tags: [cli, scaffolding, templates]
timestamp: 2026-09-15T00:00:00Z
---
# Project templates with a shared overlay

**Date:** 2026-09-15
**Status:** Accepted

## Context

`dreego new` and `dreego init` copied one of two fixed blueprints from
`cmd/dreego/blueprints/{default,landing}`. Both blueprints carried their own
`main.go.tmpl`, `Taskfile.yml`, `.gitignore`, and `www/dreego.config.json`, so
the files were duplicated and drifted independently. The blueprint name was
also hardcoded in the command implementations, so an additional scaffold would
have meant another full copy of every common file.

The layout is borrowed from Wails v3 (`internal/templates` plus a `_common`
overlay and per-template metadata), adapted to Dreego's module boundaries.

## Decision

Scaffolding moves into `cmd/dreego/internal/templates/`:

```
cmd/dreego/internal/templates/
├── templates.go              ← embedded filesystem, listing, metadata, install
├── templates_test.go
├── _common/
│   ├── main.go.tmpl          ← SSR entrypoint
│   ├── Taskfile.yml
│   ├── .gitignore.tmpl
│   └── www/dreego.config.json
└── web-minimal/
    ├── template.json
    └── www/
        ├── layouts/default.dreego
        └── routes/+page.dreego
```

`cmd/dreego/blueprints/` is deleted.

### `_common` overlay semantics

Installation copies `_common` first and then overlays the selected template on
top of the same destination tree. A relative path that exists in both layers is
overwritten by the template's file. A template does not have to repeat any
`_common` file it accepts unchanged. This is the extension point for a future
template that needs its own entrypoint, Taskfile, or configuration: it ships
only the files that differ.

### `template.json` instead of YAML

Each template carries its metadata in a `template.json` file with the fields
`name`, `title`, `description`, `type`, `adapter`, and `extraRequires`. The
module that ships the templates is CI-gated to the standard library and the
Go-project modules under `golang.org/x/` only. `gopkg.in/yaml.v3` would violate
that boundary, while `encoding/json` is already available. JSON is used only
for build-time metadata; it is not part of the generated application.

For `web-minimal` the fields are `name` `web-minimal`, `type` `web`,
`adapter` `ssr`, and an empty `extraRequires` list. `adapter` and
`extraRequires` are the declared extension point for templates that need an
additional Dreego module.

### `.tmpl` strips, `§$name$§` substitutes

The `.tmpl` suffix means only that the suffix is removed when the file is
written. The file body is never run through `text/template`. The literal token
`§$name$§` is the single placeholder and is replaced with the project or module
name via a plain string replacement.

`text/template` is deliberately not used because scaffolded files legitimately
contain `{{ }}` Go template syntax. A template engine would try to evaluate
those expressions and corrupt the generated application.

### Shipped templates

| Template | Purpose |
|----------|---------|
| `web-minimal` | Smallest SSR application: SSR entrypoint, config, one layout, one route, and one style block. |

Installation skips `template.json`, `.DS_Store`, and generated files
(`dree.go`, `*_dreego.go`, `handle_*.go`).

### CLI surface

- `dreego init <path> [-t <template>]` — default template `web-minimal`.
- `dreego new <name> [-t <template>]` — default template `web-minimal`.
- `-l` / `--list` — lists the available templates and exits without
  scaffolding.
- `dreego task [args...]` — forwards to the external `task` binary and fails
  with a clear message when that binary is absent.

## Consequences

- The template filesystem is embedded with an explicit list:
  `//go:embed all:_common all:web-minimal`. Adding a template therefore
  requires adding its directory to that list in `templates.go`; a new directory
  on disk is not shipped automatically.
- A registry test guards the list: `TestListSorted` pins the returned names and
  `TestEveryDiskTemplateRegistered` fails when a directory with a
  `template.json` is not registered by `List()`. `TestMetaOfValid` checks the
  `web-minimal` metadata.
- A new template is one directory plus its `template.json`; common files are
  never duplicated.
- Removing `cmd/dreego/blueprints/` is a breaking change for anything that
  referenced the old scaffold paths or the `default`/`landing` names.
- No third-party YAML dependency is added to a module that is standard-library
  and `golang.org/x/` only.
- `dreego task` depends on the external `task` binary being on `PATH`; Dreego
  does not embed a task runner. With no arguments it runs `task --list`.
- Template substitution stays intentionally narrow. A file that needs
  conditional or repeated content must ship as a concrete file instead.

## Alternatives considered

### YAML metadata like Wails v3

Rejected because it adds a third-party dependency to a gated module for a small
amount of build-time metadata.

### Generate files with `text/template`

Rejected because generated applications contain `{{ }}` expressions that the
template engine would evaluate.

### Ship each template standalone without `_common`

Rejected because the shared entrypoint, Taskfile, config, and ignore rules
would be duplicated in every template.

## Not decided or built yet

The following are explicitly open and not part of this change:

- No `doctor` command for checking a scaffolded project.
- No interactive template picker.
- No remote or git-hosted template sources.
- No template versioning or migration of already-scaffolded projects.
- No post-generation hooks.
- No embedded task runner; `dreego task` only forwards.
- Additional templates, such as a Wails desktop template, are deliberately
  deferred. The `adapter` and `extraRequires` fields of `template.json` are the
  extension point for a template that depends on another Dreego module.
