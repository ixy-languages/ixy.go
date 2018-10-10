package driver

import (
	"fmt"
	"time"
)

//DeviceStats holds stats
type DeviceStats struct {
	device IxyInterface
	rxPackets, txPackets, rxBytes, txBytes uint64
}

//PrintStats prints stats
func (stats *DeviceStats) PrintStats() {
	dev := stats.device.getIxyDev()
	var addr string
	if stats.device != nil {
		addr = dev.PciAddr
	} else {
		addr = "???"
	}
	fmt.Printf("[%v] RX: %v bytes %v packets\n", addr, stats.rxBytes, stats.rxPackets)
	fmt.Printf("[%v] TX: %v bytes %v packets\n", addr, stats.txBytes, stats.txPackets)
}

func diffMpps(pktsNew, pktsOld uint64, nanos time.Duration) float64 { //get duration by start := time.Now();t := time.Now();elapsed := t.Sub(start)
	return float64(pktsNew-pktsOld) / 1000000.0 / (float64(nanos) / 1000000000.0)
}

func diffMbit(bytesNew, bytesOld, pktsNew, pktsOld uint64, nanos time.Duration) uint32 {
	// take stuff on the wire into account, i.e., the preamble, SFD and IFG (20 bytes)
	// otherwise it won't show up as 10000 mbit/s with small packets which is confusing
	return uint32((float64(bytesNew-bytesOld)/1000000.0/(float64(nanos)/1000000000.0))*8 + diffMpps(pktsNew, pktsOld, nanos)*20*8)
}

//PrintStatsDiff get difference between reciever and previous stats
func (stats *DeviceStats) PrintStatsDiff(statsOld *DeviceStats, nanos time.Duration) {
	oldDev := statsOld.device.getIxyDev()
	newDev := stats.device.getIxyDev()
	var addr string
	if statsOld.device != nil {
		addr = oldDev.PciAddr
	} else {
		addr = "???"
	}
	fmt.Printf("[%v] RX: %v Mbit/s %.2f Mpps\n", addr, diffMbit(stats.rxBytes, statsOld.rxBytes, stats.rxPackets, statsOld.rxPackets, nanos), diffMpps(stats.rxPackets, statsOld.rxPackets, nanos))
	//fmt.Printf("[%v] DMA Good Packets: %v\n", addr, stats.rxDmaPackets)
	if stats.device != nil {
		addr = newDev.PciAddr
	} else {
		addr = "???"
	}
	fmt.Printf("[%v] TX: %v Mbit/s %.2f Mpps\n", addr, diffMbit(stats.txBytes, statsOld.txBytes, stats.txPackets, statsOld.txPackets, nanos), diffMpps(stats.txPackets, statsOld.txPackets, nanos))
}

//StatsInit initialize device stats
func (stats *DeviceStats) StatsInit(dev IxyInterface) {
	stats.rxPackets = 0
	stats.txPackets = 0
	stats.rxBytes = 0
	stats.txBytes = 0
	stats.device = dev
	if dev != nil {
		dev.ReadStats(nil)
	}
}
