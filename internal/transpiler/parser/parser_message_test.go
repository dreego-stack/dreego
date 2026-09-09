package parser

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func TestParseMessageExpression(t *testing.T) {
	toks, err := lexer.Lex(`<body><p>[[ cart.items count=len(items) owner=user.Name ]]</p></body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(toks).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var message *ir.TemplateNode
	for i := range file.Body.Nodes {
		if file.Body.Nodes[i].Type == ir.NodeMessage {
			message = &file.Body.Nodes[i]
		}
	}
	if message == nil {
		t.Fatalf("message node not found: %+v", file.Body.Nodes)
	}
	if message.MessageKey != "cart.items" {
		t.Errorf("key = %q", message.MessageKey)
	}
	want := []ir.MessageArgument{{Name: "count", Expression: "len(items)"}, {Name: "owner", Expression: "user.Name"}}
	if len(message.MessageArgs) != len(want) {
		t.Fatalf("arguments = %+v", message.MessageArgs)
	}
	for i := range want {
		if message.MessageArgs[i] != want[i] {
			t.Errorf("argument %d = %+v, want %+v", i, message.MessageArgs[i], want[i])
		}
	}
}

func TestParseMessageExpressionRejectsInvalidSyntax(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want string
	}{
		{"empty", `<body>[[ ]]</body>`, "message key is required"},
		{"invalid key", `<body>[[ cart..items ]]</body>`, `invalid message key "cart..items"`},
		{"missing value", `<body>[[ cart.items count= ]]</body>`, `message argument "count" requires an expression`},
		{"duplicate argument", `<body>[[ cart.items count=1 count=2 ]]</body>`, `duplicate message argument "count"`},
		{"positional argument", `<body>[[ cart.items count ]]</body>`, `invalid message argument "count"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			toks, err := lexer.Lex(tc.src)
			if err != nil {
				t.Fatalf("lex: %v", err)
			}
			_, err = NewParser(toks).Parse()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestParseMessageExpressionAllowsGoExpressionsWithSpaces(t *testing.T) {
	key, arguments, err := ParseMessageExpression(`home.greeting name=fmt.Sprintf("%s %s", user.First, user.Last) count=len(items)`, 0)
	if err != nil {
		t.Fatal(err)
	}
	if key != "home.greeting" || len(arguments) != 2 {
		t.Fatalf("key = %q, arguments = %+v", key, arguments)
	}
	if arguments[0].Expression != `fmt.Sprintf("%s %s", user.First, user.Last)` {
		t.Fatalf("expression = %q", arguments[0].Expression)
	}
}
