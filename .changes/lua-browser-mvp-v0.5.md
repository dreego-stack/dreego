---
version: minor
---

- Feat: compile root and inline browser Lua into JavaScript without external dependencies or a VM.
- Feat: generate one deterministic Lua runtime asset containing only the semantic helpers used by the application.
- Feat: compile local functions, closures, and browser event callbacks from Lua.
- Docs: add equivalent Lua and JavaScript browser demos for direct comparison.
- Docs: publish client-language version targets and plan isolated client Go and Starlark experiments.
- Fix: isolate Lua blocks and prevent client strings from terminating generated script elements.
- Fix: preserve Lua truthiness and modulo behavior at JavaScript boundaries without identifier or runtime namespace collisions.
- Fix: reject unsafe Markdown code-fence language labels before generating HTML attributes.
