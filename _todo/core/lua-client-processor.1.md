---
area: transpiler
phase: multi-language-dreego
---
# Lua client processor

Implement first-party `<client lang="lua">` after TypeScript validates the
shared JavaScript output stage. It compiles Lua source to browser JavaScript;
Lua-to-Go, an embedded interpreter, and a runtime business-logic store are not
part of the work.
