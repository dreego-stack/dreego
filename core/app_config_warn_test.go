package core

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestWarnMissingSessionStore(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(old)

	app := New()
	app.warnMissingSessionStore()

	out := buf.String()
	if !strings.Contains(out, "no session store is configured") {
		t.Errorf("expected warning about missing session store, got %q", out)
	}
}
