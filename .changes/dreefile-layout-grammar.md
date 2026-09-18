---
version: minor
---

- Feat: explicit Dreefile header grammar with `DREEFILE component|layout|page`, `LAYOUT`, `COMPONENT "path" IMPORT { Names, Name as Alias }` and `GOIMPORT { paths }`; `LAYOUT` and `GOIMPORT` are parsed and reserved, their codegen consumers follow in later slices
- Breaking: the legacy `Component X (props)`, bare `import "path"` and `from "path" import { ... }` headers are rejected at `dreego generate` with a `file:line:col` diagnostic naming the replacement
- Breaking: a component's name now comes from its filename (`Card.dreego` -> `<@Card>`), so a component filename must be an exported Go identifier
- Bug: the head merge deduplicates `<title>` and meta description for body-level layouts on both sides of `{#head}`, so a page carries exactly one of each
- Docs: the Dreefile grammar, the migration, and the header ADR
