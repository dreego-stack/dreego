package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentForwardReferenceBindsPropsByName(t *testing.T) {
	t.Parallel()
	generated := dreegotest.Build(t, map[string]string{
		"www/components/Alpha.dreego": "Component Alpha ()\n<body><@Zulu second=\"B\" first=\"A\"/></body>",
		"www/components/Zulu.dreego":  "Component Zulu (first string, second string)\n<body><p>{{ first }}{{ second }}</p></body>",
		"www/routes/get.dreego":       `<body><@Alpha/></body>`,
	})
	components := generated["www/components/dree.go"]
	if !strings.Contains(components, `Zulu("A", "B")`) {
		t.Fatalf("forward component props were not bound by name, got:\n%s", components)
	}
}
