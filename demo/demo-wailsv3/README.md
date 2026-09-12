# Wails v3 timer

This experimental reference application renders Dreego pages in a native Wails
v3 window without starting a Dreego HTTP listener. Its timer state belongs to
an explicitly registered Go service. Wails-generated JavaScript calls that
service while committed TypeScript declarations describe the same contract.
Navigation and binding modules are served as in-process assets.

From the repository root, generate and build the demo through Task:

```sh
smd sh -c 'cd cmd/dreego && go build -o /tmp/dreego .'
smd sh -c 'cd demo/demo-wailsv3 && task generate'
smd sh -c 'cd demo/demo-wailsv3 && task build'
```

Regenerate bindings after changing `TimerService`:

```sh
smd sh -c 'cd demo/demo-wailsv3 && task bindings'
```

No `package.json`, npm install, Wails dev server, or Dreego HTTP server is part
of this workflow. Development reload is deliberately deterministic: stop the
binary, run `task generate`, then run `task build` and restart.

Platform tasks run on their native hosts:

```sh
task darwin:build
task windows:build
task ios:run
```

## Accessibility checks

- Tab reaches Start/Pause, Reset, About, and Back in document order.
- Enter and Space activate both timer buttons; Enter activates links.
- Focus remains visible through the high-contrast `:focus-visible` outline.
- The timer has a stable accessible name and does not announce every tick.
- The polite status region announces start, pause, reset, and completion.
- Page titles and one `h1` identify both routes to screen-reader users.

The application should be checked with the platform screen reader and keyboard
before distribution. Dreego supplies accessible reference semantics but cannot
make arbitrary application content accessible automatically.

## Packaging and troubleshooting

Phase 1 proves a native binary; signed installers and distribution packaging
are deferred until Wails v3 is stable. Build on the target operating system
with its documented Wails dependencies. If the window is blank, rerun
`dreego generate` and confirm the generated binding modules exist. The
application owns any Wails development-server choice; this reference app uses
only the in-process handler.

The Wails v3 native runtime still requires a graphical desktop. The container
proves generation and compilation; run the packaged application on a supported
desktop platform for visual and assistive-technology checks.
