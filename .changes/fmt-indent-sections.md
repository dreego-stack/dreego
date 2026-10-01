---
version: patch
---

- Feat: `dreego fmt` now indents markup by element nesting depth (four spaces per level) and code sections one level under their root tag
- Feat: `{#if}`/`{#each}` blocks indent their children; closing tags and block terminators align with what opened them
- Fix: keep `<pre>` and `{#verbatim}` content verbatim, leave `lang="md"` bodies unindented, and keep single-line sections inline while formatting
- Feat: `dreego docs skill` prints the bundled agent skill (`/_docs/skill.md`) so coding agents read the docs instead of the framework source
- Docs: document `dreego fmt` flags and the agent skill command in the CLI reference
