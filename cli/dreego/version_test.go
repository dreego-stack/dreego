package main

import (
	"testing"
)

// TestDreegoVersionInjected verifies that the version injected via -ldflags
// wins over everything else.
func TestDreegoVersionInjected(t *testing.T) {
	version = "v9.9.9"
	if got := dreegoVersion(); got != "v9.9.9" {
		t.Errorf("dreegoVersion() = %q, want %q (injected)", got, "v9.9.9")
	}
}

// TestDreegoVersionFallbackDev verifies the dev fallback when no version is
// injected and no build-info version is available (plain local build).
func TestDreegoVersionFallbackDev(t *testing.T) {
	version = ""
	got := dreegoVersion()
	if got == "" {
		t.Fatal("expected a non-empty version")
	}
	if got == "(devel)" {
		t.Errorf("expected resolved version, got build-info placeholder %q", got)
	}
}
