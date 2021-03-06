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
	ND_VAR                     // variable
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
	ND_INDEX                   // x[y]
)

var nodeKindName = map[NodeKind]string{
	ND_ADD:     "ND_ADD",
	ND_SUB:     "ND_SUB",
	ND_MUL:     "ND_MUL",
	ND_DIV:     "ND_DIV",
	ND_ASSIGN:  "ND_ASSIGN",
	ND_VAR:     "ND_LVAR",
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
	ND_INDEX:   "ND_INDEX",
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
	Locals       *VarList
	Block        *Node

	// var
	Var *Var
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

type Labeler struct {
	Counter int
}

func (l *Labeler) New() string {
	label := fmt.Sprintf(".L.data.%d", l.Counter)
	l.Counter++
	return label
}

var labeler = &Labeler{}

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
		tok := expectIdent()
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
	globals = make(map[string]*Var)
	for !token.atEof() {
		switch token.str {
		case "func":
			function()
		case "var":
			gvar()
		default:
			panic(fmt.Sprintf("expected declaration, found %s", token.str))
		}
		consume(";")
	}
}

func function() {
	expect("func")
	locals = nil
	tok := expectIdent()
	node := &Node{
		Kind:         ND_FUNC,
		FunctionName: tok.str,
		Args:         definedArgs(),
	}
	if !consume("{") {
		node.Type = parseType()
		expect("{")
	}
	node.Block = block()
	node.Locals = locals
	node.addType()
	code = append(code, node)
}

func gvar() {
	expect("var")
	tok := expectIdent()
	if globals[tok.str] != nil {
		panic(fmt.Sprintf("%s redeclared in this block", tok.str))
	}
	gvar := newGVar(tok.str, parseType())
	globals[gvar.Name] = gvar
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
		return newNode(ND_INC, node, nil)
	} else if consume("--") {
		return newNode(ND_DEC, node, nil)
	} else if peek("[") {
		return index(node)
	}
	return node
}

func index(base *Node) *Node {
	if !consume("[") {
		return base
	}
	i := primary()
	expect("]")
	node := newNode(ND_INDEX, base, i)
	return index(node)
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
				Kind: ND_VAR,
				Var:  lvar,
				Type: lvar.Type,
			}
			return node
		} else if gvar := globals[tok.str]; gvar != nil {
			node := &Node{
				Kind: ND_VAR,
				Var:  gvar,
				Type: gvar.Type,
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

	if token.kind == TK_STR {
		ty := arrayOf(byteType, uint(token.len))
		v := newGVar(labeler.New(), ty)
		v.Content = token.str
		v.Len = token.len
		token = token.next
		return newVarNode(v)
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

func newVarNode(v *Var) *Node {
	node := &Node{
		Kind: ND_VAR,
		Var:  v,
		Type: v.Type,
	}
	return node
}

func newGVarNode(name string, ty *Type) *Node {
	v := newGVar(name, ty)
	return newVarNode(v)
}

func newGVar(name string, ty *Type) *Var {
	gvar := &Var{
		Name: name,
		Type: ty,
	}
	globals[name] = gvar
	return gvar
}

func newLVarNode(name string, ty *Type) *Node {
	v := newLVar(name, ty)
	return newVarNode(v)
}

func newLVar(name string, ty *Type) *Var {
	lvar := &Var{
		Name:    name,
		Type:    ty,
		IsLocal: true,
	}
	locals = &VarList{locals, lvar}
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
	kind := expectType()
	return &Type{kind, nil, 0}
}
