package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	dreego "codeberg.org/dreego/dreego/core"
)

// detectChanges scans dir for .dreego files and compares their modtimes
// against the previous map. It returns the changed files (relative paths)
// plus the updated mtime map. A file is "changed" when it is new, its
// modtime moved, or it disappeared from the previous map.
func detectChanges(dir string, mtimes map[string]time.Time) (changed []string, updated map[string]time.Time) {
	updated = make(map[string]time.Time)
	seen := make(map[string]bool)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		rel := strings.TrimPrefix(path, dir+string(os.PathSeparator))
		seen[rel] = true

		prev, ok := mtimes[rel]
		if !ok || !prev.Equal(info.ModTime()) {
			changed = append(changed, rel)
		}
		updated[rel] = info.ModTime()
		return nil
	})

	for name := range mtimes {
		if !seen[name] {
			changed = append(changed, name)
		}
	}

	return changed, updated
}

// shouldRestart reports whether a server restart is required. Any .dreego
// change needs codegen + rebuild, so it is true whenever changed is non-empty.
func shouldRestart(changed []string) bool {
	return len(changed) > 0
}

func cmdDev(args []string) {
	if err := dreego.Run(false); err != nil {
		fmt.Fprintf(os.Stderr, "generate error: %v\n", err)
		os.Exit(1)
	}
	if err := cmdBuildE(nil); err != nil {
		fmt.Fprintf(os.Stderr, "build error: %v\n", err)
		os.Exit(1)
	}

	projDir, _, name := findMain()
	bin := filepath.Join(projDir, "build", "bin", name)
	wd := wd()

	fmt.Printf("dreego dev: watching %s for .dreego changes (Ctrl-C to stop)\n", wd)

	var server *exec.Cmd
	server, err := startServer(bin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start error: %v\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Prime the mtime map before the loop so the first tick has a real
	// diff baseline instead of treating every .dreego file as new.
	_, mtimes := detectChanges(wd, map[string]time.Time{})
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			if server != nil && server.Process != nil {
				server.Process.Signal(syscall.SIGTERM)
				server.Process.Wait()
			}
			fmt.Println("\ndreego dev: stopped")
			return
		case <-ticker.C:
			changed, updated := detectChanges(wd, mtimes)
			mtimes = updated
			if !shouldRestart(changed) {
				continue
			}
			fmt.Printf("change detected: %v\n", changed)
			if err := cmdBuildE(nil); err != nil {
				fmt.Fprintf(os.Stderr, "build error: %v\n", err)
				continue
			}
			restartServer(&server, bin)
		}
	}
}

func startServer(bin string) (*exec.Cmd, error) {
	c := exec.Command(bin)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return nil, err
	}
	fmt.Printf("dreego dev: server started (%s)\n", bin)
	return c, nil
}

func restartServer(server **exec.Cmd, bin string) {
	if *server != nil && (*server).Process != nil {
		// Graceful stop: send TERM and reap the process with Wait(). This
		// blocks until the server exits; a misbehaving server could hang the
		// watcher, but that is acceptable for the dev tool.
		(*server).Process.Signal(syscall.SIGTERM)
		(*server).Wait()
	}
	c, err := startServer(bin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "restart error: %v\n", err)
		return
	}
	*server = c
}
