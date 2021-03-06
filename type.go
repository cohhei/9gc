package main

import "fmt"

type TypeKind int

const (
	TY_BYTE TypeKind = iota
	TY_INT
	TY_POINTER
	TY_ARRAY
)

var typeKindString = map[TypeKind]string{
	TY_BYTE:    "byte",
	TY_INT:     "int",
	TY_POINTER: "pointer",
	TY_ARRAY:   "array",
}

var typeNames = map[string]TypeKind{
	"int":  TY_INT,
	"byte": TY_BYTE,
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
	case TY_BYTE:
		return 1
	case TY_INT, TY_POINTER:
		return 8
	case TY_ARRAY:
		return t.Ref.size() * t.ArrayLen
	default:
		panic("unknown type")
	}
}

func (t *Type) String() string {
	return fmt.Sprintf("%+v", *t)
}

func arrayOf(ty *Type, len uint) *Type {
	return &Type{
		Kind:     TY_ARRAY,
		Ref:      ty,
		ArrayLen: len,
	}
}

var intType = &Type{Kind: TY_INT}
var byteType = &Type{Kind: TY_BYTE}

func (t *Type) isInt() bool {
	return t.Kind == TY_INT
}

func (t *Type) isPointer() bool {
	return t.Kind == TY_POINTER
}

func (t *Type) isArray() bool {
	return t.Kind == TY_ARRAY
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
	case ND_DEREF, ND_INDEX:
		ref := n.Lhs.Type.Ref
		if ref == nil {
			errorRef(n)
		}
		n.Type = ref
	}
}

func errorRef(n *Node) {
	switch n.Kind {
	case ND_DEREF:
		panic("invalid pointer dereference")
	case ND_INDEX:
		errorIndexing(n.Lhs.Type.Kind)
	}
}

func errorIndexing(kind TypeKind) {
	panic(fmt.Sprintf("invalid operation (type %s does not support indexing)", kind))
}
