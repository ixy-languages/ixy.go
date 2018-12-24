package driver

import (
	"log"
	"unsafe"
)

const maxQueues = 64

var isBig = true

//IxyInterface is the interface that has to be implemented for all substrates such as the ixgbe or virtio
type IxyInterface interface {
	RxBatch(uint16, []*PktBuf) uint32
	TxBatch(uint16, []*PktBuf) uint32
	ReadStats(*DeviceStats)
	setPromisc(bool)
	getLinkSpeed() uint32
	getIxyDev() IxyDevice
}

//IxyDevice contains information common across all substrates
type IxyDevice struct {
	PciAddr     string
	DriverName  string
	NumRxQueues uint16
	NumTxQueues uint16
}

//IxyInit initializes the driver and hands back the interface
func IxyInit(pciAddr string, rxQueues, txQueues uint16) IxyInterface {
	//get endianness of the machine
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		isBig = false
	}
	//Read PCI configuration space
	config := pciOpenResource(pciAddr, "config")
	vendorID := readIo16(config, 0)
	deviceID := readIo16(config, 2)
	classID := readIo32(config, 8) >> 24
	config.Close()
	if classID != 2 {
		log.Fatalf("Device %v is not a NIC", pciAddr)
	}
	if vendorID == 0x1af4 && deviceID >= 0x1000 {
		log.Fatalln("Virtio not supported")
		return nil
	}
	//probably an ixgbe device
	return ixgbeInit(pciAddr, rxQueues, txQueues)
}

//IxyTxBatchBusyWait calls dev.TxBatch until all packets are queued with busy waiting
func IxyTxBatchBusyWait(dev IxyInterface, queueID uint16, bufs []*PktBuf) {
	numBufs := uint32(len(bufs))
	for numSent := uint32(0); numSent != numBufs; numSent += dev.TxBatch(queueID, bufs[numSent:]) {
	} //busy wait
}
