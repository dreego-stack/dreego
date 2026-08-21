package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentCallStringLiteralProp(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (message string)
<div><p>{{ message }}</p></div>`,
		"www/routes/get.dreego": `<div><@Greet message="hi"/></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<p>hi</p>") {
		t.Fatalf("expected string literal prop rendered as <p>hi</p>, got: %s", body)
	}
}

func TestComponentCallIntLiteralProp(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Counter.dreego": `Component Counter (count int)
<div><p>{{ count }}</p></div>`,
		"www/routes/get.dreego": `<div><@Counter count={42}/></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<p>42</p>") {
		t.Fatalf("expected int literal prop rendered as <p>42</p>, got: %s", body)
	}
}

func TestComponentCallWrongTypeLiteralProp(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (message string)
<div><p>{{ message }}</p></div>`,
		"www/routes/get.dreego": `<div><@Greet message={42}/></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for wrong-type literal prop, got success: %s", out)
	}
	if !strings.Contains(out, "Greet") {
		t.Fatalf("error must name component Greet, got: %s", out)
	}
	if !strings.Contains(out, "message") {
		t.Fatalf("error must name prop message, got: %s", out)
	}
	if !strings.Contains(out, "expected string, got int") {
		t.Fatalf("error must report expected string and got int, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") {
		t.Fatalf("error must reference the calling source path, got: %s", out)
	}
}

func TestComponentCallExprPropHTTP(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Score.dreego": `Component Score (value int)
<div><p>score: {{ value }}</p></div>`,
		"www/routes/get.dreego": `<go>value := 99</go>
<div><@Score value={value}/></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<p>score: 99</p>") {
		t.Fatalf("expected typed expression prop rendered as <p>score: 99</p>, got: %s", body)
	}
}
