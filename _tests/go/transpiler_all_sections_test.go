package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerAllSections(t *testing.T) {
	dreegotest.MustCompile(t, `<head><title>All</title></head>
<go>x := "ok"</go>
<div><p>{x}</p></div>
<script>console.log("js")</script>
<style>p { color: red; }</style>`)
}
