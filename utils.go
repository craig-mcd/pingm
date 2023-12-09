package main

import (
	"fmt"
	"os"
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

			color.Set(p.bgColor)
			color.Set(p.fgColor)

			// Don't print newline due to how colour output works
			fmt.Print(p.message)
			color.Unset()
			fmt.Println()

		} else {
			fmt.Println(p.message)
		}
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
func fmtInvalidHosts(hosts []string) string {

	var sb strings.Builder
	sb.WriteString("Invalid hosts: ")
	sb.WriteString(strings.Join(hosts, " "))
	sb.WriteString("\n")

	return sb.String()
}

// signalHandler catch the signal, mark the sentinel value as done for the main loop
func signalHandler(printChan chan<- printDetails, signals <-chan os.Signal, keepRunning *bool) {

	// catch once, finish batch then quit
	<-signals
	printChan <- printDetails{message: "\rFinishing batch (ctrl-c to kill)", fgColor: color.FgBlue}
	*keepRunning = false
	// catch again, immediate quit
	<-signals
	os.Exit(1)
}
