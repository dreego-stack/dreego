---
version: patch
---

- Fix: propagate directory-walk errors in generate check, layout discovery, and static-asset scan (no silent skip)
- Fix: make `dreego generate --force` actually force a full rewrite by bypassing the up-to-date check
- Fix: reject a second concurrent `App.Listen` instead of silently corrupting server state
- Fix: `CookieStore.Destroy` clears the in-request session map so a post-logout write cannot resurrect stale state
- Fix: CSRF middleware no longer sets the readable token cookie when the session write failed
- Docs: state that request-body limits are the application's responsibility (no implicit cap)
- Docs: document rewrite ordering relative to `app.Use` middleware (original path is seen)
- Docs: remove speculative Mailer/Cache plugin interfaces that contradict the no-speculative-exports rule
- Docs: complete the v0.1 compatibility list with `ServerConfig`, `DefaultServerConfig`, and error sentinels
