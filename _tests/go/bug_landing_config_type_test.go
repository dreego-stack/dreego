package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLandingConfigType(t *testing.T) {
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	config, err := os.ReadFile(filepath.Join(dir, "testapp/dreego/config.json"))
	if err != nil {
		t.Fatalf("read config.json: %v", err)
	}
	if !strings.Contains(string(config), `"logging": {`) {
		t.Fatalf("landing config.json logging field has wrong type (B8): %s", config)
	}
}
