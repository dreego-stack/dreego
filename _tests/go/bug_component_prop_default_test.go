package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentPropDefault(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Card.dreego": `Component Card (title string = "Default Title")
<body><h1>{{ title }}</h1></body>`,
		"www/routes/get.dreego": `<body><@Card/></body>`,
	})
}
