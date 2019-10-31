package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
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
					kind: ND_ASSIGN,
					lhs: &Node{
						kind:   ND_LVAR,
						offset: 8,
					},
					rhs: &Node{
						kind: ND_NUM,
						val:  18,
					},
				},

				&Node{
					kind: ND_ASSIGN,
					lhs: &Node{
						kind:   ND_LVAR,
						offset: 16,
					},
					rhs: &Node{
						kind: ND_NUM,
						val:  3,
					},
				},

				&Node{
					kind: ND_RETURN,
					lhs: &Node{
						kind: ND_MUL,
						lhs: &Node{
							kind:   ND_LVAR,
							offset: 8,
						},
						rhs: &Node{
							kind:   ND_LVAR,
							offset: 16,
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
					kind: ND_IF,
					init: &Node{kind: ND_ASSIGN, lhs: &Node{kind: ND_LVAR, offset: 8}, rhs: &Node{kind: ND_NUM, val: 0}},
					cond: &Node{
						kind: ND_EQ,
						lhs:  &Node{kind: ND_LVAR, offset: 8},
						rhs:  &Node{kind: ND_NUM, val: 1},
					},
					then: &Node{
						kind: ND_RETURN,
						lhs:  &Node{kind: ND_LVAR, offset: 8},
					},
					els: &Node{
						kind: ND_IF,
						cond: &Node{
							kind: ND_EQ,
							lhs:  &Node{kind: ND_LVAR, offset: 8},
							rhs:  &Node{kind: ND_NUM, val: 2},
						},
						then: &Node{
							kind: ND_RETURN,
							lhs:  &Node{kind: ND_SUB, lhs: &Node{kind: ND_NUM, val: 0}, rhs: &Node{kind: ND_NUM, val: 1}},
						},
					},
				},

				&Node{
					kind: ND_RETURN,
					lhs:  &Node{kind: ND_NUM, val: 100},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for i = 1; i < 10; i++ { 1 }",
			expected: []*Node{
				{
					kind: ND_FOR,
					init: &Node{kind: ND_ASSIGN, lhs: &Node{kind: ND_LVAR, offset: 8}, rhs: &Node{kind: ND_NUM, val: 1}},
					cond: &Node{kind: ND_LT, lhs: &Node{kind: ND_LVAR, offset: 8}, rhs: &Node{kind: ND_NUM, val: 10}},
					inc:  &Node{kind: ND_INC, lhs: &Node{kind: ND_LVAR, offset: 8}},
					then: &Node{kind: ND_NUM, val: 1},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for i < 10 { 1 }",
			expected: []*Node{
				{
					kind: ND_FOR,
					cond: &Node{kind: ND_LT, lhs: &Node{kind: ND_LVAR, offset: 8}, rhs: &Node{kind: ND_NUM, val: 10}},
					then: &Node{kind: ND_NUM, val: 1},
				},
			},
		},
		{
			desc:  "ForStatement",
			input: "for { i-- }",
			expected: []*Node{
				{
					kind: ND_FOR,
					then: &Node{kind: ND_DEC, lhs: &Node{kind: ND_LVAR, offset: 8}},
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
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Fatalf("Parse '%s' failed.\nactual:\n%+v\nexpected:\n%+v\n", tC.input, showNodes(actual), showNodes(tC.expected))
			}
		})
	}
}

func showNodes(nodes []*Node) string {
	strs := make([]string, len(nodes))
	for i, n := range nodes {
		strs[i] = showHands(n, "\t")
	}
	return strings.Join(strs, "\n")
}

func showHands(node *Node, tabs string) string {
	str := fmt.Sprintf("%s%+v", tabs, node)
	if node.lhs != nil {
		str = fmt.Sprintf("%s\n%s", str, showHands(node.lhs, tabs+"\t"))
	}
	if node.rhs != nil {
		str = fmt.Sprintf("%s\n%s", str, showHands(node.rhs, tabs+"\t"))
	}
	if node.cond != nil {
		str = fmt.Sprintf("%s\n%s", str, showHands(node.cond, tabs+"\t"))
	}
	if node.then != nil {
		str = fmt.Sprintf("%s\n%s", str, showHands(node.then, tabs+"\t"))
	}
	return str
}
