package main

import (
	"fmt"
	"github.com/go-ping/ping"
)

func main() {
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {

		panic(err)

	}
	pinger.Count = 3

	pinger.Run() // blocks until finished

	stats := pinger.Statistics() // get send/receive/rtt stats
	fmt.Print(stats)
}
