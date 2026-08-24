
---
type: Decision
title: No `<catch>` Tag — Errors via Go Idioms
description: Error handling via Go idioms in the go block instead of a special catch tag
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# No `<catch>` Tag — Errors via Go Idioms

**Date:** 2026-07-28
**Status:** Accepted

## Context

Originally, a `<catch>` block was proposed to catch errors in the template (inspired by JavaScript's try/catch).

## Problem

`<catch>` is a concept from languages with exceptions (JavaScript, Java) and feels foreign in the Go world.

Go developers handle errors explicitly (`if err != nil`). A separate error handling tag would violate Go idioms and unnecessarily increase the learning curve for Go developers.

## Decision

**No `<catch>` tag.** Errors are handled in the `<server>` block and passed to the template as variables:

```html
<server>
    user, err := db.GetUser(id)
    hasError := err != nil
</server>

<body class="profile">
    {#if hasError}
        <p class="error">User could not be loaded.</p>
    {#else}
        <h1>Hello, {user.Name}!</h1>
    {/if}
</body>
```

## Rationale

1. Keeps the template language extremely slim
2. Go developers don't need to learn anything new
3. No framework magic — just variables and `{#if}`
4. Explicit error handling is idiomatic Go

## Consequences

- Template logic (`{#if}`) also handles error rendering
- Developers must handle errors themselves (as is common in Go)
- No implicit error handling in the framework
