---
type: Decision
title: First-party internationalization
description: Dreego provides replaceable, compile-time-validated localization
tags: [i18n, transpiler, runtime, accessibility]
timestamp: 2026-09-09T00:00:00Z
---
# First-party internationalization

**Date:** 2026-09-09
**Status:** Accepted

## Context

Internationalization affects template syntax, request state, generated code,
HTML language metadata, catalogs, diagnostics, routing, and application caches.
Preparing translated strings in each method-specific server section would
duplicate work and make translations unavailable to layouts and components.

Dreego already distinguishes a section processor language such as
`<body lang="md">` from the language of rendered content. These concepts must
remain separate.

## Decision

Dreego provides a complete first-party internationalization implementation.
Applications can disable it or replace the localizer through an explicit App
method. No external plugin is required for the standard experience.

Templates use a dedicated message expression:

```html
<body>
    <h1>[[ home.title ]]</h1>
    <p>[[ cart.items count=len(items) ]]</p>
</body>
```

Messages live outside `.dreego` files:

```text
locales/
├── de/
│   ├── common.ftl
│   └── cart.ftl
└── en/
    ├── common.ftl
    └── cart.ftl
```

`dreego generate` validates message syntax, keys, arguments, locale coverage,
and fallback cycles. Catalogs are compiled into deterministic generated Go;
production rendering does not read catalog files from disk.

The `lang` attribute on semantic sections continues to select their source
processor. I18n configuration uses `locale` terminology. When i18n is enabled,
Dreego supplies the selected BCP 47 locale to the rendered document's `<html
lang>` attribute and supplies `dir="rtl"` when required.

## Configuration

```json
{
  "i18n": {
    "enabled": true,
    "defaultLocale": "de",
    "locales": ["de", "en"],
    "urlStrategy": "none",
    "detection": ["account", "cookie", "browser", "custom", "default"]
  }
}
```

Supported URL strategies are `none`, `prefix`, and `domain`. URL selection is
independent from locale detection. An explicit URL locale, when configured,
wins over implicit detection.

## Resolution

The default order without localized URLs is:

1. explicit user choice for the current request;
2. application account resolver;
3. Dreego locale cookie;
4. browser `Accept-Language`;
5. ordered custom resolvers;
6. website default locale.

Core does not perform IP geolocation. A GeoIP integration is a custom resolver
placed before the default. Automatic resolvers return only supported locales
and never overwrite an explicit user choice.

An unsupported explicit selection is rejected instead of silently pretending
that the fallback was selected. Stale cookies and automatic unsupported hints
are ignored, allowing later resolvers to continue.

## Optional negotiation

Locale resolution and user interaction are separate. An optional selection
handler may render a region chooser or expose a suggested locale to an
application banner before the website default is committed. It must preserve
only local return paths, prevent redirect loops, work without JavaScript, and
remain operable by keyboard and screen reader.

Core does not render a branded popup. `dreego-ui` may provide an accessible
language and region selector while other UI libraries provide alternatives.

## Replacement model

The built-in localizer is the default. The App owns its configuration and may
disable or replace it before build, following the session and error-handler
configuration model. Replacement contracts remain provisional until their
first real implementation validates them.

The public contract must preserve typed message arguments in generated code.
A string key plus `map[string]any` may exist as an explicit dynamic boundary,
but it is not the generated template path.

## Dependencies

Core and the transpiler prefer the Go standard library. Modules maintained by
the Go project under `golang.org/x/` are allowed without requiring the feature
to move into a plugin. Third-party dependencies remain outside Core.

The built-in implementation uses a Go-version-compatible release of
`golang.org/x/text` for BCP 47 matching and Unicode locale behavior. Dependency
versions are pinned and updated through the normal review and test workflow.

## Quality requirements

The first stable release must cover BCP 47 matching, explicit fallbacks,
cardinal and ordinal plurals, select variants, parameters, numbers, percent,
currency display formatting, date and time, time zones, RTL metadata,
pseudolocalization, missing message diagnostics, catalog namespaces,
HTML-context escaping, and extraction for translation-management systems.

Currency localization formats an amount and its existing ISO 4217 currency. It
never converts an amount, selects another currency, fetches exchange rates, or
applies financial rounding policy. Exchange-rate conversion is dynamic domain
logic owned by the application or a financial integration.

Locale-sensitive caches must include the effective locale in their key. The
`none` URL strategy documents its CDN and search-indexing consequences;
localized public content recommends `prefix` or `domain` without requiring it.

## Not doing

- No `<lang>` root section.
- No translation preparation in `<server>` sections.
- No built-in IP database or geolocation service.
- No flag-only language control.
- No unescaped translated HTML in the initial contract.
- No silent fallback for a user's invalid explicit choice.

## Initial implementation choices

- Catalog behavior is based on `golang.org/x/text`; Dreego does not implement a
  partial Fluent parser and does not add a third-party Fluent dependency.
- Parameterized expressions use `[[ cart.items count=len(items) ]]`.
- A missing default-locale message is an error. An explicitly configured
  non-default fallback is allowed with a diagnostic.
- Locale resolvers only resolve. A separate selection handler owns an optional
  banner or region-selection response.
- Generated formatting accepts concrete Go values. Currency requires an amount
  and an explicit ISO 4217 code; it performs no conversion.
