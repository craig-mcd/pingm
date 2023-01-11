package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
	probing "github.com/prometheus-community/pro-bing"
)

// cleanHosts filter out valid and invalid hosts returning both
func cleanHosts(dirtyHosts []string) ([]string, []string) {

	hosts := []string{}
	invalidHosts := []string{}

	for _, host := range dirtyHosts {

		_, err := net.LookupIP(host)

		if err == nil {
			hosts = append(hosts, host)
		} else {
			invalidHosts = append(invalidHosts, host)
		}

	}

	return hosts, invalidHosts
}

// processHost check if a host is reachable and send results to the print channel
func processHost(host string, wg *sync.WaitGroup, timeout time.Duration, printChan chan<- printDetails) {

	var message printDetails

	// handle any panic, send any message to the 'print channel', update the wait group
	defer func() {

		if r := recover(); r != nil {

			message = printDetails{
				message: fmt.Sprintf("panic on %s goroutine: %s", host, r),
				fgColor: color.FgRed,
			}
		}

		printChan <- message

		wg.Done()
	}()

	pinger, err := probing.NewPinger(host)

	// return early if failed to create Pinger
	if err != nil {

		message = printDetails{
			message: fmt.Sprintf("%-30s %s", host, err),
			fgColor: color.FgRed,
		}

		return
	}

	pinger.Count = 1           // we only want a single response
	pinger.SetPrivileged(true) // send ICMP and not UDP
	pinger.Timeout = timeout
	err = pinger.Run()

	if err == nil {

		stats := pinger.Statistics()

		if stats.PacketsRecv > 0 {
			message = printDetails{
				message: fmt.Sprintf("%-30s %dms", host, stats.AvgRtt/time.Millisecond),
				fgColor: color.FgGreen,
			}
		} else {
			message = printDetails{
				message: fmt.Sprintf("%-30s timed out", host),
				fgColor: color.FgCyan,
			}
		}

	} else {
		message = printDetails{
			message: fmt.Sprintf("%-30s %s", host, err.Error()),
			fgColor: color.FgRed,
		}
	}
}
