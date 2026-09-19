---
version: patch
---

- Refactor: move `internal/transpiler` to `internal/dreefile` and regroup the compiler into `ir`, `dreecode`, `gogen`, `codegen`, `jsoutput`, `tokens`, `lexer`, `parser`, and `sections/{head,style,body/{html,md},client/{js,ts,lua}}`
- Refactor: move the `<md>`-tag region scanning out of the parser into `sections/body/md`, so the parser imports no Markdown package
- Refactor: split the mini-template-language semantics (`ParseExpression`, `FindExprEnd`, `FindMessageEnd`, `ParseMessageExpression`) into `dreecode`, and Go emission and source positions into `gogen`
- Refactor: make `sections/client` the explicit orchestrator over the `client/{js,ts,lua}` language children
- Refactor: add `internal/md` as the shared runtime Markdown string path so `core` no longer imports the compiler
- Docs: correct the layering ADR and phase plan to the verified dependency edges (`body/html -> head, style, client`; `lexer`/`parser -> dreecode`; `jsoutput -> gogen`; test files excluded from the structural check) and re-point stale docs from `internal/transpiler` to `internal/dreefile`
- Docs: record the deliberate scope deviation that `route/`, `component/`, `layout/`, and `assets/` stayed in the `dreefile` root this phase instead of becoming own packages
- Chore: add `_tests/sh/check-layering.sh` and wire it into `Taskfile.yml` and `_tests/test.sh` so the dreefile layering rule is enforced in CI
- Chore: assert the compiler dependency rule over non-test files (sections, shared packages, front-end direction) and that no compiler package imports `core`
