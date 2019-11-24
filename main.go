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

	codegen(code)
}
