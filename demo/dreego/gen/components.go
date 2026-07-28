package gen

import (
	"fmt"
	"html"
	"strings"

	core "codeberg.org/dreego/dreego/dreego-core"
)

func FeatureCard(icon string, title string, text string) core.Component {
	return core.ComponentFunc(func(ctx *core.SSRContext) (string, error) {
		var b strings.Builder

		b.WriteString("<div data-scope=\"0f3bf4da8b93\">")
		b.WriteString(`
    `)
		b.WriteString(`<div class="bg-white p-8 group hover:bg-slate-50 transition-colors">`)
		b.WriteString(`
        `)
		b.WriteString(`<div class="text-3xl mb-4">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", icon)))
		b.WriteString(`</div>`)
		b.WriteString(`
        `)
		b.WriteString(`<h3 class="text-lg font-bold text-slate-900 mb-2">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", title)))
		b.WriteString(`</h3>`)
		b.WriteString(`
        `)
		b.WriteString(`<p class="text-slate-500 text-sm leading-relaxed">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", text)))
		b.WriteString(`</p>`)
		b.WriteString(`
    `)
		b.WriteString(`</div>`)
		b.WriteString(`
`)
		b.WriteString("</div>")

		return b.String(), nil
	})
}

func Hero(badge string, title string, subtitle string, cta1 string, cta2 string) core.Component {
	return core.ComponentFunc(func(ctx *core.SSRContext) (string, error) {
		var b strings.Builder

		b.WriteString("<div data-scope=\"5b854fdbcd3b\">")
		b.WriteString(`
    `)
		b.WriteString(`<main class="pt-28 pb-20 px-6 max-w-4xl mx-auto text-center relative">`)
		b.WriteString(`
        `)
		b.WriteString(`<div class="absolute inset-0 bg-grid opacity-40 -z-10">`)
		b.WriteString(`</div>`)
		b.WriteString(`
        `)
		b.WriteString(`<div class="inline-flex items-center gap-2 px-4 py-1.5 bg-blue-50 text-blue-700 rounded-full text-sm font-semibold mb-8 ring-1 ring-blue-200/60">`)
		b.WriteString(`✨ `)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", badge)))
		b.WriteString(`</div>`)
		b.WriteString(`
        `)
		b.WriteString(`<h1 class="text-5xl md:text-7xl font-extrabold text-slate-900 tracking-tight leading-[1.05] mb-6 max-w-3xl mx-auto">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", title)))
		b.WriteString(`</h1>`)
		b.WriteString(`
        `)
		b.WriteString(`<p class="text-lg md:text-xl text-slate-500 max-w-2xl mx-auto mb-10 leading-relaxed">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", subtitle)))
		b.WriteString(`</p>`)
		b.WriteString(`
        `)
		b.WriteString(`<div class="flex flex-wrap justify-center gap-4 mb-16">`)
		b.WriteString(`
            `)
		b.WriteString(`<a href="/signup" class="px-8 py-3.5 bg-slate-900 text-white font-semibold rounded-xl hover:bg-slate-800 transition-colors text-lg shadow-xl shadow-slate-900/20">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", cta1)))
		b.WriteString(`</a>`)
		b.WriteString(`
            `)
		b.WriteString(`<a href="/docs" class="px-8 py-3.5 bg-white text-slate-700 font-semibold rounded-xl hover:bg-slate-50 border border-slate-200 transition-colors text-lg shadow-sm">`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", cta2)))
		b.WriteString(`</a>`)
		b.WriteString(`
        `)
		b.WriteString(`</div>`)
		b.WriteString(`

        `)
		b.WriteString(`<div class="inline-block bg-slate-900 rounded-2xl p-0.5 shadow-2xl shadow-slate-900/20">`)
		b.WriteString(`
            `)
		b.WriteString(`<div class="bg-slate-950 rounded-2xl px-6 py-5 text-left">`)
		b.WriteString(`
                `)
		b.WriteString(`<div class="flex items-center gap-2 mb-4">`)
		b.WriteString(`
                    `)
		b.WriteString(`<div class="w-3 h-3 rounded-full bg-red-400">`)
		b.WriteString(`</div>`)
		b.WriteString(`
                    `)
		b.WriteString(`<div class="w-3 h-3 rounded-full bg-amber-400">`)
		b.WriteString(`</div>`)
		b.WriteString(`
                    `)
		b.WriteString(`<div class="w-3 h-3 rounded-full bg-emerald-400">`)
		b.WriteString(`</div>`)
		b.WriteString(`
                `)
		b.WriteString(`</div>`)
		b.WriteString(`
                `)
		b.WriteString(`<code class="text-sm text-emerald-400 font-mono leading-relaxed">`)
		b.WriteString(`$ dreego new myapp`)
		b.WriteString(`<br>`)
		b.WriteString(`$ dreego generate &amp;&amp; go run .`)
		b.WriteString(`<br>`)
		b.WriteString(`<span class="text-slate-500">`)
		b.WriteString(`→ Server listening on :8080`)
		b.WriteString(`</span>`)
		b.WriteString(`</code>`)
		b.WriteString(`
            `)
		b.WriteString(`</div>`)
		b.WriteString(`
        `)
		b.WriteString(`</div>`)
		b.WriteString(`
    `)
		b.WriteString(`</main>`)
		b.WriteString(`
`)
		b.WriteString("</div>")
		b.WriteString("<style>")
		b.WriteString(`[data-scope=5b854fdbcd3b] .bg-grid { background-image: radial-gradient(circle, #e2e8f0 1px, transparent 1px); background-size: 24px 24px;}`)
		b.WriteString("</style>")

		return b.String(), nil
	})
}

