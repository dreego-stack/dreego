//go:build linux && !race

package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	wailsadapter "github.com/dreego-stack/dreego/adapter/wails"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const nativeWailsHelper = "DREEGO_WAILS_NATIVE_NO_TCP_HELPER"

func TestWailsNativeProcessHasNoTCPListener(t *testing.T) {
	if os.Getenv(nativeWailsHelper) == "1" {
		runNativeWailsHelper(t)
		return
	}
	display := startVirtualDisplay(t)
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve test executable: %v", err)
	}
	trace := filepath.Join(t.TempDir(), "listen.trace")
	var output lockedBuffer
	command := exec.Command("strace", "-f", "-e", "trace=listen", "-o", trace,
		executable, "-test.run=^TestWailsNativeProcessHasNoTCPListener$")
	command.Env = append(os.Environ(),
		nativeWailsHelper+"=1",
		"DISPLAY="+display,
		"GSK_RENDERER=cairo",
		"LIBGL_ALWAYS_SOFTWARE=1",
		"WEBKIT_DISABLE_DMABUF_RENDERER=1",
	)
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatalf("start native Wails helper: %v", err)
	}
	t.Cleanup(func() { _ = command.Process.Kill() })
	waitForWindow(t, display, "Dreego No TCP Test", command, &output)
	listeners := processTreeTCPListeners(t, command.Process.Pid)
	if len(listeners) != 0 {
		t.Fatalf("native Wails process opened TCP listeners: %v", listeners)
	}
	helperPID := directChildPID(t, command.Process.Pid)
	if err := syscall.Kill(helperPID, syscall.SIGINT); err != nil {
		t.Fatalf("stop native Wails helper: %v", err)
	}
	waitForProcess(t, command, &output)
	contents, err := os.ReadFile(trace)
	if err != nil {
		t.Fatalf("read listener trace: %v", err)
	}
	if strings.Contains(string(contents), "listen(") {
		t.Fatalf("native Wails process called listen:\n%s", contents)
	}
}

func runNativeWailsHelper(t *testing.T) {
	app := dreego.New()
	err := app.RegisterRender("/", dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte("<main><h1>Dreego No TCP Test</h1></main>")}, nil
	}))
	if err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	handler, err := wailsadapter.New(app)
	if err != nil {
		t.Fatalf("wails adapter: %v", err)
	}
	wailsApp := application.New(application.Options{Name: "Dreego No TCP Test", Assets: application.AssetOptions{Handler: handler}})
	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Dreego No TCP Test", URL: "/", Width: 640, Height: 480})
	if err := wailsApp.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
}

func startVirtualDisplay(t *testing.T) string {
	t.Helper()
	display := ":97"
	command := exec.Command("Xvfb", display, "-screen", "0", "1024x768x24", "-nolisten", "tcp")
	var output lockedBuffer
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatalf("start Xvfb: %v", err)
	}
	t.Cleanup(func() {
		_ = command.Process.Kill()
		_, _ = command.Process.Wait()
	})
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat("/tmp/.X11-unix/X97"); err == nil {
			return display
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("Xvfb did not become ready: %s", output.String())
	return ""
}

func waitForWindow(t *testing.T, display, title string, process *exec.Cmd, output *lockedBuffer) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		windowTree, _ := exec.Command("xwininfo", "-display", display, "-root", "-tree").CombinedOutput()
		if strings.Contains(string(windowTree), title) {
			return
		}
		if err := process.Process.Signal(syscall.Signal(0)); err != nil {
			t.Fatalf("native Wails helper exited before opening a window: %s", output.String())
		}
		time.Sleep(50 * time.Millisecond)
	}
	_ = process.Process.Kill()
	t.Fatalf("native Wails window did not appear: %s", output.String())
}

func directChildPID(t *testing.T, parent int) int {
	t.Helper()
	entries, err := os.ReadDir("/proc")
	if err != nil {
		t.Fatalf("read /proc: %v", err)
	}
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		contents, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "stat"))
		if err != nil {
			continue
		}
		fields := strings.Fields(string(contents)[strings.LastIndexByte(string(contents), ')')+1:])
		if len(fields) > 1 {
			processParent, _ := strconv.Atoi(fields[1])
			if processParent == parent {
				return pid
			}
		}
	}
	t.Fatalf("native Wails helper process not found")
	return 0
}

func waitForProcess(t *testing.T, command *exec.Cmd, output *lockedBuffer) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			return
		}
		exitError, ok := err.(*exec.ExitError)
		if ok {
			status, ok := exitError.ProcessState.Sys().(syscall.WaitStatus)
			if ok && status.Signaled() && status.Signal() == syscall.SIGINT {
				return
			}
		}
		t.Fatalf("native Wails helper shutdown: %v: %s", err, output.String())
	case <-time.After(5 * time.Second):
		_ = command.Process.Kill()
		t.Fatal("native Wails helper did not shut down")
	}
}

func processTreeTCPListeners(t *testing.T, root int) []string {
	t.Helper()
	pids := map[int]bool{root: true}
	for changed := true; changed; {
		changed = false
		entries, err := os.ReadDir("/proc")
		if err != nil {
			t.Fatalf("read /proc: %v", err)
		}
		for _, entry := range entries {
			pid, err := strconv.Atoi(entry.Name())
			if err != nil || pids[pid] {
				continue
			}
			contents, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "stat"))
			if err != nil {
				continue
			}
			fields := strings.Fields(string(contents)[strings.LastIndexByte(string(contents), ')')+1:])
			if len(fields) > 1 {
				parent, _ := strconv.Atoi(fields[1])
				if pids[parent] {
					pids[pid] = true
					changed = true
				}
			}
		}
	}
	listeners := tcpListenerInodes(t, "/proc/net/tcp")
	for inode := range tcpListenerInodes(t, "/proc/net/tcp6") {
		listeners[inode] = true
	}
	var found []string
	for pid := range pids {
		entries, _ := os.ReadDir(fmt.Sprintf("/proc/%d/fd", pid))
		for _, entry := range entries {
			target, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", pid, entry.Name()))
			if err != nil || !strings.HasPrefix(target, "socket:[") {
				continue
			}
			inode := strings.TrimSuffix(strings.TrimPrefix(target, "socket:["), "]")
			if listeners[inode] {
				found = append(found, fmt.Sprintf("pid=%d inode=%s", pid, inode))
			}
		}
	}
	return found
}
