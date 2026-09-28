//go:build linux && !race

package tests

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
)

type lockedBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (b *lockedBuffer) Write(content []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.Write(content)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.String()
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
