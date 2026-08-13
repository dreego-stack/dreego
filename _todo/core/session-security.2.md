---
area: security
phase: v0.1-blocker
depends_on: session-security.1
---
# Secure session and cookie defaults

## Goal
Make sessions secure by default while allowing applications and auth plugins to replace storage through the core interface.

## Acceptance criteria
- Session cookies use a documented `SameSite` default.
- Empty, short, or otherwise unsafe signing and encryption secrets fail during construction or startup.
- Secure-cookie behavior works behind explicitly trusted TLS-terminating proxies without trusting arbitrary forwarded headers.
- Session-store read and write failures are observable and cannot silently report success.
- CSRF session writes use the same secure cookie policy.
- Set, delete, destroy, and CSRF operations preserve one app- or store-level cookie policy and cannot silently drop `Secure`, `HttpOnly`, `SameSite`, or encryption.
- Cookie sessions have a documented size limit and an observable error when encoded state is too large.
- Applications can configure cookie policy and replace the store without disabling core security accidentally.
- Unit and HTTP-level regression tests cover direct TLS, trusted proxy, untrusted proxy, and store-failure cases.
