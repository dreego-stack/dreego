---
version: patch
---

- Feat: named component props are bound by name, order-independent, with source-aware errors for missing, unknown, and duplicate props
- Feat: simple literal expression props are type-checked against declared prop types at generate time
- Feat: self-closing component calls reject children; empty slots render as empty string
- Feat: named slots are validated against declared slots; slot content is scoped per invocation and cannot leak to sibling calls
- Bug: nested slot declarations inside slot content are rejected at generate time
