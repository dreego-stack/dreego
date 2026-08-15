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
<go>msg := "hello"</go>
<div><h1>{{ msg }}</h1></div>`,
		"each": `<go>items := []string{"a", "b"}</go>
<div><ul>{#each items as item}<li>{{ $loop.Index }}: {{ item }}</li>{/each}</ul></div>`,
		"if-else": `<go>show := true</go>
<div>{#if show}<p>yes</p>{#else}<p>no</p>{/if}</div>`,
		"style-script": `<head><title>S</title></head>
<style>.x { color: red; }</style>
<script>const a = 1;</script>
<div><p>styled</p></div>`,
	}
	for name, src := range fixtures {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			gen := dreegotest.Build(t, map[string]string{
				"dreego/routes/get.dreego": src,
			})
			cliOut := gen["dreego/gen/routes.go"]
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
<div class="badge">{{ label }}</div>`,
		"Status": `Component Status (ok bool)
<div>{#if ok}<span>on</span>{#else}<span>off</span>{/if}</div>`,
		"Card": `Component Card (title string)
<style>.card { padding: 1rem; }</style>
<div class="card"><h2>{{ title }}</h2></div>`,
	}
	for name, src := range fixtures {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			gen := dreegotest.Build(t, map[string]string{
				"dreego/components/" + name + ".dreego": src,
			})
			cliOut := gen["dreego/gen/components.go"]
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
