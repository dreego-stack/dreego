---
version: patch
---

- Fix: make the module boundary test version-agnostic so coordinated module releases can be tagged again after `release-prep.py` rewrites the internal versions
- Test: assert one coordinated `vX.Y.Z` version across every internal `go.mod` requirement and `go.work` replacement instead of a hardcoded version literal
