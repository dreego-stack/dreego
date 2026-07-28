---
id: template-filters.1
title: Template-Filter ({var|raw}, {var|upper})
status: 23
phase: v0.0.9
requires:
  - transpiler.1
  - xss.1
created: 2026-07-26
changed: 2026-07-28
---

Pipe-Syntax für Template-Ausdrücke: {name|upper}, {date|format:"2006-01-02"}, {html|raw} (Escape-Hatch für vertrauenswürdiges HTML). Parser erweitern, CodeGen generiert Filter-Chain.
