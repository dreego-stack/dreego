# Output Safety

Dreego renders every dynamic value through a safe-value rule that matches its
HTML context. The rule is chosen at generate time from the position of the
`{{ expression }}` placeholder, so a value that is safe as text cannot remain
dangerous in a URL, attribute, script, or style context.

## Context Rules

| Context | Example | Rule | Effect |
|---------|---------|------|--------|
| Text | `<p>{{ v }}</p>` | `SafeText` | HTML-escapes markup and quotes |
| Attribute | `<div title="{{ v }}">` | `SafeAttr` | HTML-escapes markup and quotes |
| URL | `<a href="{{ v }}">` | `SafeURL` | Rejects unsafe schemes, then HTML-escapes |
| Script | `<button onclick="{{ v }}">` | `SafeScript` | JSON-encodes the value, then HTML-escapes |
| Style | `<div style="{{ v }}">` | `SafeStyle` | HTML-escapes, so `</style>` and `<!--` are inert in any case |
| Meta refresh | `<meta http-equiv="refresh" content="{{ v }}">` | `SafeRefresh` | Validates the `url=` value, then HTML-escapes |
| Nested document | `<iframe srcdoc="{{ v }}">` | `SafeSrcdoc` | Escapes twice so markup remains text when the browser parses the nested document |
| Raw | `{{ v\|raw }}` | `SafeRaw` | No escaping — explicit opt-in for trusted values |

The URL rule applies to `href`, `src`, `action`, `srcset`, `poster`, `cite`,
`formaction`, `data`, `background`, `longdesc`, `usemap`, and `xlink:href`
(SVG `use`). The script rule applies to every attribute whose name starts with
`on` followed by a letter (`onclick`, `onload`, `onerror`, …) and to the
Alpine/HTMX event directives `x-on:*`, `@*`, and `hx-on:*` (`x-on:click`,
`@click`, `hx-on:click`, `hx-on::before-request`). The Alpine evaluator
directives `x-data`, `x-init`, `x-effect`, `x-html`, `x-show`, `x-model`,
`x-text`, and `x-transition` are also treated as script context, because
Alpine evaluates their values as JavaScript. The style
rule applies to the `style` attribute and to the Alpine style bindings
`x-bind:style` and `:style`.

The `<head>` section is a string, not parsed markup, so head expressions are
classified by scanning the surrounding tag with the same rules as body
attributes: `<link href="{{ v }}">` uses the URL rule, `<meta
http-equiv="refresh" content="{{ v }}">` validates the embedded `url=` value,
and `<title>{{ v }}</title>` uses the text rule. Head classification follows
HTML attribute syntax: attribute names, `=`, and values are matched
case-insensitively with whitespace tolerance, and unquoted attribute values are
recognized, so `<meta HTTP-EQUIV = refresh content={{ v }}>` is still
classified as a refresh meta tag.

## URL Scheme Validation

`SafeURL` allows relative URLs, fragment-only values (`#fragment`),
protocol-relative URLs (`//host/path`), and the `http`, `https`, `mailto`,
and `tel` schemes. Any other scheme — including `javascript:`, `data:`,
`vbscript:`, and `file:` — is replaced with `#`. The check is case-insensitive
and rejects obfuscated schemes that embed whitespace or control characters
(`\x00`, `\x0b`, `\x0c`, space, tab, newline, carriage return) before the
colon, because browsers strip those characters while parsing URLs. `srcset`
values are validated as a whole: a value that contains an unsafe scheme
anywhere (`a.jpg 1x, javascript:alert(1)`) is rejected.

`SafeRefresh` extracts the URL portion after `url=` with the same tolerance
that browsers apply: the `url` keyword is matched case-insensitively and
whitespace around `url`, `=`, and the URL itself is ignored. The extracted
URL is validated with the same scheme rules as `SafeURL`. Unsafe variants such
as `0; url = javascript:alert(1)`, `0;url =javascript:alert(1)`, and
`0;URL = javascript:alert(1)` are all rejected and replaced with `url=#`.
The extraction spans newlines and carriage returns, so obfuscated schemes
such as `0;url=java\nscript:alert(1)` and `0;url=java\rscript:alert(1)` are
rejected too — browsers strip those characters before parsing the scheme.

## Raw Opt-In

`{{ expression|raw }}` emits the value without any escaping. It is the visible,
documented opt-in for trusted HTML, URLs, or scripts. Using `|raw` with
untrusted input reintroduces the XSS risk that the context rules prevent.

## Guarantees and Limitations

- Every `{{ expression }}` placeholder is escaped for its context by default.
- `|raw` is the only way to bypass escaping, and it is explicit in the source.
- The rules protect the generated HTML document. They do not protect values
  that are written directly with `c.Write`, `c.JSON`, or `c.XML`, or values
  stored in sessions and re-emitted through other channels.
- `SafeURL` validates the scheme, not the destination. A `https:` URL can still
  point anywhere; applications that need an allowlist of hosts must validate
  the parsed URL themselves.
- `SafeScript` JSON-encodes the value, so it is safe as a JS string literal.
  It does not make arbitrary JavaScript safe to execute.
- `SafeStyle` HTML-escapes the value, so `</style>` and `<!--` cannot break
  out of the style element in any casing. It does not neutralize embedded URLs
  such as `background: url(javascript:…)`; style values that must contain
  dynamic URLs should be validated before being rendered.
- The script rule treats any attribute whose name starts with `on` followed by
  a letter as an event handler. Attributes such as `once` or `only` therefore
  get the conservative `SafeScript` treatment even though they are not event
  handlers; this is safe (the value is JSON-encoded, never executed) and
  avoids false negatives for unknown event handler names.
- Alpine value bindings other than `x-bind:style`/`:style` (for example
  `x-bind:href="{{ v }}"` or `:href="{{ v }}"`) fall back to the attribute rule,
  which escapes the value but does not validate its scheme. The URL rule only
  applies to the static attribute names listed above; validate bound URLs in
  the `<server>` section before rendering them.
- Static markup in `.dreego` templates is emitted verbatim. Only dynamic
  `{{ expression }}` values are subject to the context rules.
