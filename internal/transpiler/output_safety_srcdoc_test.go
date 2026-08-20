package transpiler

import "testing"

func TestSrcdocUsesNestedDocumentContext(t *testing.T) {
	if got := attrContext("srcdoc"); got != "SafeSrcdoc" {
		t.Fatalf("attrContext(srcdoc) = %q, want SafeSrcdoc", got)
	}
}
