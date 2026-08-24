package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugMultilineComponentPropsCompile(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": "Component Card (title string, count int)\n<body><p>{{ title }} {{ count }}</p></body>",
		"www/routes/get.dreego": `<body><@Card
			title="Hello"
			count={2}
		/></body>`,
	})
}
