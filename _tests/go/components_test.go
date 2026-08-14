package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Cmp.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2></article></div>",
		"dreego/routes/get.dreego":     `<div><@Card title="Hello"/></div>`,
	})
}

func TestComponentEmptyProps(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Empty.dreego": "Component Empty ()\n<div><p>no props</p></div>",
		"dreego/routes/get.dreego":       `<div><@Empty/></div>`,
	})
}

func TestComponentMultiProps(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Profile.dreego": `Component Profile (name string, role string, email string)
<div><h2>{{ name }}</h2><p>{{ role }}</p><a href="mailto:{{ email }}">{{ email }}</a></div>`,
		"dreego/routes/get.dreego": `<div><@Profile name="Ada" role="Admin" email="ada@example.com"/></div>`,
	})
}

func TestComponentNameClash(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/button.dreego":        "Component FlatButton (label string)\n\n<div><button>{{ label }}</button></div>",
		"dreego/components/button/button.dreego": "Component NestedButton (label string)\n\n<div><button class=\"nested\">{{ label }}</button></div>",
		"dreego/routes/get.dreego":               `<div><@FlatButton label="Click"/><@NestedButton label="Go"/></div>`,
	})
}

func TestComponentNamedSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article>{#slot header}{/slot}<h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      `<div><@Card title="Hi">{#slot header}<strong>HEADER</strong>{/slot}<p>body</p></@Card></div>`,
	})
}

func TestComponentNamedSlotEmpty(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card ()\n<div><article>{#slot header}{/slot}</article></div>",
		"dreego/routes/get.dreego":      `<div><@Card><p>only default slot</p></@Card></div>`,
	})
}

func TestComponentNamedSlotExpr(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Greet.dreego": "Component Greet ()\n<div><p>{#slot header}{/slot} {#slot}</p></div>",
		"dreego/routes/get.dreego": `<go>name := "World"</go>
<div><@Greet>{#slot header}<strong>{{ name }}</strong>{/slot}!</@Greet></div>`,
	})
}

func TestComponentNamedSlotMulti(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Page.dreego": "Component Page (title string)\n<div><header>{#slot header}{/slot}</header><main>{#slot}</main><footer>{#slot footer}{/slot}</footer></div>",
		"dreego/routes/get.dreego": `<div><@Page title="Multi">
{#slot header}<nav>menu</nav>{/slot}
content
{#slot footer}<small>2026</small>{/slot}
</@Page></div>`,
	})
}

func TestComponentNested(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Inner.dreego": "Component Inner ()\n<div><span>inner</span></div>",
		"dreego/components/Outer.dreego": "Component Outer ()\n<div><article><@Inner/></article></div>",
		"dreego/routes/get.dreego":       `<div><@Outer/></div>`,
	})
}

func TestComponentNotFound(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><@Missing/></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if dreegotest.BuildInDirOK(t, dir) {
		t.Fatal("expected build failure but succeeded")
	}
}

func TestComponentPropExpr(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Greet.dreego": "Component Greet (name string)\n<div><p>Hello {{ name }}</p></div>",
		"dreego/routes/get.dreego": `<go>n:="World"</go>
<div><@Greet name={n}/></div>`,
	})
}

func TestComponentPropExpression(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2></article></div>",
		"dreego/routes/get.dreego": `<go>type User struct { Name string }
user := User{Name: "Ada"}</go>
<div><@Card title={user.Name}/></div>`,
	})
}

func TestComponentScopedStyle(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Box.dreego": "Component Box ()\n<div><div class=\"box\"><p>scoped</p></div></div>\n<style>.box{border:1px solid red}</style>",
		"dreego/routes/get.dreego": `<head><title>T</title></head>
<div><@Box/><p class="box">unscoped</p></div>
<style>.box{color:blue}</style>`,
	})
}

func TestComponentSelfClosing(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Cmp.dreego": "Component Icon (name string)\n<div><i class=\"icon\">{{ name }}</i></div>",
		"dreego/routes/get.dreego":     `<div><@Icon name="star"/></div>`,
	})
}

func TestComponentWithGo(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Greeting.dreego": "Component Greeting (name string)\n<go>msg := \"Hi \" + name</go>\n<div><p>{{ msg }}</p></div>",
		"dreego/routes/get.dreego":          `<div><@Greeting name="World"/></div>`,
	})
}

func TestComponentWithSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>",
		"dreego/routes/get.dreego":      `<div><@Card title="Welcome"><p>body text</p></@Card></div>`,
	})
}
