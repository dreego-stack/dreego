package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugTextBeforeSection(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/app/routes/+page.dreego": `<!doctype html>
<html lang="en">
<server>msg := "hi"</server>
<body><p>{{ msg }}</p></body>`,
	})
}

func TestBugRootComponentCallRejected(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/components/Card.dreego": "DREEFILE component ()\n<body>Card</body>",
		"www/app/routes/+page.dreego":    `<@Card />`,
	})
}
