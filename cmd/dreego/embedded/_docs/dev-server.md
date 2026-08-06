# Dev Server

`dreego dev` is a development helper that watches your `.dreego` files and rebuilds + restarts the server automatically on every change.

## Usage

```bash
dreego dev
```

On start it runs `generate` + `build` and launches the server, then watches the working directory for `.dreego` changes:

- It polls every 500 ms and compares file modification times.
- When a `.dreego` file is added, changed, or removed it regenerates, rebuilds, and **gracefully restarts** the server (SIGTERM + reap).
- A build error does **not** kill the watcher — the previously running server stays up and the next change retries the build.
- Press `Ctrl-C` to stop both the watcher and the server.

## What it is not

`dreego dev` is a **dev tool, not for production**. For hot reloading with live browser reload, see [Hot Reload](https://codeberg.org/dreego/dreego/src/branch/main/_docs/hot-reload.md) (Air / `entr`).

## See Also

- [CLI Reference](https://codeberg.org/dreego/dreego/src/branch/main/_docs/cli.md)
- [Hot Reload](https://codeberg.org/dreego/dreego/src/branch/main/_docs/hot-reload.md)
- [Docs Index](https://codeberg.org/dreego/dreego/src/branch/main/_docs/index.md)
