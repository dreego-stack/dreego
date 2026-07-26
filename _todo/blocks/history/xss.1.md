---
id: xss.1
title: XSS Auto-Escaping ({variable} → html.EscapeString)
status: 9
phase: v0.0.2
requires:
  - transpiler.1
created: 2026-07-25
changed: 2026-07-26
---

Template-Expression-Nodes generieren html.EscapeString(fmt.Sprintf("%v", expr)). Conditional "html" Import nur bei Expressions.
