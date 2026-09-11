package core

import (
	"errors"
	"testing"
)

func TestStaticAssetReturnsIndependentCopy(t *testing.T) {
	app := New()
	if err := app.RegisterStatic("/assets/app.js", "application/javascript", []byte("ready()")); err != nil {
		t.Fatalf("RegisterStatic: %v", err)
	}
	asset, err := app.StaticAsset("/assets/app.js")
	if err != nil {
		t.Fatalf("StaticAsset: %v", err)
	}
	asset.Content[0] = 'X'
	again, err := app.StaticAsset("/assets/app.js")
	if err != nil {
		t.Fatalf("StaticAsset again: %v", err)
	}
	if string(again.Content) != "ready()" || again.MIME != "application/javascript" {
		t.Fatalf("asset = %#v", again)
	}
}

func TestStaticAssetRejectsUnknownAndTraversalPaths(t *testing.T) {
	app := New()
	for _, assetPath := range []string{"/missing", "/../secret", `/assets\secret`, "/assets/%2e%2e/secret"} {
		if _, err := app.StaticAsset(assetPath); !errors.Is(err, ErrStaticAssetNotFound) {
			t.Errorf("StaticAsset(%q) error = %v", assetPath, err)
		}
	}
}

func TestRegisterStaticRejectsTraversalPaths(t *testing.T) {
	app := New()
	for _, assetPath := range []string{"../secret", "/../secret", `/assets\secret`, "/assets/%2e%2e/secret"} {
		if err := app.RegisterStatic(assetPath, "text/plain", []byte("secret")); err == nil {
			t.Errorf("RegisterStatic(%q) succeeded", assetPath)
		}
	}
}
