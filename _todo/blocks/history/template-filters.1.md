---
id: template-filters.1
title: Template Filters ({var|raw}, {var|upper})
status: 23
phase: v0.0.9
requires:
  - transpiler.1
  - xss.1
created: 2026-07-26
changed: 2026-07-28
---

Pipe syntax for template expressions: {name|upper}, {date|format:"2006-01-02"}, {html|raw} (escape hatch for trusted HTML). Extend parser, CodeGen generates filter chain.
