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
<go>msg := "hi"</go>
<div><p>{{ msg }}</p></div>`,
	})
}

func TestBugRootComponentCallRejected(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/components/Card.dreego": "Component Card ()\n<div>Card</div>",
		"www/routes/get.dreego":      `<@Card />`,
	})
}
