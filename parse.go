package main

// NodeKind is a type for the kind of Node
type NodeKind int

const (
	ND_ADD    = iota // +
	ND_SUB           // -
	ND_MUL           // *
	ND_DIV           // /
	ND_ASSIGN        // =
	ND_LVAR          // local variable
	ND_EQ            // ==
	ND_NE            // !=
	ND_LT            // <
	ND_LE            // <=
	ND_INC           // ++
	ND_DEC           // --
	ND_NUM           // number
	ND_RETURN        // return
	ND_IF            // if
	ND_FOR           // for
)

// Node is a type for the abstract syntax tree
type Node struct {
	kind   NodeKind // The kind of the node
	lhs    *Node    // left-hand side
	rhs    *Node    // right-hand side
	val    int      // The value of ND_NUM
	offset int      // The velue of ND_LVAR

	// "if" and "for"
	cond *Node
	then *Node
	els  *Node
	init *Node
	inc  *Node
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
			kind: ND_RETURN,
			lhs:  expr(),
		}
	} else if consume("if") {
		node = ifstmt()
	} else if consume("for") {
		node = &Node{kind: ND_FOR}
		if consume("{") { // for {}
			node.then = stmt()
			expect("}")
		} else {
			unknown := expr()
			if consume(";") { // for i=0;i<N;i++ {}
				node.init = unknown
				node.cond = expr()
				expect(";")
				node.inc = expr()
			} else { // for i<N {}
				node.cond = unknown
			}
			expect("{")
			node.then = stmt()
			expect("}")
		}

	} else {
		node = expr()
	}

	consume(";")
	return node
}

func ifstmt() *Node {
	node := &Node{kind: ND_IF}
	unknown := expr()
	if consume(";") { // if i:=0; i<N {}
		node.init = unknown
		node.cond = expr()
	} else { // if i<N {}
		node.cond = unknown
	}
	expect("{")
	node.then = stmt()
	expect("}")
	if consume("else") {
		if consume("if") {
			node.els = ifstmt()
		} else if consume("{") {
			node.els = stmt()
			expect("}")
		}
	}
	return node
}

func program() {
	for !token.atEof() {
		code = append(code, stmt())
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
	}
	return postfix()
}

func postfix() *Node {
	node := primary()
	if consume("++") {
		node = &Node{
			kind: ND_INC,
			lhs:  node,
		}
	} else if consume("--") {
		node = &Node{
			kind: ND_DEC,
			lhs:  node,
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
		var node Node
		node.kind = ND_LVAR

		lvar := tok.findLVar()
		if lvar != nil {
			node.offset = lvar.offset
		} else {
			lvar := &LVar{
				next:   locals,
				name:   tok.str,
				len:    tok.len,
				offset: getOffset(),
			}
			node.offset = lvar.offset
			locals = lvar
		}
		return &node
	}

	// If not so, it should be a number
	val, err := expectNumber()
	if err != nil {
		panic(err)
	}
	return newNodeNum(val)
}

func getOffset() int {
	if locals == nil {
		return 8
	}
	return locals.offset + 8
}
