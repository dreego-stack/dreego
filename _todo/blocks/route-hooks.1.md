---
id: route-hooks.1
title: Plugin-Route-Registration
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Eintrittspunkt für Plugin-Routen. Plugins registrieren eigene URL-Pfade (z.B. /admin/*, /api/auth/*). Keine Dateisystem-basierte Registrierung — programmatisch via Plugin-API. Generierte dreego/gen/dree.go sammelt alle Plugin-Routen.
