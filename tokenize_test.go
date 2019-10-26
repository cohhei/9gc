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
		{
			"abc=1;abc",
			&Token{
				kind: TK_IDENT,
				str:  "abc",
				len:  3,
				next: &Token{
					kind: TK_RESERVED,
					str:  "=",
					len:  1,
					next: &Token{
						kind: TK_NUM,
						val:  1,
						str:  "1",
						len:  1,
						next: &Token{
							kind: TK_RESERVED,
							str:  ";",
							len:  1,
							next: &Token{
								kind: TK_IDENT,
								str:  "abc",
								len:  3,
								next: tokenEof,
							},
						},
					},
				},
			},
		},
		{
			"return 5;",
			&Token{
				kind: TK_RETURN,
				str:  "return",
				len:  6,
				next: &Token{
					kind: TK_NUM,
					str:  "5",
					len:  1,
					val:  5,
					next: &Token{
						kind: TK_RESERVED,
						str:  ";",
						len:  1,
						next: tokenEof,
					},
				},
			},
		},
		{
			"returned;",
			&Token{
				kind: TK_IDENT,
				str:  "returned",
				len:  8,
				next: &Token{
					kind: TK_RESERVED,
					str:  ";",
					len:  1,
					next: tokenEof,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.str, func(t *testing.T) {
			if err := tokenize(tC.str); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(token, tC.expected) {
				t.Fatalf("Tokenizing '%s' failed.\nactual:\n%+v\nexpected:\n%+v\n", tC.str, showTokens(token), showTokens(tC.expected))
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
