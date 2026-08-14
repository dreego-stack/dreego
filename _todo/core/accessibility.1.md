---
area: accessibility
phase: v0.1-blocker
---
# Establish accessibility as a framework quality gate

## Goal
Make accessibility a documented and testable Dreego design principle across the CLI, diagnostics, documentation, generated projects, and official components without claiming that every user application is automatically accessible.

## Acceptance criteria
- CLI output remains linear and understandable with a screen reader.
- Meaning is never communicated through color alone; color-disabled output remains complete.
- Errors lead with file, source position when available, cause, and a practical next action.
- Interactive CLI behavior has non-interactive equivalents and stable exit codes.
- Documentation uses descriptive headings and links, short navigable sections, and copyable commands.
- Generated blueprints use semantic HTML, associated form labels, keyboard-operable controls, visible focus styles, and accessible contrast.
- Generated blueprints and any Core-owned example components document keyboard and screen-reader behavior and ship accessible defaults.
- The transpiler evaluates useful diagnostics for common issues such as missing image alternatives or form labels, with documented limits and escape hatches for valid advanced markup.
- Accessibility checks combine automated tests with documented manual keyboard and screen-reader verification.
- Public wording promises accessible tools and defaults, not automatic conformance of arbitrary applications.
- Accessibility requirements for future UI plugins remain in the plugin ecosystem rather than blocking Core v0.1.
