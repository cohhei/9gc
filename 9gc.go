package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "The number of arguments is incorrect.\n")
		return
	}
	
	a, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, "The number of arguments is incorrect.\n")
		return
	}

	fmt.Printf(".intel_syntax noprefix\n");
  fmt.Printf(".global main\n");
  fmt.Printf("main:\n");
  fmt.Printf("  mov rax, %d\n", a);
  fmt.Printf("  ret\n");
}