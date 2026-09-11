# Phase: Wails v3 target

## Release strategy

Wails support is split because Dreego begins its integration while Wails v3 is
still beta. Phase 1 is planned as the v0.8 release slice. It pins one upstream
beta version and proves a small, opt-in desktop host. The release remains
provisional and does not imply support for the complete Wails API.

Phase 2 has no Dreego version assignment. It starts only after Wails v3 reaches
an upstream stable release and Phase 1 has produced compatibility evidence from
tests and a maintained reference application.

## Phase 1 goal

Render Dreego applications inside Wails without running a local HTTP server and
without requiring application developers to operate an npm build pipeline.
The same components and target-neutral App remain usable for web targets.

## Target package

The intended package is:

```text
github.com/dreego-stack/dreego/target/wails
```

The first-party target lives in the monorepo because it must coordinate render
results, embedded assets, navigation, DreeJS, diagnostics, and generated host
bindings. Provider-like desktop features may remain separate plugins.

Wails depends on the completed render foundation and follows the TypeScript and
Lua client processors. That sequence proves a shared JavaScript output and
asset pipeline before the desktop bridge adds another host lifecycle. Wails has
no dependency on static generation.

## Host model

The target owns:

- initial document rendering;
- asset delivery through Wails-supported mechanisms;
- navigation between Dreego routes;
- lifecycle integration with window startup and shutdown;
- generated, typed Go-to-client bindings;
- target capability declarations;
- development reload behavior.

It does not emulate HTTP with a hidden localhost server. HTTP request methods,
cookies, response headers, and redirects cannot silently retain SSR semantics.
Navigation, persistence, and identity receive explicit desktop equivalents.

## Client bridge

Wails expects browser-side bridge code even when business logic is written in
Go. Dreego should generate and manage that bridge. Developers may use
`<client lang="js">` or a processor such as TypeScript for presentation logic,
but no project-owned npm configuration is required for the standard path.

Bindings must be generated from explicit exported contracts. Avoid dynamic
string calls where Go types can produce client declarations and diagnostics.
External values crossing the bridge are validated at the boundary.

## Context and capabilities

Wails can provide:

- rendered HTML and client assets;
- a desktop host bridge;
- navigation;
- local persistence through explicit APIs;
- window and application lifecycle.

It does not inherently provide:

- HTTP request or response access;
- SSR middleware;
- browser cookies with server semantics;
- public server routes;
- multi-instance distributed state.

Components that require unavailable capabilities fail generation with a
specific diagnostic. A plugin can provide Wails-only functionality without
making unrelated web builds depend on Wails.

## Phase 1 implementation slices

1. Render one component as an initial Wails document without HTTP.
2. Load scoped styles and embedded static assets.
3. Navigate between literal Dreego routes.
4. Generate one typed Go method binding and validate boundary errors.
5. Establish the client asset and bridge boundary that DreeJS can use later.
6. Add development reload and a reference desktop application.

## Phase 1 acceptance criteria

- The reference application opens and renders without a listening TCP socket.
- The same component is exercised under SSR and Wails.
- Styles, head entries, and assets behave deterministically.
- Typed bindings surface incompatible values at generation or build time.
- Shutdown releases renderer, plugin, and bridge resources.
- A developer can build the reference application without creating a
  `package.json` or invoking npm manually.
- Keyboard navigation, focus management, and screen-reader semantics are part
  of the reference application's quality gate.

## Deferred Phase 2

Phase 2 may harden and expand:

- coverage of stable Wails v3 application and window APIs;
- broader typed bridge integration without granting ambient privileges;
- packaging, distribution, platform compatibility, and performance;
- development tooling informed by Phase 1 usage;
- Wails-specific live bridge updates where DreeJS proves a shared need.

Before Phase 2 starts, the pinned dependency must move to a stable Wails v3
release, the complete Phase 1 suite must pass against it, and an architecture
review must record any required migration.

## Risks

- Wails version changes can create a moving host boundary. Keep version-specific
  code inside the target package and document the supported range.
- Browser navigation assumptions may not match desktop history and window
  behavior. Specify the navigation contract before exposing it.
- A broad bridge can expose dangerous desktop capabilities. Generate only
  explicitly registered methods and preserve least privilege.
- A beta dependency can change before general availability. Keep Phase 1
  opt-in, pin its exact version, and contain version-specific code in the host.

## Not in this phase

- Replacing Wails or hiding all of its application lifecycle.
- Mobile support.
- An HTTP compatibility server.
- Automatic access to filesystem, shell, clipboard, or window APIs without
  explicit application registration.
