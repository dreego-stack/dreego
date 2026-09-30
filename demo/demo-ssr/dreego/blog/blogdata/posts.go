package blogdata

type Post struct {
	Title   string
	Slug    string
	Date    string
	Excerpt string
}

func Posts() []Post {
	return []Post{
		{Title: "Hello, Dreego", Slug: "hello-dreego", Date: "2026-08-24", Excerpt: "A first look at the v0.2 and v0.3 releases: typed SSR, the render foundation, and Markdown bodies."},
		{Title: "The Render Foundation", Slug: "rendering-foundation", Date: "2026-08-18", Excerpt: "How the v0.2 render foundation keeps generated inputs typed and opens the door to SSG and Wails."},
		{Title: "Markdown in the Body", Slug: "markdown-in-the-body", Date: "2026-08-10", Excerpt: "Write real Markdown in a <body lang=\"md\"> section and mix in Dreego constructs."},
		{Title: "Interactive Markdown", Slug: "interactive-markdown", Date: "2026-08-08", Excerpt: "HTMX-powered interactivity inside a Markdown body — components, styles, and scripts coexist."},
	}
}
