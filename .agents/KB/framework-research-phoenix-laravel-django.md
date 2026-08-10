
---
type: Reference
title: Framework Comparison: Phoenix, Laravel, Django — Relevance for Dreego
description: Systematic comparison of Phoenix, Laravel, and Django — adoptable features and anti-patterns for Dreego
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Framework Comparison: Phoenix, Laravel, Django — Relevance for Dreego

**Date:** 2026-07-28
**Purpose:** Systematic analysis of which features of the three established frameworks Dreego should adopt, adapt, or deliberately avoid.

---

## 1. Phoenix (Elixir)

### 1.1 LiveView — Architecture (WebSockets)

LiveView works via **WebSockets** (not SSE). Architecture:

1. Initial HTTP request renders static HTML (SSR).
2. Client establishes WebSocket connection to the server.
3. Server spawns a **stateful Elixir process** (GenServer-like) per client.
4. Each state change on the server computes a **minimal diff** and pushes only the changed HTML over the WebSocket to the client.
5. Client-side, a small JS selectively replaces the DOM (morphdom-like).
6. On connection loss: **graceful reconnection** with state recovery.

**Important Details:**
- `mount/3` is called **twice**: once for static HTML, once on WebSocket connect. `connected?(socket)` distinguishes the modes.
- `assign_async/3` — load asynchronous data in tasks, with `AsyncResult` for Loading/Failed/OK states.
- `start_async/3` + `handle_async/3` — low-level async control.
- `stream/4` — efficient list management (Insert/Delete/Update without full re-render).
- **Hibernation:** LiveView can compress its process state after 15s of inactivity (`hibernate_after`).
- **Flash Messages:** Built-in `put_flash/3`, `clear_flash/1` — also work over WebSockets.
- **File Uploads:** `allow_upload/3`, `consume_uploaded_entries/3` — fully integrated with progress tracking.
- **Lifecycle Hooks:** `attach_hook/4` for `:handle_params`, `:handle_event`, `:handle_info`, `:after_render`.
- **LiveComponents:** Stateful sub-components with their own `mount`, `update`, `handle_event`.

#### For Dreego adopt / adapt:

| LiveView Feature | Dreego Adaptation | Priority |
|---|---|---|
| SSR → WebSocket Upgrade | HTMX partials + SSE via Datastar. No WebSocket requirement, HTMX can do SSE/WS per use case | ✅ V1 |
| Stateful Server Process | Go goroutine per session (lightweight). Model: `connected?()` check | ✅ V2 |
| `assign_async` | `<go>` block: async functions with `loading`/`error`/`ok` states in template | ✅ V2 |
| `stream/4` (efficient lists) | `{#each}` with SSE-based streaming updates | ✅ V3 |
| Graceful Reconnection | Datastar's built-in SSE Reconnect | ✅ V1 |
| Flash Messages | Built-in Flash in `dreego.Context` | ✅ V1 |
| File Uploads | `dreego-storage` plugin, model: `allow_upload/3` API | ✅ V2 |
| Lifecycle Hooks | `dreego.Plugin` Interface → Middleware hooks + transpiler hooks | ✅ V1 |
| LiveComponents | `.dreego` components with their own state (<go> per component) | ✅ V2 |

### 1.2 Phoenix PubSub

Phoenix PubSub is a **cluster-wide Publish/Subscribe system**:

```elixir
PubSub.subscribe(:my_pubsub, "user:123")
PubSub.broadcast(:my_pubsub, "user:123", {:user_update, %{name: "Shane"}})
```

- **Adapter Architecture:** PG2 (Distributed Erlang, Default), Redis, or custom adapter.
- **`broadcast_from/5`:** Broadcast to all except the sender.
- **`direct_broadcast/5`:** Broadcast to a specific node.
- **`local_broadcast/4`:** Local node only.
- **Custom Dispatcher:** "Fastlane" — for channels, messages are encoded once and written directly into all sockets instead of per channel.
- **Safe Pool-Size Migration:** `broadcast_pool_size` allows rolling updates without message loss.

#### For Dreego adopt:

- **Event Bus in Core:** `dreego.Emit("user:123", payload)` → all subscribers (SSE/WS) receive update.
- **Don't build yourself:** Either NATS embedded, Redis PubSub, or Go channel-based event bus.
- **Plugin:** `dreego-pubsub` with Redis/NATS backend.

### 1.3 Ecto — Why It Is Good

Ecto is **not an ORM** but a **Database Wrapper + Query Builder**:

1. **Explicit Separation:** Schema (mapping) ≠ Query (queries) ≠ Repo (database access) ≠ Changeset (validation).
2. **Changesets:** Data validation and casting as a **separate object**, not on the model. `Ecto.Changeset` holds the state (valid?, errors, changes).
3. **Composable Queries:** Queries are immutable and chainable:
   ```elixir
   query = from u in User, where: u.age > 18
   query = from u in query, where: u.active == true
   ```
4. **Multi:** Transactions with multiple operations `Ecto.Multi` — Rollback on error.
5. **Embedded Schemas:** DB-less schemas for form validation or API responses.
6. **Migrations:** Generated and versioned.
7. **No Magic:** No lazy loading, no implicit N+1. Everything explicit.

#### For Dreego adopt:

| Ecto Feature | Dreego Adaptation |
|---|---|
| Changesets | `dreego-db` plugin: Validation object independent of DB model |
| Composable Queries | Go: Query builder with chainable methods (similar to squirrel) |
| Multi (Transactions) | Go: `db.Transaction(func(tx *sql.Tx) error { ... })` |
| Embedded Schemas | Go structs with tags for validation without DB binding |
| `Ecto.Enum` | Go: `enum:"admin,user,mod"` struct tag |

### 1.4 Generators (`mix phx.gen.*`)

Phoenix has a tiered generator strategy:

| Generator | Creates |
|---|---|
| `phx.gen.embedded` | Only schema (no migration, no context) |
| `phx.gen.schema` | Schema + Migration |
| `phx.gen.context` | Schema + Migration + Context Module |
| `phx.gen.live` | Context + LiveView (CRUD) |
| `phx.gen.json` | Context + JSON API Controller + Views |
| `phx.gen.html` | Context + HTML Controller + Templates |

Templates can be **overridden per project** (`priv/templates/`).

#### For Dreego adopt:

- `dreego generate page` — create a `.dreego` file with `<go>` + Template (analogous to `phx.gen.live`)
- `dreego generate resource` — CRUD scaffolding (Route + Form + Validation)
- `dreego generate schema` — Go struct with validation tags
- **Template Overrides** via `.dreego/templates/`

### 1.5 Phoenix Template Engine (HEEx / EEx)

HEEx is Phoenix' **template engine with component model**:

- **`~H` Sigil:** Compile-time-checked HTML with components
- **Components:** `<.component_name attr="val">` → Function calls
- **Slots:** `<:slot_name>` for multi-slot components
- **Directives:** `:if`, `:for`, `:let`
- **HEEx Expressions:** `{@assign}` for values
- **HTML validation at compile time**

#### For Dreego adopt:

- `{#if}`, `{#each}`, `{#switch}`, `{#slot}` are already planned → ✅
- HEEx-like compile-time validation: `.dreego` → Transpiler checks for valid HTML → ✅ V2
- Component model: `<dreego:map />`, `<x-card>` (analogous to Blade Components) → ✅

---

### 1.6 What Makes Phoenix "On Par with JS Frameworks"

1. **LiveView replaces SPA JS:** No React/Vue needed for reactive UIs.
2. **Erlang/OTP Concurrency:** Millions of parallel connections.
3. **Soft-Realtime Built-in:** Channels, PubSub, Presence (who is online?).
4. **Fault-Tolerance:** Supervisor trees, automatic recovery.
5. **Single Binary Deployment:** Mix Releases.

### 1.7 Deployment (Releases)

- **Mix Releases:** Elixir's built-in release system → Single Binary (similar to Go!).
- **Phoenix Endpoint Configuration:** Everything in `config/runtime.exs` — environment-specific.
- **Docker-friendly:** Minimal Alpine image with release.

#### For Dreego: Go does single-binary deployment natively better than any other stack. This is one of the main advantages—absolutely emphasize and expand.

---

### 1.8 Anti-Patterns of Phoenix Dreego Should Avoid

| Anti-Pattern | Why |
|---|---|
| **Magic Macros:** `use Phoenix.LiveView` injects implicit callbacks | Dreego: Everything explicit in Go, no magic |
| **Too many conventions:** Phoenix is very opinionated | Dreego: Conventions yes (file-based routing), but override paths for everything |
| **Elixir learning curve:** Functional language, pattern matching, OTP | Dreego: Go is simpler, flatter learning curve |
| **WebSocket-First:** Doesn't work well behind some proxies | Dreego: SSE-First (more compatible, HTTP-based) |
| **LiveView State Bloat:** Large LiveViews can consume a lot of RAM | Dreego: Flat state in `<go>` — no nested processes per user |

---

## 2. Laravel (PHP)

### 2.1 Eloquent ORM — What Makes It Popular

Eloquent is an **ActiveRecord ORM**:

1. **Convention over Configuration:**
   - `Flight` Model → `flights` table (plural, snake_case)
   - `id` as primary key (auto-increment)
   - `created_at`, `updated_at` automatic
2. **Relationship API:**
   ```php
   $user->posts()->where('active', 1)->get();
   $post->user->name; // Dynamic Property
   ```
3. **Eager Loading:** `User::with('posts.comments')->get()` — prevents N+1.
4. **Mass Assignment Protection:** `$fillable` / `$guarded` — explicit, secure.
5. **Query Scopes:** `User::active()->verified()->get()`
6. **Accessors/Mutators:** `getNameAttribute()`, `setPasswordAttribute()`
7. **Model Events:** `creating`, `updated`, `deleting` — Observer pattern.
8. **Soft Deletes:** `deleted_at` timestamp, `withTrashed()`, `onlyTrashed()`.
9. **Casting:** `$casts` property for JSON, Array, Date, Enum, etc.
10. **Pruning:** `MassPrunable` trait — automatically delete old models.

#### For Dreego adopt:

| Eloquent Feature | Dreego Adaptation |
|---|---|
| Convention over Config | `dreego-db`: Go struct → table name via pluralizer |
| Relationships | `dreego-db`: Struct tags `dreego:"has_many:Posts"` → Preload system |
| Eager Loading / N+1 Prevention | `dreego-db`: `db.With("Posts.Comments").Find(&user)` |
| Mass Assignment (`$fillable`) | Go struct tags: `json:"name" db:"name" fillable:"true"` |
| Query Scopes | Go method chaining: `db.Scope(ActiveUsers).Find(&users)` |
| Model Events / Observers | Go: `AfterCreate`, `BeforeDelete` Interface |
| Soft Deletes | `deleted_at` nullable, `db.WithTrashed()` |
| Casting | Go struct tags with type mapping |
| `firstOrCreate`, `updateOrCreate` | Convenience methods in `dreego-db` |

### 2.2 Blade Templating — Features

Blade compiles to **plain PHP** and caches the results. Features:

1. **Template Inheritance:**
   ```blade
   @extends('layouts.app')
   @section('content')
   @endsection
   @yield('content')
   ```
2. **Components & Slots:**
   ```blade
   <x-alert type="error" :message="$message">
       <x-slot name="title">Error!</x-slot>
   </x-alert>
   ```
3. **Component Methods:** `$isSelected($value)` callable in template.
4. **Anonymous Components:** Only template file, no PHP class needed.
5. **Conditional Class Merging (`$attributes->merge()`, `@class([])`).**
6. **Stacks:** `@push('scripts')` / `@stack('scripts')` — push assets from child components into the layout head.
7. **Service Injection:** `@inject('metrics', 'App\Services\MetricsService')`.
8. **`@once` Directive:** Code is output only once per render cycle.
9. **Blade Fragments:** `@fragment` for AJAX partial loading (perfect with HTMX!).
10. **Custom Directives:** `Blade::directive('datetime', fn($exp) => ...)`.
11. **Custom If Statements:** `Blade::if('cloud', fn() => ...)`.
12. **`@verbatim`:** Leave JS template syntax (Vue, Alpine) untouched in Blade.
13. **`$loop` Variable:** Index, Iteration, First, Last, Even, Odd, Depth, Parent.

#### For Dreego adopt:

| Blade Feature | Dreego Adaptation |
|---|---|
| Components (`<x-alert>`) | `<dreego:alert>` or as plugin tags (already planned) |
| Conditional Classes | `<div class:active={isActive}>` — Svelte-Style (planned) |
| Stacks (`@push`/`@stack`) | `<slot name="head">` or `<head>` block CSS/JS Injection |
| `@once` | `{#once}` block → V2 |
| Blade Fragments + HTMX | HTMX `hx-select` + partial templates → great combination |
| `$loop` Variable | `{#each items as item, index}` → `{index}`, `{first}`, `{last}` Built-in → V1 |
| Template Inheritance | Layouts via `layout.dreego` (planned) |

### 2.3 Artisan CLI — Generator Overview

Artisan has **30+ make commands**:

```
make:model       make:controller    make:middleware
make:migration   make:seeder        make:factory
make:request     make:policy        make:command
make:job         make:notification  make:mail
make:event       make:listener      make:component
make:test        make:rule          make:observer
make:provider    make:cast          make:channel
make:enum        make:exception     make:resource
```

Each command creates a file in the right place with the right boilerplate.

**Additional Artisan Features:**
- **Tinker:** REPL with access to the entire framework (Eloquent, Jobs, Events).
- **Stub Customization:** `php artisan stub:publish` → Stubs overridable in the project.
- **Interactive Prompts:** `ask()`, `secret()`, `confirm()`, `anticipate()`, `choice()`.
- **Tables & Progress Bars:** `$this->table()`, `$this->withProgressBar()`.
- **Signal Handling:** `$this->trap(SIGTERM, ...)`.
- **Programmatic Execution:** `Artisan::call('mail:send', [...])` from code.
- **Queueing Commands:** `Artisan::queue('mail:send', [...])`.

#### For Dreego adopt:

- `dreego new` — Project scaffolding
- `dreego generate page` — New page
- `dreego generate resource` — CRUD
- `dreego generate middleware` — Middleware
- `dreego routes` — Show all routes
- `dreego tinker` — Go REPL (via yaegi or go-pry)
- **Stub Customization** → `.dreego/templates/`

### 2.4 Ecosystem: Forge, Vapor, Nova, Spark

| Product | Description | Dreego Equivalent |
|---|---|---|
| **Forge** | Server management (provisioning, deployment) | `dreego deploy` — via SSH or Docker |
| **Vapor** | Serverless Laravel on AWS Lambda | Not relevant for Go (Single Binary) |
| **Nova** | Admin panel generator | `dreego-admin` plugin |
| **Spark** | SaaS starter kit (Billing, Teams) | `dreego-saas` plugin |
| **Envoyer** | Zero-downtime deployment | `dreego deploy` with blue-green |
| **Horizon** | Queue monitoring dashboard | `dreego-jobs` dashboard |
| **Telescope** | Debugging & Monitoring | `dreego-devtools` |
| **Pennant** | Feature Flags | `dreego-features` plugin |
| **Pulse** | Performance Monitoring | Not V1 |
| **Reverb** | WebSocket server (first-party) | Not needed (SSE) |
| **Echo** | Client-side WebSocket library | Datastar (SSE) suffices |

**Takeaway:** Laravel's Ecosystem is the **greatest strength** of the framework. Dreego doesn't need to replicate this in V1, but the **plugin architecture must be open enough** that such tools can emerge as plugins.

### 2.5 Queue/Job System

Laravel's Queue system is extremely mature:

1. **Unified API across backends:** Database, Redis, SQS, Beanstalkd, MongoDB.
2. **Job classes** with `handle()` method.
3. **Dependency Injection** in `handle()` via Service Container.
4. **Eloquent Model Serialization:** Models are serialized as identifiers and reloaded during processing.
5. **Unique Jobs:** `ShouldBeUnique` — no duplicate in the queue.
6. **Job Middleware:** Rate Limiting, Overlap Prevention, Exception Throttling.
7. **Job Batching:** Execute multiple jobs as a batch, with callbacks (then/catch/finally).
8. **Job Chaining:** Execute jobs sequentially.
9. **Delayed Dispatching:** `->delay(now()->addMinutes(10))`.
10. **`dispatchAfterResponse()`:** Execute job only after HTTP response.
11. **Failed Jobs:** `failed_jobs` table, automatic retries, `retryUntil()`.
12. **Queue Priorities:** `--queue=high,default` — worker prioritizes.
13. **Horizon Dashboard:** Redis queue monitoring.

#### For Dreego adopt:

- `dreego-jobs` plugin with:
  - Interface: `type Job interface { Handle() error }`
  - Backends: Redis, PG, NATS
  - Delayed Jobs: `dreego.Dispatch(job).Delay(10 * time.Minute)`
  - Job Middleware (Rate Limiting, Retry)
  - Failed Job Logging
- Go goroutines are more native than PHP workers — a job system in Go is **significantly more performant**.

### 2.6 Notification System

Laravel Notifications:
- **Multiple Channels:** Mail, SMS (Vonage), Slack, Database, Broadcast, custom.
- **`Notification::send($users, new InvoicePaid($invoice))`**.
- **`Notifiable` Trait:** `$user->notify(new WelcomeNotification())`.
- **Database Notifications:** Stored in DB, retrievable via API.
- **Markdown Mail Templates:** Pre-formatted email components.

#### For Dreego:
- `dreego-notify` plugin with channel interface:
  ```go
  type Channel interface {
      Send(to User, notification Notification) error
  }
  ```

### 2.7 Validation Approach

Laravel's Validation is the **best of all frameworks**:

1. **Request Validation:** `$request->validate([...])` → automatic redirect/JSON error message.
2. **Form Requests:** Own class with `rules()`, `authorize()`, `messages()`, `after()`.
3. **`prepareForValidation()`:** Normalize data before validation.
4. **`passedValidation()`:** Transform data after successful validation.
5. **`stopOnFirstFailure`:** Abort on first error.
6. **Validated/Safe API:** `$request->safe()->only(['email'])`.
7. **Named Error Bags:** `validateWithBag('login', [...])`.
8. **50+ Built-in Rules:** `required`, `email`, `unique:table`, `exists:table`, `confirmed`, `date`, `file`, `image`, `array`, `json`, `uuid`, custom Rule objects, Closures.
9. **Custom Messages:** Per field + rule, in language files or inline.
10. **Error Response Format:** JSON with `message` + `errors` (422).

#### For Dreego adopt (highest priority!):

Dreego already has `go-playground/validator` in the stack. This needs to be extended:

| Laravel Feature | Dreego Adaptation |
|---|---|
| Form Requests | `type LoginForm struct { ... }` in `<go>` block with struct tags |
| Auto-Error-Redirect | `dreego.FormHandler` with automatic flash-back |
| Field-Level Errors | `c.Errors.Get("email")` in template |
| `@error('email')` | `{#if errors.email}` in template → V1 |
| `old('email')` | `c.Old("email")` → V1 |
| Custom Rules | `dreego.RegisterRule("strong_password", func) + struct tag |
| JSON Error Format | Standardized 422 JSON → V1 |
| `prepareForValidation` | `func (f *LoginForm) Prepare() { ... }` |
| Safe/Only/Except | `dreego.SafeInput(r, "email", "name")` |

### 2.8 Laravel's "Batteries Included" Philosophy

Laravel ships **everything with it**:
- Auth (Login, Register, Password Reset, Email Verification)
- Session (File, Cookie, Database, Redis)
- Cache (File, Redis, Memcached, DynamoDB)
- Filesystem (Local, S3, FTP, SFTP)
- Mail (SMTP, Mailgun, Postmark, SES)
- Broadcasting (Pusher, Ably, Redis/Soketi)
- Scheduling (Cron expression in PHP, no crontab)
- Localization (Translation files per language)
- Logging (Monolog, Channels)
- Encryption (AES-256)
- Hashing (bcrypt, argon2)
- CSRF (automatic)
- Collections (Array helper with 100+ methods)

#### For Dreego:
- **Core must remain minimal** (Go Standard Library philosophy)
- **Plugins as "Batteries"**: `dreego-auth`, `dreego-cache`, `dreego-mail`, `dreego-storage`
- **Don't ship everything, but make it easy to install** (`dreego add auth`)

---

### 2.9 Anti-Patterns of Laravel Dreego Should Avoid

| Anti-Pattern | Why |
|---|---|
| **Facades:** Static proxies, hide real dependencies | Dreego: Dependency injection via struct fields |
| **Magic Properties:** `$user->posts` (Eloquent Dynamic Property) | Dreego: Explicit methods |
| **ActiveRecord ORM:** DB logic and business logic in the same object | Dreego: Separation of Schema + Query + Repo (Ecto pattern) |
| **PHP's Shared-Nothing:** Every request starts from zero | Dreego: Go's request pool, connection pool natively |
| **Too many conventions:** Everything in `app/`, everything in `Http/Controllers` | Dreego: Free project structure |
| **Lazy Loading trap:** `$post->author->name` triggers N+1 | Dreego: Explicit preloads |
| **String-based validation:** `'required|email|unique:users'` | Dreego: Struct tags (type-safe) |

---

## 3. Django (Python)

### 3.1 Django Admin — Auto-generated Admin Panels

The Django Admin is **Django's killer feature**:

1. **Auto-Discovery:** `admin.py` in each app is automatically found.
2. **Model Registration:**
   ```python
   @admin.register(Author)
   class AuthorAdmin(admin.ModelAdmin):
       list_display = ["name", "title"]
       list_filter = ["status"]
       search_fields = ["name"]
   ```
3. **CRUD out of the box:** Create, Read, Update, Delete for each registered model.
4. **Customization without writing code:**
   - `list_display`, `list_filter`, `search_fields`
   - `date_hierarchy` (auto. calendar drilldown)
   - `ordering`, `list_per_page`, `list_editable` (Inline editing in the list!)
   - `readonly_fields`, `fields`, `fieldsets` (Field grouping with collapse)
   - `prepopulated_fields` (Generate slug from title)
   - `autocomplete_fields` (Select2 with async search)
   - `raw_id_fields` (Popup selector for ForeignKeys)
   - `filter_horizontal` / `filter_vertical` (Two-box widget for ManyToMany)
   - `inlines` (Edit child models in parent form — TabularInline, StackedInline)
5. **Admin Actions:** Batch operations (e.g. "Mark as published" for selected objects).
6. **Custom Templates:** Admin pages fully overridable.
7. **Facet Counts:** Shows count per filter option.
8. **Responsive Design:** Works on mobile.

#### For Dreego:

- **`dreego-admin` is a MUST-HAVE plugin.**
- Django Admin is **the standard** every admin generator tool must measure against.
- Generate automatic CRUD UI from Go structs.
- `list_display`, `list_filter`, `search_fields` as struct tag options:
  ```go
  type Post struct {
      Title string `admin:"list_display,search"`
      Status string `admin:"list_filter"`
  }
  ```

### 3.2 Django ORM and Migrations

Django ORM Features:
- **Model Definition:** Python class → DB table automatically.
- **Field Types:** 40+ Built-in Fields (`CharField`, `TextField`, `IntegerField`, `FileField`, `ImageField`, `JSONField`, etc.).
- **Relationships:** `ForeignKey`, `ManyToManyField`, `OneToOneField`.
- **Migrations:** `makemigrations` (auto-detect Changes), `migrate` (apply), `sqlmigrate` (SQL preview).
- **Manager:** `objects` — Query interface per model.
- **QuerySet:** Lazy, chainable, immutable:
  ```python
  Post.objects.filter(status='published').exclude(deleted=True).order_by('-created_at')
  ```
- **Q-Objects:** Complex queries with `|` (OR), `&` (AND), `~` (NOT).
- **F-Expressions:** DB-level operations `Post.objects.update(views=F('views') + 1)`.
- **Aggregation:** `Count`, `Sum`, `Avg`, `Max`, `Min`.
- **select_related (JOIN) / prefetch_related (separate Query):** Eager Loading.
- **Database Router:** Multi-DB support.
- **Transactions:** `transaction.atomic()`.
- **Signals:** `pre_save`, `post_save`, `pre_delete`, `post_delete`.

#### For Dreego adopt:

| Django Feature | Dreego Adaptation |
|---|---|
| Auto-Migration Detection | `dreego-db`: Go struct changes → generate migration SQL |
| QuerySet (lazy, chainable) | Query builder with method chaining |
| F-Expressions | DB-level updates |
| select_related / prefetch_related | Explicit preloads |
| Signals | Model events (`AfterCreate`, `BeforeDelete`) |
| Multi-DB | Connection string per model |

### 3.3 Django Template Engine (DTL)

DTL Features:
- **Variables:** `{{ variable }}` — auto-escaping (XSS protection).
- **Filters:** `{{ name|lower }}`, `{{ date|date:"Y-m-d" }}`, `{{ text|truncatechars:100 }}`.
- **Tags:** `{% if %}`, `{% for %}`, `{% block %}`, `{% extends %}`, `{% include %}`.
- **Template Inheritance:**
  ```django
  {% extends "base.html" %}
  {% block content %}
  {% endblock %}
  ```
- **Context Processors:** Request data automatically available in template.
- **Custom Tags & Filters:** Extensible.
- **Template Loaders:** Filesystem, App directories, Cached.
- **Jinja2 Support:** Django can also use Jinja2 as engine.

**BUT:** DTL is deliberately **restricted** — no complex logic in templates (philosophy: logic belongs in views, not templates). Less powerful than Blade or HEEx. **No components, no slots.**

#### For Dreego: Dreego's template engine (Svelte-inspired) is already more advanced than DTL. Django's Template Engine is not a model for Dreego — Blade and HEEx are the better references.

### 3.4 "Batteries Included"

Django's Built-in Apps (in `INSTALLED_APPS`):
- `django.contrib.admin` — Admin Panel
- `django.contrib.auth` — User, Groups, Permissions
- `django.contrib.contenttypes` — Generic Relations
- `django.contrib.sessions` — Session Management
- `django.contrib.messages` — Flash Messages
- `django.contrib.staticfiles` — Static File Management
- `django.contrib.sites` — Multi-Site Framework
- `django.contrib.sitemaps` — Sitemap Generation
- `django.contrib.syndication` — RSS Feed Generation
- `django.contrib.gis` — GeoDjango (PostGIS)
- `django.contrib.postgres` — PostgreSQL-specific Features

#### For Dreego:
- **Core:** Router, Context, Middleware, Template Engine.
- **Official Plugins:** Auth, Admin, Sessions, Storage, Cache, Mail, Jobs.
- **Community Plugins:** Everything else.

### 3.5 Authentication System

Django Auth — **most complete built-in auth system of all frameworks**:

1. **User Model:** `AbstractUser` (username, password, email, is_staff, is_superuser) — extensible.
2. **Permissions:** `User.has_perm('app.action_model')`, Groups.
3. **Authentication Backends:** ModelBackend (Username+Password), RemoteUserBackend, custom Backends.
4. **Session Authentication:** Cookie-based.
5. **Login/Logout Views:** `LoginView`, `LogoutView` (ready to use).
6. **Password Management:**
   - `PasswordChangeView`, `PasswordChangeDoneView`
   - `PasswordResetView`, `PasswordResetDoneView`, `PasswordResetConfirmView`, `PasswordResetCompleteView`
   - Password Hashers: PBKDF2 (default), Argon2, BCrypt
7. **Password Validators:** Minimum length, Common-Password-Check, Numeric-Only-Check, User-Attribute-Similarity-Check.
8. **Login Rate Limiting:** `AXES` (third-party, but standard).
9. **User Creation Forms:** `UserCreationForm`, `UserChangeForm`.
10. **`login_required` Decorator:** `@login_required` on views.
11. **`user_passes_test`:** Custom permission checks.
12. **Mixin Classes:** `LoginRequiredMixin`, `PermissionRequiredMixin`, `UserPassesTestMixin`.

#### For Dreego:

- **`dreego-auth` plugin with:**
  - Go Interface: `type User interface { GetID() string; HasRole(string) bool }`
  - Session-based auth (Cookie)
  - Middleware: `auth.Required`, `auth.HasRole("admin")`
  - Password Hashing: Argon2id via `golang.org/x/crypto`
  - Login/Register/Reset routes
  - OAuth2/OIDC via third-party (`dreego-oauth`)
  - CSRF protection in core, used by auth

### 3.6 Django REST Framework (DRF)

DRF is the **standard for API development in Django**:

1. **Serializers:** Model → JSON/Dict, with validation:
   ```python
   class UserSerializer(serializers.ModelSerializer):
       class Meta:
           model = User
           fields = ['id', 'username', 'email']
   ```
2. **ViewSets:** CRUD logic in one class:
   ```python
   class UserViewSet(viewsets.ModelViewSet):
       queryset = User.objects.all()
       serializer_class = UserSerializer
   ```
3. **Routers:** Automatic URL configuration from ViewSets:
   ```python
   router = routers.DefaultRouter()
   router.register(r'users', UserViewSet)
   ```
4. **Browsable API:** Interactive API docs in browser (GET/POST/PUT/DELETE).
5. **Authentication:** SessionAuth, TokenAuth, OAuth1/2, JWT (via third-party).
6. **Permissions:** `IsAuthenticated`, `IsAdminUser`, `DjangoModelPermissions`, custom.
7. **Throttling:** Rate Limiting (`AnonRateThrottle`, `UserRateThrottle`).
8. **Pagination:** `PageNumberPagination`, `LimitOffsetPagination`, `CursorPagination`.
9. **Filtering:** `SearchFilter`, `OrderingFilter`, `DjangoFilterBackend`.
10. **Versioning:** URL-Path, Query-Parameter, Accept-Header.
11. **Renderers:** JSON, Browsable API, JSONP, XML, YAML.
12. **Parsers:** JSON, Form, MultiPart, FileUpload.
13. **Schemas:** OpenAPI/Swagger Auto-Generation.

#### For Dreego:

- **`dreego-api` plugin:**
  - ViewSet-like pattern for REST endpoints:
    ```go
    type UserAPI struct {
        dreego.API
    }
    func (api *UserAPI) List() ([]User, error) { ... }
    func (api *UserAPI) Get(id string) (*User, error) { ... }
    ```
  - Auto-Router: `dreego.RegisterAPI(router, "/api/users", &UserAPI{})`
  - OpenAPI generation from Go structs + routes
  - Browsable API? Optional, but nice-to-have.

### 3.7 What Keeps Django Relevant After 20+ Years

1. **Stability:** No breaking changes without long deprecation paths.
2. **Documentation:** Gold standard — every class, every method documented.
3. **Batteries Included:** For 80% of web apps, Django's built-in features suffice.
4. **Admin:** Still unbeaten.
5. **Community:** Django Packages, DjangoCon, DSF.
6. **Security:** XSS, CSRF, SQL Injection, Clickjacking — all built-in protected.
7. **Migrations:** Was the first framework with a robust migration system.
8. **ORM:** Among the best ORMs (even if Ecto is better).

#### For Dreego:
- **Aim for longevity:** Not hype-driven, but solid architecture.
- **Documentation from Day 1:** Go doc + detailed guides.
- **Stable API:** `dreego.Context` Interface, Plugin Interface → Versioning.

---

### 3.8 Anti-Patterns of Django Dreego Should Avoid

| Anti-Pattern | Why |
|---|---|
| **Monolithic `settings.py`:** Everything in one file, quickly becomes confusing | Dreego: `dreego.config.json` per project, modular |
| **Django ORM is slow on complex queries:** suboptimal SQL on `annotate()` + `aggregate()` | Dreego: Raw SQL escape hatch, query builder that produces good SQL |
| **Template engine is too restrictive:** No components, no functions with arguments | Dreego: Modern template engine with components and slots |
| **Python GIL:** Limits concurrency | Dreego: Go's goroutines are superior |
| **"Django way or the highway":** Very opinionated | Dreego: Conventions yes, but override options |
| **Migrations can become slow on large projects:** 1000+ migrations | Dreego: Support squash migrations from the start |
| **No async support in ORM (until Django 4.1):** async ORM is not yet mature | Dreego: Go is naturally concurrent |

---

## Summary: What Dreego Should Learn from Each Framework

### From Phoenix (Elixir):
1. **LiveView Architecture** → SSE-First reactivity via Datastar
2. **Ecto Pattern** (Schema ≠ Query ≠ Repo ≠ Changeset) → Separation of concerns
3. **PubSub System** → Event bus for real-time updates
4. **Generator Tiering** (embedded → schema → context → live) → `dreego generate`
5. **Compile-Time Template Validation** → Transpiler checks HTML & Go syntax
6. **Graceful Reconnection** → SSE with auto-reconnect

### From Laravel (PHP):
1. **Eloquent's Developer Experience** → Naming conventions, Eager Loading, Scopes
2. **Validation System** → Form Requests, Field-Level Errors, `safe()` / `old()` API
3. **Artisan CLI** → 30+ `make` commands, stub customization
4. **Blade Components & Slots** → `<x-*>` = `<dreego:*>`
5. **Queue/Job System** → Job Middleware, Batching, Chaining, Delayed Dispatch
6. **Notification System** → Multi-Channel (Mail, DB, Slack)
7. **Ecosystem-First Thinking** → Plugin interface must enable plugin ecosystem

### From Django (Python):
1. **Admin Panel** → `dreego-admin` plugin (highest priority)
2. **Auth System** → Password hashing, reset flow, permission system
3. **DRF's ViewSet/Router Pattern** → `dreego-api` plugin
4. **Migrations** → Auto-detection from Go structs
5. **Documentation** → Gold standard as model
6. **Stability & Longevity** → No breaking changes without deprecation

### Anti-Patterns (consistent across all frameworks):
1. **Magic / Implicitness** → Dreego: Everything explicit (Go philosophy)
2. **Too opinionated** → Override paths for conventions
3. **Monolithic Core** → Core minimal, plugins for batteries
4. **String-based APIs** → Type safety via Go structs
5. **Lazy Loading / Implicit N+1** → Explicit preloads
