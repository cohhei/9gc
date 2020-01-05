package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected []*Node
		globals  map[string]*Var
	}{
		{
			desc:  "Function",
			input: "func add(a int,b int) { return a + b } func main() { return add(1,2) }",
			expected: []*Node{
				{
					Kind:         ND_FUNC,
					FunctionName: "add",
					Args: []*Node{
						{Kind: ND_VAR, Type: intType, Var: lvarInt("a")}, {Kind: ND_VAR, Type: intType, Var: lvarInt("b")},
					},
					Locals: &VarList{
						Var: lvarInt("b"),
						Next: &VarList{
							Var: lvarInt("a"),
						},
					},
					Block: &Node{
						Kind: ND_BLOCK, Body: []*Node{
							{
								Kind: ND_RETURN,
								Lhs: &Node{
									Kind: ND_ADD,
									Type: intType,
									Lhs:  &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("a")},
									Rhs:  &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("b")},
								},
							},
						},
					},
				},
				{
					Kind:         ND_FUNC,
					FunctionName: "main",
					Args:         []*Node{},
					Block: &Node{
						Kind: ND_BLOCK, Body: []*Node{
							{
								Kind: ND_RETURN,
								Lhs: &Node{
									Kind: ND_FUNCALL, FunctionName: "add",
									Type: intType,
									Args: []*Node{
										{Kind: ND_NUM, Type: intType, Val: 1}, {Kind: ND_NUM, Type: intType, Val: 2},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc:  "Address and dereference operator",
			input: "func main() {var x int;var y *int;y = &x;*y=3}",
			expected: []*Node{
				{
					Kind:         ND_FUNC,
					FunctionName: "main",
					Args:         []*Node{},
					Locals: &VarList{
						Next: &VarList{Var: lvarInt("x")},
						Var:  lvarPointerInt("y"),
					},
					Block: &Node{
						Kind: ND_BLOCK, Body: []*Node{
							{Kind: ND_VAR, Type: intType, Var: lvarInt("x")},
							{Kind: ND_VAR, Type: &Type{TY_POINTER, intType, 0}, Var: lvarPointerInt("y")},
							{Kind: ND_ASSIGN,
								Lhs: &Node{Kind: ND_VAR, Type: &Type{TY_POINTER, intType, 0}, Var: lvarPointerInt("y")},
								Rhs: &Node{Kind: ND_ADDR, Type: &Type{TY_POINTER, intType, 0}, Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("x")}},
							},
							{Kind: ND_ASSIGN,
								Lhs: &Node{Kind: ND_DEREF, Type: intType, Lhs: &Node{Kind: ND_VAR, Type: &Type{TY_POINTER, intType, 0}, Var: lvarPointerInt("y")}},
								Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 3},
							},
						},
					},
				},
			},
		},
		{
			desc:    "Global variable",
			input:   "var i int",
			globals: map[string]*Var{"i": {Name: "i", Type: intType}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			token, code, locals = nil, nil, nil
			if err := tokenize(tC.input); err != nil {
				t.Fatal(err)
			}
			program()
			actual := code

			if diff := cmp.Diff(actual, tC.expected); diff != "" {
				t.Errorf("Hogefunc differs: (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(globals, tC.globals); tC.globals != nil && diff != "" {
				t.Errorf("Hogefunc differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestStmt(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected []*Node
	}{
		{
			desc:  "LocalVariable",
			input: "a:=18;triple:=3;return a*triple;",
			expected: []*Node{
				&Node{
					Kind: ND_ASSIGN,
					Lhs: &Node{
						Kind: ND_VAR,
						Var:  lvarInt("a"),
						Type: intType,
					},
					Rhs: &Node{
						Kind: ND_NUM,
						Val:  18,
						Type: intType,
					},
				},

				&Node{
					Kind: ND_ASSIGN,
					Lhs: &Node{
						Kind: ND_VAR,
						Var:  lvarInt("triple"),
						Type: intType,
					},
					Rhs: &Node{
						Kind: ND_NUM,
						Val:  3,
						Type: intType,
					},
				},

				&Node{
					Kind: ND_RETURN,
					Lhs: &Node{
						Kind: ND_MUL,
						Type: intType,
						Lhs: &Node{
							Kind: ND_VAR,
							Var:  lvarInt("a"),
							Type: intType,
						},
						Rhs: &Node{
							Kind: ND_VAR,
							Var:  lvarInt("triple"),
							Type: intType,
						},
					},
				},
			},
		},
		{
			desc:  "IfStatements",
			input: "if a := 0; a==1 { return a } else if a == 2 { return -1 }; return 100",
			expected: []*Node{
				{
					Kind: ND_IF,
					Init: &Node{Kind: ND_ASSIGN, Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("a")}, Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 0}},
					Cond: &Node{
						Kind: ND_EQ,
						Lhs:  &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("a")},
						Rhs:  &Node{Kind: ND_NUM, Type: intType, Val: 1},
						Type: intType,
					},
					Then: &Node{
						Kind: ND_BLOCK,
						Body: []*Node{{
							Kind: ND_RETURN,
							Lhs:  &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("a")},
						}},
					},
					Els: &Node{
						Kind: ND_IF,
						Cond: &Node{
							Kind: ND_EQ,
							Type: intType,
							Lhs:  &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("a")},
							Rhs:  &Node{Kind: ND_NUM, Type: intType, Val: 2},
						},
						Then: &Node{
							Kind: ND_BLOCK,
							Body: []*Node{{
								Kind: ND_RETURN,
								Lhs:  &Node{Kind: ND_SUB, Type: intType, Lhs: &Node{Kind: ND_NUM, Type: intType, Val: 0}, Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 1}},
							}},
						},
					},
				},

				&Node{
					Kind: ND_RETURN,
					Lhs:  &Node{Kind: ND_NUM, Type: intType, Val: 100},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for i := 1; i < 10; i++ { 1 }",
			expected: []*Node{
				{
					Kind: ND_FOR,
					Init: &Node{Kind: ND_ASSIGN, Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("i")}, Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 1}},
					Cond: &Node{
						Kind: ND_LT, Type: intType,
						Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("i")},
						Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 10},
					},
					Inc:  &Node{Kind: ND_INC, Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("i")}},
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_NUM, Type: intType, Val: 1}}},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "var i int;for i < 10 { 1 }",
			expected: []*Node{
				{Kind: ND_VAR, Type: intType, Var: lvarInt("i")},
				{
					Kind: ND_FOR,
					Cond: &Node{
						Kind: ND_LT, Type: intType,
						Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("i")},
						Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 10},
					},
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_NUM, Type: intType, Val: 1}}},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "var i int;for { i-- }",
			expected: []*Node{
				{Kind: ND_VAR, Type: intType, Var: lvarInt("i")},
				{
					Kind: ND_FOR,
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_DEC, Lhs: &Node{Kind: ND_VAR, Type: intType, Var: lvarInt("i")}}}},
				},
			},
		},
		{
			desc:  "Array",
			input: "var i [10][2]int",
			expected: []*Node{
				{Kind: ND_VAR, Type: arrayOf(arrayOf(intType, 2), 10), Var: lvarPointerPoinsterInt("i")},
			},
		},
		{
			desc:  "Index",
			input: "var i [10][2]int;i[1][1]",
			expected: []*Node{
				{Kind: ND_VAR, Type: arrayOf(arrayOf(intType, 2), 10), Var: lvarPointerPoinsterInt("i")},
				{Kind: ND_INDEX,
					Type: intType,
					Lhs: &Node{
						Kind: ND_INDEX,
						Type: arrayOf(intType, 2),
						Lhs:  &Node{Kind: ND_VAR, Type: arrayOf(arrayOf(intType, 2), 10), Var: lvarPointerPoinsterInt("i")},
						Rhs:  &Node{Kind: ND_NUM, Type: intType, Val: 1},
					},
					Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 1},
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			token, code, locals = nil, nil, nil
			globals = make(map[string]*Var)
			if err := tokenize(tC.input); err != nil {
				t.Fatal(err)
			}
			var actual []*Node
			for !token.atEof() {
				node := stmt()
				node.addType()
				actual = append(actual, node)
			}

			if diff := cmp.Diff(actual, tC.expected); diff != "" {
				t.Errorf("Hogefunc differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func lvarInt(s string) *Var {
	return &Var{Name: s, Type: intType, IsLocal: true}
}

func lvarPointerInt(s string) *Var {
	return &Var{Name: s, Type: &Type{TY_POINTER, intType, 0}, IsLocal: true}
}

func lvarPointerPoinsterInt(s string) *Var {
	return &Var{Name: s, Type: arrayOf(arrayOf(intType, 2), 10), IsLocal: true}
}
