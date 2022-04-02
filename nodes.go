package main

import (
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

func processNode(node string, wg *sync.WaitGroup, timeout time.Duration) {

	defer wg.Done()

	pinger, err := ping.NewPinger(node)

	if err != nil {
		panic(err)
	}

	pinger.Count = 1
	pinger.SetPrivileged(true)
	pinger.Timeout = timeout
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
}
