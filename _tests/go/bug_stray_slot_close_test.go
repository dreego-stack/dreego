package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugStraySlotClose(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/components/Card.dreego": `Component Card ()
<div><article>{#slot header}</article></div>`,
		"www/routes/get.dreego": `<div><@Card>{/slot}</@Card></div>`,
	})
}
