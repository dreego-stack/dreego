---
version: patch
---

- Chore: migrate remaining ~120 shell tests (`_tests/core/{Components,Template,Imports,Static,Session,Routing,Layout,ContentType,FormActions,Middleware,Config,Deployment,CLI,Bugs}`) to Go in `_tests/go/` via `dreegotest`
- Feat: extend `dreegotest` with CLI/project helpers (`CLIBin`, `ProjectDir`, `RunCLI`, `BuildInDirOK`, `LatestTag`) and a cookie-jar HTTP client (`ServeSetup`, `Client.Request/Cookie`)
