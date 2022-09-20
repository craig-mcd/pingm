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

	dirtyHosts := flag.Args()
	hosts, invalidHosts := cleanHosts(dirtyHosts)
	colorOutput := !noColor // this is required due to how bool flags works

	// No valid hosts supplied, exit
	if len(hosts) == 0 {
		fmt.Println("No valid hosts supplied.")
		os.Exit(0)
	}

	// Display invalid suppled hosts
	if len(invalidHosts) > 0 {
		printInvalidHosts(invalidHosts)
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	var wg sync.WaitGroup

	// This channel is this single place to print output
	printChan := make(chan printDetails, len(hosts))
	go printManager(printChan, colorOutput)

	// Handle SIGINT, SIGTERM
	keepRunning := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// used to catch the signal, mark the sentinel value as done for the main loop
	go func() {
		<-signals
		printChan <- printDetails{message: "\rFinishing batch (ctrl-c to kill)", fgColor: color.FgRed}
		keepRunning = false
		<-signals
		os.Exit(1)
	}()

	for keepRunning {

		// used to force each iteration to wait for the timeout
		wg.Add(1)
		go func() {
			// 100ms added to the batch timeout so not the exact value of the ping timeout
			time.Sleep(timeoutDuration + 100*time.Millisecond)
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
