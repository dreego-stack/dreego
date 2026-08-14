---
version: patch
---

- Bug: lexer treats < in Go blocks, > in quoted attributes, and < > in script/style as tags
- Bug: || in template expressions is misparsed as filter pipeline
- Bug: unknown template filters are silently ignored instead of failing
