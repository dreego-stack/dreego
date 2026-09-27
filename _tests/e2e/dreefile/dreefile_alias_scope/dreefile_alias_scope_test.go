package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDreefileAliasResolvesAcrossFiles(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Card.dreego": `DREEFILE component (title string)
<body><article><h2>{{ title }}</h2></article></body>`,
		"www/routes/a.dreego": `COMPONENT "www/components" IMPORT { Card as ProductCard }
<body><@ProductCard title="from a"/></body>`,
		"www/routes/b.dreego": `<body><@ProductCard title="from b"/></body>`,
	})
	for _, tc := range []struct {
		path string
		want string
	}{
		{"/a", "<h2>from a</h2>"},
		{"/b", "<h2>from b</h2>"},
	} {
		code, body := c.Get(t, tc.path)
		if code != 200 {
			t.Fatalf("GET %s status = %d, want 200", tc.path, code)
		}
		if !strings.Contains(body, tc.want) {
			t.Fatalf("GET %s body missing %q, got: %s", tc.path, tc.want, body)
		}
	}
}
