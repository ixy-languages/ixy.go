package driver

import (
	"fmt"
	"time"
)

// DeviceStats holds stats for a device
type DeviceStats struct {
	device IxyInterface
	RXPackets, TXPackets, RXBytes, TXBytes uint64
}

// Reset resets the packet and byte counters
func (stats *DeviceStats) Reset() {
	stats.RXPackets = 0
	stats.TXPackets = 0
	stats.RXBytes = 0
	stats.TXBytes = 0
}

// Diff subtracts the given statistics and returns the difference.
// The difference is not linked to any device, and cannot be used further.
func (stats *DeviceStats) Diff(old *DeviceStats) *DeviceStats {
	return &DeviceStats {
		RXPackets: stats.RXPackets - old.RXPackets,
		TXPackets: stats.TXPackets - old.TXPackets,
		RXBytes:   stats.RXBytes - old.RXBytes,
		TXBytes:   stats.TXBytes - old.TXBytes,
	}
}

// Rate returns the traffic rates given old statistics and the time between the old and current statistics.
func (stats *DeviceStats) Rate(old *DeviceStats, dt time.Duration) (rxpps, txpps, rxBps, txBps float64) {
	diff := stats.Diff(old)
	rxpps = float64(diff.RXPackets) / dt.Seconds()
	txpps = float64(diff.TXPackets) / dt.Seconds()
	rxBps = float64(diff.RXBytes) / dt.Seconds()
	txBps = float64(diff.TXBytes) / dt.Seconds()
	return
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
	fmt.Printf("[%v] RX: %v bytes %v packets\n", addr, stats.RXBytes, stats.RXBytes)
	fmt.Printf("[%v] TX: %v bytes %v packets\n", addr, stats.TXBytes, stats.TXPackets)
}

//PrintStatsDiff get difference between reciever and previous stats
func (stats *DeviceStats) PrintStatsDiff(statsOld *DeviceStats, nanos time.Duration) {
	rxpps, txpps, rxBps, txBps := stats.Rate(statsOld, nanos)
	oldDev := statsOld.device.getIxyDev()
	newDev := stats.device.getIxyDev()
	var addr string
	if statsOld.device != nil {
		addr = oldDev.PciAddr
	} else {
		addr = "???"
	}
	fmt.Printf("[%v] RX: %v Mbit/s %.2f Mpps\n", addr, rxBps/1e6, rxpps/1e6)

	if stats.device != nil {
		addr = newDev.PciAddr
	} else {
		addr = "???"
	}
	fmt.Printf("[%v] TX: %v Mbit/s %.2f Mpps\n", addr, txBps/1e6, txpps/1e6)
}

//StatsInit initialize device stats
func (stats *DeviceStats) StatsInit(dev IxyInterface) {
	stats.Reset()
	stats.device = dev
	if dev != nil {
		dev.ReadStats(nil)
	}
}
