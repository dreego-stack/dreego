# Wails v3 timer

This experimental reference application renders Dreego pages in a native Wails
v3 window without starting a Dreego HTTP listener. Its timer state belongs to
an explicitly registered Go service. Wails-generated JavaScript calls that
service while committed TypeScript declarations describe the same contract.
Navigation and binding modules are served as in-process assets.

From the repository root, generate the demo and build it through `smd`:

```sh
smd go build -o /tmp/dreego ./cli/dreego
smd sh -c 'cd demo && /tmp/dreego generate'
smd sh -c 'cd demo && go build -o /tmp/dreego-wails-timer ./cmd/wails-v3'
```

Regenerate bindings after changing `TimerService`:

```sh
smd sh -c 'cd demo && wails3 generate bindings -d wails-v3/bindings -ts -i -b ./wails-v3 ./cmd/wails-v3'
smd sh -c 'cd demo && wails3 generate bindings -d wails-v3/static/bindings -b -noevents ./wails-v3 ./cmd/wails-v3'
```

No `package.json`, npm install, Wails dev server, or Dreego HTTP server is part
of this workflow. Development reload is deliberately deterministic: stop the
binary, regenerate Dreego output and changed bindings, then rebuild and restart.

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
`dreego generate` and confirm the generated binding modules exist. If a service
fails before the window opens, follow the `Fix:` text in the binding diagnostic.
Unset `FRONTEND_DEVSERVER_URL`; Dreego rejects it to preserve the no-listener
contract.

The Wails v3 native runtime still requires a graphical desktop. The container
proves generation and compilation; run the packaged application on a supported
desktop platform for visual and assistive-technology checks.
