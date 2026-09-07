package ts

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestDeclarationsMapGoComponentProps(t *testing.T) {
	component := &ir.ComponentDef{Name: "Profile", Props: []ir.Prop{
		{Name: "name", Type: "string"},
		{Name: "count", Type: "int"},
		{Name: "active", Type: "bool"},
		{Name: "scores", Type: "[]float64"},
		{Name: "labels", Type: "map[string]string"},
		{Name: "nickname", Type: "*string"},
		{Name: "identifier", Type: "uint64"},
	}}
	got, err := Declarations(component, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"interface ProfileProps {",
		"name: string;",
		"count: number;",
		"active: boolean;",
		"scores: Array<number>;",
		"labels: Record<string, string>;",
		"nickname: string | null;",
		"identifier: number;",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("declarations missing %q:\n%s", want, got)
		}
	}
}

func TestDeclarationsMapGoModels(t *testing.T) {
	got, err := Declarations(nil, []ir.ServerSection{{Code: "type User struct { Name string `json:\"name\"`; Secret string `json:\"-\"`; Age int; Tags []string; Parent *User }"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"interface User {", "name: string;", "Age: number;", "Tags: Array<string>;", "Parent: User | null;"} {
		if !strings.Contains(got, want) {
			t.Fatalf("declarations missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Secret") {
		t.Fatalf("ignored field was emitted:\n%s", got)
	}
}
