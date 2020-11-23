package main

import (
	"flag"
	"fmt"
	"github.com/go-ping/ping"
)


func main() {
	var host = flag.String("h", "", "Host address")
	flag.Parse()

	pinger, err := ping.NewPinger(*host)
	pinger.Count = 10

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}
