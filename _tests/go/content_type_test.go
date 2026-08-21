package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestContentTypeAcceptFallback(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>
    msg := "hello"
</go>
<go type="json">
    c.JSON(200, map[string]string{"msg": msg})
</go>
<div><h1>{{ msg }}</h1></div>`,
	})
	code, body := c.Get(t, "/")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "<h1>hello</h1>")
}

func TestContentTypeAcceptJSON(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go type="json">
    c.JSON(200, map[string]string{"ok": "true"})
</go>`,
	})
	code, body, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/json"})
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, `"ok"`)
}

func TestContentTypeAcceptXML(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go type="xml">
    user := struct{XMLName struct{} ` + "`xml:\"user\"`" + `; Name string ` + "`xml:\"name\"`" + `}{Name: "Lukas"}
    c.XML(200, user)
</go>`,
	})
	code, body, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/xml"})
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "<user>")
}

func TestContentTypeBindError(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/post.dreego": `<go type="json">
    var input map[string]any
    err := c.Bind(&input)
    if err != nil {
        c.JSON(400, map[string]string{"error": err.Error()})
    } else {
        c.JSON(200, input)
    }
</go>`,
	})
	code, body, _ := c.Request(t, "POST", "/", "not json", map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	dreegotest.MustStatus(t, code, 400)
	dreegotest.MustContainBody(t, body, `"error"`)
}

func TestContentTypeBindPost(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/post.dreego": `<go type="json">
    var input map[string]any
    c.Bind(&input)
    input["echo"] = true
    c.JSON(200, input)
</go>`,
	})
	code, body, _ := c.Request(t, "POST", "/", `{"name":"Lukas"}`, map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, `"name":"Lukas"`)
	dreegotest.MustContainBody(t, body, `"echo":true`)
}

func TestContentTypeCustomBasic(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go type="custom">
    msg := []byte("hello world")
    c.Write(200, "text/plain", msg)
</go>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "text/plain")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "c.Write")
}

func TestContentTypeHTMLDefault(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go>
    msg := "hello"
</go>
<div><h1>{{ msg }}</h1></div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "text/html")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "b.WriteString")
}

func TestContentTypeJSONAutoImports(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/post.dreego": `<go type="json">
    var input map[string]any
    c.Bind(&input)
    input["echo"] = true
    c.JSON(200, input)
</go>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "c.JSON")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "c.Bind")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "application/json")
}

func TestContentTypeJSONBasic(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go type="json">
    user := map[string]string{"name": "Lukas"}
    c.JSON(200, user)
</go>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "application/json")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "c.JSON")
}

func TestContentTypeJSONShared(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go>
    msg := "Lukas"
</go>

<go type="json">
    c.JSON(200, map[string]string{"name": msg})
</go>

<div>
    <h1>{{ msg }}</h1>
</div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "application/json")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "c.JSON")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "Lukas")
}

func TestContentTypeMultiTyped(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>
    name := "Lukas"
</go>
<go type="json">
    c.JSON(200, map[string]string{"format": "json", "name": name})
</go>
<go type="xml">
    user := struct{XMLName struct{} ` + "`xml:\"user\"`" + `; Name string ` + "`xml:\"name\"`" + `}{Name: name}
    c.XML(200, user)
</go>
<div><h1>{{ name }}</h1></div>`,
	})
	_, jbody, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/json"})
	dreegotest.MustContainBody(t, jbody, `"format":"json"`)
	_, xbody, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/xml"})
	dreegotest.MustContainBody(t, xbody, "<user>")
	_, hbody := c.Get(t, "/")
	dreegotest.MustContainBody(t, hbody, "<h1>Lukas</h1>")
}

func TestContentTypeXMLBasic(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go type="xml">
    user := struct{XMLName struct{} ` + "`xml:\"user\"`" + `; Name string ` + "`xml:\"name\"`" + `}{Name: "Lukas"}
    c.XML(200, user)
</go>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "application/xml")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "c.XML")
}
