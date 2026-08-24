package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugUnclosedNamedSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/components/Card.dreego": `Component Card ()
<body><article>{#slot header}</article></body>`,
		"www/routes/get.dreego": `<body><@Card>{#slot header}no close</@Card></body>`,
	})
}
