package lua

import (
	"sort"
	"strings"
)

type helper struct {
	dependencies []string
	code         string
}

var helpers = map[string]helper{
	"truthy": {code: "truthy: value => value !== false && value !== null && value !== undefined"},
	"and": {
		dependencies: []string{"truthy"},
		code:         "and: (left, right) => globalThis.dreegoLua.truthy(left) ? right() : left",
	},
	"or": {
		dependencies: []string{"truthy"},
		code:         "or: (left, right) => globalThis.dreegoLua.truthy(left) ? left : right()",
	},
	"print":       {code: "print: (...values) => console.log(...values)"},
	"number":      {code: "number: value => { if (typeof value !== 'number') throw new TypeError('Lua arithmetic requires numbers'); return value }"},
	"compare":     {code: "compare: (left, right) => { if (typeof left !== typeof right || !['number','string'].includes(typeof left)) throw new TypeError('Lua comparison requires two numbers or two strings'); return [left, right] }"},
	"add":         {dependencies: []string{"number"}, code: "add: (left, right) => globalThis.dreegoLua.number(left) + globalThis.dreegoLua.number(right)"},
	"sub":         {dependencies: []string{"number"}, code: "sub: (left, right) => globalThis.dreegoLua.number(left) - globalThis.dreegoLua.number(right)"},
	"mul":         {dependencies: []string{"number"}, code: "mul: (left, right) => globalThis.dreegoLua.number(left) * globalThis.dreegoLua.number(right)"},
	"div":         {dependencies: []string{"number"}, code: "div: (left, right) => globalThis.dreegoLua.number(left) / globalThis.dreegoLua.number(right)"},
	"neg":         {dependencies: []string{"number"}, code: "neg: value => -globalThis.dreegoLua.number(value)"},
	"lt":          {dependencies: []string{"compare"}, code: "lt: (left, right) => { [left, right] = globalThis.dreegoLua.compare(left, right); return left < right }"},
	"le":          {dependencies: []string{"compare"}, code: "le: (left, right) => { [left, right] = globalThis.dreegoLua.compare(left, right); return left <= right }"},
	"gt":          {dependencies: []string{"compare"}, code: "gt: (left, right) => { [left, right] = globalThis.dreegoLua.compare(left, right); return left > right }"},
	"ge":          {dependencies: []string{"compare"}, code: "ge: (left, right) => { [left, right] = globalThis.dreegoLua.compare(left, right); return left >= right }"},
	"concatValue": {code: "concatValue: value => { if (typeof value !== 'string' && typeof value !== 'number') throw new TypeError('Lua concatenation requires strings or numbers'); return String(value) }"},
	"concat": {
		dependencies: []string{"concatValue"},
		code:         "concat: (left, right) => globalThis.dreegoLua.concatValue(left) + globalThis.dreegoLua.concatValue(right)",
	},
	"mod": {dependencies: []string{"number"}, code: "mod: (left, right) => { left = globalThis.dreegoLua.number(left); right = globalThis.dreegoLua.number(right); if (right === 0) throw new RangeError('Lua modulo by zero'); return ((left % right) + right) % right }"},
}

func Bundle(features string) string {
	if features == "" {
		return ""
	}
	selected := map[string]bool{}
	var include func(string)
	include = func(name string) {
		if selected[name] {
			return
		}
		selected[name] = true
		for _, dependency := range helpers[name].dependencies {
			include(dependency)
		}
	}
	for _, feature := range strings.Split(features, ",") {
		include(feature)
	}
	names := make([]string, 0, len(selected))
	for name := range selected {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, helpers[name].code)
	}
	return "globalThis.dreegoLua = Object.freeze({...globalThis.dreegoLua," + strings.Join(parts, ",") + "});"
}

func FeatureList(features string) []string {
	if features == "" {
		return nil
	}
	return strings.Split(features, ",")
}

func BundleFeatures(features map[string]bool) string {
	names := make([]string, 0, len(features))
	for name := range features {
		names = append(names, name)
	}
	sort.Strings(names)
	return Bundle(strings.Join(names, ","))
}
