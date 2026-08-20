package validate

import "testing"

func TestSaveOldWithFormTag(t *testing.T) {
	c := &testSetter{}
	f := bindFormStruct{Name: "Bob"}
	SaveOld(c, f)
	if c.data["old_name"] != "Bob" {
		t.Errorf("expected 'Bob', got '%v'", c.data["old_name"])
	}
}

func TestSaveErrorsEmpty(t *testing.T) {
	c := &testSetter{}
	SaveErrors(c, map[string]string{})
	if len(c.data) != 0 {
		t.Error("expected no data for empty errors")
	}
}
