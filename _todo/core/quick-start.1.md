---
area: developer-experience
phase: v0.1-blocker
depends_on:
  - api-freeze.1
  - app-runtime.1
---
# Make the five-minute quick start reliable

## Goal
Provide one canonical path from installation to a running Dreego application using the explicit App API.

## Acceptance criteria
- README and Getting Started teach one consistent project-creation command.
- Instructions contain no repository-local replace directive or outdated module path.
- A fresh project generates, builds, starts, and answers an HTTP request using only published public APIs.
- CI runs the documented commands as a black-box test.
- Errors explain missing Go, unsupported versions, invalid project names, and dependency-resolution failures.
