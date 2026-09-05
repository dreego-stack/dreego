package routes

import "encoding/json"

type NewsPost struct {
	Title   string
	Date    string
	Content string
}

var newsPosts = []NewsPost{}

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
