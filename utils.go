package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func printTimestamp() {

	color.Set(color.BgHiYellow)
	color.Set(color.FgBlack)
	defer color.Unset()

	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05"))
}

func printInvalidNodes(nodes []string) {

	fmt.Print("Invalid nodes: ")

	for _, n := range nodes {
		fmt.Print(n, " ")
	}

	fmt.Printf("\n\n")
}
