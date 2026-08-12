---
version: patch
---

- Bug: fix flaky `TestCookieStoreEncryptValueNotPlaintext` — the byte-level substring check for `user_id`/`42` in the decoded cookie value was statistically unsound: the random AES-GCM ciphertext (base64-encoded) can coincidentally contain those bytes (~1.5% per run), causing intermittent `make test` failures. The test now asserts the encryption marker and checks the decrypted payload instead, making it deterministic (verified 300/300 green).
