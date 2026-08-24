package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Cmp.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2></article></body>",
		"www/routes/get.dreego":     `<body><@Card title="Hello"/></body>`,
	})
}

func TestComponentEmptyProps(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Empty.dreego": "Component Empty ()\n<body><p>no props</p></body>",
		"www/routes/get.dreego":       `<body><@Empty/></body>`,
	})
}

func TestComponentMultiProps(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Profile.dreego": `Component Profile (name string, role string, email string)
<body><h2>{{ name }}</h2><p>{{ role }}</p><a href="mailto:{{ email }}">{{ email }}</a></body>`,
		"www/routes/get.dreego": `<body><@Profile name="Ada" role="Admin" email="ada@example.com"/></body>`,
	})
}

func TestComponentNameClash(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/button.dreego":        "Component FlatButton (label string)\n\n<body><button>{{ label }}</button></body>",
		"www/components/button/button.dreego": "Component NestedButton (label string)\n\n<body><button class=\"nested\">{{ label }}</button></body>",
		"www/routes/get.dreego":               `<body><@FlatButton label="Click"/><@NestedButton label="Go"/></body>`,
	})
}

func TestComponentNamedSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article>{#slot header}{/slot}<h2>{{ title }}</h2><div>{#slot}</div></article></body>",
		"www/routes/get.dreego":      `<body><@Card title="Hi">{#slot header}<strong>HEADER</strong>{/slot}<p>body</p></@Card></body>`,
	})
}

func TestComponentNamedSlotEmpty(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card ()\n<body><article>{#slot header}{/slot}</article></body>",
		"www/routes/get.dreego":      `<body><@Card><p>only default slot</p></@Card></body>`,
	})
}

func TestComponentNamedSlotExpr(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Greet.dreego": "Component Greet ()\n<body><p>{#slot header}{/slot} {#slot}</p></body>",
		"www/routes/get.dreego": `<server>name := "World"</server>
<body><@Greet>{#slot header}<strong>{{ name }}</strong>{/slot}!</@Greet></body>`,
	})
}

func TestComponentNamedSlotMulti(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Page.dreego": "Component Page (title string)\n<body><header>{#slot header}{/slot}</header><main>{#slot}</main><footer>{#slot footer}{/slot}</footer></body>",
		"www/routes/get.dreego": `<body><@Page title="Multi">
{#slot header}<nav>menu</nav>{/slot}
content
{#slot footer}<small>2026</small>{/slot}
</@Page></body>`,
	})
}

func TestComponentNested(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Inner.dreego": "Component Inner ()\n<body><span>inner</span></body>",
		"www/components/Outer.dreego": "Component Outer ()\n<body><article><@Inner/></article></body>",
		"www/routes/get.dreego":       `<body><@Outer/></body>`,
	})
}

func TestComponentNotFound(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body><@Missing/></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an unknown component:\n%s", out)
	}
	if !strings.Contains(out, "unknown component Missing") {
		t.Fatalf("unexpected diagnostic:\n%s", out)
	}
}

func TestComponentPropExpr(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Greet.dreego": "Component Greet (name string)\n<body><p>Hello {{ name }}</p></body>",
		"www/routes/get.dreego": `<server>n:="World"</server>
<body><@Greet name={n}/></body>`,
	})
}

func TestComponentPropExpression(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2></article></body>",
		"www/routes/get.dreego": `<server>type User struct { Name string }
user := User{Name: "Ada"}</server>
<body><@Card title={user.Name}/></body>`,
	})
}

func TestComponentScopedStyle(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Box.dreego": "Component Box ()\n<body><div class=\"box\"><p>scoped</p></div></body>\n<style>.box{border:1px solid red}</style>",
		"www/routes/get.dreego": `<head><title>T</title></head>
<body><@Box/><p class="box">unscoped</p></body>
<style>.box{color:blue}</style>`,
	})
}

func TestComponentSelfClosing(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Cmp.dreego": "Component Icon (name string)\n<body><i class=\"icon\">{{ name }}</i></body>",
		"www/routes/get.dreego":     `<body><@Icon name="star"/></body>`,
	})
}

func TestComponentWithGo(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Greeting.dreego": "Component Greeting (name string)\n<server>msg := \"Hi \" + name</server>\n<body><p>{{ msg }}</p></body>",
		"www/routes/get.dreego":          `<body><@Greeting name="World"/></body>`,
	})
}

func TestComponentWithSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2><div>{#slot}</div></article></body>",
		"www/routes/get.dreego":      `<body><@Card title="Welcome"><p>body text</p></@Card></body>`,
	})
}
