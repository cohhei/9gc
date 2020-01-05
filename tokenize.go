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
	TK_STR
	TK_NUM
	TK_EOF
)

// Token
type Token struct {
	str  string    // Token string
	len  int       // Token length
	val  int       // The value of TK_NUM
	kind TokenKind // The kind of the token
	next *Token    // The next token
}

// Current token
var token *Token

type Var struct {
	Name    string
	Type    *Type
	IsLocal bool

	// Local variables
	Offset uint // Offset from RBP

	// Global variables
	Content string
	Len     int
}

type VarList struct {
	Next *VarList
	Var  *Var
}

var locals *VarList         // Local variables
var globals map[string]*Var // Global variables

func (t *Token) findLVar() *Var {
	for v := locals; v != nil; v = v.Next {
		if t.str == v.Var.Name {
			return v.Var
		}
	}
	return nil
}

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
	if !token.isReserved() || len(op) != token.len || op != token.str {
		return false
	}
	token = token.next
	return true
}

func peek(op string) bool {
	return token.isReserved() && len(op) == token.len && op == token.str
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
func expect(op string) {
	if !token.isReserved() || len(op) != token.len || op != token.str {
		panic(fmt.Errorf("expected '%s', found '%s'", op, token.str))
	}
	token = token.next
}

// expectNumber returns the value and read the next token, otherwise panic.
func expectNumber() int {
	if token.kind != TK_NUM {
		panic(fmt.Sprintf("'%s' is not a number.", string(token.str)))
	}
	val := token.val
	token = token.next
	return val
}

func expectIdent() *Token {
	tok := consumeIdent()
	if tok == nil {
		panic(fmt.Sprintf("expected 'IDENT', found %s", token.str))
	}
	return tok
}

func expectType() TypeKind {
	kind, isType := typeNames[token.str]
	if !isType {
		panic(fmt.Sprintf("expected 'TYPE', found %s", token.str))
	}
	token = token.next
	return kind
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

		// String literals
		if str[0] == '"' {
			var content string
			for i, s := range str[1:] {
				if s == '"' {
					content = str[1 : i+1]
					break
				}
			}
			if len(content) == 0 {
				panic("string literal not terminated")
			}
			cur = cur.newToken(TK_STR, content, len(content))
			str = str[len(content)+2:]
			continue
		}

		// Multi-letter punctuator
		if startswitch(str, "==") || startswitch(str, "!=") ||
			startswitch(str, "<=") || startswitch(str, ">=") ||
			startswitch(str, "++") || startswitch(str, "--") ||
			startswitch(str, ":=") {
			cur = cur.newToken(TK_RESERVED, str[:2], 2)
			str = str[len(cur.str):]
			continue
		}

		if strings.Contains("+-*/()<>;={},&[]", str[0:1]) {
			cur = cur.newToken(TK_RESERVED, str[:1], 1)
			str = next(str)
			continue
		}

		if isDigit(str[0]) {
			var err error
			cur, err = cur.readDigit(str)
			if err != nil {
				return err
			}
			str = str[len(cur.str):]
			continue
		}

		if k := startWithReserved(str); k != "" {
			len := len(k)
			cur = cur.newToken(TK_RESERVED, str[:len], len)
			str = str[len:]
			continue
		}

		if isIdent(str[0]) {
			cur = cur.readIdent(str)
			str = str[len(cur.str):]
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

var keywords = []string{
	"return", "if", "else", "for", "func", "var", "int", "byte",
}

func startWithReserved(str string) string {
	for _, k := range keywords {
		l := len(k)
		if startswitch(str, k) && (len(str) == l || !isDigit(str[l]) && !isIdent(str[l])) {
			return k
		}
	}
	return ""
}

func isDigit(s byte) bool {
	return '0' <= s && s <= '9'
}

func (t *Token) readDigit(str string) (*Token, error) {
	for i := 0; i < len(str); i++ {
		if isDigit(str[i]) {
			continue
		}
		dig, err := strconv.Atoi(str[:i])
		if err != nil {
			return nil, err
		}
		tok := t.newToken(TK_NUM, str[:i], i)
		tok.val = dig
		return tok, nil
	}

	dig, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	tok := t.newToken(TK_NUM, str, len(str))
	tok.val = dig
	return tok, nil
}

func isIdent(s byte) bool {
	return ('a' <= s && s <= 'z') || ('A' <= s && s <= 'Z') || s == '_'
}

func (t *Token) readIdent(str string) *Token {
	for i := 0; i < len(str); i++ {
		if isIdent(str[i]) {
			continue
		}
		return t.newToken(TK_IDENT, str[:i], i)
	}

	return t.newToken(TK_IDENT, str, len(str))
}

func (t *Token) isReserved() bool {
	return t.kind == TK_RESERVED
}
