package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	timeout := flag.Int("t", 5_000, "timeout value in millis")
	flag.Parse()

	dirtyNodes := flag.Args()
	nodes, invalid := cleanNodes(dirtyNodes)

	if len(nodes) == 0 {
		fmt.Println("No valid nodes supplied.")
		os.Exit(0)
	}

	if len(invalid) > 0 {
		fmt.Println("Invalid nodes:", invalid)
	}

	fmt.Println("Nodes supplied:", nodes)
	fmt.Printf("Timeout %dms\n", *timeout)
}
