
---
type: Reference
title: C# Blazor & ASP.NET Core Research — Dreego Relevance
description: Systematic analysis of C# Blazor & ASP.NET Core features transferable to Dreego
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# C# Blazor & ASP.NET Core Research — Dreego Relevance

**Date:** 2026-07-28
**Purpose:** Systematic analysis of the C#/.NET web framework world. C# is like Go — compiled and statically typed — the architectural lessons are directly transferable.

## Key Findings

### What Dreego Can Learn from Blazor

Blazor's rendering modes (`Static SSR`, `Interactive Server`, `Interactive WebAssembly`, `Interactive Auto`) are a direct model for Dreego. Dreego does SSR-First, which Blazor calls "Static SSR", and HTMX/Alpine.js correspond to "Interactive Server" (server-side events, DOM updates without full reload).

### Middleware Pipeline (ASP.NET Core)
One of the cleanest middleware architectures:
- `Use()` — chainable, `await next(context)` for pre/post
- `Run()` — terminal, short-circuits
- `Map()` — branch by path prefix
- `UseWhen()` — conditional branch, rejoins main pipeline
- Documented ordering: Exception → HTTPS → Static → Routing → CORS → Auth → Endpoints

### Configuration System
Hierarchical key-value pairs with clear priority:
1. CLI args (highest)
2. Environment Variables
3. User Secrets (Dev only)
4. `appsettings.{Env}.json`
5. `appsettings.json` (lowest)

Plus `reloadOnChange` for live config updates.

### Form Handling
- `EditForm` + `DataAnnotationsValidator` — declarative
- `OnValidSubmit`/`OnInvalidSubmit` — clean lifecycle
- `ValidationSummary` + `ValidationMessage` — component-based
- `[SupplyParameterFromForm]` — auto-map POST body to properties
- Antiforgery: Auto-added, opt-out possible

### Dependency Injection
Built-in, no external framework:
- Transient (new per inject), Scoped (per request/circuit), Singleton
- Constructor Injection (preferred)
- Keyed Services (.NET 8+)
- `IOptions<T>` Pattern — typed config

### Component Model
- Single-file: `.razor` = Markup + `@code` = Logic
- `[Parameter]`, `[CascadingParameter]` for Props
- `RenderFragment` for Children/Slots
- `@typeparam` for Generics
- Lifecycle: OnInit, OnParamsSet, OnAfterRender, ShouldRender, Dispose
- `@key` for Identity in Lists
- `@attributes` for Attribute Splatting
- `@rendermode` for mixing Static + Interactive

### Hot Reload
Changes to C#, Razor, CSS without app restart. State is preserved.

### Scaffolding
`dotnet new blazor` — interactive (Server/WebAssembly/Auto/None)
ASP.NET Core Identity Scaffolder: Login, Register, Manage Pages
Entity Framework Scaffolder: Models from existing DB

## For Dreego: Directly Transferable

| .NET Feature | Dreego Adaptation |
|---|---|
| Convention over Configuration | Directory structure = Routing, Namespace = file path |
| `_Imports.razor` | `_imports.dreego` per directory |
| Middleware order + documentation | Go middleware chain with clear ordering |
| Configuration priority system | CLI > Env > `.dreego.yaml` > Defaults |
| `EditForm` pattern | `g-submit` + struct validation |
| Component Lifecycle | `<go>` block = OnInit, no further lifecycle phases needed |
| `@key` | Not needed (HTMX does DOM replacement) |
| DI: Keyed Services | Go: map-based registration |
| `IOptions<T>` | Go struct with `dreego:"config"` tags |

*See thinking-list.md for the full detailed list of Blazor features.*
