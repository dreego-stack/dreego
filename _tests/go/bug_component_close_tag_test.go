package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentCloseTag(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": `Component Card (title string)
<body><article><h2>{{ title }}</h2><div>{#slot}</div></article></body>`,
		"www/routes/get.dreego": `<body><@Card title="Hi">text</@Card></body>`,
	})
}
