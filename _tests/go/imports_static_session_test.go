package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestImportsBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "DREEFILE component (title string)\n<body><article><h2>{{ title }}</h2></article></body>",
		"www/routes/+page.dreego": `COMPONENT "www/components" IMPORT { Card }
<body><@Card title="Imported!"/></body>`,
	})
}

func TestImportsMissing(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `COMPONENT "www/components" IMPORT { Nope }
<body><p>hi</p></body>`,
	})
}

func TestImportsMultiFile(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/button/Login.dreego": "DREEFILE component ()\n<body><button>Login</button></body>",
		"www/routes/+page.dreego": `COMPONENT "www/components/button" IMPORT { Login }
<body><@Login/></body>`,
	})
}

func TestStaticBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<body><p>hello</p></body>`,
	})
}

func TestStaticCollision(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/about/+page.dreego": `<body><p>about</p></body>`,
	})
}

func TestStaticSubdir(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<body><p>hello</p></body>`,
	})
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
c.SetSessionVal("key","val")
c.DelSessionVal("key")
v:=c.SessionVal("key")
</server>
<body><p>{{ v }}</p></body>`,
	})
}

func TestSessionDestroy(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
c.SetSessionVal("a","1")
c.DestroySession()
v:=c.SessionVal("a")
</server>
<body><p>{{ v }}</p></body>`,
	})
}

func TestSessionNoStore(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>v:=c.SessionVal("x")</server>
<body><p>{{ v }}</p></body>`,
	})
}

func TestSessionSetGet(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
    c.SetSessionVal("key", "val")
    v := c.SessionVal("key")
    _ = v
</server>
<body><p>session set/get</p></body>`,
	})
}
