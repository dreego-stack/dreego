package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		cmdNew(os.Args[2:])
	case "init":
		cmdInit(os.Args[2:])
	case "generate":
		cmdGenerate(os.Args[2:])
	case "build":
		cmdBuild(os.Args[2:])
	case "run":
		cmdRun(os.Args[2:])
	case "dev":
		cmdDev(os.Args[2:])
	case "docs":
		cmdDocs(os.Args[2:])
	case "fmt":
		cmdFmt(os.Args[2:])
	case "feedback":
		cmdFeedback()
	case "version", "--version", "-v":
		fmt.Println(dreegoVersion())
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`dreego — Go Web Framework CLI (dev tools, not for production)

usage: dreego <command> [flags]

commands:
  new <name>             create a new project from landing template
  init <path>            create a minimal dreego project from blueprint
  generate [--force] [--check] transpile .dreego files to Go code
  fmt [--check] [--stdout] [path]  format .dreego files (like gofmt)
  build [--target <os/arch>] generate + go build → build/bin/<name>
  run [-d] [-t <seconds>] build + start server (dev only)
  dev                    watch .dreego files, rebuild + restart on change
  docs [-p <name>] [--web] [--json] [--dump] [--list] [path]  local docs (default: core /_docs/index.md)
  feedback               open browser to submit feedback/issue
  version, --version, -v  show the dreego CLI version
  help                   show this help

flags:
  --force                force regeneration of all files
  --target <os/arch>     cross-compile target (e.g. linux/amd64, darwin/arm64)
  -p <name>              docs of a dreego plugin (github.com/dreego-stack/<name>)
  --web                  open docs in browser instead of terminal
  --list                 list all core + plugin doc pages
  -d                     debug mode: write logs to build/logs/<utc>.log
  -t <seconds>           auto-stop server after N seconds (timer)

examples:
  dreego new myapp            create project with landing page
  dreego generate             transpile changed .dreego files
  dreego generate --force     force full regeneration
  dreego build                generate + build binary (local platform)
  dreego build --target linux/amd64  cross-compile for Docker
  dreego run                  build + start server (foreground)
  dreego run -d               build + start + log to file
  dreego run -t 60            build + start + stop after 60s
  dreego run -d -t 60         debug log + 60s timer
  dreego dev                  watch + rebuild + restart on change
  dreego docs                 show core docs index (terminal)
  dreego docs -p plugin-sse   show plugin docs index
  dreego docs --list          list all core + plugin pages
  dreego docs --web           open docs index in browser
  dreego docs --json          structured JSON for AI agents
  dreego docs --dump          all docs for LLM context
  dreego feedback             submit issue / feedback
`)
}

func cmdGenerate(args []string) {
	force := false
	check := false
	for _, a := range args {
		if a == "--force" {
			force = true
		}
		if a == "--check" {
			check = true
		}
	}
	if !check {
		if err := dreego.Run(force); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	if check {
		var genFile string
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if genFile != "" || err != nil || info.IsDir() {
				return nil
			}
			if strings.Contains(path, "/gen/routes.go") {
				genFile = path
			}
			return nil
		})
		if genFile == "" {
			fmt.Fprintf(os.Stderr, "no generated files found, run dreego generate first\n")
			os.Exit(1)
		}
		genInfo, _ := os.Stat(genFile)
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".dreego") {
				return nil
			}
			if info.ModTime().After(genInfo.ModTime()) {
				fmt.Fprintf(os.Stderr, "stale: %s is newer than %s\n", path, genFile)
				os.Exit(1)
			}
			return nil
		})
		fmt.Println("generated code is up-to-date")
	}
}

// cmdBuild wraps cmdBuildE for the CLI, printing the error and exiting.
func cmdBuild(args []string) {
	if err := cmdBuildE(args); err != nil {
		fmt.Fprintf(os.Stderr, "build error: %v\n", err)
		os.Exit(1)
	}
}

// cmdBuildE runs generation and the go build, returning an error instead of
// calling os.Exit. Callers in the dev watcher use it so a build failure does
// not kill the process.
func cmdBuildE(args []string) error {
	target := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--target" && i+1 < len(args) {
			target = args[i+1]
			i++
		} else if strings.HasPrefix(args[i], "--target=") {
			target = strings.TrimPrefix(args[i], "--target=")
		}
	}

	if err := dreego.Run(false); err != nil {
		return err
	}

	projDir, pkg, name := findMain()
	outDir := filepath.Join(projDir, "build", "bin")
	os.MkdirAll(outDir, 0755)
	out := filepath.Join(outDir, name)

	if target != "" {
		parts := strings.SplitN(target, "/", 2)
		out += "-" + strings.ReplaceAll(target, "/", "-")
		fmt.Printf("cross-compiling %s → %s (target %s)\n", pkg, out, target)
		c := exec.Command("go", "build", "-o", out, "./"+pkg)
		c.Env = os.Environ()
		if len(parts) == 2 {
			c.Env = append(c.Env, "GOOS="+parts[0], "GOARCH="+parts[1], "CGO_ENABLED=0")
		}
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return err
		}
		fmt.Println("build ok")
		return nil
	}

	fmt.Printf("building %s → %s\n", pkg, out)
	c := exec.Command("go", "build", "-o", out, "./"+pkg)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	fmt.Println("build ok")
	return nil
}

func cmdRun(args []string) {
	debug := false
	timer := 0

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d":
			debug = true
		case "-t":
			if i+1 < len(args) {
				timer, _ = strconv.Atoi(args[i+1])
				i++
			}
		}
	}

	cmdBuild(nil)

	projDir, _, name := findMain()
	bin := filepath.Join(projDir, "build", "bin", name)

	fmt.Printf("starting %s", bin)
	if timer > 0 {
		fmt.Printf(" (auto-stop in %ds)", timer)
	}
	fmt.Println()

	c := exec.Command(bin)
	c.Stderr = os.Stderr

	if debug {
		logDir := filepath.Join(projDir, "build", "logs")
		os.MkdirAll(logDir, 0755)
		logFile := filepath.Join(logDir, time.Now().UTC().Format("2006-01-02T150405")+".log")
		f, err := os.Create(logFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "log error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		c.Stdout = f
		c.Stderr = f
		fmt.Printf("logging to %s\n", logFile)
	} else {
		c.Stdout = os.Stdout
	}

	if err := c.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "start error: %v\n", err)
		os.Exit(1)
	}

	if timer > 0 {
		go scheduleStop(c.Process, time.Duration(timer)*time.Second)
		c.Wait()
		return
	}

	c.Wait()
}

// scheduleStop sends SIGTERM to the process after the given delay (graceful
// shutdown, bug B20) and falls back to SIGKILL if signaling fails. Extracted
// from cmdRun so the timer behavior is testable without a full server run.
func scheduleStop(proc *os.Process, after time.Duration) {
	time.Sleep(after)
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		fmt.Fprintf(os.Stderr, "timer: signal error: %v\n", err)
		proc.Kill()
	} else {
		fmt.Println("timer: server stopped")
	}
}

func findMain() (projDir, pkg, name string) {
	return findMainIn(wd())
}

func findMainIn(dir string) (projDir, pkg, name string) {
	if _, err := os.Stat(filepath.Join(dir, "main.go")); err == nil {
		return ".", ".", filepath.Base(dir)
	}

	for _, d := range []string{"demo", "cmd"} {
		mp := filepath.Join(dir, d, "main.go")
		if _, err := os.Stat(mp); err == nil {
			return d, d, d
		}
	}

	return ".", ".", "server"
}

func wd() string {
	d, _ := os.Getwd()
	return d
}
