---
area: core
phase: after-multi-language-processors
---
# First-party internationalization

## Goal

Provide a complete, replaceable i18n implementation with compile-time-checked
template messages, project catalogs, request locale resolution, and accessible
document metadata without requiring an external plugin.

## Acceptance criteria

- `[[ key ]]` and parameterized message expressions work in HTML bodies,
  attributes, heads, layouts, components, and every HTTP method body.
- Catalogs under `locales/` are validated and compiled by `dreego generate`.
- The runtime supports the locale features and quality requirements recorded in
  `_docs/decisions/i18n-core.md`.
- The standard resolver supports cookies, `Accept-Language`, fallback, optional
  localized URLs, ordered custom resolvers, and an optional selection handler.
- IP geolocation is not built in and can be registered before fallback.
- Generated documents expose correct `lang` and text direction metadata.
- The built-in implementation can be disabled or replaced before App build.
- Unit, integration, accessibility, fuzz, race, and deterministic-generation
  tests cover the public behavior.
- Public documentation explains URL, caching, privacy, accessibility, and SEO
  tradeoffs without claiming that one routing strategy fits every application.

## Dependency

- Semantic sections and first-party language processors are complete.
- The proposed ADR is accepted and its open questions are resolved.
