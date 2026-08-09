---
version: patch
---

- Chore: rewrite CHANGELOG.md to the new line-based standard (Feat/Fix/Chore prefix per line, no section titles)
- Fix: run-timer-sigterm test — 10x faster (timer 30s → 3s), uses go run per how-to standard, no more flaky port race