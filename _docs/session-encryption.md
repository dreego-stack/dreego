# Session Encryption & Security

Dreego sessions are HMAC-signed by default. You can also encrypt the session payload with AES-256-GCM so cookie contents are not readable by the client.

## Enable Encryption

Pass `Encrypt: true` in session options:

```go
store := dreego.NewCookieStore(secret)
store.Set(w, r, "user_id", "42", &dreego.Options{Encrypt: true})
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
- Decryption verifies the HMAC first; any tampering or wrong key rejects the cookie and returns an error from `Get`.

## Secret Requirements

The secret must be at least 32 bytes. `NewCookieStore` panics during construction when the secret is empty or shorter than 32 bytes. `App.Build()` also validates the store via the `Validate()` method and panics when an unsafe secret reaches startup.

## Cookie Defaults

Session cookies use secure defaults:

- `HttpOnly: true`
- `SameSite: Lax` (documented default)
- `Secure: true` when the request uses TLS
- `Path: "/"`

Pass `&dreego.Options{Encrypt: true, Secure: true, HttpOnly: true}` to override per call. A configurable `CookiePolicy` on the store can set app-wide defaults that cannot be silently downgraded: `Set`, `Delete`, `Destroy`, and CSRF writes all preserve the policy and cannot drop `Secure`, `HttpOnly`, `SameSite`, or encryption. `SetCookiePolicy` merges partial policies with the secure defaults — for example, passing only `SameSite: Strict` keeps `HttpOnly: true` and `Path: "/"`.

Cookie paths are app-wide policy rather than a per-call option. Set
`CookiePolicy.Path` once. Passing a different `Options.Path` to `Set` returns
`ErrCookiePathOverride`, preventing cookies that `Delete` or `Destroy` cannot
expire at the same path.

## Trusted TLS Proxies

Behind a TLS-terminating reverse proxy, set `Secure: true` only when the request arrives from an explicitly trusted proxy. Configure trusted proxy addresses on the store:

```go
store := dreego.NewCookieStore(secret)
store.SetTrustedProxies([]string{"10.0.0.1"})
```

When `X-Forwarded-Proto: https` arrives from a trusted proxy, the cookie is marked `Secure`. Requests from untrusted addresses, or when no proxies are configured, ignore the forwarded header and fall back to direct TLS detection (`r.TLS != nil`). This prevents arbitrary clients from forcing `Secure` cookies over plain HTTP.

CSRF cookies follow the same policy: `CSRF(store)` reads the store's trusted-proxy list and applies the same `Secure` decision as the session cookie. The readable CSRF cookie keeps `SameSite: Strict`; the token itself lives in the (HttpOnly) session cookie, so the readable cookie never needs cross-site sends.

## Size Limit

Cookie sessions have a documented size limit of 4096 bytes for the encoded cookie value. `Set` returns `ErrSessionTooLarge` when the encoded state exceeds the limit. Use a server-side store plugin for larger state.

## Observable Store Failures

`CookieStore.Get` returns an error when the cookie fails integrity verification (tampered, wrong key, corrupt JSON, invalid base64). Write paths recover instead of failing: `Set` (and `Delete`, which delegates to it) treats a failed verification as an empty session and overwrites the invalid cookie, so login and CSRF keep working after key rotation. `SetSessionVal`, `DelSessionVal`, and `DestroySession` surface genuine store errors as HTTP 500 with a generic body; the internal cause is logged server-side. Handlers must stop rendering after a failed session write, because the framework cannot interrupt them. Custom stores that implement `Get`/`Set` with errors get the same treatment.

## Authentication-State Invalidation

Auth plugins replace or invalidate pre-authentication state by calling `Set` and `Delete` on the same store:

```go
func login(w http.ResponseWriter, r *http.Request, store dreego.Store) {
    store.Set(w, r, "user_id", user.ID, nil)
    store.Delete(w, r, "pending_oauth_state")
}
```

`Destroy` invalidates the complete authentication state: it clears every key and sends a single expired cookie that preserves `Secure`, `HttpOnly`, `SameSite`, and `Path` from the active policy. Logout must call `Destroy`, not `Delete` on one key.

The cookie store does not rotate a server-side identifier because it has none. A stolen cookie can be replayed until it expires or the signing secret rotates. To invalidate a session immediately, rotate the secret or use a server-side store plugin that tracks identifiers.

## Key Rotation

Rotate the signing secret by deploying a new `NewCookieStore` with the new secret. Cookies signed or encrypted with the old secret fail verification: `Get` returns an error, and any write starts a fresh session with the old keys gone, so the user simply re-authenticates. There is no automatic dual-secret grace period — deploy rotation during a maintenance window or accept a transient re-login wave.

## Replaceable Store

Applications can replace the built-in cookie store with a custom implementation via `app.SetSessionStore(store)`. The custom store must implement the `Store` interface. When it also implements `Validate() error`, `App.Build()` calls it at startup so unsafe configuration fails early. Replacing the store does not disable CSRF or cookie-policy enforcement — `CSRF(store)` and the session middleware continue to use the new store.

## Security Notes

- Encrypted cookies are still readable in size and timing; do not store large secrets in sessions.
- Encryption only protects the cookie contents. Use `Secure: true` and `HttpOnly: true` to protect transport and JavaScript access.
- The cookie store is stateless. Replay protection and immediate invalidation require a server-side store.
