---
version: patch
---

- Bug: `dreego fmt` is now semantics-preserving for body-level layouts; a document-level `<head>` nested in `<body>` is no longer hoisted and trailing `</html>`/`</body>` are no longer dropped
- Bug: `dreego fmt` reorders only whitespace-separated root sections and preserves every token and the surrounding text
- Bug: `dreego fmt` no longer rewrites string-literal contents; whitespace and `|` normalization stays outside `"…"`, `'…'`, and `` `…` `` values such as `{{ "a  b" }}` or `{#if x == "a | b"}`
- Bug: `dreego fmt` preserves `<server>`, `<client>`, and `<style>` sections byte for byte, including Go raw-string contents and alignment spacing
- Test: round-trip property matrix over the scaffold layouts and body-level layout variants that asserts fmt never changes the lexed section structure or skeleton tags
- Test: literal guards and an exact code-section comparison that fail against the previous whitespace collapsing
- Test: `dreego fmt --check` never writes, including on a body-level layout
