---
version: patch
---

- Bug: emit the components import in generated layouts so component calls in a layout compile
- Bug: parse header import directives in layout files instead of rejecting them as body text
- Bug: track layout source positions so component errors point at the layout file
