//go:build linux

package tests

import (
	"os"
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/target/wails"
)

func TestWailsRenderDoesNotOpenTCPListener(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	app := dreego.New()
	err := app.RegisterRender("/", dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		close(entered)
		<-release
		return dreego.Result{HTML: []byte("<main>Timer</main>")}, nil
	}))
	if err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	host, err := wails.New(app)
	if err != nil {
		t.Fatalf("wails.New: %v", err)
	}
	before := processTCPListeners(t)
	done := make(chan error, 1)
	go func() {
		_, renderErr := host.Render("/")
		done <- renderErr
	}()
	<-entered
	during := func() map[string]bool {
		defer close(release)
		return processTCPListeners(t)
	}()
	if err := <-done; err != nil {
		t.Fatalf("Render: %v", err)
	}
	for inode := range during {
		if !before[inode] {
			t.Fatalf("Wails render opened TCP listener socket inode %s", inode)
		}
	}
}

func processTCPListeners(t *testing.T) map[string]bool {
	t.Helper()
	listeners := tcpListenerInodes(t, "/proc/net/tcp")
	for inode := range tcpListenerInodes(t, "/proc/net/tcp6") {
		listeners[inode] = true
	}
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		t.Fatalf("read process descriptors: %v", err)
	}
	owned := make(map[string]bool)
	for _, entry := range entries {
		target, err := os.Readlink("/proc/self/fd/" + entry.Name())
		if err != nil || !strings.HasPrefix(target, "socket:[") {
			continue
		}
		inode := strings.TrimSuffix(strings.TrimPrefix(target, "socket:["), "]")
		if listeners[inode] {
			owned[inode] = true
		}
	}
	return owned
}

func tcpListenerInodes(t *testing.T, path string) map[string]bool {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	listeners := make(map[string]bool)
	for _, line := range strings.Split(string(contents), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 9 && fields[3] == "0A" {
			listeners[fields[9]] = true
		}
	}
	return listeners
}
