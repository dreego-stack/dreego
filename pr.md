---
version: patch
---

- Docs: mark superseded architecture decisions (Chi, validator, Tailwind-core, monorepo-plugin, target, runtime) as historical in .agents/decisions/
- Docs: restructure the decision index (.agents/index.md) into current, provisional, research, and superseded sections
- Docs: add progressive enhancement guide (_docs/progressive-enhancement.md) with a worked HTMX/Alpine.js/plain-JS example and security, failure, and no-JavaScript behavior
- Docs: link the progressive enhancement guide from _docs/index.md and _docs/sitemap.json
- Docs: fix the worked example in the progressive enhancement guide — add the missing g-action, use the current filename-based routing scheme, and use hx-boost instead of the unimplemented fragment-swap claim
- Docs: correct the HTMX CSRF claims in the progressive enhancement guide and _docs/forms.md (HTMX does not send the token automatically; use the hidden field or hx-headers)
- Docs: align the AGENTS.md "Architecture Guarantees" section with the superseded SSG/Wails decisions (SSR-only until v1, no reserved CLI flags)
- Docs: mark the dreego-architecture and signals-and-runes concepts as historical where superseded (monorepo plugin layout, Tailwind core dependency, Datastar/SSE client options)
