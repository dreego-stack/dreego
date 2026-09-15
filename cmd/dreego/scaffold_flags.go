package main

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/cmd/dreego/internal/templates"
)

type scaffoldFlags struct {
	template string
	list     bool
	args     []string
}

func parseScaffoldFlags(args []string) (scaffoldFlags, error) {
	var flags scaffoldFlags
	for i := 0; i < len(args); i++ {
		switch arg := args[i]; {
		case arg == "-l" || arg == "--list":
			flags.list = true
		case arg == "-t" || arg == "--template":
			if i+1 >= len(args) || args[i+1] == "" {
				return flags, fmt.Errorf("flag %s requires a template name", arg)
			}
			flags.template = args[i+1]
			i++
		case strings.HasPrefix(arg, "-t=") || strings.HasPrefix(arg, "--template="):
			flag, value, _ := strings.Cut(arg, "=")
			if value == "" {
				return flags, fmt.Errorf("flag %s requires a template name", flag)
			}
			flags.template = value
		default:
			flags.args = append(flags.args, arg)
		}
	}
	return flags, nil
}

func resolveTemplate(name string) (templates.Meta, error) {
	if name == "" {
		name = templates.DefaultName
	}
	meta, err := templates.MetaOf(name)
	if err != nil {
		return templates.Meta{}, fmt.Errorf("unknown template %q\n  available templates: %s", name, templateNames())
	}
	return meta, nil
}

func templateNames() string {
	metas := templates.List()
	names := make([]string, 0, len(metas))
	for _, meta := range metas {
		names = append(names, meta.Name)
	}
	return strings.Join(names, ", ")
}

func printTemplates() {
	for _, meta := range templates.List() {
		fmt.Printf("%s — %s\n  %s\n", meta.Name, meta.Title, meta.Description)
	}
}
