package lexer

import (
	"reflect"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func FuzzLexer(f *testing.F) {
	seeds := []string{
		"",
		"<body>hello</body>",
		"<server>msg := \"hi\"</server>\n<body><p>{{ msg }}</p></body>",
		`<server method="post">x := 1</server><body>{#if x}<p>yes</p>{/if}</body>`,
		`<body>{#each items as item}<li>{{ item.name }}</li>{/each}</body>`,
		`<head><title>T</title></head><body>body</body>`,
		`<client>let a = 1 < 2;</client><style>.a > .b { color: red }</style><body>x</body>`,
		`<body><@Card title="Hi" count={n}/><@Card>slot</@Card></body>`,
		`<body>{#verbatim}{{ raw }}{/verbatim}</body>`,
		`<body>{#slot name}content{/slot}{#slot}default{/slot}</body>`,
		`<body><input value="x>y" data-x='a{b'><a href="/p/{{ id }}">link</a></body>`,
		`<body><p>{{ user.name | upper }}</p></body>`,
		`<body><div><p>nested</p></div></body>`,
		`<server>if a < b { return a }</server>`,
		`<@Comp/>`,
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 64<<10 {
			t.Skip("input too large")
		}

		toks, err := Lex(input)

		if err == nil {
			lexProperties(t, input, toks)
		}

		again, err2 := Lex(input)
		if err != nil || err2 != nil {
			return
		}
		if !reflect.DeepEqual(toks, again) {
			t.Fatalf("non-deterministic lex output for %q", input)
		}
	})
}

func lexProperties(t *testing.T, input string, toks []tokens.Token) {
	t.Helper()

	maxTokens := 2*len(input) + 2
	if len(toks) > maxTokens {
		t.Fatalf("token count %d exceeds bound %d for input of length %d",
			len(toks), maxTokens, len(input))
	}

	if n := len(toks); n == 0 || toks[n-1].Type != tokens.TokenEOF {
		t.Fatalf("expected EOF as last token, got %d tokens", n)
	}

	last := 0
	for _, tok := range toks {
		if tok.Pos < last {
			t.Fatalf("token positions not monotonic: %d after %d", tok.Pos, last)
		}
		last = tok.Pos
		if tok.Pos < 0 || tok.Pos > len(input) {
			t.Fatalf("token position %d out of bounds for input of length %d",
				tok.Pos, len(input))
		}
		if tok.Type == tokens.TokenText {
			if tok.Pos+len(tok.Value) > len(input) {
				t.Fatalf("text token %q at %d overruns input of length %d",
					tok.Value, tok.Pos, len(input))
			}
			if input[tok.Pos:tok.Pos+len(tok.Value)] != tok.Value {
				t.Fatalf("text token %q does not match input at %d", tok.Value, tok.Pos)
			}
		}
	}
}
