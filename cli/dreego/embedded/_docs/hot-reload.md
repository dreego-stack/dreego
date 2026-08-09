# Hot Reload (with Air)

Dreego uses [Air](https://github.com/air-verse/air) for hot reload during development. Air watches your `.dreego` files, runs `dreego generate`, rebuilds, and restarts the server automatically.

## Setup

```bash
go install github.com/air-verse/air@latest
```

## Config

Create `.air.toml` in your project root:

```toml
root = "."
tmp_dir = "build/air"

[build]
  cmd = "dreego generate && go build -o ./build/air/main . && ./build/air/main"
  bin = "./build/air/main"
  full_bin = "APP_ENV=dev ./build/air/main"
  watch_dir = "dreego"
  include_ext = ["go", "dreego", "json"]
  exclude_dir = ["dreego/gen", "dreego/static", "build"]
  delay = 500

[misc]
  clean_on_exit = true
```

## Run

```bash
air
```

On any `.dreego` file change, Air will:
1. Run `dreego generate` (transpile changed files)
2. Run `go build` (compile)
3. Kill the old server process
4. Start the new server process

## Without Air

For quick ad-hoc reload without installing Air:

```bash
while true; do dreego generate && go build -o ./tmp/srv . && ./tmp/srv; sleep 2; done
```

Or use `dreego run -t 0` with an external file watcher like `entr`:

```bash
find dreego -name '*.dreego' | entr -r dreego run
```
