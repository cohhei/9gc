package main

import "fmt"

var label int

func seq() int {
	s := label
	label++
	return s
}
func genLval(node *Node) {
	if node.kind != ND_LVAR {
		panic("Not valiable")
	}

	fmt.Printf("  mov rax, rbp\n")
	fmt.Printf("  sub rax, %d\n", node.offset)
	fmt.Printf("  push rax\n")
}
func gen(node *Node) {
	switch node.kind {
	case ND_NUM:
		fmt.Printf("  push %d\n", node.val)
		return
	case ND_LVAR:
		genLval(node)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov rax, [rax]\n")
		fmt.Printf("  push rax\n")
		return
	case ND_ASSIGN:
		genLval(node.lhs)
		gen(node.rhs)

		fmt.Printf("  pop rdi\n")
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov [rax], rdi\n")
		fmt.Printf("  push rdi\n")
		return
	case ND_RETURN:
		gen(node.lhs)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov rsp, rbp\n")
		fmt.Printf("  pop rbp\n")
		fmt.Printf("  ret\n")
		return
	case ND_ADD, ND_SUB, ND_MUL, ND_DIV, ND_EQ, ND_NE, ND_LT, ND_LE:
		genBinary(node)
	case ND_IF:
		gen(node.cond)
		s := seq()
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		if node.els != nil {
			fmt.Printf("  je .L.else.%d\n", s)
			gen(node.then)
			fmt.Printf(".L.else.%d:\n", s)
			gen(node.els)
		} else {
			fmt.Printf("  je .L.end.%d\n", s)
			gen(node.then)
		}
		if node.els != nil {
		}
		fmt.Printf(".L.end.%d:\n", s)
	}
}

func genBinary(node *Node) {
	gen(node.lhs)
	gen(node.rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.kind {
	case ND_ADD:
		fmt.Printf("  add rax, rdi\n")
	case ND_SUB:
		fmt.Printf("  sub rax, rdi\n")
	case ND_MUL:
		fmt.Printf("  imul rax, rdi\n")
	case ND_DIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
	case ND_EQ:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_NE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_LT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzb rax, al\n")
	case ND_LE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzb rax, al\n")
	}

	fmt.Printf("  push rax\n")
}
