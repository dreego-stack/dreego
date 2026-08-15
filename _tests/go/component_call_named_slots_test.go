package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentCallNamedSlotRender(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article>{#slot header}{/slot}<h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      `<div><@Card title="Hi">{#slot header}<strong>HEADER</strong>{/slot}<p>body</p></@Card></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<strong>HEADER</strong>") {
		t.Fatalf("named header slot not rendered, got: %s", body)
	}
	if !strings.Contains(body, "<p>body</p>") {
		t.Fatalf("default slot not rendered, got: %s", body)
	}
	if !strings.Contains(body, "<h2>Hi</h2>") {
		t.Fatalf("title prop not rendered, got: %s", body)
	}
}

func TestComponentCallNamedSlotUnknownError(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card ()\n<div><article>{#slot header}{/slot}<div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      "<div>\n  <@Card>{#slot footer}<p>extra</p>{/slot}</@Card>\n</div>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate error for unknown slot, got success: %s", out)
	}
	want := "dreego/routes/get.dreego:2:3: Card: unknown slot \"footer\""
	if !strings.Contains(out, want) {
		t.Fatalf("expected error\n%s\ngot:\n%s", want, out)
	}
}

func TestComponentCallNamedSlotNestedDeclarationError(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card ()\n<div><article>{#slot}</article></div>",
		"dreego/routes/get.dreego":      `<div><@Card>{#slot header}{#slot footer}<p>x</p>{/slot}{/slot}</@Card></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate error for nested slot declaration, got success: %s", out)
	}
	want := "dreego/routes/get.dreego:1:27: Card: nested slot declaration is not allowed"
	if !strings.Contains(out, want) {
		t.Fatalf("expected error\n%s\ngot:\n%s", want, out)
	}
}

func TestComponentCallNamedSlotSiblingIsolation(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2>{#slot header}{/slot}<div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego": `<div>
<@Card title="First">{#slot header}<strong>only first</strong>{/slot}<p>first body</p></@Card>
<@Card title="Second"/>
</div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<h2>First</h2>") {
		t.Fatalf("first card title missing, got: %s", body)
	}
	if !strings.Contains(body, "<h2>Second</h2>") {
		t.Fatalf("second card title missing, got: %s", body)
	}
	if !strings.Contains(body, "<strong>only first</strong>") {
		t.Fatalf("first card named slot missing, got: %s", body)
	}
	if !strings.Contains(body, "<p>first body</p>") {
		t.Fatalf("first card default slot missing, got: %s", body)
	}
	before := strings.Index(body, "<h2>Second</h2>")
	if before < 0 {
		t.Fatalf("second card title missing, got:\n%s", body)
	}
	after := strings.Index(body[before:], "</article></div>") + before
	if after <= before {
		t.Fatalf("could not locate second card body, got:\n%s", body)
	}
	second := body[before:after]
	if strings.Contains(second, "<strong>only first</strong>") {
		t.Fatalf("second card leaked named slot from first, got:\n%s", second)
	}
	if strings.Contains(second, "<p>first body</p>") {
		t.Fatalf("second card leaked default slot from first, got:\n%s", second)
	}
}

func TestComponentCallNestedComponentInNamedSlot(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Icon.dreego": "Component Icon (name string)\n<div><span class=\"icon\">{{ name }}</span></div>",
		"dreego/components/Card.dreego": "Component Card ()\n<div><article>{#slot header}{/slot}</article></div>",
		"dreego/routes/get.dreego":      `<div><@Card>{#slot header}<@Icon name="star"/>{/slot}</@Card></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, `<span class="icon">star</span>`) {
		t.Fatalf("nested component inside named slot not rendered, got: %s", body)
	}
}

func TestComponentCallNamedSlotHTTP(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Page.dreego": "Component Page ()\n<div><header>{#slot header}{/slot}</header><main>{#slot}</main></div>",
		"dreego/routes/get.dreego":      `<div><@Page>{#slot header}<nav>menu</nav>{/slot}<p>content</p></@Page></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<nav>menu</nav>") {
		t.Fatalf("named slot output missing in HTTP body, got: %s", body)
	}
	if !strings.Contains(body, "<header>") || !strings.Contains(body, "</header>") {
		t.Fatalf("named slot wrapper missing in HTTP body, got: %s", body)
	}
	if !strings.Contains(body, "<p>content</p>") {
		t.Fatalf("default slot output missing in HTTP body, got: %s", body)
	}
}
