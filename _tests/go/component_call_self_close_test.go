package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentCallSelfCloseNoChildren(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      `<div><@Card title="Hello"/></div>`,
	})
}

func TestComponentCallSelfCloseWithChildrenError(t *testing.T) {
	t.Parallel()
	out, err := dreegotest.RunCLI(t, dreegotest.ProjectDir(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      "<div>\n  <@Card title=\"Hello\"/>\n  <p>extra</p>\n</div>",
	}), "generate")
	if err == nil {
		t.Fatal("expected generate error, got none")
	}
	want := "dreego/routes/get.dreego:2:3: Card: self-closing call must not contain children"
	if !strings.Contains(out, want) {
		t.Fatalf("expected error\n%s\ngot:\n%s", want, out)
	}
}

func TestComponentCallDefaultSlotFallback(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      `<div><@Card title="Hello"></@Card></div>`,
	})
}

func TestComponentCallSelfCloseHTTP(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      `<div><@Card title="Hello"/></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("expected status 200, got %d", code)
	}
	for _, want := range []string{"<h2>Hello</h2>", "<div></div>"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %q, got:\n%s", want, body)
		}
	}
}

func TestComponentCallSelfCloseWhitespaceOnly(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      "<div><@Card title=\"Hello\"/>   \n\t\n</div>",
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("expected status 200, got %d", code)
	}
	if strings.Contains(body, "extra") {
		t.Fatalf("body should not contain unexpected text, got:\n%s", body)
	}
	if !strings.Contains(body, "<div></div>") {
		t.Fatalf("expected empty default slot, got:\n%s", body)
	}
}

func TestComponentCallSelfCloseNested(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Outer.dreego": "Component Outer ()\n<div><section>{#slot}</section></div>",
		"dreego/components/Inner.dreego": "Component Inner ()\n<div><span>inner</span></div>",
		"dreego/routes/get.dreego":       "<div><@Outer><@Inner/></@Outer></div>",
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("expected status 200, got %d", code)
	}
	if !strings.Contains(body, "<span>inner</span>") {
		t.Fatalf("expected inner self-closing component rendered, got:\n%s", body)
	}
}

func TestComponentCallSelfCloseNamedSlot(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Panel.dreego": "Component Panel ()\n<div><header>{#slot header}{/slot}</header><main>{#slot}</main></div>",
		"dreego/routes/get.dreego":       "<div><@Panel/></div>",
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("expected status 200, got %d", code)
	}
	if !strings.Contains(body, "<header></header>") {
		t.Fatalf("expected empty named header slot, got:\n%s", body)
	}
	if !strings.Contains(body, "<main></main>") {
		t.Fatalf("expected empty default slot, got:\n%s", body)
	}
}
