package lua

import (
	"strings"
	"testing"
)

func TestCompileAcceptsDocumentedLiteralAndOperatorMatrix(t *testing.T) {
	t.Parallel()
	programs := map[string]string{
		"nil":            `local value = nil`,
		"single quote":   `local value = 'text'`,
		"comparison":     `local value = 1 <= 2`,
		"inequality":     `local value = 1 ~= 2`,
		"multiplication": `local value = 2 * 3`,
		"division":       `local value = 6 / 2`,
		"logical and":    `local value = true and "yes"`,
		"logical not":    `local value = not false`,
		"line comment":   "-- ignored\nlocal value = true",
		"method call":    `document:querySelector("main")`,
		"table":          `local value = { key = "value" }`,
		"while loop":     `while false do print("loop") end`,
		"numeric for":    `for i = 1, 3 do print(i) end`,
		"local restricted call": `local setmetatable = function(value) return value end
local value = setmetatable("safe")`,
		"local restricted root": `local coroutine = document
coroutine.createElement("div")`,
	}
	for name, source := range programs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			artifact, err := Compile(source)
			if err != nil {
				t.Fatalf("Compile() error = %v", err)
			}
			if artifact.Code == "" {
				t.Fatal("Compile() returned empty JavaScript")
			}
		})
	}
}

func TestCompileRejectsDocumentedUnsupportedSyntax(t *testing.T) {
	t.Parallel()
	programs := map[string]string{
		"generic for":     `for key, value in pairs(items) do print(key) end`,
		"varargs":         `local function values(...) return ... end`,
		"multiple return": `local function pair() return 1, 2 end`,
		"metatable":       `setmetatable(value, meta)`,
		"coroutine":       `coroutine.create(function() end)`,
		"coroutine value": `local value = coroutine`,
		"loadfile":        `loadfile("module.lua")`,
		"collectgarbage":  `collectgarbage()`,
	}
	for name, source := range programs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := Compile(source); err == nil {
				t.Fatalf("Compile(%q) succeeded", source)
			}
		})
	}
}

func TestRuntimeBundleIsStableForFeatureOrderAndDuplicates(t *testing.T) {
	t.Parallel()
	first := Bundle("truthy,print,concat,print")
	second := Bundle("concat,print,truthy")
	if first != second {
		t.Fatalf("bundle depends on feature order or duplicates:\nfirst: %s\nsecond: %s", first, second)
	}
	for _, helper := range []string{"console.log", "concatValue", "truthy"} {
		if !strings.Contains(first, helper) {
			t.Fatalf("bundle missing %q: %s", helper, first)
		}
	}
}
