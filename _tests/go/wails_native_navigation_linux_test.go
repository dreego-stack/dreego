//go:build linux

package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	wailsadapter "github.com/dreego-stack/dreego/adapter/wails"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const nativeNavigationHelper = "DREEGO_WAILS_NATIVE_NAVIGATION_HELPER"
const nativeNavigationTrace = "DREEGO_WAILS_NATIVE_NAVIGATION_TRACE"

func TestWailsNativeNavigationUsesLiteralHistory(t *testing.T) {
	if os.Getenv(nativeNavigationHelper) == "1" {
		runNativeNavigationHelper(t)
		return
	}
	display := startVirtualDisplay(t)
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve test executable: %v", err)
	}
	trace := filepath.Join(t.TempDir(), "navigation.trace")
	var output bytes.Buffer
	command := exec.Command(executable, "-test.run=^TestWailsNativeNavigationUsesLiteralHistory$")
	command.Env = append(os.Environ(),
		nativeNavigationHelper+"=1",
		nativeNavigationTrace+"="+trace,
		"DISPLAY="+display,
		"WEBKIT_DISABLE_SANDBOX_THIS_IS_DANGEROUS=1",
	)
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatalf("start native navigation helper: %v", err)
	}
	t.Cleanup(func() { _ = command.Process.Kill() })
	waitForNavigationTrace(t, trace, []string{"root", "settings", "root-back"}, command, &output)
	if err := command.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("stop native navigation helper: %v", err)
	}
	waitForProcess(t, command, &output)
}

func runNativeNavigationHelper(t *testing.T) {
	trace := os.Getenv(nativeNavigationTrace)
	app := dreego.New()
	registerNavigationPage(t, app, "/", trace, "root", `
<script>
if (sessionStorage.getItem("visited-settings") === null) {
  window.addEventListener("pageshow", (event) => {
    if (event.persisted) fetch("/root-back")
  })
  sessionStorage.setItem("visited-settings", "yes")
  setTimeout(() => { location.href = "/settings" }, 100)
} else {
  fetch("/root-back")
}
</script>`)
	registerNavigationPage(t, app, "/settings", trace, "settings", `<script>setTimeout(() => history.back(), 500)</script>`)
	registerNavigationPage(t, app, "/root-back", trace, "root-back", "")
	handler, err := wailsadapter.New(app)
	if err != nil {
		t.Fatalf("wails adapter: %v", err)
	}
	wailsApp := application.New(application.Options{Name: "Dreego Navigation Test", Assets: application.AssetOptions{Handler: handler}})
	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Dreego Navigation Test", URL: "/", Width: 640, Height: 480})
	if err := wailsApp.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
}

func registerNavigationPage(t *testing.T, app *dreego.App, routePath, trace, marker, script string) {
	t.Helper()
	err := app.RegisterRender(routePath, dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		file, err := os.OpenFile(trace, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return dreego.Result{}, err
		}
		if _, err := file.WriteString(marker + "\n"); err != nil {
			_ = file.Close()
			return dreego.Result{}, err
		}
		if err := file.Close(); err != nil {
			return dreego.Result{}, err
		}
		return dreego.Result{HTML: []byte("<!doctype html><title>" + marker + "</title><main>" + marker + "</main>" + script)}, nil
	}))
	if err != nil {
		t.Fatalf("RegisterRender %s: %v", routePath, err)
	}
}

func waitForNavigationTrace(t *testing.T, trace string, want []string, process *exec.Cmd, output *bytes.Buffer) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		contents, _ := os.ReadFile(trace)
		if markersInOrder(string(contents), want) {
			return
		}
		if err := process.Process.Signal(syscall.Signal(0)); err != nil {
			t.Fatalf("native navigation helper exited early: %s", output.String())
		}
		time.Sleep(50 * time.Millisecond)
	}
	contents, _ := os.ReadFile(trace)
	t.Fatalf("native navigation trace = %q, want %q: %s", contents, want, output.String())
}

func markersInOrder(trace string, markers []string) bool {
	lines := strings.Fields(trace)
	next := 0
	for _, line := range lines {
		if next < len(markers) && line == markers[next] {
			next++
		}
	}
	return next == len(markers)
}
