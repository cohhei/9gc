package main

import "fmt"

// NodeKind is a type for the kind of Node
type NodeKind int

const (
	ND_ADD     NodeKind = iota // +
	ND_SUB                     // -
	ND_MUL                     // *
	ND_DIV                     // /
	ND_ASSIGN                  // =
	ND_LVAR                    // local variable
	ND_EQ                      // ==
	ND_NE                      // !=
	ND_LT                      // <
	ND_LE                      // <=
	ND_INC                     // ++
	ND_DEC                     // --
	ND_NUM                     // number
	ND_RETURN                  // return
	ND_IF                      // if
	ND_FOR                     // for
	ND_BLOCK                   // { ... }
	ND_FUNCALL                 // Function call
	ND_FUNC                    // Function
	ND_ADDR                    // &
	ND_DEREF                   // *
)

var nodeKindName = map[NodeKind]string{
	ND_ADD:     "ND_ADD",
	ND_SUB:     "ND_SUB",
	ND_MUL:     "ND_MUL",
	ND_DIV:     "ND_DIV",
	ND_ASSIGN:  "ND_ASSIGN",
	ND_LVAR:    "ND_LVAR",
	ND_EQ:      "ND_EQ",
	ND_NE:      "ND_NE",
	ND_LT:      "ND_LT",
	ND_LE:      "ND_LE",
	ND_INC:     "ND_INC",
	ND_DEC:     "ND_DEC",
	ND_NUM:     "ND_NUM",
	ND_RETURN:  "ND_RETURN",
	ND_IF:      "ND_IF",
	ND_FOR:     "ND_FOR",
	ND_BLOCK:   "ND_BLOCK",
	ND_FUNCALL: "ND_FUNCALL",
	ND_FUNC:    "ND_FUNC",
	ND_ADDR:    "ND_ADDR",
	ND_DEREF:   "ND_DEREF",
}

func (nk NodeKind) String() string {
	return nodeKindName[nk]
}

// Node is a type for the abstract syntax tree
type Node struct {
	Kind NodeKind // The kind of the node
	Type *Type    // Type, e.g. int or pointer to int
	Lhs  *Node    // left-hand side
	Rhs  *Node    // right-hand side
	Val  int      // The value of ND_NUM

	// "if" and "for"
	Cond *Node
	Then *Node
	Els  *Node
	Init *Node
	Inc  *Node

	// block
	Body []*Node

	// function
	FunctionName string
	Args         []*Node
	Locals       *LVarList
	Block        *Node

	// var
	LVar *LVar
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	node := &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
	return node
}

func newNodeNum(val int) *Node {
	node := &Node{
		Kind: ND_NUM,
		Val:  val,
	}
	return node
}

var code []*Node

func assign() *Node {
	node := equality()
	if consume("=") || consume(":=") {
		node = newNode(ND_ASSIGN, node, assign())
	}
	return node
}

func expr() *Node {
	return assign()
}

func stmt() *Node {
	var node *Node

	if consume("return") {
		node = &Node{
			Kind: ND_RETURN,
			Lhs:  equality(),
		}
	} else if consume("if") {
		node = ifstmt()
	} else if consume("for") {
		node = &Node{Kind: ND_FOR}
		if consume("{") { // for {}
			node.Then = block()
		} else {
			unknown := expr()
			if consume(";") { // for i=0;i<N;i++ {}
				node.Init = unknown
				node.Cond = expr()
				expect(";")
				node.Inc = expr()
			} else { // for i<N {}
				node.Cond = unknown
			}
			expect("{")
			node.Then = block()
		}
	} else if consume("{") {
		node = block()
	} else if consume("var") {
		tok := consumeIdent()
		if tok == nil {
			panic("expected 'IDENT'")
		}

		lvar := tok.findLVar()
		if lvar != nil {
			panic(fmt.Sprintf("%s redeclared in this block", tok.str))
		}

		node = newLVarNode(tok.str, parseType())
	} else {
		node = expr()
	}

	consume(";")
	return node
}

func ifstmt() *Node {
	node := &Node{Kind: ND_IF}
	unknown := expr()
	if consume(";") { // if i:=0; i<N {}
		node.Init = unknown
		node.Cond = expr()
	} else { // if i<N {}
		node.Cond = unknown
	}
	expect("{")
	node.Then = block()
	if consume("else") {
		if consume("if") {
			node.Els = ifstmt()
		} else {
			expect("{")
			node.Els = block()
		}
	}
	return node
}

func block() *Node {
	node := &Node{Kind: ND_BLOCK}
	for !consume("}") {
		node.Body = append(node.Body, stmt())
	}
	return node
}

func program() {
	for !token.atEof() {
		if consume("func") {
			locals = nil
			tok := consumeIdent()
			if tok == nil {
				panic("expected 'IDENT'")
			}
			node := &Node{
				Kind:         ND_FUNC,
				FunctionName: tok.str,
				Args:         definedArgs(),
			}
			if !consume("{") {
				node.Type = parseType()
			}
			node.Block = block()
			node.Locals = locals
			node.addType()
			code = append(code, node)
			continue
		}

		panic(fmt.Sprintf("expected declaration, found %s", token.str))
	}
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
	} else if consume("&") {
		return newNode(ND_ADDR, unary(), nil)
	} else if consume("*") {
		return newNode(ND_DEREF, unary(), nil)
	}
	return postfix()
}

func postfix() *Node {
	node := primary()
	if consume("++") {
		node = &Node{
			Kind: ND_INC,
			Lhs:  node,
		}
	} else if consume("--") {
		node = &Node{
			Kind: ND_DEC,
			Lhs:  node,
		}
	}
	return node
}

func primary() *Node {
	// If the next token is '(', it shouled be '(' expr ')'
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	if tok := consumeIdent(); tok != nil {
		// Function call
		if consume("(") {
			node := Node{
				Kind:         ND_FUNCALL,
				FunctionName: tok.str,
				Args:         args(),
			}
			return &node
		}

		// Variables
		lvar := tok.findLVar()
		if lvar != nil {
			node := &Node{
				Kind: ND_LVAR,
				LVar: lvar,
				Type: lvar.Type,
			}
			return node
		} else {
			if !consume(":=") {
				panic(fmt.Sprintf("undeclared name: %s %s %s", tok.str, tok.next.str, tok.next.next.str))
			}
			rhs := equality()
			rhs.addType()
			node := newNode(ND_ASSIGN, newLVarNode(tok.str, rhs.Type), rhs)
			return node
		}
	}

	// If not so, it should be a number
	return newNodeNum(expectNumber())
}

func args() []*Node {
	args := []*Node{}
	for !consume(")") {
		args = append(args, assign())
		if consume(",") {
			continue
		}
	}
	return args
}

func definedArgs() []*Node {
	expect("(")
	args := []*Node{}
	for !consume(")") {
		if tok := consumeIdent(); tok != nil {
			args = append(args, newLVarNode(tok.str, parseType()))
		}
		consume(",")
	}
	return args
}

func newLVarNode(name string, ty *Type) *Node {
	node := &Node{
		Kind: ND_LVAR,
		LVar: newLVar(name, ty),
		Type: ty,
	}
	return node
}

func newLVar(name string, ty *Type) *LVar {
	lvar := &LVar{
		Name: name,
		Type: ty,
	}
	locals = &LVarList{locals, lvar}
	return lvar
}

func array() *Type {
	expect("[")
	l := expectNumber()
	expect("]")
	return arrayOf(parseType(), uint(l))
}

func parseType() *Type {
	if peek("[") {
		return array()
	}
	if consume("*") {
		ty := parseType()
		return &Type{TY_POINTER, ty, 0}
	}
	expect("int")
	return &Type{TY_INT, nil, 0}
}
