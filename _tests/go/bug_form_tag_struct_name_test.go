package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugFormTagStructName(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-search.dreego": `<go>
type LoginForm struct {
}
type SearchQuery struct {
	Query string ` + "`form:\"q\"`" + `
}
func Search(c *dreego.SSRContext, form SearchQuery) error {
	return nil
}
</go>
<div>
  <form g-action="Search" method="post">
    <input name="q">
    <button>Search</button>
  </form>
</div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "SearchQuery")
}
