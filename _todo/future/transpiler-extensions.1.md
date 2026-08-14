---
area: architecture
phase: v2-or-later
---
# Transpiler extension points

## Goal
Allow proven plugins to add sections such as Markdown, SVG, or charts through a formal processor interface.

Markdown and MDX are examples for evaluating the boundary, not committed
features. A proposal must first explain styling, component composition, i18n,
escaping, and conflicts with the single `<div>` template root.

## Acceptance criteria
- The lexer-to-codegen pipeline is formalized before extension hooks are added.
- Extensions cannot bypass validation or target boundaries.
- At least two real processors validate the contract.
