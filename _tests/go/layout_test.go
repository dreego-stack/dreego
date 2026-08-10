package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestLayoutNoLayout(t *testing.T) {
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hello no layout</p></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "hello no layout") {
		t.Fatalf("route content missing in body: %s", body)
	}
	if strings.Contains(body, "<html>") {
		t.Fatalf("body unexpectedly contains layout html: %s", body)
	}
}

func TestLayoutWithHead(t *testing.T) {
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestLayoutWithSlot(t *testing.T) {
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}
