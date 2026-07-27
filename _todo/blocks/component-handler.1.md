---
id: component-handler.1
title: ComponentHandler (Buffered Mode + Functional Options)
status: planned
phase: v0.0.4
requires:
  - context-refactoring.1
created: 2026-07-27
changed: 2026-07-27
---

Templ-inspiriert: `ComponentHandler` mit Buffer-first-Strategie (erst in Buffer rendern, bei Erfolg Status+Headers+Body). Functional Options: `WithStatus(int)`, `WithContentType(string)`, `WithErrorHandler(fn)`, `WithStreaming()`. Entkoppelt Rendering von HTTP — `Render(ctx, io.Writer) error` statt `http.ResponseWriter`. Ermöglicht SSG/SSR/Desktop ohne HTTP-Abhängigkeit.
