package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugGenerateWithoutRoutesBuilds(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/dreego.config.json": `{}`,
		"www/app/routes/+page.dreego": `<body><p>hello</p></body>`,
	})
}
