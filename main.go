package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
)

func main() {

	timeout := flag.Int64("t", 5_000, "timeout value in millis")
	flag.Parse()

	dirtyNodes := flag.Args()
	nodes, invalidNodes := cleanNodes(dirtyNodes)

	if len(nodes) == 0 {
		fmt.Println("No valid nodes supplied.")
		os.Exit(0)
	}

	if len(invalidNodes) > 0 {
		printInvalidNodes(invalidNodes)
	}

	timeoutDuration := time.Duration(*timeout) * time.Millisecond
	var wg sync.WaitGroup

	for {
		// used to force each iteration to wait for the timeout
		// TODO check if this is idiomatic Go
		wg.Add(1)
		go func() {
			time.Sleep(timeoutDuration)
			wg.Done()
		}()

		diplayTimestamp()

		for _, node := range nodes {

			wg.Add(1)

			go func(node string) {

				defer wg.Done()

				pinger, err := ping.NewPinger(node)

				if err != nil {
					panic(err)
				}

				pinger.Count = 1
				pinger.SetPrivileged(true)
				pinger.Timeout = timeoutDuration
				err = pinger.Run()

				if err == nil {

					stats := pinger.Statistics()

					if stats.PacketsRecv > 0 {
						color.Green("%-30s %dms\n", node, *&stats.AvgRtt/time.Millisecond)
					} else {
						color.Cyan("%-30s timed out\n", node)
					}
				} else {
					color.Red("%-30s %s\n", node, err.Error())
				}

			}(node)
		}

		wg.Wait()
		fmt.Println()
	}
}
