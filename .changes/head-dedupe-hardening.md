---
version: none
---

- Hardening: head dedupe detects `<title>` and meta description case-insensitively
- Hardening: a non-literal layout head prefix emits a generate diagnostic instead of silently skipping dedupe
- Chore: use the `ir.HeadPlaceholder` constant instead of hardcoded `{#head}`
