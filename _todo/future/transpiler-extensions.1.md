---
area: architecture
phase: v2-or-later
---
# Transpiler extension points

## Goal
Allow proven plugins to add sections such as Markdown, SVG, or charts through a formal processor interface.

## Acceptance criteria
- The lexer-to-codegen pipeline is formalized before extension hooks are added.
- Extensions cannot bypass validation or target boundaries.
- At least two real processors validate the contract.
