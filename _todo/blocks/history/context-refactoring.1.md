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

Context-Interface mit embedding context.Context, Param, Data, Render. SSRContext als konkrete Impl. Codegen generiert *context.SSRContext + context.NewSSR. Kein map[string]string mehr.
