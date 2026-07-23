# Framework-Vergleich: Phoenix, Laravel, Django — Relevanz für Dreego

**Datum:** 23.07.2026
**Zweck:** Systematische Analyse, welche Features der drei etablierten Frameworks Dreego übernehmen, adaptieren oder bewusst weglassen sollte.

---

## 1. Phoenix (Elixir)

### 1.1 LiveView — Architektur (WebSockets)

LiveView funktioniert über **WebSockets** (nicht SSE). Architektur:

1. Initialer HTTP-Request rendert statisches HTML (SSR).
2. Client stellt WebSocket-Verbindung zum Server her.
3. Server spawn einen **stateful Elixir-Prozess** (GenServer-ähnlich) pro Client.
4. Jeder State-Change auf dem Server berechnet ein **Minimal-Diff** und pusht nur das geänderte HTML über den WebSocket an den Client.
5. Client-seitig ersetzt ein kleines JS das DOM selektiv (morphdom-ähnlich).
6. Bei Verbindungsabbruch: **graceful reconnection** mit State-Recovery.

**Wichtige Details:**
- `mount/3` wird **zweimal** aufgerufen: einmal für statisches HTML, einmal beim WebSocket-Connect. `connected?(socket)` unterscheidet die Modi.
- `assign_async/3` — asynchrone Daten in Tasks laden, mit `AsyncResult` für Loading/Failed/OK States.
- `start_async/3` + `handle_async/3` — niedrig-level async control.
- `stream/4` — effiziente Listen-Verwaltung (Insert/Delete/Update ohne kompletten Re-Render).
- **Hibernation:** LiveView kann nach 15s Inaktivität seinen Prozess-State komprimieren (`hibernate_after`).
- **Flash Messages:** Built-in `put_flash/3`, `clear_flash/1` — funktionieren auch über WebSockets.
- **File Uploads:** `allow_upload/3`, `consume_uploaded_entries/3` — komplett integriert mit Progress-Tracking.
- **Lifecycle Hooks:** `attach_hook/4` für `:handle_params`, `:handle_event`, `:handle_info`, `:after_render`.
- **LiveComponents:** Stateful sub-components mit eigenem `mount`, `update`, `handle_event`.

#### Für Dreego adoptieren / adaptieren:

| LiveView-Feature | Dreego-Adaption | Priorität |
|---|---|---|
| SSR → WebSocket Upgrade | HTMX partials + SSE via Datastar. Kein WebSocket-Zwang, HTMX kann SSE/WS je nach use case | ✅ V1 |
| Stateful Server-Prozess | Go-Goroutine pro Session (leichtgewichtig). Vorbild: `connected?()` Check | ✅ V2 |
| `assign_async` | `<go>`-Block: async functions mit `loading`/`error`/`ok` States im Template | ✅ V2 |
| `stream/4` (effiziente Listen) | `{#each}` mit SSE-basierten Streaming-Updates | ✅ V3 |
| Graceful Reconnection | Datastar's built-in SSE Reconnect | ✅ V1 |
| Flash Messages | Built-in Flash in `dreego.Context` | ✅ V1 |
| File Uploads | `dreego-storage` Addon, Vorbild: `allow_upload/3` API | ✅ V2 |
| Lifecycle Hooks | `dreego.Plugin` Interface → Middleware-Hooks + Transpiler-Hooks | ✅ V1 |
| LiveComponents | `.dreego`-Komponenten mit eigenem State (<go> pro Komponente) | ✅ V2 |

### 1.2 Phoenix PubSub

Phoenix PubSub ist ein **cluster-weites Publish/Subscribe-System**:

```elixir
PubSub.subscribe(:my_pubsub, "user:123")
PubSub.broadcast(:my_pubsub, "user:123", {:user_update, %{name: "Shane"}})
```

- **Adapter-Architektur:** PG2 (Distributed Erlang, Default), Redis, oder custom Adapter.
- **`broadcast_from/5`:** Broadcast an alle außer den Sender.
- **`direct_broadcast/5`:** Broadcast an einen bestimmten Node.
- **`local_broadcast/4`:** Nur lokaler Node.
- **Custom Dispatcher:** "Fastlane" — bei Channels werden Nachrichten einmal encodiert und direkt in alle Sockets geschrieben, statt pro Channel.
- **Safe Pool-Size-Migration:** `broadcast_pool_size` erlaubt Rolling-Updates ohne Nachrichtenverlust.

#### Für Dreego adoptieren:

- **Event-Bus im Core:** `dreego.Emit("user:123", payload)` → alle Subscriber (SSE/WS) bekommen Update.
- **Nicht selbst bauen:** Entweder NATS embedded, Redis PubSub, oder Go-Kanäle-basierter Event-Bus.
- **Addon:** `dreego-pubsub` mit Redis/NATS-Backend.

### 1.3 Ecto — Warum es gut ist

Ecto ist **kein ORM** sondern ein **Database-Wrapper + Query-Builder**:

1. **Explizite Trennung:** Schema (Mapping) ≠ Query (Abfragen) ≠ Repo (Datenbank-Zugriff) ≠ Changeset (Validierung).
2. **Changesets:** Daten-Validierung und Casting als **eigenes Objekt**, nicht am Model. `Ecto.Changeset` hält den State (valid?, errors, changes).
3. **Composable Queries:** Queries sind immutable und chainable:
   ```elixir
   query = from u in User, where: u.age > 18
   query = from u in query, where: u.active == true
   ```
4. **Multi:** Transaktionen mit mehreren Operations `Ecto.Multi` — Rollback bei Fehler.
5. **Embedded Schemas:** DB-lose Schemas für Form-Validierung oder API-Responses.
6. **Migrations:** Generiert und versioniert.
7. **No Magic:** Kein lazy loading, kein implicit N+1. Alles explizit.

#### Für Dreego adoptieren:

| Ecto-Feature | Dreego-Adaption |
|---|---|
| Changesets | `dreego-db` Addon: Validierungsobjekt unabhängig vom DB-Model |
| Composable Queries | Go: Query-Builder mit Chainable-Methoden (ähnlich squirrel) |
| Multi (Transaktionen) | Go: `db.Transaction(func(tx *sql.Tx) error { ... })` |
| Embedded Schemas | Go-Structs mit Tags für Validierung ohne DB-Bindung |
| `Ecto.Enum` | Go: `enum:"admin,user,mod"` struct tag |

### 1.4 Generatoren (`mix phx.gen.*`)

Phoenix hat eine gestaffelte Generator-Strategie:

| Generator | Erstellt |
|---|---|
| `phx.gen.embedded` | Nur Schema (keine Migration, kein Context) |
| `phx.gen.schema` | Schema + Migration |
| `phx.gen.context` | Schema + Migration + Context-Modul |
| `phx.gen.live` | Context + LiveView (CRUD) |
| `phx.gen.json` | Context + JSON-API-Controller + Views |
| `phx.gen.html` | Context + HTML-Controller + Templates |

Templates können **pro Projekt überschrieben** werden (`priv/templates/`).

#### Für Dreego adoptieren:

- `dreego generate page` — eine `.dreego`-Datei mit `<go>` + Template erstellen (analog zu `phx.gen.live`)
- `dreego generate resource` — CRUD-Gerüst (Route + Form + Validierung)
- `dreego generate schema` — Go-Struct mit Validierungs-Tags
- **Template-Overrides** via `.dreego/templates/`

### 1.5 Phoenix Template Engine (HEEx / EEx)

HEEx ist Phoenix' **Template-Engine mit Component-Modell**:

- **`~H` Sigil:** Kompilierungszeit-geprüftes HTML mit Komponenten
- **Components:** `<.component_name attr="val">` → Funktionsaufrufe
- **Slots:** `<:slot_name>` für Multi-Slot-Komponenten
- **Directives:** `:if`, `:for`, `:let`
- **HEEx Expressions:** `{@assign}` für Werte
- **HTML-Validierung zur Compile-Zeit**

#### Für Dreego adoptieren:

- `{#if}`, `{#each}`, `{#switch}`, `{#slot}` sind bereits geplant → ✅
- HEEx-ähnliche Compile-Time-Validierung: `.dreego` → Transpiler prüft auf valides HTML → ✅ V2
- Komponenten-Modell: `<dreego:map />`, `<x-card>` (analog zu Blade Components) → ✅

---

### 1.6 Was Phoenix "on par mit JS-Frameworks" macht

1. **LiveView ersetzt SPA-JS:** Kein React/Vue nötig für reaktive UIs.
2. **Erlang/OTP-Concurrency:** Millionen paralleler Verbindungen.
3. **Soft-Realtime Built-in:** Channels, PubSub, Presence (wer ist online?).
4. **Fault-Tolerance:** Supervisor-Trees, automatische Recovery.
5. **Single Binary Deployment:** Mix Releases.

### 1.7 Deployment (Releases)

- **Mix Releases:** Elixir's Built-in Release-System → Single Binary (ähnlich Go!).
- **Phoenix-Endpoint-Konfiguration:** Alles in `config/runtime.exs` — Environment-spezifisch.
- **Docker-freundlich:** Minimales Alpine-Image mit Release.

#### Für Dreego: Go macht Single-Binary-Deployment nativ besser als jeder andere Stack. Das ist einer der Hauptvorteile—unbedingt betonen und ausbauen.

---

### 1.8 Anti-Patterns von Phoenix, die Dreego vermeiden sollte

| Anti-Pattern | Warum |
|---|---|
| **Magische Macros:** `use Phoenix.LiveView` injiziert implizit Callbacks | Dreego: Alles explizit in Go, keine Magic |
| **Zu viele Konventionen:** Phoenix ist sehr opinionated | Dreego: Konventionen ja (file-based routing), aber Override-Pfade für alles |
| **Elixir-Lernkurve:** Funktionale Sprache, Pattern Matching, OTP | Dreego: Go ist simpler, flachere Lernkurve |
| **WebSocket-First:** Funktioniert nicht gut hinter manchen Proxies | Dreego: SSE-First (kompatibler, HTTP-basiert) |
| **LiveView State Bloat:** Große LiveViews können viel RAM verbrauchen | Dreego: Flat state in `<go>` — keine verschachtelten Prozesse pro User |

---

## 2. Laravel (PHP)

### 2.1 Eloquent ORM — Was es beliebt macht

Eloquent ist ein **ActiveRecord-ORM**:

1. **Convention over Configuration:**
   - `Flight` Model → `flights` Tabelle (Plural, snake_case)
   - `id` als Primary Key (auto-increment)
   - `created_at`, `updated_at` automatisch
2. **Relationship API:**
   ```php
   $user->posts()->where('active', 1)->get();
   $post->user->name; // Dynamic Property
   ```
3. **Eager Loading:** `User::with('posts.comments')->get()` — verhindert N+1.
4. **Mass Assignment Protection:** `$fillable` / `$guarded` — explizit, sicher.
5. **Query Scopes:** `User::active()->verified()->get()`
6. **Accessors/Mutators:** `getNameAttribute()`, `setPasswordAttribute()`
7. **Model Events:** `creating`, `updated`, `deleting` — Observer-Pattern.
8. **Soft Deletes:** `deleted_at` timestamp, `withTrashed()`, `onlyTrashed()`.
9. **Casting:** `$casts` Property für JSON, Array, Date, Enum, etc.
10. **Pruning:** `MassPrunable` Trait — automatisch alte Models löschen.

#### Für Dreego adoptieren:

| Eloquent-Feature | Dreego-Adaption |
|---|---|
| Convention over Config | `dreego-db`: Go-Struct → Tabellenname via Pluralizer |
| Relationships | `dreego-db`: Struct-Tags `dreego:"has_many:Posts"` → Preload-System |
| Eager Loading / N+1 Prevention | `dreego-db`: `db.With("Posts.Comments").Find(&user)` |
| Mass Assignment (`$fillable`) | Go-Struct Tags: `json:"name" db:"name" fillable:"true"` |
| Query Scopes | Go-Method-Chaining: `db.Scope(ActiveUsers).Find(&users)` |
| Model Events / Observers | Go: `AfterCreate`, `BeforeDelete` Interface |
| Soft Deletes | `deleted_at` nullable, `db.WithTrashed()` |
| Casting | Go-Struct-Tags mit Type-Mapping |
| `firstOrCreate`, `updateOrCreate` | Convenience-Methoden in `dreego-db` |

### 2.2 Blade Templating — Features

Blade kompiliert zu **Plain PHP** und cached die Ergebnisse. Features:

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
3. **Component Methods:** `$isSelected($value)` im Template aufrufbar.
4. **Anonymous Components:** Nur Template-Datei, kein PHP-Class nötig.
5. **Conditional Class Merging (`$attributes->merge()`, `@class([])`).**
6. **Stacks:** `@push('scripts')` / `@stack('scripts')` — Assets von Child-Komponenten ins Layout-Head schieben.
7. **Service Injection:** `@inject('metrics', 'App\Services\MetricsService')`.
8. **`@once` Directive:** Code wird nur einmal pro Render-Zyklus ausgegeben.
9. **Blade Fragments:** `@fragment` für AJAX-Teilladungen (perfekt mit HTMX!).
10. **Custom Directives:** `Blade::directive('datetime', fn($exp) => ...)`.
11. **Custom If Statements:** `Blade::if('cloud', fn() => ...)`.
12. **`@verbatim`:** JS-Template-Syntax (Vue, Alpine) im Blade unberührt lassen.
13. **`$loop` Variable:** Index, Iteration, First, Last, Even, Odd, Depth, Parent.

#### Für Dreego adoptieren:

| Blade-Feature | Dreego-Adaption |
|---|---|
| Components (`<x-alert>`) | `<dreego:alert>` oder als Addon-Tags (bereits geplant) |
| Conditional Classes | `<div class:active={isActive}>` — Svelte-Style (geplant) |
| Stacks (`@push`/`@stack`) | `<slot name="head">` oder `<head>`-Block CSS/JS Injection |
| `@once` | `{#once}` block → V2 |
| Blade Fragments + HTMX | HTMX `hx-select` + Partial-Templates → super Kombination |
| `$loop` Variable | `{#each items as item, index}` → `{index}`, `{first}`, `{last}` Built-in → V1 |
| Template Inheritance | Layouts via `layout.dreego` (geplant) |

### 2.3 Artisan CLI — Generator-Übersicht

Artisan hat **30+ make-Kommandos**:

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

Jedes Kommando erzeugt eine Datei an der richtigen Stelle mit der richtigen Boilerplate.

**Weitere Artisan-Features:**
- **Tinker:** REPL mit Zugriff auf das gesamte Framework (Eloquent, Jobs, Events).
- **Stub Customization:** `php artisan stub:publish` → Stubs im Projekt überschreibbar.
- **Interactive Prompts:** `ask()`, `secret()`, `confirm()`, `anticipate()`, `choice()`.
- **Tables & Progress Bars:** `$this->table()`, `$this->withProgressBar()`.
- **Signal Handling:** `$this->trap(SIGTERM, ...)`.
- **Programmatic Execution:** `Artisan::call('mail:send', [...])` aus Code.
- **Queueing Commands:** `Artisan::queue('mail:send', [...])`.

#### Für Dreego adoptieren:

- `dreego new` — Projekt-Scaffolding
- `dreego generate page` — Neue Seite
- `dreego generate resource` — CRUD
- `dreego generate middleware` — Middleware
- `dreego routes` — Alle Routen anzeigen
- `dreego tinker` — Go-REPL (via yaegi oder go-pry)
- **Stub Customization** → `.dreego/templates/`

### 2.4 Ecosystem: Forge, Vapor, Nova, Spark

| Produkt | Beschreibung | Dreego-Äquivalent |
|---|---|---|
| **Forge** | Server-Management (provisioning, deployment) | `dreego deploy` — via SSH oder Docker |
| **Vapor** | Serverless Laravel auf AWS Lambda | Nicht relevant für Go (Single Binary) |
| **Nova** | Admin-Panel-Generator | `dreego-admin` Addon |
| **Spark** | SaaS-Starter-Kit (Billing, Teams) | `dreego-saas` Addon |
| **Envoyer** | Zero-Downtime Deployment | `dreego deploy` mit Blue-Green |
| **Horizon** | Queue-Monitoring-Dashboard | `dreego-jobs` Dashboard |
| **Telescope** | Debugging & Monitoring | `dreego-devtools` |
| **Pennant** | Feature Flags | `dreego-features` Addon |
| **Pulse** | Performance Monitoring | Nicht V1 |
| **Reverb** | WebSocket-Server (First-Party) | Nicht nötig (SSE) |
| **Echo** | Client-seitige WebSocket-Bibliothek | Datastar (SSE) reicht |

**Takeaway:** Laravel's Ecosystem ist die **größte Stärke** des Frameworks. Dreego muss das nicht in V1 replizieren, aber die **Plugin-Architektur muss offen genug sein**, dass solche Tools als Addons entstehen können.

### 2.5 Queue/Job-System

Laravel's Queue-System ist extrem ausgereift:

1. **Unified API über Backends:** Database, Redis, SQS, Beanstalkd, MongoDB.
2. **Job-Klassen** mit `handle()`-Methode.
3. **Dependency Injection** in `handle()` via Service Container.
4. **Eloquent Model Serialization:** Models werden als Identifier serialisiert und beim Verarbeiten neu geladen.
5. **Unique Jobs:** `ShouldBeUnique` — kein Duplikat in der Queue.
6. **Job Middleware:** Rate Limiting, Overlap Prevention, Exception Throttling.
7. **Job Batching:** Mehrere Jobs als Batch ausführen, mit Callbacks (then/catch/finally).
8. **Job Chaining:** Jobs sequenziell ausführen.
9. **Delayed Dispatching:** `->delay(now()->addMinutes(10))`.
10. **`dispatchAfterResponse()`:** Job erst nach HTTP-Response ausführen.
11. **Failed Jobs:** `failed_jobs`-Tabelle, automatische Retries, `retryUntil()`.
12. **Queue Priorities:** `--queue=high,default` — Worker priorisiert.
13. **Horizon Dashboard:** Redis-Queue-Monitoring.

#### Für Dreego adoptieren:

- `dreego-jobs` Addon mit:
  - Interface: `type Job interface { Handle() error }`
  - Backends: Redis, PG, NATS
  - Delayed Jobs: `dreego.Dispatch(job).Delay(10 * time.Minute)`
  - Job Middleware (Rate Limiting, Retry)
  - Failed Job Logging
- Go-Goroutinen sind nativer als PHP-Worker — ein Job-System in Go ist **deutlich performanter**.

### 2.6 Notification System

Laravel Notifications:
- **Mehrere Channels:** Mail, SMS (Vonage), Slack, Database, Broadcast, custom.
- **`Notification::send($users, new InvoicePaid($invoice))`**.
- **`Notifiable` Trait:** `$user->notify(new WelcomeNotification())`.
- **Database Notifications:** In der DB gespeichert, per API abrufbar.
- **Markdown Mail Templates:** Vorformatierte E-Mail-Komponenten.

#### Für Dreego:
- `dreego-notify` Addon mit Channel-Interface:
  ```go
  type Channel interface {
      Send(to User, notification Notification) error
  }
  ```

### 2.7 Validation Approach

Laravel's Validation ist **beste aller Frameworks**:

1. **Request-Validation:** `$request->validate([...])` → automatische Redirect/JSON-Fehlermeldung.
2. **Form Requests:** Eigene Klasse mit `rules()`, `authorize()`, `messages()`, `after()`.
3. **`prepareForValidation()`:** Daten vor Validierung normalisieren.
4. **`passedValidation()`:** Daten nach erfolgreicher Validierung transformieren.
5. **`stopOnFirstFailure`:** Abbruch beim ersten Fehler.
6. **Validated/Safe API:** `$request->safe()->only(['email'])`.
7. **Named Error Bags:** `validateWithBag('login', [...])`.
8. **50+ Built-in Rules:** `required`, `email`, `unique:table`, `exists:table`, `confirmed`, `date`, `file`, `image`, `array`, `json`, `uuid`, custom Rule-Objekte, Closures.
9. **Custom Messages:** Per Feld + Regel, in Language-Files oder inline.
10. **Error-Response-Format:** JSON mit `message` + `errors` (422).

#### Für Dreego adoptieren (höchste Priorität!):

Dreego hat bereits `go-playground/validator` im Stack. Das muss ergänzt werden:

| Laravel-Feature | Dreego-Adaption |
|---|---|
| Form Requests | `type LoginForm struct { ... }` im `<go>`-Block mit Struct-Tags |
| Auto-Error-Redirect | `dreego.FormHandler` mit automatischem Flash-Back |
| Field-Level Errors | `c.Errors.Get("email")` im Template |
| `@error('email')` | `{#if errors.email}` im Template → V1 |
| `old('email')` | `c.Old("email")` → V1 |
| Custom Rules | `dreego.RegisterRule("strong_password", func) + struct tag |
| JSON Error Format | Standardisiertes 422 JSON → V1 |
| `prepareForValidation` | `func (f *LoginForm) Prepare() { ... }` |
| Safe/Only/Except | `dreego.SafeInput(r, "email", "name")` |

### 2.8 Laravel's "Batteries Included"-Philosophie

Laravel liefert **alles mit**:
- Auth (Login, Register, Passwort-Reset, Email-Verification)
- Session (File, Cookie, Database, Redis)
- Cache (File, Redis, Memcached, DynamoDB)
- Filesystem (Local, S3, FTP, SFTP)
- Mail (SMTP, Mailgun, Postmark, SES)
- Broadcasting (Pusher, Ably, Redis/Soketi)
- Scheduling (Cron-Expression in PHP, kein Crontab)
- Localization (Übersetzungs-Dateien pro Sprache)
- Logging (Monolog, Channels)
- Encryption (AES-256)
- Hashing (bcrypt, argon2)
- CSRF (automatisch)
- Collections (Array-Helper mit 100+ Methoden)

#### Für Dreego:
- **Core muss minimal bleiben** (Go-Standard-Library-Philosophie)
- **Addons als "Batteries"**: `dreego-auth`, `dreego-cache`, `dreego-mail`, `dreego-storage`
- **Nicht alles mitliefern, aber einfach installierbar machen** (`dreego add auth`)

---

### 2.9 Anti-Patterns von Laravel, die Dreego vermeiden sollte

| Anti-Pattern | Warum |
|---|---|
| **Facades:** Statische Proxies, verstecken echte Abhängigkeiten | Dreego: Dependency Injection via Struct-Felder |
| **Magic Properties:** `$user->posts` (Eloquent Dynamic Property) | Dreego: Explizite Methoden |
| **ActiveRecord ORM:** DB-Logik und Business-Logik im selben Objekt | Dreego: Trennung von Schema + Query + Repo (Ecto-Muster) |
| **PHP's Shared-Nothing:** Jeder Request startet von Null | Dreego: Go's Request-Pool, Connection-Pool nativ |
| **Zu viele Konventionen:** Alles in `app/`, alles in `Http/Controllers` | Dreego: Freie Projekt-Struktur |
| **Lazy Loading Falle:** `$post->author->name` löst N+1 aus | Dreego: Explizite Preloads |
| **String-basierte Validierung:** `'required|email|unique:users'` | Dreego: Struct-Tags (type-safe) |

---

## 3. Django (Python)

### 3.1 Django Admin — Auto-generierte Admin-Panels

Das Django Admin ist **das Killer-Feature von Django**:

1. **Auto-Discovery:** `admin.py` in jeder App wird automatisch gefunden.
2. **Model-Registrierung:**
   ```python
   @admin.register(Author)
   class AuthorAdmin(admin.ModelAdmin):
       list_display = ["name", "title"]
       list_filter = ["status"]
       search_fields = ["name"]
   ```
3. **CRUD out of the box:** Create, Read, Update, Delete für jedes registrierte Model.
4. **Customization ohne Code schreiben:**
   - `list_display`, `list_filter`, `search_fields`
   - `date_hierarchy` (autom. Kalender-Drilldown)
   - `ordering`, `list_per_page`, `list_editable` (Inline-Editing in der Liste!)
   - `readonly_fields`, `fields`, `fieldsets` (Field-Grouping mit Collapse)
   - `prepopulated_fields` (Slug aus Title generieren)
   - `autocomplete_fields` (Select2 mit Async-Suche)
   - `raw_id_fields` (Popup-Selector für ForeignKeys)
   - `filter_horizontal` / `filter_vertical` (Zwei-Box-Widget für ManyToMany)
   - `inlines` (Child-Models im Parent-Form bearbeiten — TabularInline, StackedInline)
5. **Admin Actions:** Batch-Operationen (z.B. "Mark as published" für ausgewählte Objekte).
6. **Custom Templates:** Admin-Seiten komplett überschreibbar.
7. **Facetten-Counts:** Zeigt Anzahl pro Filter-Option an.
8. **Responsive Design:** Funktioniert auf Mobile.

#### Für Dreego:

- **`dreego-admin` ist ein MUST-HAVE-Addon.**
- Django Admin ist **der Standard**, an dem sich jedes Admin-Generator-Tool messen muss.
- Automatische CRUD-UI aus Go-Structs generieren.
- `list_display`, `list_filter`, `search_fields` als Struct-Tag-Optionen:
  ```go
  type Post struct {
      Title string `admin:"list_display,search"`
      Status string `admin:"list_filter"`
  }
  ```

### 3.2 Django ORM und Migrations

Django ORM Features:
- **Model Definition:** Python-Klasse → DB-Tabelle automatisch.
- **Field Types:** 40+ Built-in Fields (`CharField`, `TextField`, `IntegerField`, `FileField`, `ImageField`, `JSONField`, etc.).
- **Relationships:** `ForeignKey`, `ManyToManyField`, `OneToOneField`.
- **Migrations:** `makemigrations` (auto-detect Changes), `migrate` (apply), `sqlmigrate` (SQL-Vorschau).
- **Manager:** `objects` — Query-Interface pro Model.
- **QuerySet:** Lazy, chainable, immutable:
  ```python
  Post.objects.filter(status='published').exclude(deleted=True).order_by('-created_at')
  ```
- **Q-Objects:** Komplexe Queries mit `|` (OR), `&` (AND), `~` (NOT).
- **F-Expressions:** DB-Level-Operationen `Post.objects.update(views=F('views') + 1)`.
- **Aggregation:** `Count`, `Sum`, `Avg`, `Max`, `Min`.
- **select_related (JOIN) / prefetch_related (separate Query):** Eager Loading.
- **Database Router:** Multi-DB-Support.
- **Transactions:** `transaction.atomic()`.
- **Signals:** `pre_save`, `post_save`, `pre_delete`, `post_delete`.

#### Für Dreego adoptieren:

| Django-Feature | Dreego-Adaption |
|---|---|
| Auto-Migration Detection | `dreego-db`: Go-Struct-Änderungen → Migration-SQL generieren |
| QuerySet (lazy, chainable) | Query-Builder mit Method-Chaining |
| F-Expressions | DB-Level-Updates |
| select_related / prefetch_related | Explizite Preloads |
| Signals | Model-Events (`AfterCreate`, `BeforeDelete`) |
| Multi-DB | Connection-String pro Model |

### 3.3 Django Template Engine (DTL)

DTL-Features:
- **Variables:** `{{ variable }}` — auto-escaping (XSS-Schutz).
- **Filters:** `{{ name|lower }}`, `{{ date|date:"Y-m-d" }}`, `{{ text|truncatechars:100 }}`.
- **Tags:** `{% if %}`, `{% for %}`, `{% block %}`, `{% extends %}`, `{% include %}`.
- **Template Inheritance:**
  ```django
  {% extends "base.html" %}
  {% block content %}
  {% endblock %}
  ```
- **Context Processors:** Request-Daten automatisch im Template verfügbar.
- **Custom Tags & Filters:** Erweiterbar.
- **Template Loaders:** Dateisystem, App-Verzeichnisse, Cached.
- **Jinja2 Support:** Django kann auch Jinja2 als Engine nutzen.

**ABER:** DTL ist bewusst **eingeschränkt** — keine komplexe Logik in Templates (Philosophie: Logik gehört in Views, nicht Templates). Weniger mächtig als Blade oder HEEx. **Keine Components, keine Slots.**

#### Für Dreego: Dreego's Template-Engine (Svelte-inspiriert) ist bereits fortschrittlicher als DTL. Django's Template Engine ist kein Vorbild für Dreego — Blade und HEEx sind die besseren Referenzen.

### 3.4 "Batteries Included"

Django's Built-in Apps (in `INSTALLED_APPS`):
- `django.contrib.admin` — Admin-Panel
- `django.contrib.auth` — User, Groups, Permissions
- `django.contrib.contenttypes` — Generic Relations
- `django.contrib.sessions` — Session-Management
- `django.contrib.messages` — Flash-Messages
- `django.contrib.staticfiles` — Static-File-Management
- `django.contrib.sites` — Multi-Site-Framework
- `django.contrib.sitemaps` — Sitemap-Generierung
- `django.contrib.syndication` — RSS-Feed-Generierung
- `django.contrib.gis` — GeoDjango (PostGIS)
- `django.contrib.postgres` — PostgreSQL-spezifische Features

#### Für Dreego:
- **Core:** Router, Context, Middleware, Template-Engine.
- **Offizielle Addons:** Auth, Admin, Sessions, Storage, Cache, Mail, Jobs.
- **Community-Addons:** Alles andere.

### 3.5 Authentication System

Django Auth — **komplettestes Built-in Auth-System aller Frameworks**:

1. **User Model:** `AbstractUser` (username, password, email, is_staff, is_superuser) — erweiterbar.
2. **Permissions:** `User.has_perm('app.action_model')`, Groups.
3. **Authentication Backends:** ModelBackend (Username+Password), RemoteUserBackend, custom Backends.
4. **Session-Authentication:** Cookie-basiert.
5. **Login/Logout Views:** `LoginView`, `LogoutView` (fertig).
6. **Password Management:**
   - `PasswordChangeView`, `PasswordChangeDoneView`
   - `PasswordResetView`, `PasswordResetDoneView`, `PasswordResetConfirmView`, `PasswordResetCompleteView`
   - Passwort-Hasher: PBKDF2 (default), Argon2, BCrypt
7. **Password Validators:** Minimum-Länge, Common-Password-Check, Numeric-Only-Check, User-Attribute-Similarity-Check.
8. **Login-Rate-Limiting:** `AXES` (Drittanbieter, aber Standard).
9. **User Creation Forms:** `UserCreationForm`, `UserChangeForm`.
10. **`login_required` Decorator:** `@login_required` auf Views.
11. **`user_passes_test`:** Custom Permission-Checks.
12. **Mixin-Klassen:** `LoginRequiredMixin`, `PermissionRequiredMixin`, `UserPassesTestMixin`.

#### Für Dreego:

- **`dreego-auth` Addon mit:**
  - Go-Interface: `type User interface { GetID() string; HasRole(string) bool }`
  - Session-basierte Auth (Cookie)
  - Middleware: `auth.Required`, `auth.HasRole("admin")`
  - Password-Hashing: Argon2id via `golang.org/x/crypto`
  - Login/Register/Reset-Routen
  - OAuth2/OIDC via Drittanbieter (`dreego-oauth`)
  - CSRF-Schutz im Core, von Auth genutzt

### 3.6 Django REST Framework (DRF)

DRF ist der **Standard für API-Entwicklung in Django**:

1. **Serializers:** Model → JSON/Dict, mit Validierung:
   ```python
   class UserSerializer(serializers.ModelSerializer):
       class Meta:
           model = User
           fields = ['id', 'username', 'email']
   ```
2. **ViewSets:** CRUD-Logik in einer Klasse:
   ```python
   class UserViewSet(viewsets.ModelViewSet):
       queryset = User.objects.all()
       serializer_class = UserSerializer
   ```
3. **Routers:** Automatische URL-Konfiguration aus ViewSets:
   ```python
   router = routers.DefaultRouter()
   router.register(r'users', UserViewSet)
   ```
4. **Browsable API:** Interaktive API-Doku im Browser (GET/POST/PUT/DELETE).
5. **Authentication:** SessionAuth, TokenAuth, OAuth1/2, JWT (via Drittanbieter).
6. **Permissions:** `IsAuthenticated`, `IsAdminUser`, `DjangoModelPermissions`, custom.
7. **Throttling:** Rate Limiting (`AnonRateThrottle`, `UserRateThrottle`).
8. **Pagination:** `PageNumberPagination`, `LimitOffsetPagination`, `CursorPagination`.
9. **Filtering:** `SearchFilter`, `OrderingFilter`, `DjangoFilterBackend`.
10. **Versioning:** URL-Path, Query-Parameter, Accept-Header.
11. **Renderers:** JSON, Browsable API, JSONP, XML, YAML.
12. **Parsers:** JSON, Form, MultiPart, FileUpload.
13. **Schemas:** OpenAPI/Swagger Auto-Generation.

#### Für Dreego:

- **`dreego-api` Addon:**
  - ViewSet-ähnliches Pattern für REST-Endpoints:
    ```go
    type UserAPI struct {
        dreego.API
    }
    func (api *UserAPI) List() ([]User, error) { ... }
    func (api *UserAPI) Get(id string) (*User, error) { ... }
    ```
  - Auto-Router: `dreego.RegisterAPI(router, "/api/users", &UserAPI{})`
  - OpenAPI-Generation aus Go-Structs + Routes
  - Browsable API? Optional, aber nice-to-have.

### 3.7 Was Django nach 20+ Jahren relevant hält

1. **Stabilität:** Keine Breaking Changes ohne lange Deprecation-Pfade.
2. **Dokumentation:** Goldstandard — jede Klasse, jede Methode dokumentiert.
3. **Batteries Included:** Für 80% der Web-Apps reicht Django's Built-in-Features.
4. **Admin:** Ist immer noch ungeschlagen.
5. **Community:** Django Packages, DjangoCon, DSF.
6. **Sicherheit:** XSS, CSRF, SQL-Injection, Clickjacking — alles Built-in geschützt.
7. **Migrations:** War das erste Framework mit robustem Migration-System.
8. **ORM:** Gehört zu den besten ORMs (auch wenn Ecto besser ist).

#### Für Dreego:
- **Langlebigkeit anstreben:** Nicht hype-driven, sondern solide Architektur.
- **Dokumentation von Tag 1:** Go-doc + ausführliche Guides.
- **Stabile API:** `dreego.Context` Interface, Plugin-Interface → Versionierung.

---

### 3.8 Anti-Patterns von Django, die Dreego vermeiden sollte

| Anti-Pattern | Warum |
|---|---|
| **Monolithisches `settings.py`:** Alles in einer Datei, wird schnell unübersichtlich | Dreego: `dreego.config.json` pro Projekt, modular |
| **Django ORM ist langsam bei komplexen Queries:** suboptimales SQL bei `annotate()` + `aggregate()` | Dreego: Raw SQL Escape-Hatch, Query-Builder der gutes SQL erzeugt |
| **Template-Engine ist zu restriktiv:** Keine Components, keine Funktionen mit Argumenten | Dreego: Moderne Template-Engine mit Components und Slots |
| **Python GIL:** Begrenzt Concurrency | Dreego: Go's Goroutinen sind überlegen |
| **"Django way or the highway":** Sehr opinionated | Dreego: Konventionen ja, aber Override-Möglichkeiten |
| **Migrations können bei großen Projekten langsam werden:** 1000+ Migrationen | Dreego: Squash-Migrations von Anfang an unterstützen |
| **Kein Async-Support im ORM (bis Django 4.1):** async ORM ist noch nicht ausgereift | Dreego: Go ist von Natur aus concurrent |

---

## Zusammenfassung: Was Dreego von jedem Framework lernen sollte

### Von Phoenix (Elixir):
1. **LiveView-Architektur** → SSE-First Reaktivität via Datastar
2. **Ecto-Pattern** (Schema ≠ Query ≠ Repo ≠ Changeset) → Trennung von Concerns
3. **PubSub-System** → Event-Bus für Echtzeit-Updates
4. **Generator-Staffelung** (embedded → schema → context → live) → `dreego generate`
5. **Compile-Time Template Validation** → Transpiler prüft HTML & Go-Syntax
6. **Graceful Reconnection** → SSE mit auto-reconnect

### Von Laravel (PHP):
1. **Eloquent's Developer Experience** → Naming-Conventions, Eager Loading, Scopes
2. **Validation-System** → Form Requests, Field-Level Errors, `safe()` / `old()` API
3. **Artisan CLI** → 30+ `make`-Kommandos, Stub-Customization
4. **Blade Components & Slots** → `<x-*>` = `<dreego:*>`
5. **Queue/Job-System** → Job-Middleware, Batching, Chaining, Delayed Dispatch
6. **Notification-System** → Multi-Channel (Mail, DB, Slack)
7. **Ecosystem-First-Denken** → Plugin-Interface muss Addon-Ökosystem ermöglichen

### Von Django (Python):
1. **Admin-Panel** → `dreego-admin` Addon (höchste Priorität)
2. **Auth-System** → Password-Hashing, Reset-Flow, Permission-System
3. **DRF's ViewSet/Router-Pattern** → `dreego-api` Addon
4. **Migrations** → Auto-Detection aus Go-Structs
5. **Dokumentation** → Goldstandard als Vorbild
6. **Stabilität & Langlebigkeit** → Keine Breaking Changes ohne Deprecation

### Anti-Patterns (konsistent über alle Frameworks):
1. **Magic / Implicitness** → Dreego: Alles explizit (Go-Philosophie)
2. **Too opinionated** → Override-Pfade für Konventionen
3. **Monolithischer Core** → Core minimal, Addons für Batteries
4. **String-basierte APIs** → Type-Safety via Go-Structs
5. **Lazy Loading / Implicit N+1** → Explizite Preloads
