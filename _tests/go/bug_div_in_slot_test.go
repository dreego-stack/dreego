package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDivInSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": `Component Card ()
<body><article>{#slot}</article></body>`,
		"www/routes/get.dreego": `<body><@Card><div class="inner">hi</div></@Card></body>`,
	})
}
