---
id: named-slots.1
title: Named Slots ({#slot header}...{/slot})
status: 19
phase: v0.0.8
requires:
  - component-handler.1
created: 2026-07-28
changed: 2026-07-28
---

Named Slots: `{#slot header}...{/slot}` Block-Syntax. Component: `{#slot header}{/slot}` als Platzhalter. Route: `{#slot header}content{/slot}` als Definition. Mehrere benannte Slots pro Component. TokenSlotOpen/TokenSlotClose im Lexer. parseSlotNodes() im Parser. Codegen: `c.Set("slot_name", content)` in Route, `c.Get("slot_name")` in Component.
