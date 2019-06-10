package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ixy-languages/ixy.go/driver"
)

const batchSize = 256


func forward(rxDev, txDev driver.IxyInterface, rxQueue, txQueue uint16, bufs []*driver.PktBuf) {
	numRx := rxDev.RxBatch(rxQueue, bufs)
	if numRx > 0 {
		//touch all packets, otherwise it's a completely unrealistic workload if the packet just stays in L3
		for i := uint32(0); i < numRx; i++ {
			bufs[i].Pkt[64]++
		}
		numTx := txDev.TxBatch(txQueue, bufs[:numRx])
		//there are two ways to handle the case that packets are not being sent out:
		//either wait on tx or drop them; in this case it's better to drop them, otherwise we accumulate latency
		for i := numTx; i < numRx; i++ {
			driver.PktBufFree(bufs[i])
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("%v forwards packets between two ports.\n", os.Args[0])
		fmt.Printf("Usage: %v <pci bus id2> <pci bus id1>\n", os.Args[0])
		return
	}
	
	dev1 := driver.IxyInit(os.Args[1], 1, 1)
	dev2 := driver.IxyInit(os.Args[2], 1, 1)
	if dev1 == nil || dev2 == nil {
		fmt.Print("trying to start driver with unsupported device, exiting now\n")
		return
	}

	lastStatsPrinted := time.Now() //includes monotonic clock reading
	stats1 := driver.DeviceStats{}
	stats1old := driver.DeviceStats{}
	stats2 := driver.DeviceStats{}
	stats2old := driver.DeviceStats{}
	stats1.StatsInit(dev1)
	stats1old.StatsInit(dev1)
	stats2.StatsInit(dev2)
	stats2old.StatsInit(dev2)

	bufs := make([]*driver.PktBuf, batchSize)
	counter := uint64(0)
	for {
		forward(dev1, dev2, 0, 0, bufs)
		forward(dev2, dev1, 0, 0, bufs)

		//don't poll the time unnecessarily
		counter++
		if counter&0xfff == 0 {
			t := time.Now()
			if t.Sub(lastStatsPrinted) > time.Second {
				//every second
				dev1.ReadStats(&stats1)
				stats1.PrintStatsDiff(&stats1old, t.Sub(lastStatsPrinted))
				stats1old = stats1
				if dev1 != dev2 {
					dev2.ReadStats(&stats2)
					stats2.PrintStatsDiff(&stats2old, t.Sub(lastStatsPrinted))
					stats2old = stats2
				}
				lastStatsPrinted = t
			}
		}
	}
}
