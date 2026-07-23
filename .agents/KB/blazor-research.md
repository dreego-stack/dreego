# C# Blazor & ASP.NET Core Research — Dreego-Relevanz

**Datum:** 23.07.2026
**Zweck:** Systematische Analyse der C#/.NET Webframework-Welt. C# ist wie Go kompiliert und statisch typisiert — die architektonischen Lehren sind direkt übertragbar.

## Kernergebnisse

### Was Dreego von Blazor lernen kann

Blazor's Rendering-Modes (`Static SSR`, `Interactive Server`, `Interactive WebAssembly`, `Interactive Auto`) sind ein direktes Vorbild für Dreego. Dreego macht SSR-First was Blazor "Static SSR" nennt, und HTMX/Alpine.js entsprechen "Interactive Server" (Server-seitige Events, DOM-Updates ohne Full Reload).

### Middleware-Pipeline (ASP.NET Core)
Eine der saubersten Middleware-Architekturen:
- `Use()` — chainable, `await next(context)` für pre/post
- `Run()` — terminal, short-circuits
- `Map()` — branch by path prefix
- `UseWhen()` — conditional branch, rejoins main pipeline
- Dokumentierte Reihenfolge: Exception → HTTPS → Static → Routing → CORS → Auth → Endpoints

### Configuration System
Hierarchische Key-Value Pairs mit klarer Priorität:
1. CLI-Args (höchste)
2. Environment Variables
3. User Secrets (Dev only)
4. `appsettings.{Env}.json`
5. `appsettings.json` (niedrigste)

Plus `reloadOnChange` für Live-Config-Updates.

### Form-Handling
- `EditForm` + `DataAnnotationsValidator` — deklarativ
- `OnValidSubmit`/`OnInvalidSubmit` — clean lifecycle
- `ValidationSummary` + `ValidationMessage` — component-based
- `[SupplyParameterFromForm]` — auto-map POST body zu Properties
- Antiforgery: Auto-added, opt-out möglich

### Dependency Injection
Built-in, kein externes Framework:
- Transient (new per inject), Scoped (per request/circuit), Singleton
- Constructor Injection (preferred)
- Keyed Services (.NET 8+)
- `IOptions<T>` Pattern — typed config

### Component Model
- Single-file: `.razor` = Markup + `@code` = Logic
- `[Parameter]`, `[CascadingParameter]` für Props
- `RenderFragment` für Children/Slots
- `@typeparam` für Generics
- Lifecycle: OnInit, OnParamsSet, OnAfterRender, ShouldRender, Dispose
- `@key` für Identity in Lists
- `@attributes` für Attribute Splatting
- `@rendermode` für Mix von Static + Interactive

### Hot Reload
Änderungen an C#, Razor, CSS ohne App-Neustart. State bleibt erhalten.

### Scaffolding
`dotnet new blazor` — interaktiv (Server/WebAssembly/Auto/None)
ASP.NET Core Identity Scaffolder: Login, Register, Manage Pages
Entity Framework Scaffolder: Models aus existierender DB

## Für Dreego: Direkt Übertragbar

| .NET Feature | Dreego-Adaption |
|---|---|
| Convention over Configuration | Ordnerstruktur = Routing, Namespace = Dateipfad |
| `_Imports.razor` | `_imports.dreego` pro Ordner |
| Middleware-Ordnung + Dokumentation | Go-Middleware-Chain mit klarer Reihenfolge |
| Configuration-Priority-System | CLI > Env > `.dreego.yaml` > Defaults |
| `EditForm`-Pattern | `g-submit` + Struct-Validierung |
| Component Lifecycle | `<go>`-Block = OnInit, keine weiteren Lifecycle-Phases nötig |
| `@key` | Nicht nötig (HTMX macht DOM-Replacement) |
| DI: Keyed Services | Go: Map-basierte Registrierung |
| `IOptions<T>` | Go-Struct mit `dreego:"config"` Tags |

*Siehe thinking-list.md für die vollständige Detail-Liste der Blazor-Features.*
