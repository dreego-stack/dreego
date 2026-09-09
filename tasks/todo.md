# First-party i18n tasks

- [ ] Tokenize `[[ key ]]` and parameterized message expressions.
  - Acceptance: lexer preserves content and source position in body and head.
  - Verify: targeted lexer tests in `smd`.
  - Files: token definitions, lexer, lexer tests.

- [ ] Parse message tokens into dedicated IR nodes.
  - Acceptance: keys and Go argument expressions are structured and malformed
    syntax reports source positions.
  - Verify: parser tests and parser fuzz tests.
  - Files: IR node, body parser, head parser, parser tests.

- [ ] Discover and validate project locale configuration and catalogs.
  - Acceptance: default locale, supported locales, duplicate keys, missing
    defaults, fallback cycles, and deterministic discovery are validated.
  - Verify: config and generator integration tests.
  - Files: config types, catalog loader, generator wiring, tests.

- [ ] Add the built-in catalog runtime.
  - Acceptance: BCP 47 matching and fallback use pinned `golang.org/x/text`;
    the App owns immutable built state.
  - Verify: runtime unit, race, and dependency-gate tests.
  - Files: i18n runtime package, App configuration, facade, tests.

- [ ] Generate escaped message output end to end.
  - Acceptance: messages render in bodies, attributes, heads, layouts,
    components, Markdown, and all HTTP methods without server preparation.
  - Verify: output unit tests and `_tests/go/i18n_messages_test.go`.
  - Files: HTML output, generated catalog, integration test, fixture.

- [ ] Resolve request locales.
  - Acceptance: explicit choice, account, cookie, browser, custom resolvers, and
    default follow documented precedence; unsupported hints safely pass.
  - Verify: table-driven HTTP, race, and proxy tests.
  - Files: resolver, middleware, App API, tests.

- [ ] Add selection and document metadata.
  - Acceptance: local return paths, loop prevention, `html[lang]`, RTL direction,
    and non-JavaScript operation are tested.
  - Verify: security and accessibility integration tests.
  - Files: selection handler, document output, integration tests, docs.

- [ ] Add advanced formatting and translation workflow.
  - Acceptance: plural, ordinal, select, number, percent, currency display,
    date, time, time zone, pseudolocale, and catalog export are deterministic.
  - Verify: locale matrix tests and CLI tests.
  - Files: formatter packages and focused test files, kept below size limits.

- [ ] Complete release documentation and quality gates.
  - Acceptance: guides cover configuration, fallbacks, caching, privacy, SEO,
    accessibility, replacement, and the no-conversion currency boundary.
  - Verify: full suite, coverage gate, `dreego generate --check`, file lengths.
  - Files: docs, reference app, sitemap, change entry.
