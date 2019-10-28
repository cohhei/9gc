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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
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
	return str
}
