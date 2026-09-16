package main

import "testing"

func TestParseScaffoldFlagsEmptyTemplateValue(t *testing.T) {
	cases := [][]string{
		{"-t", ""},
		{"--template", ""},
		{"-t="},
		{"--template="},
	}
	for _, args := range cases {
		if _, err := parseScaffoldFlags(args); err == nil {
			t.Errorf("parseScaffoldFlags(%q) must fail for an empty template value", args)
		}
	}
}

func TestParseScaffoldFlagsMissingTemplateValue(t *testing.T) {
	for _, args := range [][]string{{"-t"}, {"--template"}} {
		if _, err := parseScaffoldFlags(args); err == nil {
			t.Errorf("parseScaffoldFlags(%q) must fail for a missing template value", args)
		}
	}
}

func TestParseScaffoldFlagsValidTemplate(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"-t", "web-minimal"}, "web-minimal"},
		{[]string{"--template", "web-minimal"}, "web-minimal"},
		{[]string{"-t=web-minimal"}, "web-minimal"},
		{[]string{"--template=web-minimal"}, "web-minimal"},
	}
	for _, tc := range cases {
		flags, err := parseScaffoldFlags(tc.args)
		if err != nil {
			t.Errorf("parseScaffoldFlags(%q) unexpected error: %v", tc.args, err)
			continue
		}
		if flags.template != tc.want {
			t.Errorf("parseScaffoldFlags(%q) template = %q, want %q", tc.args, flags.template, tc.want)
		}
	}
}
