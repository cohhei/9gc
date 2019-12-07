package main

type TypeKind int

const (
	TY_INT TypeKind = iota
	TY_POINTER
	TY_ARRAY
)

var typeKindString = map[TypeKind]string{
	TY_INT:     "int",
	TY_POINTER: "pointer",
	TY_ARRAY:   "array",
}

func (tk TypeKind) String() string {
	return typeKindString[tk]
}

type Type struct {
	Kind     TypeKind
	Ref      *Type // The pointer's reference
	ArrayLen uint
}

func (t *Type) size() uint {
	switch t.Kind {
	case TY_INT, TY_POINTER:
		return 8
	case TY_ARRAY:
		return t.Ref.size() * t.ArrayLen
	default:
		panic("unknown type")
	}
}

func arrayOf(ty *Type, len uint) *Type {
	return &Type{
		Kind:     TY_ARRAY,
		Ref:      ty,
		ArrayLen: len,
	}
}

var intType = &Type{Kind: TY_INT}

func (t *Type) isInt() bool {
	return t.Kind == TY_INT
}

func (t *Type) isPointer() bool {
	return t.Kind == TY_POINTER
}

func (n *Node) addType() {
	if n == nil || n.Type != nil {
		return
	}

	n.Lhs.addType()
	n.Rhs.addType()
	n.Cond.addType()
	n.Then.addType()
	n.Els.addType()
	n.Init.addType()
	n.Inc.addType()
	n.Block.addType()
	for _, stmt := range n.Body {
		stmt.addType()
	}
	for _, a := range n.Args {
		a.addType()
	}

	switch n.Kind {
	case ND_ADD, ND_SUB, ND_MUL, ND_DIV, ND_EQ, ND_NE, ND_LT, ND_LE, ND_FUNCALL, ND_NUM:
		n.Type = intType
	case ND_ADDR:
		n.Type = &Type{TY_POINTER, n.Lhs.Type, 0}
	case ND_DEREF:
		ref := n.Lhs.Type.Ref
		if ref == nil {
			panic("nvalid pointer dereference")
		}
		n.Type = ref
	}
}
