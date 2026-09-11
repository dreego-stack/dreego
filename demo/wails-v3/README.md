# Wails v3 timer

This experimental reference application renders one Dreego page in a native
Wails v3 window without starting a Dreego HTTP listener. The timer works with
keyboard-accessible native buttons and announces only start, pause, reset, and
completion state changes.

From the repository root, generate the demo and build it through `smd`:

```sh
smd sh -c 'go build -o /tmp/dreego ./cli/dreego && cd demo && /tmp/dreego generate && go build -o /tmp/dreego-wails-timer ./cmd/wails-v3'
```

The Wails v3 native runtime still requires a graphical desktop. The container
proves generation and compilation; run the packaged application on a supported
desktop platform for visual and assistive-technology checks.
