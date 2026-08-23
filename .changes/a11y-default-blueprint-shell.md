---
version: patch
---

- Feat: add accessible layout shell to default (init) blueprint with skip link, main landmark, and lang attribute
- Feat: add aria-label and skip link to PageShell component fixture
- Docs: add Accessibility section to components.md covering nav labels, landmarks, skip links, and decorative elements
- Test: add a11y_runtime_test.go asserting rendered HTML retains accessibility attributes
- Test: strengthen cli_accessibility_test to verify default blueprint skip-link and main landmark