---
version: patch
---

- Bug: apply `|raw` and other expression filters in attribute, URL, script, and style contexts instead of emitting invalid Go (`undefined: raw`)
- Bug: allow `{#if}` conditions on strings, numbers, and slices by routing them through a truthiness helper (empty string, zero, and empty collections are false), fixing the non-compiling `_docs/forms.md` example
