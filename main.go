package main

import (
	"flag"
	"fmt"
)

func main() {

	timeout := flag.Int("t", 5_000, "timeout value in millis")
	flag.Parse()

	dirtyNodes := flag.Args()

	fmt.Println("Nodes supplied:", dirtyNodes)
	fmt.Printf("Timeout %dms\n", *timeout)
}
