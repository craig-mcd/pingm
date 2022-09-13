package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// printDetails is used for passing standard information into the print channel
type printDetails struct {
	message string
	bgColor color.Attribute
	fgColor color.Attribute
}

// printManager used as single place to control and print output
func printManager(printChan <-chan printDetails, colorOutput bool) {

	for p := range printChan {

		if colorOutput {

			if p.bgColor != 0 {
				color.Set(p.bgColor)
			}

			if p.fgColor != 0 {
				color.Set(p.fgColor)
			}
		}

		fmt.Println(p.message)

		color.Unset()
	}
}

// timestamp used to print time info before each new batch runs
func timestamp() printDetails {

	now := time.Now().Format("2006-01-02 15:04:05")

	return printDetails{
		message: now,
		bgColor: color.BgHiYellow,
		fgColor: color.FgBlack,
	}
}

// printInvalidHosts helper function to display invalid supplied hosts
func printInvalidHosts(hosts []string) {

	var sb strings.Builder
	sb.WriteString("Invalid hosts: ")
	sb.WriteString(strings.Join(hosts, " "))
	sb.WriteString("\n")

	fmt.Println(sb.String())
}
