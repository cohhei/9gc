package main

import (
	"reflect"
	"testing"
)

var tokenEof = &Token{kind: TK_EOF}

func TestTokenize(t *testing.T) {
	testCases := []struct {
		str      string
		expected *Token
	}{
		{
			" 1 ",
			&Token{TK_NUM, tokenEof, 1, "1"},
		},
		{
			"0 + 45 - 5 ",
			&Token{
				kind: TK_NUM,
				val:  0,
				str:  "0",
				next: &Token{
					kind: TK_RESERVED,
					str:  "+",
					next: &Token{
						kind: TK_NUM,
						val:  45,
						str:  "45",
						next: &Token{
							kind: TK_RESERVED,
							str:  "-",
							next: &Token{
								kind: TK_NUM,
								val:  5,
								str:  "5",
								next: tokenEof,
							},
						},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.str, func(t *testing.T) {
			actual, err := tokenize(tC.str)

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, tC.expected) {
				t.Fatalf("Tokenizing '%s' failed.\nactual:\t%+v\nexpected:\t%+v\n", tC.str, actual, tC.expected)
			}
		})
	}
}
