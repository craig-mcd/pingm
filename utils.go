package main

import (
	"fmt"
	"strings"
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

func timestamp() printDetails {

	now := time.Now().Format("2006-01-02 15:04:05")

	return printDetails{
		message: now,
		bgColor: color.BgHiYellow,
		fgColor: color.FgBlack,
	}
}

func printInvalidNodes(nodes []string) {

	var sb strings.Builder
	sb.WriteString("Invalid nodes: ")
	sb.WriteString(strings.Join(nodes, " "))
	sb.WriteString("\n")

	fmt.Println(sb.String())
}
