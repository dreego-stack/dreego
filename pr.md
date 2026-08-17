---
version: patch
---

- Feat: one canonical quick-start path in README and _docs/getting-started.md (dreego new → generate → go run .)
- Feat: dreego new validates the project name and checks for the go executable before scaffolding
- Feat: dreego init validates the target path and checks for the go executable
- Feat: scaffolded main.go reads DREEGO_PORT so the documented server can run on a free port
- Feat: go mod tidy failure now prints the underlying command output and a DREEGO_LOCAL_REPO hint
- Test: black-box quick_start_test.go exercises scaffold + generate + build + serve + HTTP request
- Test: invalid project names, missing go, existing target, and no-arg are asserted
- Docs: README and _docs/getting-started.md no longer reference a repo-local replace directive or outdated module path