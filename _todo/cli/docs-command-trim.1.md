---
area: dx
phase: pre-v0.1
---
# Evaluate trimming the docs command surface

## Goal
Evaluate trimming the docs command surface (`--web`/`--json`/`--dump`) for
pre-v0.1. If the docs say the CLI reference lists them, decide: keep+document
or remove.

## Acceptance criteria
- A decision is recorded in this todo file after implementation (mark done
  with verdict).

## Verdict: keep the full surface (option a)
`--json`, `--list`, and `-p` are referenced by integration tests
(`_tests/go/cli_scaffold_test.go`); all flags (`--web`/`--json`/`--dump`/
`--list`/`-p`) are already documented in `_docs/cli.md` and the `help` text in
`cli/dreego/main.go`. No flag removal. No code or doc change needed.
