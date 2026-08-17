package core

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultServerTimeoutsAreSecure(t *testing.T) {
	c := DefaultServerConfig()
	cases := map[string]struct {
		got  any
		zero any
	}{
		"ReadHeaderTimeout": {c.ReadHeaderTimeout, time.Duration(0)},
		"ReadTimeout":       {c.ReadTimeout, time.Duration(0)},
		"WriteTimeout":      {c.WriteTimeout, time.Duration(0)},
		"IdleTimeout":       {c.IdleTimeout, time.Duration(0)},
		"MaxHeaderBytes":    {c.MaxHeaderBytes, 0},
		"ShutdownTimeout":   {c.ShutdownTimeout, time.Duration(0)},
	}
	for name, tc := range cases {
		if tc.got == tc.zero {
			t.Errorf("%s default not set", name)
		}
	}
}

func TestSetServerConfigBeforeBuild(t *testing.T) {
	app := New()
	custom := DefaultServerConfig()
	custom.ReadHeaderTimeout = 5 * time.Second
	custom.ShutdownTimeout = 3 * time.Second
	if err := app.SetServerConfig(custom); err != nil {
		t.Fatalf("SetServerConfig: %v", err)
	}
	app.Build()
	if err := app.SetServerConfig(custom); !errors.Is(err, ErrAppBuilt) {
		t.Fatalf("expected ErrAppBuilt after build, got %v", err)
	}
}
