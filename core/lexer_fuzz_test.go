package core

import (
	"reflect"
	"testing"
)

func FuzzLexer(f *testing.F) {
	seeds := []string{
		"",
		"<div>hello</div>",
		"<go>msg := \"hi\"</go>\n<div><p>{{ msg }}</p></div>",
		`<go method="post">x := 1</go><div>{#if x}<p>yes</p>{/if}</div>`,
		`<div>{#each items as item}<li>{{ item.name }}</li>{/each}</div>`,
		`<head><title>T</title></head><div>body</div>`,
		`<script>let a = 1 < 2;</script><style>.a > .b { color: red }</style><div>x</div>`,
		`<div><@Card title="Hi" count={n}/><@Card>slot</@Card></div>`,
		`<div>{#verbatim}{{ raw }}{/verbatim}</div>`,
		`<div>{#slot name}content{/slot}{#slot}default{/slot}</div>`,
		`<div><input value="x>y" data-x='a{b'><a href="/p/{{ id }}">link</a></div>`,
		`<div><p>{{ user.name | upper }}</p></div>`,
		`<div><div><p>nested</p></div></div>`,
		`<go>if a < b { return a }</go>`,
		`<@Comp/>`,
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 64<<10 {
			t.Skip("input too large")
		}

		tokens, err := Lex(input)

		if err == nil {
			lexProperties(t, input, tokens)
		}

		again, err2 := Lex(input)
		if err != nil || err2 != nil {
			return
		}
		if !reflect.DeepEqual(tokens, again) {
			t.Fatalf("non-deterministic lex output for %q", input)
		}
	})
}

func lexProperties(t *testing.T, input string, tokens []Token) {
	t.Helper()

	maxTokens := 2*len(input) + 2
	if len(tokens) > maxTokens {
		t.Fatalf("token count %d exceeds bound %d for input of length %d",
			len(tokens), maxTokens, len(input))
	}

	if n := len(tokens); n == 0 || tokens[n-1].Type != TokenEOF {
		t.Fatalf("expected EOF as last token, got %d tokens", n)
	}

	last := 0
	for _, tok := range tokens {
		if tok.Pos < last {
			t.Fatalf("token positions not monotonic: %d after %d", tok.Pos, last)
		}
		last = tok.Pos
		if tok.Pos < 0 || tok.Pos > len(input) {
			t.Fatalf("token position %d out of bounds for input of length %d",
				tok.Pos, len(input))
		}
		if tok.Type == TokenText {
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
