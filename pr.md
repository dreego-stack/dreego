---
version: none
---

- Fix: release-prep.py also updates embedded CHANGELOG copy
  (DocsSync test failed because prepare job changed CHANGELOG.md
  but not cmd/dreego/embedded/CHANGELOG.md)
