package lua

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (e *emitter) expression(value expression) (string, error) {
	switch current := value.(type) {
	case literalExpression:
		return literal(current), nil
	case nameExpression:
		if current.name == "print" && !e.isLocal("print") {
			e.use("print")
			return "globalThis.dreegoLua.print", nil
		}
		if (forbiddenCalls[current.name] || forbiddenRoots[current.name]) && !e.isLocal(current.name) {
			return "", fmt.Errorf("Lua browser MVP: %s is not available", current.name)
		}
		return jsIdentifier(current.name), nil
	case memberExpression:
		if root, ok := memberRoot(current); ok && forbiddenRoots[root] && !e.isLocal(root) {
			return "", fmt.Errorf("Lua browser MVP: %s is not available", root)
		}
		object, err := e.expression(current.object)
		return object + "." + current.field, err
	case callExpression:
		return e.call(current)
	case indexExpression:
		object, err := e.expression(current.object)
		if err != nil {
			return "", err
		}
		index, err := e.expression(current.index)
		e.use("get")
		return fmt.Sprintf("globalThis.dreegoLua.get(%s, %s)", object, index), err
	case tableExpression:
		return e.table(current)
	case unaryExpression:
		return e.unary(current)
	case binaryExpression:
		return e.binary(current)
	case functionExpression:
		e.pushScope(current.params...)
		params := make([]string, 0, len(current.params))
		for _, param := range current.params {
			params = append(params, jsIdentifier(param))
		}
		body, err := e.statements(current.body, 1)
		e.popScope()
		if err != nil {
			return "", err
		}
		return "(" + strings.Join(params, ", ") + ") => {\n" + body + "}", nil
	default:
		return "", fmt.Errorf("unsupported Lua expression %T", value)
	}
}

func literal(value literalExpression) string {
	switch value.kind {
	case tokenString:
		return strconv.Quote(value.value)
	case tokenTrue:
		return "true"
	case tokenFalse:
		return "false"
	case tokenNil:
		return "null"
	default:
		return value.value
	}
}

func (e *emitter) call(value callExpression) (string, error) {
	callee, err := e.expression(value.callee)
	if err != nil {
		return "", err
	}
	args := make([]string, 0, len(value.args))
	for _, argument := range value.args {
		code, err := e.expression(argument)
		if err != nil {
			return "", err
		}
		args = append(args, code)
	}
	return callee + "(" + strings.Join(args, ", ") + ")", nil
}

func (e *emitter) unary(value unaryExpression) (string, error) {
	right, err := e.expression(value.right)
	if err != nil {
		return "", err
	}
	if value.op == tokenNot {
		e.use("truthy")
		return "!globalThis.dreegoLua.truthy(" + right + ")", nil
	}
	if value.op == tokenHash {
		e.use("length")
		return "globalThis.dreegoLua.length(" + right + ")", nil
	}
	e.use("neg")
	return "globalThis.dreegoLua.neg(" + right + ")", nil
}

func (e *emitter) table(value tableExpression) (string, error) {
	fields := make([]string, 0, len(value.fields))
	for _, field := range value.fields {
		key, err := e.expression(field.key)
		if err != nil {
			return "", err
		}
		item, err := e.expression(field.value)
		if err != nil {
			return "", err
		}
		fields = append(fields, "["+key+", "+item+"]")
	}
	e.use("table")
	return "globalThis.dreegoLua.table([" + strings.Join(fields, ", ") + "])", nil
}

func (e *emitter) binary(value binaryExpression) (string, error) {
	left, err := e.expression(value.left)
	if err != nil {
		return "", err
	}
	right, err := e.expression(value.right)
	if err != nil {
		return "", err
	}
	switch value.op {
	case tokenAnd:
		e.use("and")
		return fmt.Sprintf("globalThis.dreegoLua.and(%s, () => %s)", left, right), nil
	case tokenOr:
		e.use("or")
		return fmt.Sprintf("globalThis.dreegoLua.or(%s, () => %s)", left, right), nil
	case tokenConcat:
		e.use("concat")
		return fmt.Sprintf("globalThis.dreegoLua.concat(%s, %s)", left, right), nil
	case tokenPlus, tokenMinus, tokenStar, tokenSlash, tokenPercent:
		feature := map[tokenKind]string{tokenPlus: "add", tokenMinus: "sub", tokenStar: "mul", tokenSlash: "div", tokenPercent: "mod"}[value.op]
		e.use(feature)
		return fmt.Sprintf("globalThis.dreegoLua.%s(%s, %s)", feature, left, right), nil
	case tokenLess, tokenLessEqual, tokenGreater, tokenGreaterEqual:
		feature := map[tokenKind]string{tokenLess: "lt", tokenLessEqual: "le", tokenGreater: "gt", tokenGreaterEqual: "ge"}[value.op]
		e.use(feature)
		return fmt.Sprintf("globalThis.dreegoLua.%s(%s, %s)", feature, left, right), nil
	}
	operator := map[tokenKind]string{tokenEqual: "===", tokenNotEqual: "!=="}[value.op]
	return fmt.Sprintf("(%s %s %s)", left, operator, right), nil
}

var reservedJavaScriptNames = map[string]bool{
	"await": true, "break": true, "case": true, "catch": true, "class": true,
	"const": true, "continue": true, "debugger": true, "default": true,
	"delete": true, "do": true, "else": true, "export": true, "extends": true,
	"finally": true, "for": true, "function": true, "if": true, "import": true,
	"in": true, "instanceof": true, "let": true, "new": true, "return": true,
	"static": true, "super": true, "switch": true, "this": true, "throw": true,
	"try": true, "typeof": true, "var": true, "void": true, "while": true,
	"with": true, "yield": true,
}

func jsIdentifier(name string) string {
	if reservedJavaScriptNames[name] {
		return "$lua_" + name
	}
	return name
}

func (e *emitter) sortedFeatures() []string {
	features := make([]string, 0, len(e.features))
	for feature := range e.features {
		features = append(features, feature)
	}
	sort.Strings(features)
	return features
}
