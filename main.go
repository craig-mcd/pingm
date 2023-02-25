package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

func main() {

	// Flag parsing
	var timeout int64
	var noColor bool

	flag.Int64Var(&timeout, "t", 5, "timeout value in seconds")
	flag.Int64Var(&timeout, "timeout", 5, "timeout value in seconds")
	flag.BoolVar(&noColor, "nc", false, "disable color output")
	flag.BoolVar(&noColor, "nocolor", false, "disable color output")
	flag.Parse()
	colorOutput := !noColor // this is required due to how bool flags works

	dirtyHosts := flag.Args()

	// no hosts provided, exit
	if len(dirtyHosts) == 0 {

		// print channel is not setup yet
		if colorOutput {
			color.Set(color.FgRed)
		}

		fmt.Println("No hosts supplied.")
		os.Exit(0)
	}

	// split valid and invalid hosts
	hosts, invalidHosts := cleanHosts(dirtyHosts)

	// No valid hosts supplied, exit
	if len(hosts) == 0 {

		// print channel is not setup yet
		if colorOutput {
			color.Set(color.FgRed)
		}

		fmt.Println("No valid hosts supplied.")
		os.Exit(0)
	}

	// This channel is this single place to print output
	printChan := make(chan printDetails, len(hosts))
	go printManager(printChan, colorOutput)

	// Display invalid suppled hosts
	if len(invalidHosts) > 0 {
		printChan <- printDetails{message: fmtInvalidHosts(invalidHosts), fgColor: color.FgBlue}
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	var wg sync.WaitGroup

	// Handle SIGINT, SIGTERM
	keepRunning := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// register cleanup
	go signalHandler(printChan, signals, &keepRunning)

	for keepRunning {

		// used to force each iteration to wait for the timeout
		wg.Add(1)
		go func() {
			// 10ms added to the batch timeout so not the exact value of the ping timeout
			time.Sleep(timeoutDuration + 10*time.Millisecond)
			wg.Done()
		}()

		printChan <- timestamp()

		for _, host := range hosts {
			wg.Add(1)
			go processHost(host, &wg, timeoutDuration, printChan)
		}

		wg.Wait()
		printChan <- printDetails{message: "\n"}
	}
}
