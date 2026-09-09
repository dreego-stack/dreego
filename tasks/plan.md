# First-party i18n implementation plan

## Objective

Implement the accepted first-party i18n decision in reviewable vertical slices
without exposing a partial public feature.

## Assumptions

- `stage/v0.5` semantic sections are the implementation base.
- `[[ message.key ]]` is distinct from Go `{{ expression }}` interpolation.
- Existing `[[param]]` catch-all route filenames remain unchanged because route
  filenames and template bodies are separate grammars.
- Currency localization formats but never converts money.

## Technical order

1. Add lexer and parser representation for message expressions.
2. Add catalog configuration, discovery, validation, and deterministic output.
3. Add the built-in runtime catalog and `golang.org/x/text` foundation.
4. Generate escaped message rendering for routes, layouts, heads, attributes,
   components, and method bodies.
5. Add cookie, browser, custom resolver, and fallback locale resolution.
6. Add language selection, replacement APIs, and document metadata.
7. Add plural, select, typed formatting, pseudolocalization, and export support.
8. Complete documentation, fuzzing, race tests, accessibility checks, and the
   reference application.

## Risks and mitigations

- Template ambiguity: tokenize message expressions before parsing and fuzz
  interactions with Markdown, verbatim blocks, attributes, and route syntax.
- Public API freezing: keep the feature internal until an end-to-end slice
  proves the smallest App contract.
- Unicode correctness: use a Go-version-compatible pinned `golang.org/x/text`.
- Cache leakage: require locale-aware cache keys and test response variation.
- Unsafe translated markup: escape by HTML context and omit raw messages from
  the initial contract.
- Redirect abuse: accept only validated local return paths.

## Verification checkpoints

- Each task starts with a failing permanent test.
- Each completed slice passes its targeted tests and the dependency gate in
  `smd` before commit.
- Public slices pass the full `_tests/test.sh` suite before review.
- Every handwritten file remains at or below 300 lines.

## Boundaries

- Always preserve dependency injection and App ownership.
- Ask before changing the accepted template or catalog syntax.
- Never add IP geolocation, exchange-rate conversion, or translated raw HTML to
  Core.
