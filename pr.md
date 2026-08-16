---
version: patch
---

- Security: render every `{{ expression }}` through a context-aware safe-value rule (text, attribute, URL, script, style, meta refresh) instead of a single HTML escape
- Security: URL attributes (`href`, `src`, `action`, `srcset`, `xlink:href`, …) reject unsafe schemes such as `javascript:` and `data:` by default, replacing them with `#`
- Security: reject obfuscated URL schemes that embed whitespace or control characters (`\x00`, `\x0b`, `\x0c`, space, tab, newline, carriage return) before the colon — browsers strip them while parsing, so `java\nscript:` previously executed
- Security: validate `<meta http-equiv="refresh">` `content` URLs through `SafeRefresh`; the `url=` keyword is matched case-insensitively with whitespace tolerance (`0; url = javascript:…` is rejected), the extraction spans newlines so `0;url=java\nscript:…` is rejected too, and `<link href>` in `<head>` now uses the URL rule
- Security: script attributes (`onclick`, `hx-on:*`, `x-on:*`, `@*`, `x-data`, `x-init`, `x-effect`, `x-html`, `x-show`, `x-model`, `x-text`, `x-transition`, …) JSON-encode the value; style attributes and Alpine style bindings (`style`, `x-bind:style`, `:style`) HTML-escape so `</style>` and `<!--` are inert in any casing
- Security: head classification follows HTML attribute syntax — `http-equiv = "refresh"` with whitespace around `=` and unquoted values such as `<link href={{ u }}>` are recognized as their real context
- Feat: `|raw` remains the explicit, documented opt-in for trusted values; generated code no longer imports `html`
- Docs: add `_docs/security.md` describing the exact guarantees and limitations of the context rules
