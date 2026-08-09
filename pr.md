---
version: none
---

- Feat: auto-prepare on PR — pull_request.yml applies pr.md (changelog + version) and removes pr.md on the PR branch after checks pass
- Fix: drop manual release-prep workflow (forgettable, left pr.md in main)
- Feat: release.yml tags on push to main after merge (no bot commit on main, no bypass needed)
- Chore: branch protection requires 0 reviews — merge button only, no approve needed
