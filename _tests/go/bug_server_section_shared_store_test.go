package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// The documented shared-store pattern from _docs/server-section.md must
// compile: a leading declaration block with a mutex-guarded store, plus a
// method and a handler func, shared across route files in the same package.
func TestBugServerSectionSharedStoreCompiles(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `GOIMPORT { sync }

<server>
type Store struct {
    mu    sync.Mutex
    items []string
}

var store = &Store{}

func (s *Store) Add(item string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items = append(s.items, item)
}
</server>
<body><p>{{ len(store.items) }}</p></body>`,
		"www/routes/add.dreego": `<server>
store.Add("x")
</server>
<body><p>{{ len(store.items) }}</p></body>`,
	})
}

// A request-local var after a statement must not become shared package state.
func TestBugServerSectionRequestLocalVarStaysLocal(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
count := 3
var product string
product = "mug"
</server>
<body><p>{{ product }}{{ count }}</p></body>`,
	})
}
