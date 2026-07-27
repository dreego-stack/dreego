---
id: ci-check.1
title: dreego generate --check (CI Mode)
status: 16
phase: v0.0.6
requires:
  - cli.1
created: 2026-07-27
changed: 2026-07-28
---

Templ-inspiriert: `dreego generate --check` validiert, dass alle generierten Dateien aktuell sind. Exit non-zero wenn veraltet. Ermöglicht CI-Pipeline: "Fail if generated code is stale". Kein Re-Generate, nur Validation.
