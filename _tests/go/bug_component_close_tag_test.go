package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentCloseTag(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": `Component Card (title string)
<div><article><h2>{{ title }}</h2><div>{#slot}</div></article></div>`,
		"www/routes/get.dreego": `<div><@Card title="Hi">text</@Card></div>`,
	})
}
