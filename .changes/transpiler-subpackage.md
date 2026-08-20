---
version: none
---

- Refactor: move the transpiler out of `core/` into `internal/transpiler/` (CLI and dreegotest are its only consumers; `core/` is now the runtime framework only, public import path unchanged)
