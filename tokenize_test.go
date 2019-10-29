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
		expected []*Token
	}{
		{
			" 1 ",
			[]*Token{{"1", 1, 1, TK_NUM, tokenEof}},
		},
		{
			"0 + 45 - 5 ",
			[]*Token{
				{"0", 1, 0, TK_NUM, nil},
				{"+", 1, 0, TK_RESERVED, nil},
				{"45", 2, 45, TK_NUM, nil},
				{"-", 1, 0, TK_RESERVED, nil},
				{"5", 1, 5, TK_NUM, nil},
				tokenEof,
			},
		},
		{
			"5*(9-6)",
			[]*Token{
				{"5", 1, 5, TK_NUM, nil},
				{"*", 1, 0, TK_RESERVED, nil},
				{"(", 1, 0, TK_RESERVED, nil},
				{"9", 1, 9, TK_NUM, nil},
				{"-", 1, 0, TK_RESERVED, nil},
				{"6", 1, 6, TK_NUM, nil},
				{")", 1, 0, TK_RESERVED, nil},
				tokenEof,
			},
		},
		{
			"42!=42",
			[]*Token{
				{"42", 2, 42, TK_NUM, nil},
				{"!=", 2, 0, TK_RESERVED, nil},
				{"42", 2, 42, TK_NUM, nil},
				tokenEof,
			},
		},
		{
			"abc=1;abc",
			[]*Token{
				{"abc", 3, 0, TK_IDENT, nil},
				{"=", 1, 0, TK_RESERVED, nil},
				{"1", 1, 1, TK_NUM, nil},
				{";", 1, 0, TK_RESERVED, nil},
				{"abc", 3, 0, TK_IDENT, nil},
				tokenEof,
			},
		},
		{
			"return 5;",
			[]*Token{
				{"return", 6, 0, TK_RESERVED, nil},
				{"5", 1, 5, TK_NUM, nil},
				{";", 1, 0, TK_RESERVED, nil},
				tokenEof,
			},
		},
		{
			"returned;",
			[]*Token{
				{"returned", 8, 0, TK_IDENT, nil},
				{";", 1, 0, TK_RESERVED, nil},
				tokenEof,
			},
		},
		{
			"if a==1 { return a } else { return 0 }",
			[]*Token{
				{"if", 2, 0, TK_RESERVED, nil},
				{"a", 1, 0, TK_IDENT, nil},
				{"==", 2, 0, TK_RESERVED, nil},
				{"1", 1, 1, TK_NUM, nil},
				{"{", 1, 0, TK_RESERVED, nil},
				{"return", 6, 0, TK_RESERVED, nil},
				{"a", 1, 0, TK_IDENT, nil},
				{"}", 1, 0, TK_RESERVED, nil},
				{"else", 4, 0, TK_RESERVED, nil},
				{"{", 1, 0, TK_RESERVED, nil},
				{"return", 6, 0, TK_RESERVED, nil},
				{"0", 1, 0, TK_NUM, nil},
				{"}", 1, 0, TK_RESERVED, nil},
				tokenEof,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.str, func(t *testing.T) {
			if err := tokenize(tC.str); err != nil {
				t.Fatal(err)
			}

			expected := joinTokens(tC.expected)
			if !reflect.DeepEqual(token, expected) {
				t.Fatalf("Tokenizing '%s' failed.\nactual:\n%+v\nexpected:\n%+v\n", tC.str, showTokens(token), showTokens(expected))
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

func joinTokens(tokens []*Token) *Token {
	var head Token
	cur := &head
	for _, t := range tokens {
		cur.next = t
		cur = cur.next
	}
	return head.next
}
