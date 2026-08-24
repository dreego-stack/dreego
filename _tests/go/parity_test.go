package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestParityCLIAndDreegotestGenerate(t *testing.T) {
	t.Parallel()
	fixtures := map[string]string{
		"basic": `<head><title>Parity</title></head>
<server>msg := "hello"</server>
<body><h1>{{ msg }}</h1></body>`,
		"each": `<server>items := []string{"a", "b"}</server>
<body><ul>{#each items as item}<li>{{ $loop.Index }}: {{ item }}</li>{/each}</ul></body>`,
		"if-else": `<server>show := true</server>
<body>{#if show}<p>yes</p>{#else}<p>no</p>{/if}</body>`,
		"style-script": `<head><title>S</title></head>
<style>.x { color: red; }</style>
<client>const a = 1;</client>
<body><p>styled</p></body>`,
	}
	for name, src := range fixtures {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			gen := dreegotest.Build(t, map[string]string{
				"www/routes/get.dreego": src,
			})
			cliOut := gen["www/routes/dree.go"]
			dgtOut := dreegotest.Generate(t, src)
			if dgtOut == "" {
				t.Fatal("dreegotest.Generate returned empty output")
			}
			if !strings.Contains(cliOut, dgtOut) {
				t.Fatalf("CLI generate diverges from dreegotest.Generate\n--- dreegotest output ---\n%s\n--- CLI routes.go ---\n%s", dgtOut, cliOut)
			}
		})
	}
}

func TestParityCLIAndDreegotestGenerateComponent(t *testing.T) {
	t.Parallel()
	fixtures := map[string]string{
		"Badge": `Component Badge (label string)
<body class="badge">{{ label }}</body>`,
		"Status": `Component Status (ok bool)
<body>{#if ok}<span>on</span>{#else}<span>off</span>{/if}</body>`,
		"Card": `Component Card (title string)
<style>.card { padding: 1rem; }</style>
<body class="card"><h2>{{ title }}</h2></body>`,
	}
	for name, src := range fixtures {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			gen := dreegotest.Build(t, map[string]string{
				"www/components/" + name + ".dreego": src,
			})
			cliOut := gen["www/components/dree.go"]
			dgtOut := dreegotest.GenerateComponent(t, src)
			if dgtOut == "" {
				t.Fatal("dreegotest.GenerateComponent returned empty output")
			}
			if !strings.Contains(cliOut, dgtOut) {
				t.Fatalf("CLI generate diverges from dreegotest.GenerateComponent\n--- dreegotest output ---\n%s\n--- CLI components.go ---\n%s", dgtOut, cliOut)
			}
		})
	}
}
