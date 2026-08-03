# v0.0.22 Triple Block

Goal: Implement three planned blocks for v0.0.22, add integration tests, update docs, bump version, and commit.

1. servemux-cache.1 — cache the built middleware/router stack in core/runtime.go ✅
2. codegen-errors.1 — propagate errors in core/codegen*.go, fix nested {#if} in else bug in genTemplateNodeComp ✅
3. security-session.1 — optional AES-256-GCM encryption for session payload in core/session.go ✅

v0.0.22 was released on 2026-08-03 with tags core/v0.0.22, cmd/dreego/v0.0.22, plugins/sample/v0.0.22.
The late fixes listed below were moved to v0.0.23 (unreleased):

4. runtime.1 — exported core.Reset() helper to clear the cached middleware/router stack
5. security-session.2 — HMAC-SHA256 key derivation in core/session_keys.go
6. security-session.3 — propagate JSON marshal / encryption errors from CookieStore.Set
7. coding-standards.1 — raise max file line limit from 120 to 300

Additional deliverables:
- Integration tests under _tests/core/ for each block where applicable
- Update _docs/ and .agents/ docs if needed
- Bump VERSION file to v0.0.22 (done for v0.0.22 release)
- Update CHANGELOG.md and .agents/log.md
- Git commit via git subagent

Rules: standard library only in core, max 300 lines per file, one logical thing per file, no comments unless needed, all tests must pass via smd go test ./core/... and smd sh _tests/test.sh.

Status:
- v0.0.22 released: done
  - servemux-cache.1: done
  - codegen-errors.1: done
  - security-session.1: done
  - integration tests: done (servemux-cache, component-nested-if-else, session-encrypt)
  - docs + version bump: done
  - tags pushed: core/v0.0.22, cmd/dreego/v0.0.22, plugins/sample/v0.0.22
- v0.0.23 (unreleased) late fixes: done
  - core/runtime.go: exported Reset() added
  - core/runtime_test.go: uses Reset(); TestResetClearsCache is behavioral
  - core/session_keys.go: deriveKeys now uses HMAC-SHA256(secret, label)
  - core/session.go: sign propagates json.Marshal and encryptPayload errors
  - core/session_crypto.go: encryptPayload returns (ciphertext, error) and accepts io.Reader for nonce generation
  - core/session_encrypt_test.go: shortReader returns error; removed testEncryptReader hook
  - coding-standards: line limit raised 120 → 300 in AGENTS.md and .agents/guides/coding-standards.md; ADR line-limit-300.md created
- final test runs: core tests green, integration 147 passed, 0 failed
