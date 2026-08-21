package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugMultilineComponentPropsCompile(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string, count int)\n<div><p>{{ title }} {{ count }}</p></div>",
		"www/routes/get.dreego": `<div><@Card
			title="Hello"
			count={2}
		/></div>`,
	})
}
