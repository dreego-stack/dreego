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
	name      string
	value     expression
	recursive bool
}

type assignStatement struct {
	target expression
	value  expression
}

type expressionStatement struct{ value expression }

type returnStatement struct{ value expression }

type breakStatement struct{}

type whileStatement struct {
	condition expression
	body      []statement
}

type numericForStatement struct {
	name                 string
	initial, limit, step expression
	body                 []statement
}

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

type nameExpression struct {
	name         string
	line, column int
}

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

type indexExpression struct {
	object expression
	index  expression
}

type tableField struct {
	key   expression
	value expression
}

type tableExpression struct{ fields []tableField }

type functionExpression struct {
	params []string
	body   []statement
}

func (localStatement) statementNode()      {}
func (assignStatement) statementNode()     {}
func (expressionStatement) statementNode() {}
func (returnStatement) statementNode()     {}
func (breakStatement) statementNode()      {}
func (whileStatement) statementNode()      {}
func (numericForStatement) statementNode() {}
func (ifStatement) statementNode()         {}

func (literalExpression) expressionNode()  {}
func (nameExpression) expressionNode()     {}
func (unaryExpression) expressionNode()    {}
func (binaryExpression) expressionNode()   {}
func (memberExpression) expressionNode()   {}
func (callExpression) expressionNode()     {}
func (indexExpression) expressionNode()    {}
func (tableExpression) expressionNode()    {}
func (functionExpression) expressionNode() {}
