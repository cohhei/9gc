package main

// NodeKind is a type for the kind of Node
type NodeKind int

const (
	ND_ADD     = iota // +
	ND_SUB            // -
	ND_MUL            // *
	ND_DIV            // /
	ND_ASSIGN         // =
	ND_LVAR           // local variable
	ND_EQ             // ==
	ND_NE             // !=
	ND_LT             // <
	ND_LE             // <=
	ND_INC            // ++
	ND_DEC            // --
	ND_NUM            // number
	ND_RETURN         // return
	ND_IF             // if
	ND_FOR            // for
	ND_BLOCK          // { ... }
	ND_FUNCALL        // Function call
	ND_FUNC           // Function
	ND_ADDR           // &
	ND_DEREF          // *
)

// Node is a type for the abstract syntax tree
type Node struct {
	Kind NodeKind // The kind of the node
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
			Lhs:  expr(),
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
			expect("{")
			// locals = nil
			node.Block = block()
			node.Locals = locals
			code = append(code, node)
			continue
		}

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
		var node Node

		// Function call
		if consume("(") {
			node = Node{
				Kind:         ND_FUNCALL,
				FunctionName: tok.str,
				Args:         args(),
			}
			return &node
		}

		// Variables
		node.Kind = ND_LVAR

		lvar := tok.findLVar()
		if lvar != nil {
			node.LVar = lvar
		} else {
			lvar := &LVar{
				Name: tok.str,
				Len:  tok.len,
			}
			locals = &LVarList{
				Next: locals,
				LVar: lvar,
			}
			node.LVar = lvar
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
		args = append(args, primary())
		consume(",")
	}
	return args
}
