package lua

import (
	"fmt"
	"strings"
)

type emitter struct {
	features map[string]bool
	scopes   []map[string]bool
}

func newEmitter() *emitter {
	return &emitter{features: map[string]bool{}, scopes: []map[string]bool{{}}}
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
		if current.recursive {
			e.declare(current.name)
		}
		right, err := e.expression(current.value)
		if !current.recursive {
			e.declare(current.name)
		}
		return fmt.Sprintf("%slet %s = %s;\n", indent, jsIdentifier(current.name), right), err
	case assignStatement:
		if target, ok := current.target.(indexExpression); ok {
			object, err := e.expression(target.object)
			if err != nil {
				return "", err
			}
			index, err := e.expression(target.index)
			if err != nil {
				return "", err
			}
			right, err := e.expression(current.value)
			e.use("set")
			return fmt.Sprintf("%sglobalThis.dreegoLua.set(%s, %s, %s);\n", indent, object, index, right), err
		}
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
	case breakStatement:
		return indent + "break;\n", nil
	case whileStatement:
		return e.whileStatement(current, depth)
	case numericForStatement:
		return e.numericForStatement(current, depth)
	case ifStatement:
		return e.ifStatement(current, depth)
	default:
		return "", fmt.Errorf("unsupported Lua statement %T", value)
	}
}

func (e *emitter) whileStatement(value whileStatement, depth int) (string, error) {
	condition, err := e.expression(value.condition)
	if err != nil {
		return "", err
	}
	e.use("truthy")
	indent := strings.Repeat("  ", depth)
	e.pushScope()
	body, err := e.statements(value.body, depth+1)
	e.popScope()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%swhile (globalThis.dreegoLua.truthy(%s)) {\n%s%s}\n", indent, condition, body, indent), nil
}

func (e *emitter) numericForStatement(value numericForStatement, depth int) (string, error) {
	initial, err := e.expression(value.initial)
	if err != nil {
		return "", err
	}
	limit, err := e.expression(value.limit)
	if err != nil {
		return "", err
	}
	step, err := e.expression(value.step)
	if err != nil {
		return "", err
	}
	e.use("numericFor")
	e.pushScope(value.name)
	body, err := e.statements(value.body, depth+1)
	e.popScope()
	if err != nil {
		return "", err
	}
	indent := strings.Repeat("  ", depth)
	return fmt.Sprintf("%sfor (let %s of globalThis.dreegoLua.numericFor(%s, %s, %s)) {\n%s%s}\n", indent, jsIdentifier(value.name), initial, limit, step, body, indent), nil
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
			fmt.Fprintf(&out, "%sif (globalThis.dreegoLua.truthy(%s)) {\n", indent, condition)
		} else {
			fmt.Fprintf(&out, "%s} else if (globalThis.dreegoLua.truthy(%s)) {\n", indent, condition)
		}
		e.pushScope()
		body, err := e.statements(branch.body, depth+1)
		if err != nil {
			return "", err
		}
		out.WriteString(body)
		e.popScope()
	}
	if len(value.otherwise) > 0 {
		fmt.Fprintf(&out, "%s} else {\n", indent)
		e.pushScope()
		body, err := e.statements(value.otherwise, depth+1)
		if err != nil {
			return "", err
		}
		out.WriteString(body)
		e.popScope()
	}
	fmt.Fprintf(&out, "%s}\n", indent)
	return out.String(), nil
}

func (e *emitter) use(feature string) { e.features[feature] = true }

func (e *emitter) pushScope(names ...string) {
	scope := map[string]bool{}
	for _, name := range names {
		scope[name] = true
	}
	e.scopes = append(e.scopes, scope)
}

func (e *emitter) popScope() { e.scopes = e.scopes[:len(e.scopes)-1] }

func (e *emitter) declare(name string) { e.scopes[len(e.scopes)-1][name] = true }

func (e *emitter) isLocal(name string) bool {
	for index := len(e.scopes) - 1; index >= 0; index-- {
		if e.scopes[index][name] {
			return true
		}
	}
	return false
}
