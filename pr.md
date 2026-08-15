---
version: patch
---

- Feat: named component props are bound by name, order-independent, with source-aware errors for missing, unknown, and duplicate props
- Feat: simple literal expression props are type-checked against declared prop types at generate time
- Feat: self-closing component calls render empty slots while following sibling content remains valid
- Feat: named slots are validated against declared slots; slot content is isolated across sibling and nested calls
- Bug: nested slot declarations inside slot content are rejected at generate time
- Bug: nested component calls generate valid builder code and preserve outer slot values
- Bug: component prop binding is independent of component filename order
