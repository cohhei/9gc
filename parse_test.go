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
	}{
		{
			desc:  "Function",
			input: "func add(a int,b int) { return a + b } func main() { return add(1,2) }",
			expected: []*Node{
				{
					Kind:         ND_FUNC,
					FunctionName: "add",
					Args: []*Node{
						{Kind: ND_LVAR, Type: intType, LVar: lvarInt("a")}, {Kind: ND_LVAR, Type: intType, LVar: lvarInt("b")},
					},
					Locals: &LVarList{
						LVar: lvarInt("b"),
						Next: &LVarList{
							LVar: lvarInt("a"),
						},
					},
					Block: &Node{
						Kind: ND_BLOCK, Body: []*Node{
							{
								Kind: ND_RETURN,
								Lhs: &Node{
									Kind: ND_ADD,
									Type: intType,
									Lhs:  &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("a")},
									Rhs:  &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("b")},
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
					Locals: &LVarList{
            		Next: &LVarList{LVar: &LVar{Name: "x", Type: intType}},
								LVar: lvarPointerInt("y"),
					},
					Block: &Node{
						Kind: ND_BLOCK, Body: []*Node{
							{Kind: ND_LVAR, Type: intType, LVar: lvarInt("x")},
							{Kind: ND_LVAR, Type: &Type{TY_POINTER, intType}, LVar: lvarPointerInt("y")},
							{Kind: ND_ASSIGN,
								Lhs: &Node{Kind: ND_LVAR, Type: &Type{TY_POINTER, intType}, LVar: lvarPointerInt("y")},
								Rhs: &Node{Kind: ND_ADDR, Type: &Type{TY_POINTER, intType}, Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("x")}},
							},
							{Kind: ND_ASSIGN,
								Lhs: &Node{Kind: ND_DEREF, Type: intType, Lhs: &Node{Kind: ND_LVAR, Type: &Type{TY_POINTER, intType}, LVar: lvarPointerInt("y")}},
								Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 3},
							},
						},
					},
				},
			},
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
						Kind: ND_LVAR,
						LVar: lvarInt("a"),
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
						Kind: ND_LVAR,
						LVar: lvarInt("triple"),
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
							Kind: ND_LVAR,
							LVar: lvarInt("a"),
							Type: intType,
						},
						Rhs: &Node{
							Kind: ND_LVAR,
							LVar: lvarInt("triple"),
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
					Init: &Node{Kind: ND_ASSIGN, Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("a")}, Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 0}},
					Cond: &Node{
						Kind: ND_EQ,
						Lhs:  &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("a")},
						Rhs:  &Node{Kind: ND_NUM, Type: intType, Val: 1},
						Type: intType,
					},
					Then: &Node{
						Kind: ND_BLOCK,
						Body: []*Node{{
							Kind: ND_RETURN,
							Lhs:  &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("a")},
						}},
					},
					Els: &Node{
						Kind: ND_IF,
						Cond: &Node{
							Kind: ND_EQ,
							Type: intType,
							Lhs:  &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("a")},
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
					Init: &Node{Kind: ND_ASSIGN, Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")}, Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 1}},
					Cond: &Node{
						Kind: ND_LT, Type: intType,
						Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")},
						Rhs: &Node{Kind: ND_NUM, Type: intType, Val: 10},
					},
					Inc:  &Node{Kind: ND_INC, Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")}},
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_NUM, Type: intType, Val: 1}}},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "var i int;for i < 10 { 1 }",
			expected: []*Node{
				{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")},
				{
					Kind: ND_FOR,
					Cond: &Node{
						Kind: ND_LT, Type: intType,
						Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")},
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
				{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")},
				{
					Kind: ND_FOR,
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_DEC, Lhs: &Node{Kind: ND_LVAR, Type: intType, LVar: lvarInt("i")}}}},
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			token, code, locals = nil, nil, nil
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

func lvarInt(s string) *LVar {
	return &LVar{s, 0, intType}
}

func lvarPointerInt(s string) *LVar {
	return &LVar{s, 0, &Type{TY_POINTER, intType}}
}
