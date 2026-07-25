
---
type: Concept
title: "Dreego Ecosystem (Affons)"
description: "Tools, Services und Community-Ressourcen rund um den Dreego Core"
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Dreego Ecosystem (Affons)

## Definition

Das Dreego-Ecosystem ("Affons") umfasst alle Tools, Services und Community-Ressourcen rund um den Dreego Core. Es ist das, was Dreego von einem Tool zu einer Plattform macht.

## Ecosystem-Komponenten

### 1. CLI Tooling
| Komponente          | Beschreibung                                         | Priorität |
|---------------------|-----------------------------------------------------|-----------|
| `dreego`             | Haupt-CLI: generate, dev, build, new                | V1        |
| `dreego dev`         | Dev-Server mit Hot Reload, Tailwind JIT              | V1        |
| `dreego generate`    | Transpiler: .dreego → Go-Code                         | V1        |
| `dreego build`       | Single Binary bauen                                  | V1        |
| `dreego new`         | Neues Projekt scaffolden                             | V1        |
| `dreego add`         | Addons installieren (wrap `go get`)                  | V2        |
| `dreego upgrade`     | Framework-Version upgraden                           | V2        |

### 2. Addon Registry
| Komponente                | Beschreibung                                     | Priorität |
|---------------------------|--------------------------------------------------|-----------|
| registry.dreego.dev        | Zentrales Addon-Verzeichnis                       | V2        |
| `dreego add search`        | CLI-Suche nach Addons                             | V2        |
| Addon-Scaffolding         | `dreego add create` — neues Addon erstellen        | V2        |
| Addon-Docs                | Automatisch generierte API-Docs                   | V2        |

### 3. Dokumentation
| Komponente          | Beschreibung                                         | Priorität |
|---------------------|-----------------------------------------------------|-----------|
| docs.dreego.dev      | Offizielle Dokumentation                             | V1        |
| Tutorial            | "Build a Blog with Dreego" — Getting Started          | V1        |
| API Reference       | Go-Doc für Core-Packages                             | V1        |
| Migration Guide     | Von Next.js / SvelteKit zu Dreego                     | V2        |
| Examples Repo       | github.com/dreego-ecosystem/examples                  | V1        |

### 4. Developer Experience
| Komponente          | Beschreibung                                         | Priorität |
|---------------------|-----------------------------------------------------|-----------|
| VS Code Extension   | Syntax Highlighting für .dreego                       | V2        |
| LSP Server          | Language Server für .dreego (Autocomplete, Diagnostics)| V2        |
| Error Overlay       | Browser-Overlay bei Compile-Errors (Dev-Server)      | V1        |
| Playground          | playground.dreego.dev — .dreego im Browser testen      | V2        |

### 5. Templates & Starter
| Komponente          | Beschreibung                                         | Priorität |
|---------------------|-----------------------------------------------------|-----------|
| dreego-starter       | Minimales Starter-Projekt                            | V1        |
| dreego-blog          | Blog-Starter mit Markdown                            | V2        |
| dreego-saas          | SaaS-Starter (Auth, Payments, DB)                    | V2        |
| dreego-admin         | Admin-Dashboard-Starter                              | V2        |

### 6. Deployment & Hosting
| Komponente            | Beschreibung                                       | Priorität |
|-----------------------|----------------------------------------------------|-----------|
| Dockerfile Template    | Multi-Stage Build (test → build → deploy)          | V1        |
| Fly.io Guide           | Deployment auf Fly.io                              | V1        |
| VPS Guide              | Deployment auf Hetzner/Vultr/etc.                  | V1        |
| dreego-cloud (später)   | Managed Hosting für Dreego (optional)                | V3        |

### 7. Testing
| Komponente          | Beschreibung                                         | Priorität |
|---------------------|-----------------------------------------------------|-----------|
| `dreego test`        | Test-Runner für .dreego-Dateien                       | V2        |
| Test Helpers        | `dreegotest` Package — Request-Simulation, etc.       | V2        |
| E2E Testing         | Playwright-Integration-Guide                         | V2        |

### 8. Monitoring & Observability
| Komponente          | Beschreibung                                         | Priorität |
|---------------------|-----------------------------------------------------|-----------|
| Built-in Logging    | Strukturiertes Logging (slog)                        | V1        |
| Metrics             | Prometheus-Endpoint (optional)                       | V2        |
| Tracing             | OpenTelemetry-Integration                            | V3        |

### 9. Community
| Komponente            | Beschreibung                                       | Priorität |
|-----------------------|----------------------------------------------------|-----------|
| GitHub Discussions    | Community-Fragen & Support                         | V1        |
| Discord Server        | Echtzeit-Community                                 | V1        |
| Showcase              | dreego.dev/showcase — Projekte zeigen               | V2        |
| Newsletter            | Release-Notes & Tipps                              | V2        |

## V1 Minimal Ecosystem

Für V1 braucht Dreego nur das Nötigste:

1. **`dreego` CLI** — generate, dev, build, new
2. **docs.dreego.dev** — Getting Started + API Reference
3. **Starter-Template** — `dreego new` scaffolded
4. **GitHub Repo** — README, Contributing Guide
5. **Error Overlay** — im Dev-Server

Alles andere kommt in V2+, wenn der Core stabil ist.

## Addon-Entwicklungs-Experience

### Ein Addon erstellen (V2)
```bash
dreego add create my-addon
# Scaffolded:
# my-addon/
# ├── dreego.go          # Plugin-Interface Implementierung
# ├── assets/           # //go:embed Ressourcen
# ├── go.mod
# └── README.md
```

### Addon veröffentlichen (V2)
```bash
git push origin main
git tag v1.0.0
# Automatisch im Registry sichtbar
```

### Addon nutzen (V1)
```bash
go get github.com/dreego-ecosystem/dreego-auth
```

```go
// main.go
import "github.com/dreego-ecosystem/dreego-auth"

app.UsePlugin(auth.New("secret-key"))
```
