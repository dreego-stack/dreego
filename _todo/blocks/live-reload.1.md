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

Templ-inspiriert: Reverse-Proxy (`dreego dev --proxy`) auf eigenem Port, injiziert `<script>` vor `</body>`, SSE-Endpoint `/__dreego/reload/events`. Änderungen an `.dreego`-Dateien triggern Browser-Reload. Gzip/Brotli-transparent. Pattern aus templs generatecmd/proxy/.
