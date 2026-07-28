---
id: live-reload.1
title: Live Reload Proxy (SSE + Script Injection)
status: planned
phase: v0.1.0
requires:
  - hot-reload.1
created: 2026-07-27
changed: 2026-07-27
---

Templ-inspired: Reverse proxy (`dreego dev --proxy`) on own port, injects `<script>` before `</body>`, SSE endpoint `/__dreego/reload/events`. Changes to `.dreego` files trigger browser reload. Gzip/Brotli-transparent. Pattern from templ's generatecmd/proxy/.
