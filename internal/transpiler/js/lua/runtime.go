package lua

import (
	"maps"
	"slices"
	"strings"
)

type helper struct {
	dependencies []string
	code         string
}

var helpers = map[string]helper{
	"tableCore": {code: "tableCore: Symbol('dreegoLuaTable')"},
	"tableKey":  {code: "tableKey: key => { if (key === null || key === undefined) throw new TypeError('Lua table index is nil'); if (typeof key === 'number' && Number.isNaN(key)) throw new TypeError('Lua table index is NaN'); return key }"},
	"table": {
		dependencies: []string{"tableCore", "tableKey"},
		code:         "table: fields => { const entries = new Map(); for (const [key, value] of fields) { const checked = globalThis.dreegoLua.tableKey(key); if (value !== null && value !== undefined) entries.set(checked, value) } const target = {[globalThis.dreegoLua.tableCore]: entries}; return new Proxy(target, {get: (object, key) => key === globalThis.dreegoLua.tableCore ? entries : entries.get(globalThis.dreegoLua.tableKey(key)) ?? null, set: (object, key, value) => { const checked = globalThis.dreegoLua.tableKey(key); if (value === null || value === undefined) entries.delete(checked); else entries.set(checked, value); return true }}) }",
	},
	"get": {
		dependencies: []string{"tableCore", "tableKey"},
		code:         "get: (table, key) => { const entries = table?.[globalThis.dreegoLua.tableCore]; if (!(entries instanceof Map)) throw new TypeError('Lua table operation requires a table'); return entries.get(globalThis.dreegoLua.tableKey(key)) ?? null }",
	},
	"set": {
		dependencies: []string{"tableCore", "tableKey"},
		code:         "set: (table, key, value) => { const entries = table?.[globalThis.dreegoLua.tableCore]; if (!(entries instanceof Map)) throw new TypeError('Lua table operation requires a table'); const checked = globalThis.dreegoLua.tableKey(key); if (value === null || value === undefined) entries.delete(checked); else entries.set(checked, value); return value }",
	},
	"length": {
		dependencies: []string{"tableCore"},
		code:         "length: table => { const entries = table?.[globalThis.dreegoLua.tableCore]; if (!(entries instanceof Map)) throw new TypeError('Lua length requires a table'); let length = 0; while (entries.has(length + 1)) length++; return length }",
	},
	"numericFor": {
		code: "numericFor: function* (initial, limit, step) { if (![initial, limit, step].every(Number.isFinite)) throw new TypeError('Lua numeric for requires finite numbers'); if (step === 0) throw new RangeError('Lua numeric for step cannot be zero'); if (step > 0) { for (let value = initial; value <= limit; value += step) yield value } else { for (let value = initial; value >= limit; value += step) yield value } }",
	},
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
	names := slices.Sorted(maps.Keys(selected))
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
	names := slices.Sorted(maps.Keys(features))
	return Bundle(strings.Join(names, ","))
}
