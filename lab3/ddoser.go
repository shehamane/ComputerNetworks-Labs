package main

import (
	"flag"
	"fmt"
	"github.com/go-ping/ping"
	"time"
)

func main() {
	var host = flag.String("h", "", "Host address")
	var num = flag.Int("n", 10, "Number of pingers")
	var async = flag.Int("async", 1, "Async?")
	flag.Parse()

	pingers := make([]*ping.Pinger, *num)

	t := time.Now()

	for i := 0; i < *num; i++ {
		pingers[i], _ = ping.NewPinger(*host)
		pingers[i].SetPrivileged(true)
		pingers[i].Count = 1
		//pingers[i].Count = 2
		pingers[i].OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
			time.Sleep(time.Millisecond * 1000)
		}
		pingers[i].OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("Прошло %v\n",
				time.Since(t))
		}
		fmt.Printf("PING %s (%s):\n", pingers[i].Addr(), pingers[i].IPAddr())
	}
	for i := 0; i < 10; i++ {
		if *async == 1 {
			go pingers[i].Run()
		} else {
			pingers[i].Run()
		}
	}
	var input string
	fmt.Scanln(&input)
}
