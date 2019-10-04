package main

import (
	"fmt"
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
			&Token{TK_NUM, tokenEof, 1, "1", 1},
		},
		{
			"0 + 45 - 5 ",
			&Token{
				kind: TK_NUM,
				val:  0,
				str:  "0",
				len:  1,
				next: &Token{
					kind: TK_RESERVED,
					str:  "+",
					len:  1,
					next: &Token{
						kind: TK_NUM,
						val:  45,
						str:  "45",
						len:  2,
						next: &Token{
							kind: TK_RESERVED,
							str:  "-",
							len:  1,
							next: &Token{
								kind: TK_NUM,
								val:  5,
								str:  "5",
								len:  1,
								next: tokenEof,
							},
						},
					},
				},
			},
		},
		{
			"5*(9-6)",
			&Token{
				kind: TK_NUM,
				val:  5,
				str:  "5",
				len:  1,
				next: &Token{
					kind: TK_RESERVED,
					str:  "*",
					len:  1,
					next: &Token{
						kind: TK_RESERVED,
						str:  "(",
						len:  1,
						next: &Token{
							kind: TK_NUM,
							str:  "9",
							val:  9,
							len:  1,
							next: &Token{
								kind: TK_RESERVED,
								str:  "-",
								len:  1,
								next: &Token{
									kind: TK_NUM,
									val:  6,
									str:  "6",
									len:  1,
									next: &Token{
										kind: TK_RESERVED,
										str:  ")",
										len:  1,
										next: tokenEof,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"42!=42",
			&Token{
				kind: TK_NUM,
				val:  42,
				str:  "42",
				len:  2,
				next: &Token{
					kind: TK_RESERVED,
					str:  "!=",
					len:  2,
					next: &Token{
						kind: TK_NUM,
						val:  42,
						str:  "42",
						len:  2,
						next: tokenEof,
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
				t.Fatalf("Tokenizing '%s' failed.\nactual:\n%+v\nexpected:\n%+v\n", tC.str, showTokens(actual), showTokens(tC.expected))
			}
		})
	}
}

func showTokens(t *Token) string {
	if t.next == nil {
		return ""
	}
	return fmt.Sprintf("%+v\n%+v", t, showTokens(t.next))
}
