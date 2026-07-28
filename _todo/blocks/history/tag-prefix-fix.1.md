---
id: tag-prefix-fix.1
title: scanTag: Tag-Präfix-Matching fix (head vs header)
status: 22
phase: v0.0.9
requires:
  - transpiler.1
created: 2026-07-28
changed: 2026-07-28
---

Fixing scanTag's opening tag detection: `<header>` was incorrectly matching `<head>` because `HasPrefix` didn't check the terminator character. Added check: character after tag name must be ` `, `>`, or `/` — otherwise it's a different tag. This fixes `<header>`, `<main>`, `<footer>` in component templates.
