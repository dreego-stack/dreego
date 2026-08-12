package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentQuotedAttrs(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": `Component Card (title string)
<div><h1>{title}</h1></div>`,
		"dreego/routes/get.dreego": `<div><@Card title="Hello World"/></div>`,
	})
}
