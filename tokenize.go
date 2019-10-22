package main

import (
	"fmt"
	"strconv"
	"strings"
)

// TokenKind is a type for the kind of Token
type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_IDENT
	TK_NUM
	TK_EOF
)

// Token
type Token struct {
	kind TokenKind // The kind of the token
	next *Token    // The next token
	val  int       // The value of TK_NUM
	str  string    // Token string
	len  int       // Token length
}

// Current token
var token *Token

// atEof returns true if the token is EOF
func (t *Token) atEof() bool {
	return t.kind == TK_EOF
}

func startswitch(s1, s2 string) bool {
	return strings.HasPrefix(s1, s2)
}

// newToken creates a new token and joins it.
func (t *Token) newToken(kind TokenKind, str string, len int) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
		len:  len,
	}
	t.next = tok
	return tok
}

// consume returns true and reads the next token if it is the expected value.
// Otherwise consume returns false.
func consume(op string) bool {
	if token.kind != TK_RESERVED || len(op) != token.len || !strings.HasPrefix(op, token.str) {
		return false
	}
	token = token.next
	return true
}

// Consumes the current token if it is an identifier.
func consumeIdent() *Token {
	if token.kind != TK_IDENT {
		return nil
	}
	t := token
	token = token.next
	return t
}

// expect reads the next token if it is the expected value, otherwise returns the error.
func expect(op string) error {
	if token.kind != TK_RESERVED || len(op) != token.len || !strings.HasPrefix(op, token.str) {
		return fmt.Errorf("It is not '%s', but '%s'", op, token.str)
	}
	token = token.next
	return nil
}

// expectNumber returns the value and read the next token, otherwise returns the error.
func expectNumber() (int, error) {
	if token.kind != TK_NUM {
		return 0, fmt.Errorf("'%s' is not a number.", string(token.str[0]))
	}
	val := token.val
	token = token.next
	return val, nil
}

// tokenize tokenizes a string and returns it
func tokenize(str string) error {
	var head Token
	cur := &head

	for len(str) > 0 {
		// Skip the space
		if isSpace(str[0]) {
			str = next(str)
			continue
		}

		// Multi-letter punctuator
		if startswitch(str, "==") || startswitch(str, "!=") ||
			startswitch(str, "<=") || startswitch(str, ">=") {
			cur = cur.newToken(TK_RESERVED, str[:2], 2)
			str = str[len(cur.str):]
			continue
		}

		if strings.Contains("+-*/()<>;=", str[0:1]) {
			cur = cur.newToken(TK_RESERVED, str[:1], 1)
			str = next(str)
			continue
		}

		if isDigit(str[0]) {
			var err error
			cur, err = getDigit(cur, str)
			if err != nil {
				return err
			}
			str = str[len(cur.str):]
			continue
		}

		if 'a' <= str[0] && str[0] <= 'z' {
			cur = cur.newToken(TK_IDENT, str[:1], 1)
			str = next(str)
			continue
		}

		return fmt.Errorf("Couldn't tokenize. '%s'", str[:1])
	}

	cur.newToken(TK_EOF, str, len(str))
	token = head.next
	return nil
}

func next(str string) string {
	return str[1:]
}

func isSpace(s byte) bool {
	return s == '\t' || s == '\n' || s == '\v' || s == '\f' || s == '\r' || s == ' '
}

func isDigit(s byte) bool {
	return '0' <= s && s <= '9'
}

func getDigit(cur *Token, str string) (*Token, error) {
	for i := 0; i < len(str); i++ {
		if isDigit(str[i]) {
			continue
		}
		dig, err := strconv.Atoi(str[:i])
		if err != nil {
			return nil, err
		}
		t := cur.newToken(TK_NUM, str[:i], i)
		t.val = dig
		return t, nil
	}

	dig, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	t := cur.newToken(TK_NUM, str, len(str))
	t.val = dig
	return t, nil
}
