---
id: session.1
title: Session-Interface (Cookie Store im Core)
status: in-progress
phase: v0.0.3
requires:
  - context-refactoring.1
created: 2026-07-26
changed: 2026-07-26
---

Session-Interface im Core: Get(key)/Set(key, value)/Delete(key)/Save(). Cookie-Store als Default-Implementierung. ADR-Entscheidung session-management.md liegt vor. Kein Redis/DB — nur In-Memory + Cookie. Plugin-System erlaubt später DB-Store.
