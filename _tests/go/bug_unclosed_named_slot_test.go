package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugUnclosedNamedSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"dreego/components/Card.dreego": `Component Card ()
<div><article>{#slot header}</article></div>`,
		"dreego/routes/get.dreego": `<div><@Card>{#slot header}no close</@Card></div>`,
	})
}
