---
id: session.1
title: Session Interface (Cookie Store in Core)
status: 14
phase: v0.0.3
requires:
  - context-refactoring.1
created: 2026-07-26
changed: 2026-07-26
---

Session interface in Core: Get(key)/Set(key, value)/Delete(key)/Save(). Cookie store as default implementation. ADR decision session-management.md exists. No Redis/DB — only In-Memory + Cookie. Plugin system allows DB store later.
