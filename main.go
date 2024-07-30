package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

func main() {

	// flag parsing
	var timeout int64
	var noColor bool

	flag.Int64Var(&timeout, "t", 5, "timeout value in seconds")
	flag.Int64Var(&timeout, "timeout", 5, "timeout value in seconds")
	flag.BoolVar(&noColor, "nc", false, "disable color output")
	flag.BoolVar(&noColor, "nocolor", false, "disable color output")
	flag.Parse()
	colorOutput := !noColor // required due to how bool flags works

	suppliedHosts := flag.Args()

	// this channel is this single place to display to screen
	// just assume all hosts supplied are valid and use for
	// size of the channel buffer
	printChan := make(chan printDetails, len(suppliedHosts))
	go printManager(printChan, colorOutput)

	// no hosts provided, exit
	if len(suppliedHosts) == 0 {
		printChan <- printDetails{message: "No hosts supplied.", fgColor: color.FgRed}
		time.Sleep(50 * time.Millisecond) // gross hack
		os.Exit(0)
	}

	// split valid and invalid hosts
	hosts, invalidHosts := cleanHosts(suppliedHosts)

	// no valid hosts supplied, exit
	if len(hosts) == 0 {
		printChan <- printDetails{message: "No valid hosts supplied.", fgColor: color.FgRed}
		time.Sleep(50 * time.Millisecond) // gross hack
		os.Exit(0)
	}

	// display invalid supplied hosts
	if len(invalidHosts) > 0 {
		printChan <- printDetails{message: fmtInvalidHosts(invalidHosts), fgColor: color.FgBlue}
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	var wg sync.WaitGroup

	// handle SIGINT, SIGTERM
	keepRunning := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// register cleanup
	go signalHandler(printChan, signals, &keepRunning)

	for keepRunning {

		// used to force each iteration to wait for the timeout
		wg.Add(1)
		go func() {
			// 1ms added to the batch timeout so not the exact value of the ping timeout
			// N.B. I am not sure the timeout is required, need more testing to prove either way
			time.Sleep(timeoutDuration + 1*time.Millisecond)
			wg.Done()
		}()

		printChan <- timestamp()

		for _, host := range hosts {
			wg.Add(1)
			go processHost(host, &wg, timeoutDuration, printChan)
		}

		wg.Wait()

		// blank line between batches
		printChan <- printDetails{message: ""}
	}
}
