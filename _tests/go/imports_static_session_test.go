package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestImportsBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card (title string)\n<div><article><h2>{title}</h2></article></div>",
		"dreego/routes/get.dreego":      `import "dreego/components/Card"
<div><@Card title="Imported!"/></div>`,
	})
}

func TestImportsMissing(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `import "dreego/components/Nope"
<div><p>hi</p></div>`,
	})
}

func TestImportsMultiFile(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/button/Login.dreego": "Component Login ()\n<div><button>Login</button></div>",
		"dreego/routes/get.dreego":              `import "dreego/components/button"
<div><@Login/></div>`,
	})
}

func TestStaticBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hello</p></div>`,
	})
}

func TestStaticCollision(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/about/get.dreego": `<div><p>about</p></div>`,
	})
}

func TestStaticSubdir(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hello</p></div>`,
	})
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
c.SetSessionVal("key","val")
c.DelSessionVal("key")
v:=c.SessionVal("key")
</go>
<div><p>{v}</p></div>`,
	})
}

func TestSessionDestroy(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
c.SetSessionVal("a","1")
c.DestroySession()
v:=c.SessionVal("a")
</go>
<div><p>{v}</p></div>`,
	})
}

func TestSessionNoStore(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>v:=c.SessionVal("x")</go>
<div><p>{v}</p></div>`,
	})
}

func TestSessionSetGet(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
    c.SetSessionVal("key", "val")
    v := c.SessionVal("key")
    _ = v
</go>
<div><p>session set/get</p></div>`,
	})
}
