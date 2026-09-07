package lua

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type emitter struct {
	features map[string]bool
}

func newEmitter() *emitter {
	return &emitter{features: map[string]bool{}}
}

func (e *emitter) program(value program) (string, error) {
	return e.statements(value.statements, 0)
}

func (e *emitter) statements(values []statement, depth int) (string, error) {
	var out strings.Builder
	for _, value := range values {
		code, err := e.statement(value, depth)
		if err != nil {
			return "", err
		}
		out.WriteString(code)
	}
	return out.String(), nil
}

func (e *emitter) statement(value statement, depth int) (string, error) {
	indent := strings.Repeat("  ", depth)
	switch current := value.(type) {
	case localStatement:
		right, err := e.expression(current.value)
		return fmt.Sprintf("%slet %s = %s;\n", indent, jsIdentifier(current.name), right), err
	case assignStatement:
		left, err := e.expression(current.target)
		if err != nil {
			return "", err
		}
		right, err := e.expression(current.value)
		return fmt.Sprintf("%s%s = %s;\n", indent, left, right), err
	case expressionStatement:
		expression, err := e.expression(current.value)
		return fmt.Sprintf("%s%s;\n", indent, expression), err
	case returnStatement:
		if current.value == nil {
			return indent + "return;\n", nil
		}
		value, err := e.expression(current.value)
		return fmt.Sprintf("%sreturn %s;\n", indent, value), err
	case ifStatement:
		return e.ifStatement(current, depth)
	default:
		return "", fmt.Errorf("unsupported Lua statement %T", value)
	}
}

func (e *emitter) ifStatement(value ifStatement, depth int) (string, error) {
	indent := strings.Repeat("  ", depth)
	var out strings.Builder
	for index, branch := range value.branches {
		condition, err := e.expression(branch.condition)
		if err != nil {
			return "", err
		}
		e.use("truthy")
		if index == 0 {
			fmt.Fprintf(&out, "%sif (dreegoLua.truthy(%s)) {\n", indent, condition)
		} else {
			fmt.Fprintf(&out, "%s} else if (dreegoLua.truthy(%s)) {\n", indent, condition)
		}
		body, err := e.statements(branch.body, depth+1)
		if err != nil {
			return "", err
		}
		out.WriteString(body)
	}
	if len(value.otherwise) > 0 {
		fmt.Fprintf(&out, "%s} else {\n", indent)
		body, err := e.statements(value.otherwise, depth+1)
		if err != nil {
			return "", err
		}
		out.WriteString(body)
	}
	fmt.Fprintf(&out, "%s}\n", indent)
	return out.String(), nil
}

func (e *emitter) expression(value expression) (string, error) {
	switch current := value.(type) {
	case literalExpression:
		return literal(current), nil
	case nameExpression:
		return jsIdentifier(current.name), nil
	case memberExpression:
		object, err := e.expression(current.object)
		return object + "." + current.field, err
	case callExpression:
		return e.call(current)
	case unaryExpression:
		return e.unary(current)
	case binaryExpression:
		return e.binary(current)
	case functionExpression:
		params := make([]string, 0, len(current.params))
		for _, param := range current.params {
			params = append(params, jsIdentifier(param))
		}
		body, err := e.statements(current.body, 1)
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
	if name, ok := value.callee.(nameExpression); ok && name.name == "print" {
		e.use("print")
		callee = "dreegoLua.print"
	} else if name, ok := value.callee.(nameExpression); ok && forbiddenCalls[name.name] {
		return "", fmt.Errorf("Lua browser MVP: %s is not available", name.name)
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

var forbiddenCalls = map[string]bool{
	"collectgarbage": true,
	"dofile":         true,
	"load":           true,
	"loadfile":       true,
	"require":        true,
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
		return "_lua_" + name
	}
	return name
}

func (e *emitter) unary(value unaryExpression) (string, error) {
	right, err := e.expression(value.right)
	if err != nil {
		return "", err
	}
	if value.op == tokenNot {
		e.use("truthy")
		return "!dreegoLua.truthy(" + right + ")", nil
	}
	e.use("neg")
	return "dreegoLua.neg(" + right + ")", nil
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
		return fmt.Sprintf("dreegoLua.and(%s, () => %s)", left, right), nil
	case tokenOr:
		e.use("or")
		return fmt.Sprintf("dreegoLua.or(%s, () => %s)", left, right), nil
	case tokenConcat:
		e.use("concat")
		return fmt.Sprintf("dreegoLua.concat(%s, %s)", left, right), nil
	case tokenPlus, tokenMinus, tokenStar, tokenSlash, tokenPercent:
		feature := map[tokenKind]string{tokenPlus: "add", tokenMinus: "sub", tokenStar: "mul", tokenSlash: "div", tokenPercent: "mod"}[value.op]
		e.use(feature)
		return fmt.Sprintf("dreegoLua.%s(%s, %s)", feature, left, right), nil
	case tokenLess, tokenLessEqual, tokenGreater, tokenGreaterEqual:
		feature := map[tokenKind]string{tokenLess: "lt", tokenLessEqual: "le", tokenGreater: "gt", tokenGreaterEqual: "ge"}[value.op]
		e.use(feature)
		return fmt.Sprintf("dreegoLua.%s(%s, %s)", feature, left, right), nil
	}
	operator := map[tokenKind]string{
		tokenEqual: "===", tokenNotEqual: "!==",
	}[value.op]
	return fmt.Sprintf("(%s %s %s)", left, operator, right), nil
}

func (e *emitter) use(feature string) { e.features[feature] = true }

func (e *emitter) sortedFeatures() []string {
	features := make([]string, 0, len(e.features))
	for feature := range e.features {
		features = append(features, feature)
	}
	sort.Strings(features)
	return features
}
