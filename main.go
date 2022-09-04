package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {

	var timeout int64
	var noColor bool

	flag.Int64Var(&timeout, "t", 5, "timeout value in seconds")
	flag.Int64Var(&timeout, "timeout", 5, "timeout value in seconds")
	flag.BoolVar(&noColor, "nc", false, "disable color output")
	flag.BoolVar(&noColor, "nocolor", false, "disable color output")
	flag.Parse()

	dirtyNodes := flag.Args()
	nodes, invalidNodes := cleanNodes(dirtyNodes)
	colorOutput := !noColor // this is required due to how bool flags works

	if len(nodes) == 0 {
		fmt.Println("No valid nodes supplied.")
		os.Exit(0)
	}

	if len(invalidNodes) > 0 {
		printInvalidNodes(invalidNodes)
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	var wg sync.WaitGroup

	printChan := make(chan printDetails, len(nodes))

	go printManager(printChan, colorOutput)

	for {
		// used to force each iteration to wait for the timeout
		wg.Add(1)
		go func() {
			time.Sleep(timeoutDuration + 100*time.Millisecond)
			wg.Done()
		}()

		printChan <- timestamp()

		for _, node := range nodes {
			wg.Add(1)
			go processNode(node, &wg, timeoutDuration, printChan)
		}

		wg.Wait()
		printChan <- printDetails{message: "\n"}
	}
}
