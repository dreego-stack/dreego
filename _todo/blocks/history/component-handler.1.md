---
id: component-handler.1
title: ComponentHandler (Buffered Mode + Functional Options)
status: 18
phase: v0.0.5
requires:
  - context-refactoring.1
created: 2026-07-27
changed: 2026-07-28
---

Templ-inspired: `ComponentHandler` with buffer-first strategy (render to buffer first, on success Status+Headers+Body). Functional Options: `WithStatus(int)`, `WithContentType(string)`, `WithErrorHandler(fn)`, `WithStreaming()`. Decouples rendering from HTTP — `Render(ctx, io.Writer) error` instead of `http.ResponseWriter`. Enables SSG/SSR/Desktop without HTTP dependency.
