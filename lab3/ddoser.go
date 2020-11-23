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
	flag.Parse()

	pingers := make([]*ping.Pinger, *num)

	for i:=0; i<*num; i++{
		pingers[i], _ = ping.NewPinger(*host)
		pingers[i].SetPrivileged(true)
		//pingers[i].Count = 2
		pingers[i].OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
			time.Sleep(time.Millisecond * 1000)
		}
		fmt.Printf("PING %s (%s):\n", pingers[i].Addr(), pingers[i].IPAddr())
	}
	for i:=0; i<10; i++ {
		go pingers[i].Run()
	}
		var input string
	fmt.Scanln(&input)
}
