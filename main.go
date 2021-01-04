package main

import (
	"fmt"

	"github.com/liyiysng/scatter/node"
)

func main() {
	fmt.Printf("hello world%s", "!!!")
	n := &node.Node{Name: "foo node"}
	fmt.Print(n)
}
