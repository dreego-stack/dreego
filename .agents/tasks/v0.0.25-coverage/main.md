---
type: Task
title: v0.0.25 Coverage Improvement — Test Backlog from 3 Coverage Reviewers
status: in_progress
assign: manager
---

# v0.0.25 Coverage Improvement

Goal: Close test-coverage gaps identified by three read-only coverage reviewers. Work through the backlog item by item: test-writer writes tests (RED where a code gap exists) → coder fixes → reviewer → commit → next.

## Source reports (consolidated)

Three coverage reviewers reported gaps. Priority HIGH first, then MEDIUM, then LOW.

### HIGH — core runtime/middleware/context/response/recovery
1. runtime.go — Redirect/Rewrite/Session/Setter ungetestet (runtime_test.go erweitern)
2. middleware_csrf.go — Token-Validierung/403-Pfade ungetestet (nur Cookie-Attrs getestet)
3. context.go — komplett ohne Tests (context_test.go neu)
4. response.go — komplett ohne Tests (response_test.go neu)
5. middleware_recovery.go — komplett ohne Tests (middleware_recovery_test.go neu)

### HIGH — core lexer/parser
6. lexer.go — Lex Error-Pfade, scanTag/scanComponentTag Randfälle (lexer_test.go neu)
7. lexer_brace.go — scanBrace Token + Fehlerpfade (lexer_brace_test.go neu)
8. parser.go — Parse Fehlerpfade, parseGoAttrs (parser_test.go neu)
9. parser_section_div.go — parseTemplateNode Fehlerpfade
10. parser_template.go — parseEachClause, Control-Flow Fehler
11. codegen_template.go — NodeIf else-if chain, NodeEach else, NodeSlot, NodeComponentCall, NodeVerbatim, filters
12. codegen_component.go — compGen each/if-else-if/slot/verbatim/filters, genComponentCall non-self-close

### HIGH — CLI
13. version.go — komplett ungetestet (version_test.go neu)
14. main.go cmdGenerate --check Fehlerpfade (check-no-gen Integration)
15. main.go cmdBuildE --target parsing (build-target Integration)
16. new.go findLocalCore + cmdNew Fehlerpfade (new_test.go + new-exists/new-no-arg Integration)

### MEDIUM
17. fmt.go — Format Funktionen (fmt_test.go neu)
18. codegen_helpers.go — goLiteral, toPascalCase, scopeSelector, splitTopLevelComma, matchBrace
19. codegen.go — GenerateErrorHandler 404/500, genTypedBlocks, GenerateMethodHandler FormActions
20. codegen_head.go — genHead expression/filter/unclosed
21. codegen_layout.go — splitLayoutText, genLayoutNode
22. codegen_split.go — splitGoSections, unindent
23. middleware_health.go — funktionale Tests (SetReady)
24. middleware_requestid.go — Tests (neu)
25. middleware_compress.go — Tests (neu)
26. middleware_security.go — übrige Header
27. middleware.go RequestLogging — Tests (neu)
28. session_crypto.go — decrypt Randfälle
29. docs.go — printJSON, cmdDump, cmdDocs --json (docs_test.go erweitern)
30. init.go — cmdInit no-arg (init-no-arg Integration)
31. blueprints — Inhalts-Validierung (new-blueprint-valid Integration)

### LOW
32. config.go — config_test.go neu
33. static.go — static_test.go neu
34. session.go — verify Randfälle
35. generate.go — buildPattern, buildPageName, cleanSegment, patternSegment, errorCatchPattern
36. lexer_header.go — ParseHeader/parseProps/parseImportLine
37. main.go cmdRun -d/-t parsing (run-timer)
38. main.go version-dispatch (version Integration)
39. dev.go startServer missing binary (dev_test erweitern)

## Execution mode

- Strictly sequential, one item at a time.
- For each item: coder-testN writes tests (RED if code gap) → coder-writeN fixes code if needed → reviewer → git commit (no push).
- Agent naming: coder-testN, coder-writeN, reviewN, gitN.
- Coverage reviewers: run 3 read-only explorers per batch to keep findings current.

## Status

- [ ] Item 1: runtime.go redirect/rewrite/session tests
- [ ] Item 2: middleware_csrf.go validation tests
- [ ] ... (fill as you go)

## Parallel packages (file-disjoint)

Items 1-11 done. Remaining grouped into file-disjoint packages so they
can be written and fixed in parallel without edit conflicts:

- Pkg A (core codegen): Items 12,17,18,19,20,21,22
  Files: core/codegen_component_test.go, fmt_test.go, codegen_helpers_test.go,
  codegen_test.go, codegen_head_test.go, codegen_layout_test.go, codegen_split_test.go
- Pkg B (core middleware+session): Items 23,24,25,26,27,28
  Files: middleware_health_test.go, middleware_requestid_test.go,
  middleware_compress_test.go, middleware_security_test.go, middleware_logging_test.go,
  session_crypto_test.go
- Pkg C (core low): Items 32,33,34,35,36
  Files: config_test.go, static_test.go, session_test.go, generate_test.go,
  lexer_header_test.go
- Pkg D (CLI): Items 13,14,15,16,29,30,31,37,38,39
  Files: cmd/dreego/version_test.go, new_test.go, docs_test.go (extend),
  dev_test.go (extend) + _tests/core/CLI/check-no-gen, build-target,
  new-exists, new-no-arg, init-no-arg, new-blueprint-valid, run-timer, version

Each package: one coder-test agent writes all tests (RED if a code gap),
one coder-write agent fixes to GREEN, reviewer, commit. Agents only run
`go test <their package>` to avoid parallel smd collisions; the full
suite runs once at the end.
