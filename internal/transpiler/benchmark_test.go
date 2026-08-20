package transpiler

import "testing"

const benchPage = `<head>
    <title>Benchmark</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>

<go>
    title := "Benchmark"
    show := true
    items := []string{"a", "b", "c"}
</go>

<div class="page">
    <h1>{{ title }}</h1>
    {#if show}
        <p>visible</p>
    {/if}
    {#each items as item}
        <p>{{ item }}</p>
    {/each}
    <@Card title="Hello" />
</div>`

const benchComponent = `Component Card (title string)

<div>
    <h2>{{ title }}</h2>
</div>`

func transpilePage(src string) (string, error) {
	_, imports, body := ParseHeader(src)
	tokens, err := Lex(body)
	if err != nil {
		return "", err
	}
	p := NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return "", err
	}
	file.Imports = imports
	file.SourceContent = src
	if len(file.Go) == 0 {
		file.Go = []GoSection{{Method: "GET"}}
	}
	for i := range file.Go {
		if !file.Go[i].MethodExplicit {
			file.Go[i].Method = "GET"
		}
	}
	gen := NewGenerator()
	out, _, err := GenerateMethodHandler(gen, file, nil, "routes", "index", "/{$}", "abc")
	return out, err
}

func BenchmarkGeneratePage(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := transpilePage(benchPage); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateComponent(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		comp, _, body := ParseHeader(benchComponent)
		tokens, err := Lex(body)
		if err != nil {
			b.Fatal(err)
		}
		p := NewParser(tokens)
		file, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
		file.Component = comp
		file.SourceContent = benchComponent
		if len(file.Go) == 0 {
			file.Go = []GoSection{{Method: ""}}
		}
		gen := NewGenerator()
		if _, err := GenerateComponent(gen, file, "abc"); err != nil {
			b.Fatal(err)
		}
	}
}
