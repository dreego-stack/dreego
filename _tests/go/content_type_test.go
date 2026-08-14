package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestContentTypeAcceptFallback(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
    msg := "hello"
</go>
<go type="json">
    c.JSON(200, map[string]string{"msg": msg})
</go>
<div><h1>{{ msg }}</h1></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<h1>hello</h1>") {
		t.Fatalf("no HTML fallback, got: %s", body)
	}
}

func TestContentTypeAcceptJSON(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go type="json">
    c.JSON(200, map[string]string{"ok": "true"})
</go>`,
	})
	code, body, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/json"})
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, `"ok"`) {
		t.Fatalf("no JSON response, got: %s", body)
	}
}

func TestContentTypeAcceptXML(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go type="xml">
    user := struct{XMLName struct{} ` + "`xml:\"user\"`" + `; Name string ` + "`xml:\"name\"`" + `}{Name: "Lukas"}
    c.XML(200, user)
</go>`,
	})
	code, body, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/xml"})
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<user>") {
		t.Fatalf("no XML response, got: %s", body)
	}
}

func TestContentTypeBindError(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/post.dreego": `<go type="json">
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
	if code != 400 {
		t.Fatalf("status = %d, want 400", code)
	}
	if !strings.Contains(body, `"error"`) {
		t.Fatalf("no error response, got: %s", body)
	}
}

func TestContentTypeBindPost(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/post.dreego": `<go type="json">
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
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, `"name":"Lukas"`) {
		t.Fatalf("Bind not working, got: %s", body)
	}
	if !strings.Contains(body, `"echo":true`) {
		t.Fatalf("no echo field, got: %s", body)
	}
}

func TestContentTypeCustomBasic(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go type="custom">
    msg := []byte("hello world")
    c.Write(200, "text/plain", msg)
</go>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "text/plain")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "c.Write")
}

func TestContentTypeHTMLDefault(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
    msg := "hello"
</go>
<div><h1>{{ msg }}</h1></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "text/html")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "b.WriteString")
}

func TestContentTypeJSONAutoImports(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/post.dreego": `<go type="json">
    var input map[string]any
    c.Bind(&input)
    input["echo"] = true
    c.JSON(200, input)
</go>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "c.JSON")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "c.Bind")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "application/json")
}

func TestContentTypeJSONBasic(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go type="json">
    user := map[string]string{"name": "Lukas"}
    c.JSON(200, user)
</go>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "application/json")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "c.JSON")
}

func TestContentTypeJSONShared(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
    msg := "Lukas"
</go>

<go type="json">
    c.JSON(200, map[string]string{"name": msg})
</go>

<div>
    <h1>{{ msg }}</h1>
</div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "application/json")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "c.JSON")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Lukas")
}

func TestContentTypeMultiTyped(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
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
	if !strings.Contains(jbody, `"format":"json"`) {
		t.Fatalf("JSON broken, got: %s", jbody)
	}
	_, xbody, _ := c.Request(t, "GET", "/", "", map[string]string{"Accept": "application/xml"})
	if !strings.Contains(xbody, "<user>") {
		t.Fatalf("XML broken, got: %s", xbody)
	}
	_, hbody := c.Get(t, "/")
	if !strings.Contains(hbody, "<h1>Lukas</h1>") {
		t.Fatalf("HTML broken, got: %s", hbody)
	}
}

func TestContentTypeXMLBasic(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go type="xml">
    user := struct{XMLName struct{} ` + "`xml:\"user\"`" + `; Name string ` + "`xml:\"name\"`" + `}{Name: "Lukas"}
    c.XML(200, user)
</go>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "application/xml")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "c.XML")
}
