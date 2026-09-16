package routes

import (
	"encoding/json"
	"html"
	"log"
	"strings"
	"sync"

	dreego "github.com/dreego-stack/dreego/core"
)

type NewsPost struct {
	Title   string
	Date    string
	Content string
}

var (
	newsPosts = []NewsPost{}
	newsMu    sync.Mutex
)

const newsSeedJSON = `[
	{
		"title": "Dreego v0.3: The Render Foundation",
		"date": "2026-09-01",
		"content": "## Announcing v0.3\n\nWe are shipping the **render foundation** that powers the next phase of Dreego.\n\n- Typed SSR inputs\n- A target-neutral App\n- Runtime markdown rendering\n\nRead the [roadmap](https://github.com/dreego-stack/dreego) for what comes next."
	},
	{
		"title": "A New Logo",
		"date": "2026-08-28",
		"content": "## Fresh look\n\nDreego has a new **logo** to go with the v0.3 release.\n\n- Simpler mark\n- Slate palette\n- Built for dark and light themes"
	},
	{
		"title": "Welcome to the Dreego Blog",
		"date": "2026-08-20",
		"content": "## Welcome\n\nThis is the first post on the Dreego blog. We will share **field notes** about building a compile-time Go web framework.\n\nStay tuned for more."
	}
]`

func init() {
	_ = json.Unmarshal([]byte(newsSeedJSON), &newsPosts)
}

func RenderNewsPosts() (string, error) {
	var b strings.Builder
	newsMu.Lock()
	defer newsMu.Unlock()
	for _, p := range newsPosts {
		body, err := dreego.MarkdownToHTML(p.Content)
		if err != nil {
			log.Printf("news: skipping post %q: %v", p.Title, err)
			continue
		}
		b.WriteString(`<article class="mb-10 rounded-3xl border border-slate-200 bg-white p-7">`)
		b.WriteString(`<div class="mb-3 flex items-center justify-between">`)
		b.WriteString(`<span class="font-mono text-xs text-emerald-700">` + p.Date + `</span>`)
		b.WriteString(`<span class="font-mono text-xs uppercase tracking-[0.2em] text-slate-400">announcement</span>`)
		b.WriteString(`</div>`)
		b.WriteString(`<h2 class="mb-4 text-2xl font-semibold tracking-tight text-slate-900">` + html.EscapeString(p.Title) + `</h2>`)
		b.WriteString(body)
		b.WriteString(`</article>`)
	}
	return b.String(), nil
}
