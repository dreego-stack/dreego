
---
type: Concept
title: "Dreego Ecosystem (Affons)"
description: "Tools, services, and community resources around the Dreego Core"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Dreego Ecosystem (Affons)

## Definition

The Dreego Ecosystem ("Affons") encompasses all tools, services, and community resources around the Dreego Core. It is what turns Dreego from a tool into a platform.

## Ecosystem Components

### 1. CLI Tooling
| Component            | Description                                         | Priority |
|----------------------|-----------------------------------------------------|----------|
| `dreego`              | Main CLI: generate, dev, build, new                 | V1       |
| `dreego dev`          | Dev server with hot reload, Tailwind JIT             | V1       |
| `dreego generate`     | Transpiler: .dreego → Go code                        | V1       |
| `dreego build`        | Build single binary                                 | V1       |
| `dreego new`          | Scaffold new project                                | V1       |
| `dreego add`          | Install addons (wrap `go get`)                      | V2       |
| `dreego upgrade`      | Upgrade framework version                           | V2       |

### 2. Addon Registry
| Component                  | Description                                     | Priority |
|----------------------------|-------------------------------------------------|----------|
| registry.dreego.dev         | Central addon directory                          | V2       |
| `dreego add search`         | CLI search for addons                            | V2       |
| Addon Scaffolding          | `dreego add create` — create new addon           | V2       |
| Addon Docs                 | Automatically generated API docs                 | V2       |

### 3. Documentation
| Component           | Description                                         | Priority |
|---------------------|-----------------------------------------------------|----------|
| docs.dreego.dev      | Official documentation                               | V1       |
| Tutorial            | "Build a Blog with Dreego" — Getting Started          | V1       |
| API Reference       | Go doc for core packages                             | V1       |
| Migration Guide     | From Next.js / SvelteKit to Dreego                    | V2       |
| Examples Repo       | github.com/dreego-ecosystem/examples                  | V1       |

### 4. Developer Experience
| Component           | Description                                         | Priority |
|---------------------|-----------------------------------------------------|----------|
| VS Code Extension   | Syntax highlighting for .dreego                      | V2       |
| LSP Server          | Language Server for .dreego (Autocomplete, Diagnostics)| V2      |
| Error Overlay       | Browser overlay for compile errors (dev server)      | V1       |
| Playground          | playground.dreego.dev — test .dreego in the browser   | V2       |

### 5. Templates & Starter
| Component           | Description                                         | Priority |
|---------------------|-----------------------------------------------------|----------|
| dreego-starter       | Minimal starter project                              | V1       |
| dreego-blog          | Blog starter with Markdown                           | V2       |
| dreego-saas          | SaaS starter (Auth, Payments, DB)                    | V2       |
| dreego-admin         | Admin dashboard starter                              | V2       |

### 6. Deployment & Hosting
| Component             | Description                                       | Priority |
|-----------------------|----------------------------------------------------|----------|
| Dockerfile Template    | Multi-Stage Build (test → build → deploy)          | V1       |
| Fly.io Guide           | Deployment on Fly.io                               | V1       |
| VPS Guide              | Deployment on Hetzner/Vultr/etc.                   | V1       |
| dreego-cloud (later)   | Managed hosting for Dreego (optional)               | V3       |

### 7. Testing
| Component           | Description                                         | Priority |
|---------------------|-----------------------------------------------------|----------|
| `dreego test`        | Test runner for .dreego files                        | V2       |
| Test Helpers        | `dreegotest` package — Request simulation, etc.      | V2       |
| E2E Testing         | Playwright integration guide                         | V2       |

### 8. Monitoring & Observability
| Component           | Description                                         | Priority |
|---------------------|-----------------------------------------------------|----------|
| Built-in Logging    | Structured logging (slog)                            | V1       |
| Metrics             | Prometheus endpoint (optional)                       | V2       |
| Tracing             | OpenTelemetry integration                            | V3       |

### 9. Community
| Component             | Description                                       | Priority |
|-----------------------|----------------------------------------------------|----------|
| GitHub Discussions    | Community questions & support                      | V1       |
| Discord Server        | Real-time community                                | V1       |
| Showcase              | dreego.dev/showcase — show projects                 | V2       |
| Newsletter            | Release notes & tips                               | V2       |

## V1 Minimal Ecosystem

For V1, Dreego only needs the essentials:

1. **`dreego` CLI** — generate, dev, build, new
2. **docs.dreego.dev** — Getting Started + API Reference
3. **Starter Template** — `dreego new` scaffolded
4. **GitHub Repo** — README, Contributing Guide
5. **Error Overlay** — in the dev server

Everything else comes in V2+, once the core is stable.

## Addon Development Experience

### Create an Addon (V2)
```bash
dreego add create my-addon
# Scaffolded:
# my-addon/
# ├── dreego.go          # Plugin interface implementation
# ├── assets/           # //go:embed resources
# ├── go.mod
# └── README.md
```

### Publish an Addon (V2)
```bash
git push origin main
git tag v1.0.0
# Automatically visible in the registry
```

### Use an Addon (V1)
```bash
go get github.com/dreego-ecosystem/dreego-auth
```

```go
// main.go
import "github.com/dreego-ecosystem/dreego-auth"

app.UsePlugin(auth.New("secret-key"))
```
