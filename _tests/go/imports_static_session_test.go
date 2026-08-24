package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestImportsBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string)\n<body><article><h2>{{ title }}</h2></article></body>",
		"www/routes/get.dreego": `import "www/components/Card"
<body><@Card title="Imported!"/></body>`,
	})
}

func TestImportsMissing(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `import "www/components/Nope"
<body><p>hi</p></body>`,
	})
}

func TestImportsMultiFile(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/button/Login.dreego": "Component Login ()\n<body><button>Login</button></body>",
		"www/routes/get.dreego": `import "www/components/button"
<body><@Login/></body>`,
	})
}

func TestStaticBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<body><p>hello</p></body>`,
	})
}

func TestStaticCollision(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/about/get.dreego": `<body><p>about</p></body>`,
	})
}

func TestStaticSubdir(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<body><p>hello</p></body>`,
	})
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>
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
		"www/routes/get.dreego": `<server>
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
		"www/routes/get.dreego": `<server>v:=c.SessionVal("x")</server>
<body><p>{{ v }}</p></body>`,
	})
}

func TestSessionSetGet(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>
    c.SetSessionVal("key", "val")
    v := c.SessionVal("key")
    _ = v
</server>
<body><p>session set/get</p></body>`,
	})
}
