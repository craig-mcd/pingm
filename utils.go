package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type printDetails struct {
	message string
	bgColor color.Attribute
	fgColor color.Attribute
}

func printManager(printChan <-chan printDetails) {

	for p := range printChan {

		if p.bgColor != 0 {
			color.Set(p.bgColor)
		}

		if p.fgColor != 0 {
			color.Set(p.fgColor)
		}

		fmt.Println(p.message)

		color.Unset()
	}
}

func timestamp() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func printInvalidNodes(nodes []string) {

	fmt.Print("Invalid nodes: ")

	for _, n := range nodes {
		fmt.Print(n, " ")
	}

	fmt.Printf("\n\n")
}
