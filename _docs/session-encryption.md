# Session Encryption

Dreego sessions are HMAC-signed by default. You can also encrypt the session payload with AES-256-GCM so cookie contents are not readable by the client.

## Enable Encryption

Pass `Encrypt: true` in session options:

Use `core.Options` directly with the underlying store:

```go
store := core.NewCookieStore(secret)
store.Set(w, r, "user_id", "42", &core.Options{Encrypt: true})
```

`c.SetSessionVal` does not accept options; call `store.Set` directly when you need encryption.

Encryption applies to the whole session payload as a single JSON blob, not to individual keys.

## How It Works

- `NewCookieStore(secret)` derives two 256-bit keys from the secret via HMAC-SHA256:
  - `dreego-session-sig` — HMAC signing key
  - `dreego-session-enc` — AES-GCM encryption key
- When `Encrypt: true`, the JSON payload is encrypted first, then HMAC-signed (`encrypt-then-MAC`).
- A one-byte marker distinguishes encrypted cookies from signed-only cookies, so both formats can coexist during key rotation.
- `CookieStore.Set` returns an error if JSON marshaling or encryption fails (for example, when nonce generation fails).
- Decryption verifies the HMAC first; any tampering or wrong key rejects the cookie and returns an empty session.

## Security Notes

- The secret must be at least 32 bytes for adequate entropy.
- Encrypted cookies are still readable in size and timing; do not store large secrets in sessions.
- Encryption only protects the cookie contents. Use `Secure: true` (default is TLS-aware) and `HttpOnly: true` (default for session cookies) to protect transport and JavaScript access.

## Defaults

Session cookies use secure defaults:

- `HttpOnly: true`
- `Secure: true` when the request uses TLS
- `Path: "/"`

Pass `&core.Options{Encrypt: true, Secure: true, HttpOnly: true}` to override any of these per call.
