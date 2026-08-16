---
version: patch
---

- Bug: redirect/rewrite wildcard rules matched near-prefix paths (/api/* matched /apiary); now segment-boundary safe
- Feat: exact vs wildcard redirect/rewrite semantics documented and validated at configuration time
- Feat: invalid patterns, targets, redirect status codes, and self/wildcard loops rejected during RegisterRedirect/RegisterRewrite
- Breaking: external-URL redirect/rewrite targets (e.g. https://example.com) rejected — same-origin paths only (pre-v0.1, deliberate)