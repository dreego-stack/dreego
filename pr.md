---
version: patch
---

- Feat: CI runs core and integration tests with -race on the ubuntu runner (cgo)
- Fix: protect CookieStore policy and trustedProxies with sync.RWMutex; route isSecureForCSRF through locked accessor
- Test: add concurrent-access tests for routes, app config, session store, cookie policy, CSRF middleware, ready handler