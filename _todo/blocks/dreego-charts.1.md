---
id: dreego-charts.1
title: dreego-charts (Chart.js/Canvas Components)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - components.1
created: 2026-07-26
changed: 2026-07-26
---

Target directory: `plugins/charts/` in this repository.

Plugin/Component library for charts. Chart.js/Canvas integration. <Chart type="line" data={...} /> syntax. SSR-compatible (server renders canvas data as JSON, client hydrates). No core interface — pure UI components.
