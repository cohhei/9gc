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
			desc:  "LocalVariable",
			input: "a=18;triple=3;return a*triple;",
			expected: []*Node{
				&Node{
					Kind: ND_ASSIGN,
					Lhs: &Node{
						Kind: ND_LVAR,
						LVar: &LVar{"a", 1, 0},
					},
					Rhs: &Node{
						Kind: ND_NUM,
						Val:  18,
					},
				},

				&Node{
					Kind: ND_ASSIGN,
					Lhs: &Node{
						Kind: ND_LVAR,
						LVar: &LVar{"triple", 6, 0},
					},
					Rhs: &Node{
						Kind: ND_NUM,
						Val:  3,
					},
				},

				&Node{
					Kind: ND_RETURN,
					Lhs: &Node{
						Kind: ND_MUL,
						Lhs: &Node{
							Kind: ND_LVAR,
							LVar: &LVar{"a", 1, 0},
						},
						Rhs: &Node{
							Kind: ND_LVAR,
							LVar: &LVar{"triple", 6, 0},
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
					Init: &Node{Kind: ND_ASSIGN, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"a", 1, 0}}, Rhs: &Node{Kind: ND_NUM, Val: 0}},
					Cond: &Node{
						Kind: ND_EQ,
						Lhs:  &Node{Kind: ND_LVAR, LVar: &LVar{"a", 1, 0}},
						Rhs:  &Node{Kind: ND_NUM, Val: 1},
					},
					Then: &Node{
						Kind: ND_BLOCK,
						Body: []*Node{{
							Kind: ND_RETURN,
							Lhs:  &Node{Kind: ND_LVAR, LVar: &LVar{"a", 1, 0}},
						}},
					},
					Els: &Node{
						Kind: ND_IF,
						Cond: &Node{
							Kind: ND_EQ,
							Lhs:  &Node{Kind: ND_LVAR, LVar: &LVar{"a", 1, 0}},
							Rhs:  &Node{Kind: ND_NUM, Val: 2},
						},
						Then: &Node{
							Kind: ND_BLOCK,
							Body: []*Node{{
								Kind: ND_RETURN,
								Lhs:  &Node{Kind: ND_SUB, Lhs: &Node{Kind: ND_NUM, Val: 0}, Rhs: &Node{Kind: ND_NUM, Val: 1}},
							}},
						},
					},
				},

				&Node{
					Kind: ND_RETURN,
					Lhs:  &Node{Kind: ND_NUM, Val: 100},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for i = 1; i < 10; i++ { 1 }",
			expected: []*Node{
				{
					Kind: ND_FOR,
					Init: &Node{Kind: ND_ASSIGN, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"i", 1, 0}}, Rhs: &Node{Kind: ND_NUM, Val: 1}},
					Cond: &Node{Kind: ND_LT, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"i", 1, 0}}, Rhs: &Node{Kind: ND_NUM, Val: 10}},
					Inc:  &Node{Kind: ND_INC, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"i", 1, 0}}},
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_NUM, Val: 1}}},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for i < 10 { 1 }",
			expected: []*Node{
				{
					Kind: ND_FOR,
					Cond: &Node{Kind: ND_LT, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"i", 1, 0}}, Rhs: &Node{Kind: ND_NUM, Val: 10}},
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_NUM, Val: 1}}},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for { i-- }",
			expected: []*Node{
				{
					Kind: ND_FOR,
					Then: &Node{Kind: ND_BLOCK, Body: []*Node{{Kind: ND_DEC, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"i", 1, 0}}}}},
				},
			},
		},
		{
			desc:  "Function",
			input: "func add(a,b) { return a + b } func main() { return add(1,2) }",
			expected: []*Node{
				{
					Kind:         ND_FUNC,
					FunctionName: "add",
					Args: []*Node{
						{Kind: ND_LVAR, LVar: &LVar{"a", 1, 0}}, {Kind: ND_LVAR, LVar: &LVar{"b", 1, 0}},
					},
					Locals: &LVarList{
						LVar: &LVar{"b", 1, 0},
						Next: &LVarList{
							LVar: &LVar{"a", 1, 0},
						},
					},
					Block: &Node{
						Kind: ND_BLOCK, Body: []*Node{
							{
								Kind: ND_RETURN,
								Lhs: &Node{
									Kind: ND_ADD,
									Lhs:  &Node{Kind: ND_LVAR, LVar: &LVar{"a", 1, 0}},
									Rhs:  &Node{Kind: ND_LVAR, LVar: &LVar{"b", 1, 0}},
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
									Args: []*Node{
										{Kind: ND_NUM, Val: 1}, {Kind: ND_NUM, Val: 2},
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
			input: "&x;*x",
			expected: []*Node{
				{Kind: ND_ADDR, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"x", 1, 0}}},
				{Kind: ND_DEREF, Lhs: &Node{Kind: ND_LVAR, LVar: &LVar{"x", 1, 0}}},
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
