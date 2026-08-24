---
area: architecture
phase: planned-v0.3
---
# Transpiler extension points

## Goal
Allow external plugins to process explicit section/language pairs through a
versioned subprocess protocol. See `_plan/v0.3-language-processors.md`.

Markdown bodies and TypeScript clients are the first two protocol proofs.
Dreego retains ownership of component composition, control flow, escaping,
source positions, and the single `<body>` template root.

## Acceptance criteria
- The lexer-to-codegen pipeline is formalized before extension hooks are added.
- Extensions cannot bypass validation or target boundaries.
- At least two real processors validate the contract.
- Processor dependencies remain outside the Dreego root module.
