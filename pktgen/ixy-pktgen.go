package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"ixy.go/driver"
)

const (
	batchSize = 64
	pktSize   = 60
)

var isBig = true //Network Byte Order is always Big Endian
var check = true

//this is the packet we want to send
var pktData = []byte{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, // dst MAC
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, // src MAC
	0x08, 0x00, // ether type: IPv4
	0x45, 0x00, // Version, IHL, TOS
	(pktSize - 14) >> 8,   // ip len excluding ethernet, high byte
	(pktSize - 14) & 0xFF, // ip len exlucding ethernet, low byte
	0x00, 0x00, 0x00, 0x00, // id, flags, fragmentation
	0x40, 0x11, 0x00, 0x00, // TTL (64), protocol (UDP), checksum
	0x0A, 0x00, 0x00, 0x01, // src ip (10.0.0.1)
	0x0A, 0x00, 0x00, 0x02, // dst ip (10.0.0.2)
	0x00, 0x2A, 0x05, 0x39, // src and dst ports (42 -> 1337)
	(pktSize - 20 - 14) >> 8,   // udp len excluding ip & ethernet, high byte
	(pktSize - 20 - 14) & 0xFF, // udp len exlucding ip & ethernet, low byte
	0x00, 0x00, // udp checksum, optional
	'i', 'x', 'y', // payload
}

//calculate an IP/TCP/UDP checksum
func calcIPChecksum(data []byte) uint16 {
	if len(data)%1 != 0 {
		log.Fatal("odd-sized checksums NYI")
	}
	cs := uint32(0)
	for i := 0; i < len(data)/2; i += 2 {
		cs += uint32(binary.BigEndian.Uint16(data[i : i+2]))
		if cs > 0xFFFF {
			cs = (cs & 0xFFFF) + 1 //16 bit one's complement
		}
	}
	return ^uint16(cs)
}

func initMempool() *driver.Mempool {
	numBufs := 2048
	mempool := driver.MemoryAllocateMempool(uint32(numBufs), 0)
	//pre-fill all our packet buffers with some templates that can be modified later
	//we have to do it like this because sending is async in the hardware; we cannot re-use a buffer immediately
	bufs := make([]*driver.PktBuf, numBufs)
	for bufID := 0; bufID < numBufs; bufID++ {
		buf := driver.PktBufAlloc(mempool)
		binary.BigEndian.PutUint32(buf.Pkt[20:24], pktSize)
		copy(buf.Pkt[64:64+pktSize], pktData)
		binary.BigEndian.PutUint16(buf.Pkt[64+24:64+26], calcIPChecksum(buf.Pkt[64+14:64+34]))
		bufs[bufID] = buf
	}
	//return them all to the mempool, all future allocations will return bufs with the data set above
	for bufID := 0; bufID < numBufs; bufID++ {
		driver.PktBufFree(bufs[bufID])
	}
	return mempool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <pci bus id>\n", os.Args[0])
		return
	}

	//get endianess of the machine -> Network Byte Order is always Big Endian
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		isBig = false
	}

	mempool := initMempool()
	dev := driver.IxyInit(os.Args[1], 1, 1)

	lastStatsPrinted := time.Now()
	counter := 0
	stats := new(driver.DeviceStats)
	statsOld := new(driver.DeviceStats)
	stats.StatsInit(dev)
	statsOld.StatsInit(dev)
	seqNum := uint32(0)

	//tx loop
	for {
		//we cannot immediately recycle packets, we need to allocate new packets every time
		//the old packets might still be used by the NIC: tx is async
		bufs := driver.PktBufAllocBatch(mempool, batchSize)
		for i := 0; i < len(bufs); i++ {
			//packets can be modified here, make sure to update the checksum when changing the IP header
			binary.BigEndian.PutUint32(bufs[i].Pkt[64+pktSize-4:64+pktSize], seqNum)
			seqNum++
		}
		//the packets could be modified here to generate multiple flows
		driver.IxyTxBatchBusyWait(dev, 0, bufs)

		//don't check time for every packet, this yields +10% performance :)
		if counter&0xFFF == 0 {
			t := time.Now()
			if t.Sub(lastStatsPrinted) > time.Second {
				//every second
				dev.ReadStats(stats)
				stats.PrintStatsDiff(statsOld, t.Sub(lastStatsPrinted))
				statsOld = stats
				lastStatsPrinted = t
			}
		}
		counter++
		//track stats
	}
}
