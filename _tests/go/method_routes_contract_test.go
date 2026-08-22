package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestMethodRoutesAllVerbsRemainIsolated(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/items.dreego": `<go>value := "get"</go><div><p>{{ value }}</p></div>
<go method="post">value := "post"</go><div method="post"><p>{{ value }}</p></div>
<go method="put">value := "put"</go><div method="put"><p>{{ value }}</p></div>
<go method="delete">value := "delete"</go><div method="delete"><p>{{ value }}</p></div>`,
	})
	for _, tc := range []struct{ method, want string }{{"GET", "get"}, {"POST", "post"}, {"PUT", "put"}, {"DELETE", "delete"}} {
		code, body, _ := c.Request(t, tc.method, "/items", "", nil)
		if code != 200 || !strings.Contains(body, tc.want) {
			t.Fatalf("%s /items = %d %q, want isolated %q", tc.method, code, body, tc.want)
		}
		for _, other := range []string{"get", "post", "put", "delete"} {
			if other != tc.want && strings.Contains(body, other) {
				t.Fatalf("%s rendered %q section: %q", tc.method, other, body)
			}
		}
	}
}

func TestMethodRouteOnlyPostDoesNotRegisterGet(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/submit.dreego": `<go method="post">message := "accepted"</go><div method="post"><p>{{ message }}</p></div>`,
	})
	code, body, _ := c.Request(t, "POST", "/submit", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "accepted")
	code, _, _ = c.Request(t, "GET", "/submit", "", nil)
	if code == 200 {
		t.Fatal("GET unexpectedly matched a POST-only route")
	}
}

func TestMethodRouteCanRenderComponentAndLayout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<div><html><body>{#slot}</body></html></div>`,
		"www/components/Badge.dreego": `Component Badge ()
<div class="badge">badge</div>`,
		"www/routes/profile.dreego": `<go method="post">name := "Ada"</go><div method="post"><@Badge/> <span>{{ name }}</span></div>`,
	})
	code, body, _ := c.Request(t, "POST", "/profile", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "badge")
	dreegotest.MustContainBody(t, body, "Ada")
}

func TestMethodRouteDynamicParameter(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/users/[id].dreego": `<go method="post">id := c.Param("id")</go><div method="post"><p>saved {{ id }}</p></div>`,
	})
	code, body, _ := c.Request(t, "POST", "/users/42", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "saved 42")
}

func TestMethodRouteDuplicateGoSectionsFail(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<go method="post">x := 1</go><div method="post">{{ x }}</div>
<go method="post">x := 2</go><div method="post">{{ x }}</div>`, "duplicate")
}

func TestMethodRouteDuplicateDivSectionsFail(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<go method="post">x := 1</go><div method="post">{{ x }}</div>
<div method="post">duplicate</div>`, "duplicate")
}

func TestMethodRouteStandardExtendedMethodGenerates(t *testing.T) {
	t.Parallel()
	generated := dreegotest.Generate(t, `<go method="patch">x := 1</go><div method="patch">{{ x }}</div>`)
	dreegotest.MustContain(t, generated, "PATCH")
}

func TestMethodRouteGenerationIsDeterministic(t *testing.T) {
	t.Parallel()
	src := `<go>message := "get"</go><div><p>{{ message }}</p></div>
<go method="post">message := "post"</go><div method="post"><p>{{ message }}</p></div>`
	one := dreegotest.Generate(t, src)
	two := dreegotest.Generate(t, src)
	if one != two {
		t.Fatal("repeated generation produced different output")
	}
}

func TestMethodRouteDoesNotLeakCodeBetweenMethods(t *testing.T) {
	t.Parallel()
	out := dreegotest.Generate(t, `<go>getOnly := "get"</go><div><p>{{ getOnly }}</p></div>
<go method="post">postOnly := "post"</go><div method="post"><p>{{ postOnly }}</p></div>`)
	if strings.Count(out, "getOnly") < 1 || strings.Count(out, "postOnly") < 1 {
		t.Fatal("generated output omitted one method section")
	}
}

func TestMethodRouteHeadRequestHasDefinedResult(t *testing.T) {
	t.Parallel()
	generated := dreegotest.Generate(t, `<div><p>ok</p></div>`)
	if !strings.Contains(generated, "HandleHealth") && !strings.Contains(generated, "HandleIndex") {
		t.Fatal("GET route did not generate a handler that can define HEAD behavior")
	}
}
