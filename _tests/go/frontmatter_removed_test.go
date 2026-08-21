package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestFrontmatterIsRejected(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/routes/get.dreego": `---
title: About
---
<div><h1>About</h1></div>`,
	})
}
