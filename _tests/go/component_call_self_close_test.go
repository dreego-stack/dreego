package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentCallSelfCloseNoChildren(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2><div>{#slot}</div></article></body>",
		"www/routes/get.dreego":      `<body><@Card title="Hello"/></body>`,
	})
}

func TestComponentCallDefaultSlotFallback(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2><div>{#slot}</div></article></body>",
		"www/routes/get.dreego":      `<body><@Card title="Hello"></@Card></body>`,
	})
}

func TestComponentCallSelfCloseHTTP(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2><div>{#slot}</div></article></body>",
		"www/routes/get.dreego":      `<body><@Card title="Hello"/></body>`,
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
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2><div>{#slot}</div></article></body>",
		"www/routes/get.dreego":      "<body><@Card title=\"Hello\"/>   \n\t\n</body>",
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
		"www/components/Outer.dreego": "Component Outer ()\n<body><section>{#slot}</section></body>",
		"www/components/Inner.dreego": "Component Inner ()\n<body><span>inner</span></body>",
		"www/routes/get.dreego":       "<body><@Outer><@Inner/></@Outer></body>",
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
		"www/components/Panel.dreego": "Component Panel ()\n<body><header>{#slot header}{/slot}</header><main>{#slot}</main></body>",
		"www/routes/get.dreego":       "<body><@Panel/></body>",
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
