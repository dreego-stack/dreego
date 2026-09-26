package lexer

import (
	"errors"
	"strings"
	"testing"
)

func TestParseFileHeaderProfile(t *testing.T) {
	header, body, err := ParseFileHeaderStrict("PROFILE \"hooks\"\n\n<body>x</body>")
	if err != nil {
		t.Fatalf("PROFILE must be accepted: %v", err)
	}
	if header.Profile != "hooks" {
		t.Fatalf("expected profile hooks, got %q", header.Profile)
	}
	if body != "<body>x</body>" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestParseFileHeaderProfileAlongsideOtherDirectives(t *testing.T) {
	header, _, err := ParseFileHeaderStrict("DREEFILE page\nPROFILE \"app\"\nLAYOUT \"www/layouts/base.dreego\"\n\n<body>x</body>")
	if err != nil {
		t.Fatalf("PROFILE with other directives must be accepted: %v", err)
	}
	if header.Profile != "app" {
		t.Fatalf("expected profile app, got %q", header.Profile)
	}
	if header.Layout != "www/layouts/base.dreego" {
		t.Fatalf("expected layout preserved, got %q", header.Layout)
	}
}

func TestParseFileHeaderProfileAbsent(t *testing.T) {
	header, _, err := ParseFileHeaderStrict("<body>x</body>")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if header.Profile != "" {
		t.Fatalf("expected no profile, got %q", header.Profile)
	}
}

func TestParseFileHeaderRejectsInvalidProfile(t *testing.T) {
	for _, src := range []string{
		"PROFILE hooks\n\n<body>x</body>",
		"PROFILE \"\"\n\n<body>x</body>",
		"PROFILE\n\n<body>x</body>",
	} {
		_, _, err := ParseFileHeaderStrict(src)
		if err == nil {
			t.Fatalf("expected error for %q", src)
		}
		if !strings.Contains(err.Error(), "PROFILE") {
			t.Fatalf("error must name PROFILE, got: %v", err)
		}
		assertHeaderPosition(t, err)
	}
}

func TestParseFileHeaderRejectsDuplicateProfile(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("PROFILE \"a\"\nPROFILE \"b\"\n\n<body>x</body>")
	if err == nil {
		t.Fatal("expected duplicate PROFILE error")
	}
	if !strings.Contains(err.Error(), "duplicate PROFILE") {
		t.Fatalf("error must mention duplicate PROFILE, got: %v", err)
	}
	var he *HeaderError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HeaderError, got %T: %v", err, err)
	}
	if he.Line != 2 || he.Col != 1 {
		t.Fatalf("expected 2:1, got %d:%d (%v)", he.Line, he.Col, err)
	}
}
