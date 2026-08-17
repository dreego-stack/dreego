---
version: patch
---

- Feat: session cookies use a documented SameSite=Lax default and a configurable CookiePolicy that Set/Delete/Destroy/CSRF cannot silently downgrade (Secure, HttpOnly, SameSite, encryption)
- Feat: empty, short, or otherwise unsafe session secrets fail at construction (NewCookieStore) and at startup (App.Build validates the store)
- Feat: Secure-cookie behavior works behind explicitly trusted TLS-terminating proxies (SetTrustedProxies) without trusting arbitrary X-Forwarded-* headers
- Feat: session-store read/write failures are observable: CookieStore.Get returns errors for invalid/tampered/rotated cookies, SetSessionVal/DelSessionVal/DestroySession surface store errors as HTTP 500
- Fix: session writes (Set/Delete/CSRF) recover from tampered or rotated cookies by starting a fresh session and overwriting the invalid cookie, so login and CSRF keep working after key rotation
- Fix: session write failures return a generic 500 body and log the internal cause server-side (no error leakage to clients)
- Fix: SetCookiePolicy merges partial policies with the secure defaults, so a zero-value policy no longer silently drops HttpOnly
- Fix: CSRF cookie keeps SameSite=Strict (no silent downgrade to Lax)
- Feat: cookie sessions enforce a documented 4096-byte size limit with observable ErrSessionTooLarge
- Feat: auth-state invalidation documented for cookie sessions: login/privilege changes replace pre-auth state, Destroy invalidates the complete auth state, key rotation and replay limitations documented
- Fix: per-call Options no longer drop the app-configured SameSite policy (e.g. Strict) for SetSessionVal and CSRF session writes
