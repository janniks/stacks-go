package main

import (
	"fmt"

	stacks "github.com/janniks/stacks-go"
)

func main() {
	fmt.Println("Library version:", stacks.Version())
	fmt.Println(stacks.Hello())
}
