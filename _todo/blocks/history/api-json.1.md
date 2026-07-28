---
id: api-json.1
title: Content-Type Routing (JSON, XML, Custom)
status: 32
phase: v0.0.15
requires:
  - routing.1
  - context-refactoring.1
created: 2026-07-26
changed: 2026-07-29
---

<go type="json|xml|custom"> — Content-type routing via Accept header. JSON: c.JSON() + c.Bind(). XML: c.XML(). Custom: c.Write() + c.Wants(). Multiple <go> blocks merged into single handler: shared code runs always, typed blocks fire on Accept match. First match wins. c.Wants() returns false for empty Accept (HTML default).
