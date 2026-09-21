---
version: patch
---

- Bug: `dreego fmt` is now semantics-preserving for body-level layouts; a document-level `<head>` nested in `<body>` is no longer hoisted and trailing `</html>`/`</body>` are no longer dropped
- Bug: `dreego fmt` reorders only whitespace-separated root sections and preserves every token and the surrounding text
- Test: round-trip property matrix over the scaffold layouts and body-level layout variants that asserts fmt never changes the lexed section structure or skeleton tags
- Test: `dreego fmt --check` never writes, including on a body-level layout
