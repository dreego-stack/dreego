# Accessibility

Accessibility is a release quality gate for Dreego, not a cosmetic enhancement. CLI output, diagnostics, documentation, generated templates, and official components are designed to work without relying on sight, color, or pointer input alone.

Dreego does **not** make arbitrary user applications automatically accessible. The framework provides accessible defaults and diagnostics; applications still verify their own content and conformance.

## CLI Output

- No ANSI color codes: meaning never relies on color alone.
- Help text starts with the program name and a linear `usage:` line for screen readers.
- Generator diagnostics lead with `file:line:col`, the cause, and a practical `Fix:` action.
- Interactive workflows (`dev`, `run`) have non-interactive equivalents and stable exit codes.

## Transpiler Diagnostics

`dreego generate` emits warnings for common accessibility issues:

- `<img>` without an `alt` attribute (use `alt=""` for decorative images).
- `<input>` without an associated `<label>` (via `label[for]` matching the input `id`).

Diagnostics are printed to stderr in the form `file:line:col: cause Fix: ...`. They are warnings: generation completes so you can review them alongside other feedback, but they should be resolved before shipping.

Limits: the checker scans static markup in `.dreego` templates. Dynamically composed attributes or labels added at runtime are not detected. Treat the checker as a first-pass guard, not a conformance audit.

## Templates

Two scaffolds ship today: `web-minimal` (the default) and `web-app`.

`web-minimal` ships a minimal accessible layout in
`www/layouts/default.dreego`:

- `<html lang="en">`, a UTF-8 charset, and a viewport meta tag.
- A skip link (`<a href="#main">skip to content</a>`) that is visually hidden
  until focused, and a `<main id="main">` landmark around `{#slot}`.
- A visible focus style: `:focus-visible { outline: 3px solid #1d4ed8; }`.

The route in `www/routes/+page.dreego` stays intentionally small: one `<h1>`,
one paragraph, and a `:focus-visible` outline. It has no navigation, form
controls, or images, so it is a starting point rather than a complete
accessible application shell. Add the landmarks, labels, and alternatives your
application needs as you build it.

`web-app` ships a fuller SSR application shell in the same layout shape: a skip
link, a `<header>` with primary navigation, a labelled `<nav>` whose active link
uses a non-colour state (`aria-current="page"` plus weight/underline), a
`<main id="main">` landmark, a footer, and a visible `:focus-visible` outline.
Its `/` route has one `<h1>`, a labelled text input, and a typed server-rendered
form; `/dashboard` has one `<h1>` and a table with a caption and column headers.
Each route supplies its own `<title>` and meta description, so the layout
declares no static title. These are semantic defaults, not a conformance claim:
applications still verify their own content.

## Documentation

- Descriptive headings and links, short navigable sections, copyable commands.
- Test counts are described by layout, not by unstable numbers.
- Examples use the explicit `App` API (`app := dreego.New(); www.Register(app)`).

## Manual Verification

Automated checks cover static markup and CLI behavior. Before claiming accessibility for a real application, verify manually:

- Keyboard navigation through all interactive controls.
- Screen-reader announcement of landmarks, labels, and updates.
- Contrast ratios for text and controls.
- Reduced-motion and high-contrast preferences.

## See Also

- [Testing](../dreegotest/_docs/testing.md) — Accessibility test entries
- [CLI](../cmd/dreego/_docs/cli.md) — CLI reference
- [Getting Started](getting-started.md) — Tutorial
