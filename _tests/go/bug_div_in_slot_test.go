package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDivInSlot(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": `Component Card ()
<div><article>{#slot}</article></div>`,
		"www/routes/get.dreego": `<div><@Card><div class="inner">hi</div></@Card></div>`,
	})
}
