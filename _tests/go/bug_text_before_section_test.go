package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugTextBeforeSection(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/routes/get.dreego": `<!doctype html>
<html lang="en">
<server>msg := "hi"</server>
<body><p>{{ msg }}</p></body>`,
	})
}

func TestBugRootComponentCallRejected(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/components/Card.dreego": "Component Card ()\n<body>Card</body>",
		"www/routes/get.dreego":      `<@Card />`,
	})
}
