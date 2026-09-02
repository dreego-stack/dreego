---
area: transpiler
phase: v0.3-or-later
depends_on: none
---
# Markdown pipe tables

## Goal
Table support for the markdown processor (pipe tables per CommonMark GFM extension).

## Acceptance criteria
- Pipe table syntax `| a | b |` renders as an HTML `<table>`.
- A header separator row (`| --- | --- |`) distinguishes the header from body rows.
- Alignment variants (`:---`, `:---:`, `---:`) map to the matching `align` attributes.
- The feature has unit tests covering the base case, headerless rows, alignment, and edge inputs.
- Golden coverage is updated so the generated output is locked in.

## Dependencies
- none
