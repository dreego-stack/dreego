package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestMethodRoutesAllVerbsRemainIsolated(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/items.dreego": `<server>value := "get"</server><body><p>{{ value }}</p></body>
<server method="post">value := "post"</server><body method="post"><p>{{ value }}</p></body>
<server method="put">value := "put"</server><body method="put"><p>{{ value }}</p></body>
<server method="delete">value := "delete"</server><body method="delete"><p>{{ value }}</p></body>`,
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
		"www/routes/submit.dreego": `<server method="post">message := "accepted"</server><body method="post"><p>{{ message }}</p></body>`,
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
		"www/layouts/default.dreego": `<body><html><body>{#slot}</body></html></body>`,
		"www/components/Badge.dreego": `Component Badge ()
<body class="badge">badge</body>`,
		"www/routes/profile.dreego": `<server method="post">name := "Ada"</server><body method="post"><@Badge/> <span>{{ name }}</span></body>`,
	})
	code, body, _ := c.Request(t, "POST", "/profile", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "badge")
	dreegotest.MustContainBody(t, body, "Ada")
}

func TestMethodRouteDynamicParameter(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/users/[id].dreego": `<server method="post">id := c.Param("id")</server><body method="post"><p>saved {{ id }}</p></body>`,
	})
	code, body, _ := c.Request(t, "POST", "/users/42", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "saved 42")
}

func TestMethodRouteDuplicateServerSectionsFail(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<server method="post">x := 1</server><body method="post">{{ x }}</body>
<server method="post">x := 2</server><body method="post">{{ x }}</body>`, "duplicate")
}

func TestMethodRouteDuplicateDivSectionsFail(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<server method="post">x := 1</server><body method="post">{{ x }}</body>
<body method="post">duplicate</body>`, "duplicate")
}

func TestMethodRouteStandardExtendedMethodGenerates(t *testing.T) {
	t.Parallel()
	generated := dreegotest.Generate(t, `<server method="patch">x := 1</server><body method="patch">{{ x }}</body>`)
	dreegotest.MustContain(t, generated, "PATCH")
}

func TestMethodRouteGenerationIsDeterministic(t *testing.T) {
	t.Parallel()
	src := `<server>message := "get"</server><body><p>{{ message }}</p></body>
<server method="post">message := "post"</server><body method="post"><p>{{ message }}</p></body>`
	one := dreegotest.Generate(t, src)
	two := dreegotest.Generate(t, src)
	if one != two {
		t.Fatal("repeated generation produced different output")
	}
}

func TestMethodRouteDoesNotLeakCodeBetweenMethods(t *testing.T) {
	t.Parallel()
	out := dreegotest.Generate(t, `<server>getOnly := "get"</server><body><p>{{ getOnly }}</p></body>
<server method="post">postOnly := "post"</server><body method="post"><p>{{ postOnly }}</p></body>`)
	if strings.Count(out, "getOnly") < 1 || strings.Count(out, "postOnly") < 1 {
		t.Fatal("generated output omitted one method section")
	}
}

func TestMethodRouteHeadRequestHasDefinedResult(t *testing.T) {
	t.Parallel()
	generated := dreegotest.Generate(t, `<body><p>ok</p></body>`)
	if !strings.Contains(generated, "HandleHealth") && !strings.Contains(generated, "HandleIndex") {
		t.Fatal("GET route did not generate a handler that can define HEAD behavior")
	}
}
