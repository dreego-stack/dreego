package core

import "testing"

func TestTruthyMatrix(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"nil", nil, false},
		{"false", false, false},
		{"true", true, true},
		{"empty string", "", false},
		{"non-empty string", "x", true},
		{"whitespace string", " ", true},
		{"zero int", 0, false},
		{"non-zero int", -1, true},
		{"zero float", 0.0, false},
		{"non-zero float", 2.5, true},
		{"zero uint", uint(0), false},
		{"non-zero uint", uint(3), true},
		{"empty slice", []string{}, false},
		{"non-empty slice", []string{""}, true},
		{"empty map", map[string]string{}, false},
		{"non-empty map", map[string]string{"a": ""}, true},
		{"empty array", [0]int{}, false},
		{"non-empty array", [2]int{0, 0}, true},
		{"nil pointer", (*int)(nil), false},
		{"non-nil pointer", new(int), true},
		{"nil interface", any(nil), false},
		{"struct", struct{ A int }{}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Truthy(c.v); got != c.want {
				t.Errorf("Truthy(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}

func TestTruthyZeroValueStruct(t *testing.T) {
	type form struct {
		Email string
	}
	if !Truthy(form{}) {
		t.Error("zero-value struct must be truthy: structs have no empty state")
	}
}
