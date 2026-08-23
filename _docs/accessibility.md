# Accessibility

Accessibility is a release quality gate for Dreego, not a cosmetic enhancement. CLI output, diagnostics, documentation, generated blueprints, and official components are designed to work without relying on sight, color, or pointer input alone.

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

## Blueprints

`dreego new` scaffolds a landing-page example with:

- A layout using `<main>`, `<nav aria-label="Primary">`, a skip-link to `#main`, and `{#slot}`.
- Routes and components whose images include an `alt` attribute.
- Visible focus styles (Tailwind `focus:` variants) and semantic landmarks.

`dreego init` ships a minimal accessible layout (`www/layouts/default.dreego`) with a skip link, `<main id="main">` landmark, and `<html lang="en">`. The route stays intentionally small (a heading and a paragraph) so you can see the scaffolded structure. Add navigation, footer, and any application-specific landmarks as you build.

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

- [Testing](testing.md) — Accessibility test entries
- [CLI](cli.md) — CLI reference
- [Getting Started](getting-started.md) — Tutorial
