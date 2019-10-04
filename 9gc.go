package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	len  int       // Token length
}

// Current token
var token *Token

// consume returns true and reads the next token if it is the expected value.
// Otherwise consume returns false.
func consume(op string) bool {
	if token.kind != TK_RESERVED || len(op) != token.len || !strings.HasPrefix(op, token.str) {
		return false
	}
	token = token.next
	return true
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

		// Multi-letter punctuator
		if startswitch(str, "==") || startswitch(str, "!=") ||
			startswitch(str, "<=") || startswitch(str, ">=") {
			cur = cur.newToken(TK_RESERVED, str[:2], 2)
			str = str[len(cur.str):]
			continue
		}

		if strings.Contains("+-*/()<>", str[0:1]) {
			cur = cur.newToken(TK_RESERVED, str[:1], 1)
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

	cur.newToken(TK_EOF, str, len(str))
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
	node := expr()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	gen(node)

	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}

// NodeKind is a type for the kind of Node
type NodeKind int

const (
	ND_ADD = iota // +
	ND_SUB        // -
	ND_MUL        // *
	ND_DIV        // /
	ND_EQ         // ==
	ND_NE         // !=
	ND_LT         // <
	ND_LE         // <=
	ND_NUM        // number
)

// Node is a type for the abstract syntax tree
type Node struct {
	kind NodeKind // The kind of the node
	lhs  *Node    // left-hand side
	rhs  *Node    // right-hand side
	val  int      // The value of ND_NUM
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	node := &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
	return node
}

func newNodeNum(val int) *Node {
	node := &Node{
		kind: ND_NUM,
		val:  val,
	}
	return node
}

func expr() *Node {
	return equality()
}

func equality() *Node {
	node := relational()

	for {
		if consume("==") {
			node = newNode(ND_EQ, node, relational())
		} else if consume("!=") {
			node = newNode(ND_NE, node, relational())
		} else {
			return node
		}
	}
}

func relational() *Node {
	node := add()

	for {
		if consume("<") {
			node = newNode(ND_LT, node, add())
		} else if consume("<=") {
			node = newNode(ND_LE, node, add())
		} else if consume(">") {
			node = newNode(ND_LT, add(), node)
		} else if consume(">=") {
			node = newNode(ND_LE, add(), node)
		} else {
			return node
		}
	}
}

func add() *Node {
	node := mul()

	for {
		if consume("+") {
			node = newNode(ND_ADD, node, mul())
		} else if consume("-") {
			node = newNode(ND_SUB, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := unary()

	for {
		if consume("*") {
			node = newNode(ND_MUL, node, unary())
		} else if consume("/") {
			node = newNode(ND_DIV, node, unary())
		} else {
			return node
		}
	}
}

func unary() *Node {
	if consume("+") {
		return unary()
	} else if consume("-") {
		return newNode(ND_SUB, newNodeNum(0), unary())
	}
	return primary()
}

func primary() *Node {
	// If the next token is '(', it shouled be '(' expr ')'
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	// If not so, it should be a number
	val, err := expectNumber()
	if err != nil {
		panic(err)
	}
	return newNodeNum(val)
}

func gen(node *Node) {
	if node.kind == ND_NUM {
		fmt.Printf("  push %d\n", node.val)
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.kind {
	case ND_ADD:
		fmt.Printf("  add rax, rdi\n")
	case ND_SUB:
		fmt.Printf("  sub rax, rdi\n")
	case ND_MUL:
		fmt.Printf("  imul rax, rdi\n")
	case ND_DIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
	case ND_EQ:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_NE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_LT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_LE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzb rax, al\n")
	}

	fmt.Printf("  push rax\n")
}
