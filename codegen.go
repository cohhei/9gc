package main

import (
	"fmt"
)

var argreg = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var funcname string
var label int

func seq() int {
	s := label
	label++
	return s
}
func genAddr(node *Node) {
	switch node.Kind {
	case ND_LVAR:
		fmt.Printf("  lea rax, [rbp-%d]\n", node.LVar.Offset)
		fmt.Printf("  push rax\n")
	case ND_DEREF:
		gen(node.Lhs)
	default:
		panic("Not valiable")
	}
}

func load() {
	fmt.Printf("  pop rax\n")
	fmt.Printf("  mov rax, [rax]\n")
	fmt.Printf("  push rax\n")
}

func codegen(code []*Node) {
	fmt.Printf(".intel_syntax noprefix\n")
	var offset int
	for _, n := range code {
		switch n.Kind {
		case ND_FUNC:
			fmt.Printf(".global %s\n", n.FunctionName)
			fmt.Printf("%s:\n", n.FunctionName)
			funcname = n.FunctionName

			for _, a := range n.Args {
				offset += 8
				a.LVar.Offset = offset
			}
			for l := n.Locals; l != nil; l = l.Next {
				offset += 8
				l.LVar.Offset = offset
			}
			fmt.Printf("  push rbp\n")
			fmt.Printf("  mov rbp, rsp\n")
			fmt.Printf("  sub rsp, %d\n", offset)
			for i, a := range n.Args {
				fmt.Printf("  mov [rbp-%d], %s\n", a.LVar.Offset, argreg[i])
			}

			gen(n.Block)

			fmt.Printf(".L.return.%s:\n", funcname)
			fmt.Printf("  mov rsp, rbp\n")
			fmt.Printf("  pop rbp\n")
			fmt.Printf("  ret\n")
		default:
			panic("expected declaration")
		}
	}
}

func gen(node *Node) {
	if node == nil {
		return
	}
	switch node.Kind {
	case ND_NUM:
		fmt.Printf("  push %d\n", node.Val)
	case ND_LVAR:
		genAddr(node)
		load()
	case ND_ASSIGN:
		genAddr(node.Lhs)
		gen(node.Rhs)

		fmt.Printf("  pop rdi\n")
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov [rax], rdi\n")
		fmt.Printf("  push rdi\n")
	case ND_RETURN:
		if node.Lhs != nil {
			gen(node.Lhs)
			fmt.Printf("  pop rax\n")
		}
		fmt.Printf("  jmp .L.return.%s\n", funcname)
	case ND_ADD, ND_SUB, ND_MUL, ND_DIV, ND_EQ, ND_NE, ND_LT, ND_LE:
		genBinary(node)
	case ND_INC:
		genAddr(node.Lhs)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov rdi, [rax]\n")
		fmt.Printf("  add rdi, 1\n")
		fmt.Printf("  mov [rax], rdi\n")
	case ND_DEC:
		genAddr((node.Lhs))
		fmt.Printf("  pop rax\n")
		fmt.Printf("  mov rdi, [rax]\n")
		fmt.Printf("  sub rdi, 1\n")
		fmt.Printf("  mov [rax], rdi\n")
	case ND_IF:
		gen(node.Init)
		gen(node.Cond)
		s := seq()
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		if node.Els != nil {
			fmt.Printf("  je .L.else.%d\n", s)
			gen(node.Then)
			fmt.Printf(".L.else.%d:\n", s)
			gen(node.Els)
		} else {
			fmt.Printf("  je .L.end.%d\n", s)
			gen(node.Then)
		}
		fmt.Printf(".L.end.%d:\n", s)
	case ND_FOR:
		gen(node.Init)
		s := seq()
		fmt.Printf(".L.begin.%d:\n", s)
		gen(node.Cond)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  cmp rax, 0\n")
		fmt.Printf("  je .L.end.%d\n", s)
		gen(node.Then)
		gen(node.Inc)
		fmt.Printf("  jmp .L.begin.%d\n", s)
		fmt.Printf(".L.end.%d:\n", s)
	case ND_BLOCK:
		for _, n := range node.Body {
			gen(n)
		}
	case ND_FUNCALL:
		var nargs int
		for _, a := range node.Args {
			gen(a)
			nargs++
		}
		for i := nargs - 1; i >= 0; i-- {
			fmt.Printf("  pop %s\n", argreg[i])
		}
		// We need to align RSP to a 16 byte boundary before
		// calling a function because it is an ABI requirement.
		// RAX is set to 0 for variadic function.
		s := seq()
		fmt.Printf("  mov rax, rsp\n")
		fmt.Printf("  and rax, 15\n")
		fmt.Printf("  jnz .L.call.%d\n", s)
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  call %s\n", node.FunctionName)
		fmt.Printf("  jmp .L.end.%d\n", s)
		fmt.Printf(".L.call.%d:\n", s)
		fmt.Printf("  sub rsp, 8\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  call %s\n", node.FunctionName)
		fmt.Printf("  add rsp, 8\n")
		fmt.Printf(".L.end.%d:\n", s)
		fmt.Printf("  push rax\n")
	case ND_ADDR:
		genAddr(node.Lhs)
	case ND_DEREF:
		gen(node.Lhs)
		load()
	}
}

func genBinary(node *Node) {
	gen(node.Lhs)
	gen(node.Rhs)

	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")

	switch node.Kind {
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
