# v0.0.22-triple Remaining Fixes Plan

Source: reviewer findings + failing tests written by test-engineer.

## Fixes

1. core/runtime.go
   - Add exported `Reset()` that clears `builtHandler` (and only that).
   - Keep `runtime_test.go` using `Reset()` instead of direct `builtHandler = nil`.

2. core/session_keys.go
   - Change `deriveKeys` from SHA-256(secret $%$$%$ label) to HMAC-SHA256(secret, label) for both encryption and signing keys.

3. core/session_crypto.go + core/session.go
   - Change `encryptPayload(secret, plain string) string` to `(string, error)`.
   - Change `sign(secret, value string) string` to `(string, error)`.
   - Propagate errors through `CookieStore.Set`.
   - Remove or keep the test hook `testEncryptReader` only if needed by the failing test.

4. Verification
   - `smd go test ./core/... -count=1` must pass.
   - `smd sh _tests/test.sh` must pass.

5. Docs
   - Update CHANGELOG.md / .agents/log.md if behavior changed meaningfully.

6. Git
   - Diff review + amend commit or new commit.
