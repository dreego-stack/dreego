---
id: context-refactoring.1
title: Context Interface + SSRContext (map[string]string → any)
status: 7
phase: v0.0.2
requires:
  - transpiler.1
created: 2026-07-25
changed: 2026-07-26
---

Context interface with embedding context.Context, Param, Data, Render. SSRContext as concrete impl. Codegen generates *context.SSRContext + context.NewSSR. No more map[string]string.
