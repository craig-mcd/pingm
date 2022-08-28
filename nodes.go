package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
)

func cleanNodes(dirtyNodes []string) ([]string, []string) {

	nodes := []string{}
	invalidNodes := []string{}

	for _, node := range dirtyNodes {

		_, err := net.LookupIP(node)

		if err == nil {
			nodes = append(nodes, node)
		} else {
			invalidNodes = append(invalidNodes, node)
		}

	}

	return nodes, invalidNodes
}

func processNode(node string, wg *sync.WaitGroup, timeout time.Duration, printChan chan<- printDetails) {

	defer wg.Done()

	pinger, err := ping.NewPinger(node)
	var message printDetails

	// return early if failed to create Pinger
	if err != nil {

		message = printDetails{
			message: fmt.Sprintf("%-30s %s", node, err),
			fgColor: color.FgRed,
		}

		printChan <- message
		return
	}

	pinger.Count = 1
	pinger.SetPrivileged(true)
	pinger.Timeout = timeout
	err = pinger.Run()

	if err == nil {

		stats := pinger.Statistics()

		if stats.PacketsRecv > 0 {
			message = printDetails{
				message: fmt.Sprintf("%-30s %dms", node, stats.AvgRtt/time.Millisecond),
				fgColor: color.FgGreen,
			}
		} else {
			message = printDetails{
				message: fmt.Sprintf("%-30s timed out", node),
				fgColor: color.FgCyan,
			}
		}

	} else {
		message = printDetails{
			message: fmt.Sprintf("%-30s %s", node, err.Error()),
			fgColor: color.FgRed,
		}
	}

	printChan <- message
}
