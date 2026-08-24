package tests

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDevWatcherKillsServerThatIgnoresSIGTERM(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body>first</body>`,
	})
	mainGo := `package main
import (
	"os/signal"
	"syscall"
	"time"
)
func main() {
	signal.Ignore(syscall.SIGTERM)
	for { time.Sleep(time.Second) }
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}

	reader, writer := io.Pipe()
	cmd := exec.Command(dreegotest.CLIBin(t), "dev")
	cmd.Dir = dir
	cmd.Stdout = writer
	cmd.Stderr = writer
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		_ = writer.Close()
		close(done)
	}()
	t.Cleanup(func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	})

	lines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	waitForDevServerStart(t, lines, 5*time.Second)
	time.Sleep(750 * time.Millisecond)
	changed := filepath.Join(dir, "www/routes/get.dreego")
	if err := os.WriteFile(changed, []byte(`<body>second</body>`), 0644); err != nil {
		t.Fatal(err)
	}
	waitForDevServerStart(t, lines, 8*time.Second)
}

func waitForDevServerStart(t *testing.T, lines <-chan string, timeout time.Duration) {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case line, ok := <-lines:
			if !ok {
				t.Fatal("dreego dev exited before starting the server")
			}
			if strings.Contains(line, "dreego dev: server started") {
				return
			}
		case <-timer.C:
			t.Fatal("dreego dev did not restart within the timeout")
		}
	}
}
