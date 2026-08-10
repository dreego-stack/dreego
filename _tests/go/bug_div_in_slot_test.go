package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDivInSlot(t *testing.T) {
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Card.dreego": `Component Card ()
<div><article>{#slot}</article></div>`,
		"dreego/routes/get.dreego": `<div><@Card><div class="inner">hi</div></@Card></div>`,
	})
}
