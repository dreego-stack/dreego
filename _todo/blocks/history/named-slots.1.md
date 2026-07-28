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

Named Slots: `{#slot header}...{/slot}` block syntax. Component: `{#slot header}{/slot}` as placeholder. Route: `{#slot header}content{/slot}` as definition. Multiple named slots per component. TokenSlotOpen/TokenSlotClose in lexer. parseSlotNodes() in parser. Codegen: `c.Set("slot_name", content)` in Route, `c.Get("slot_name")` in Component.
