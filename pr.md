---
version: patch
---

- Bug: fix gzip middleware corrupting streamed responses (Flush produced an invalid gzip member without trailer)
- Bug: fix infinite recursion in gzip response writer ReadFrom causing stack overflow on io.Copy with non-WriterTo sources
- Bug: fix Accept-Encoding negotiation ignoring case and letting wildcard override an explicit gzip;q=0
- Feat: document gzip compression semantics (q-values, Vary, skip rules, Content-Length handling, optional writer interfaces, panic-safe buffering)
