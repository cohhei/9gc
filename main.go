package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "The number of arguments is incorrect.\n")
		return
	}

	if err := tokenize(os.Args[1]); err != nil {
		panic(err)
	}
	program()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	fmt.Printf("  push rbp\n")
	fmt.Printf("  mov rbp, rsp\n")
	fmt.Printf("  sub rsp, 208\n")

	for _, c := range code {
		gen(c)
		fmt.Printf("  pop rax\n")
	}

	fmt.Printf("  mov rsp, rbp\n")
	fmt.Printf("  pop rbp\n")
	fmt.Printf("  ret\n")
}
