---
version: patch
---

- Breaking: replace root `<go>`, `<div>`, and `<script>` sections with semantic `<server>`, `<body>`, and `<client>` sections before v0.1
- Feat: preserve built-in section languages in the parsed model and reject unsupported section/language pairs with source-located processor guidance
- Docs: add the semantic-sections migration guide and update scaffolds, fixtures, plans, and public documentation
