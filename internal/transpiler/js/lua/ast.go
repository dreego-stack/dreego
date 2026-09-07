package lua

type program struct {
	statements []statement
}

type statement interface {
	statementNode()
}

type expression interface {
	expressionNode()
}

type localStatement struct {
	name  string
	value expression
}

type assignStatement struct {
	target expression
	value  expression
}

type expressionStatement struct{ value expression }

type ifBranch struct {
	condition expression
	body      []statement
}

type ifStatement struct {
	branches  []ifBranch
	otherwise []statement
}

type literalExpression struct {
	kind  tokenKind
	value string
}

type nameExpression struct{ name string }

type unaryExpression struct {
	op    tokenKind
	right expression
}

type binaryExpression struct {
	left  expression
	op    tokenKind
	right expression
}

type memberExpression struct {
	object expression
	field  string
}

type callExpression struct {
	callee expression
	args   []expression
}

type functionExpression struct {
	params []string
	body   []statement
}

func (localStatement) statementNode()      {}
func (assignStatement) statementNode()     {}
func (expressionStatement) statementNode() {}
func (ifStatement) statementNode()         {}

func (literalExpression) expressionNode()  {}
func (nameExpression) expressionNode()     {}
func (unaryExpression) expressionNode()    {}
func (binaryExpression) expressionNode()   {}
func (memberExpression) expressionNode()   {}
func (callExpression) expressionNode()     {}
func (functionExpression) expressionNode() {}
