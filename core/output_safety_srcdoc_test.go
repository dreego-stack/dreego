package core

import (
	"strings"
	"testing"
)

func TestSafeSrcdocKeepsMarkupInertAfterNestedParse(t *testing.T) {
	got := SafeSrcdoc(`<script>parent.postMessage("run", "*")</script>`)
	if strings.Contains(got, "&lt;script") {
		t.Fatalf("srcdoc is escaped only once: %q", got)
	}
	if !strings.Contains(got, "&amp;lt;script&amp;gt;") {
		t.Fatalf("srcdoc does not preserve escaped markup for its nested document: %q", got)
	}
}

func TestSrcdocUsesNestedDocumentContext(t *testing.T) {
	if got := attrContext("srcdoc"); got != "SafeSrcdoc" {
		t.Fatalf("attrContext(srcdoc) = %q, want SafeSrcdoc", got)
	}
}
