package main

import (
	"fmt"
	"os"
	"strconv"
)

// TokenKind is a type for the kind of Token
type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_NUM
	TK_EOF
)

// Token
type Token struct {
	kind TokenKind // The kind of the token
	next *Token    // The next token
	val  int       // The value of TK_NUM
	str  string    // Token string
}

// Current token
var token *Token

// consume returns true and reads the next token if it is the expected value.
// Otherwise consume returns false.
func consume(op byte) bool {
	if token.kind != TK_RESERVED || token.str[0] != op {
		return false
	}
	token = token.next
	return true
}

// expect reads the next token if it is the expected value, otherwise returns the error.
func expect(op byte) error {
	if token.kind != TK_RESERVED || token.str[0] != op {
		return fmt.Errorf("It is not '%s', but '%s'", string(op), string(token.str[0]))
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

// atEof returns true if the token is EOF
func (t *Token) atEof() bool {
	return t.kind == TK_EOF
}

// newToken creates a new token and joins it.
func (t *Token) newToken(kind TokenKind, str string) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
	}
	t.next = tok
	return tok
}

// tokenize tokenizes a string and returns it
func tokenize(str string) (*Token, error) {
	var head Token
	cur := &head

	for len(str) > 0 {
		// Skip the space
		if isSpace(str[0]) {
			str = next(str)
			continue
		}

		if str[0] == '+' || str[0] == '-' {
			cur = cur.newToken(TK_RESERVED, str[:1])
			str = next(str)
			continue
		}

		if isDigit(str[0]) {
			var err error
			cur, err = getDigit(cur, str)
			if err != nil {
				return nil, err
			}
			str = str[len(cur.str):]
			continue
		}

		return nil, fmt.Errorf("Couldn't tokenize.")
	}

	cur.newToken(TK_EOF, str)
	return head.next, nil
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
		t := cur.newToken(TK_NUM, str[:i])
		t.val = dig
		return t, nil
	}

	dig, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	t := cur.newToken(TK_NUM, str)
	t.val = dig
	return t, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "The number of arguments is incorrect.\n")
		return
	}

	var err error
	token, err = tokenize(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	// The first one should be number
	n, err := expectNumber()
	if err != nil {
		panic(err)
	}
	fmt.Printf("  mov rax, %d\n", n)

	for !token.atEof() {
		if consume('+') {
			n, err := expectNumber()
			if err != nil {
				panic(err)
			}
			fmt.Printf("  add rax, %d\n", n)
			continue
		}

		if err := expect('-'); err != nil {
			panic(err)
		}

		n, err := expectNumber()
		if err != nil {
			panic(err)
		}
		fmt.Printf("  sub rax, %d\n", n)
	}

	fmt.Printf("  ret\n")
}
