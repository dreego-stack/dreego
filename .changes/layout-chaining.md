---
version: none
---

- Feat: a layout can declare `LAYOUT "path"` and is wrapped by that layout (explicit layout chaining)
- Feat: missing `LAYOUT` targets and layout cycles fail at `dreego generate` with a `file:line:col` diagnostic
- Docs: ADR records the LAYOUT codegen consumer
