---
version: patch
---

- Security: runtime markdown rendering is safe by default — raw HTML is escaped, URLs are scheme-validated, scripts are structurally impossible (dreego.MarkdownToHTML)
- Feat: dreego.MarkdownToHTMLTrusted for fully-controlled content (raw HTML passthrough, start-time warning)
- Feat: dreego.mdtohtml(...) stdlib syntax in server sections (trusted: true per call site)
- Fix: entity-encoded URL scheme bypass (java&#x09;script:) closed
- Fix: fenced-code language attribute injection closed
- Fix: data race in the markdown renderer (per-call renderer)
- Fix: news demo — title escaping, poisoned-post resilience, store mutex
- Breaking: duplicate <server> sections with the same method and content-type are now rejected at generation time (components keep concatenation)
- Tests: XSS matrix (18 hostile inputs) + fuzz target for the safe runtime path
