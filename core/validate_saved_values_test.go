package core

import "testing"

func TestSaveOldWithFormTag(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	f := bindFormStruct{Name: "Bob"}
	SaveOld(c, f)
	if c.Get("old_name") != "Bob" {
		t.Errorf("expected 'Bob', got '%v'", c.Get("old_name"))
	}
}

func TestSaveErrorsEmpty(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	SaveErrors(c, map[string]string{})
	if len(c.data) != 0 {
		t.Error("expected no data for empty errors")
	}
}
