package driver

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

const (
	driverName = "ixy-ixgbe"

	maxRxQueueEntries = 4096
	maxTxQueueEntries = 4096

	numRxQueueEntries = 512
	numTxQueueEntries = 512

	txCleanBatch = 32
)

type ixgbeDevice struct {
	ixy      IxyDevice
	addr     []byte
	rxQueues []*ixgbeRxQueue
	txQueues []*ixgbeTxQueue
}

type ixgbeRxQueue struct {
	descriptors      []IxgbeAdvRxDesc
	mempool          *Mempool
	numEntries       uint16
	rxIndex          uint16
	virtualAddresses [][]byte //[]*PktBuf		//experimental
}

type ixgbeTxQueue struct {
	descriptors      []IxgbeAdvTxDesc
	numEntries       uint16
	cleanIndex       uint16
	txIndex          uint16
	virtualAddresses [][]byte //[]*PktBuf
}

func (dev *ixgbeDevice) getIxyDev() IxyDevice {
	return dev.ixy
}

// see section 4.6.4
func (dev *ixgbeDevice) initLink() {
	// should already be set by the eeprom config, maybe we shouldn't override it here to support weirdo nics?
	setReg32(&dev.addr[0], IXGBE_AUTOC, (getReg32(&dev.addr[0], IXGBE_AUTOC)&^IXGBE_AUTOC_LMS_MASK)|IXGBE_AUTOC_LMS_10G_SERIAL)
	setReg32(&dev.addr[0], IXGBE_AUTOC, (getReg32(&dev.addr[0], IXGBE_AUTOC)&^IXGBE_AUTOC_10G_PMA_PMD_MASK)|IXGBE_AUTOC_10G_XAUI)
	// negotiate link
	setFlags32(&dev.addr[0], IXGBE_AUTOC, IXGBE_AUTOC_AN_RESTART)
	// datasheet wants us to wait for the link here, but we can continue and wait afterwards
}

func (dev *ixgbeDevice) startRxQueue(queueID int) {
	//debug info
	fmt.Printf("starting rx queue %v\n", queueID)
	queue := dev.rxQueues[queueID]
	// 2048 as pktbuf size is strictly speaking incorrect:
	// we need a few headers (1 cacheline), so there's only 1984 bytes left for the device
	// but the 82599 can only handle sizes in increments of 1 kb; but this is fine since our max packet size
	// is the default MTU of 1518
	// this has to be fixed if jumbo frames are to be supported
	// mempool should be >= the number of rx and tx descriptors for a forwarding application
	mempoolEntries := uint32(numRxQueueEntries + numTxQueueEntries)
	if mempoolEntries < 4096 {
		mempoolEntries = 4096
	}
	queue.mempool = MemoryAllocateMempool(mempoolEntries, 2048)
	if queue.numEntries&(queue.numEntries-1) == 0 {
		log.Fatal("number of queue entries must be a power of 2")
	}
	for i := uint16(0); i < queue.numEntries; i++ {
		//rxd is the struct that contains only the array raw
		rxd := queue.descriptors[i]
		buf := PktBufAlloc(queue.mempool)[0]
		if buf == nil {
			log.Fatal("failed to allocate rx descriptor")
		}
		//write into first 64 bits
		copy(rxd.raw[:8], buf[0:8])
		//write into last 64 bits
		binary.LittleEndian.PutUint64(rxd.raw[8:], uint64(0))
		queue.virtualAddresses[i] = buf //[:0]	see later
	}
	//enable queue and wait if necessary
	setFlags32(&dev.addr[0], IXGBE_RXDCTL(queueID), IXGBE_RXDCTL_ENABLE)
	waitSetReg32(&dev.addr[0], IXGBE_RXDCTL(queueID), IXGBE_RXDCTL_ENABLE)
	//rx queue starts out full
	setReg32(&dev.addr[0], IXGBE_RDH(queueID), 0)
	//was set to 0 before in the init funtion
	setReg32(&dev.addr[0], IXGBE_RDT(queueID), uint32(queue.numEntries-1))
}

func (dev *ixgbeDevice) startTxQueue(queueID int) {
	fmt.Printf("starting tx queue %v", queueID)
	queue := dev.txQueues[queueID]
	if queue.numEntries&(queue.numEntries-1) == 0 {
		log.Fatal("number of queue entries must be a power of 2")
	}
	//tx queue starts out empty
	setReg32(&dev.addr[0], IXGBE_TDH(queueID), 0)
	setReg32(&dev.addr[0], IXGBE_TDT(queueID), 0)
	// enable queue and wait if necessary
	setFlags32(&dev.addr[0], IXGBE_TXDCTL(queueID), IXGBE_TXDCTL_ENABLE)
	waitSetReg32(&dev.addr[0], IXGBE_TXDCTL(queueID), IXGBE_TXDCTL_ENABLE)
}

//see section 4.6.7
func (dev *ixgbeDevice) initRx() {
	//make sure that rx is disabled while re-configuring it
	//the datasheet also wants us to disable some crypto-offloading related rx paths (but we don't care about them)
	clearFlags32(&dev.addr[0], IXGBE_RXCTRL, IXGBE_RXCTRL_RXEN)
	// no fancy dcb or vt, just a single 128kb packet buffer for us
	setReg32(&dev.addr[0], IXGBE_RXPBSIZE(0), IXGBE_RXPBSIZE_128KB)
	for i := 1; i < 8; i++ {
		setReg32(&dev.addr[0], IXGBE_RXPBSIZE(i), 0)
	}
	// always enable CRC offloading
	setFlags32(&dev.addr[0], IXGBE_HLREG0, IXGBE_HLREG0_RXCRCSTRP)
	setFlags32(&dev.addr[0], IXGBE_RDRXCTL, IXGBE_RDRXCTL_CRCSTRIP)

	//accept broadcast packets
	setFlags32(&dev.addr[0], IXGBE_FCTRL, IXGBE_FCTRL_BAM)

	//per-queue config, same for all queues
	for i := uint16(0); i < dev.ixy.NumRxQueues; i++ {
		fmt.Printf("initializing rx queue %v\n", i)
		//enable advanced rx descriptors, we could also get away with legacy descriptors, but they aren't really easier
		setReg32(&dev.addr[0], IXGBE_SRRCTL(int(i)), (getReg32(&dev.addr[0], IXGBE_SRRCTL(int(i)))&^IXGBE_SRRCTL_DESCTYPE_MASK)|IXGBE_SRRCTL_DESCTYPE_ADV_ONEBUF)
		//drop_en causes the nic to drop packets if no rx descriptors are available instead of buffering them
		//a single overflowing queue can fill up the whole buffer and impact operations if not setting this flag
		setFlags32(&dev.addr[0], IXGBE_SRRCTL(int(i)), IXGBE_SRRCTL_DROP_EN)
		//setup descriptor ring, see section 7.1.9
		ringSizeBytes := uint32(numRxQueueEntries * 16) //unsafe.Sizeof([16]byte or IxgbeAdvRxDesc)
		mem := memoryAllocateDma(ringSizeBytes, true)
		//neat trick from Snabb: initialize 0xFF to prevent rogue memory accesses on premature DMA activation
		//copy memory -> fill in log(n) compared to iterating (even with optimization this is faster)
		//if golang mmap does zero-initialize this, then we yould probably skip this
		if len(mem.virt) != 0 {
			mem.virt[0] = 0xFF
		}
		for filled := 1; filled < len(mem.virt); filled *= 2 {
			copy(mem.virt[filled:], mem.virt[:filled])
		}
		setReg32(&dev.addr[0], IXGBE_RDBAL(int(i)), uint32(mem.phy&0xFFFFFFFF))
		setReg32(&dev.addr[0], IXGBE_RDBAH(int(i)), uint32(mem.phy>>32))
		setReg32(&dev.addr[0], IXGBE_RDLEN(int(i)), ringSizeBytes)
		fmt.Printf("rx ring %v phy addr: %+#v\n", i, mem.phy)
		fmt.Printf("rx ring %v virt addr: %+#v\n", i, uintptr(unsafe.Pointer(&mem.virt[0])))
		//set ring to empty at start
		setReg32(&dev.addr[0], IXGBE_RDH(int(i)), 0)
		setReg32(&dev.addr[0], IXGBE_RDT(int(i)), 0)
		//private data for the driver, 0-initialized
		queue := dev.rxQueues[i]
		queue.numEntries = numRxQueueEntries
		queue.rxIndex = 0
		/*
			we have the []byte mem.virt. That's where the packet descriptors (IxgbeAdvRxDesc) will be
			-> find a way so we get an []IxgbeAdvDesc that uses the same memory area
		*/
		//original: queue->descriptors = (union ixgbe_adv_tx_desc*) mem.virt;
		//test: AdvDesc hold slices
		//-> divide mem.virt into len(mem.virt)/16 subslices and collect them in a slice
		desc := make([]IxgbeAdvRxDesc, len(mem.virt)/16)
		for i := 0; i < len(mem.virt)/16; i++ {
			desc[i] = IxgbeAdvRxDesc{raw: mem.virt[i*16 : (i+1)*16]} //creating subslices should be not too inefficient since they all use the same underlying array (what we want ayways)
		}
		queue.descriptors = desc
	}
	//last step is to set some magic bits mentioned in the last sentence in 4.6.7
	setFlags32(&dev.addr[0], IXGBE_CTRL_EXT, IXGBE_CTRL_EXT_NS_DIS)
	//this flag probably refers to a broken feature: it's reserved and initialized as '1' but it must be set to '0'
	//there isn't even a constant in ixgbe_types.h for this flag
	for i := uint16(0); i < dev.ixy.NumRxQueues; i++ {
		clearFlags32(&dev.addr[0], IXGBE_DCA_RXCTRL(int(i)), 1<<12)
	}

	// start RX
	setFlags32(&dev.addr[0], IXGBE_RXCTRL, IXGBE_RXCTRL_RXEN)
}

//see section 4.6.8
func (dev *ixgbeDevice) initTx() {
	//crc offload and small packet padding
	setFlags32(&dev.addr[0], IXGBE_HLREG0, IXGBE_HLREG0_TXCRCEN|IXGBE_HLREG0_TXPADEN)

	//set default buffer size allocations
	//see also: section 4.6.11.3.4, no fancy features like DCB and VTd
	setReg32(&dev.addr[0], IXGBE_TXPBSIZE(0), IXGBE_TXPBSIZE_40KB)
	for i := 1; i < 8; i++ {
		setReg32(&dev.addr[0], IXGBE_TXPBSIZE(i), 0)
	}
	//required when not using DCB/VTd
	setReg32(&dev.addr[0], IXGBE_DTXMXSZRQ, 0xFFFF)
	clearFlags32(&dev.addr[0], IXGBE_RTTDCS, IXGBE_RTTDCS_ARBDIS)

	//per-queue config for all queues
	for i := 0; uint16(i) < dev.ixy.NumTxQueues; i++ {
		fmt.Printf("initializing tx queue %v", i)

		//setup descriptor ring, see section 7.1.9
		ringSizeBytes := uint32(numTxQueueEntries * 16) //see initRx
		mem := memoryAllocateDma(ringSizeBytes, true)
		if len(mem.virt) != 0 {
			mem.virt[0] = 0xFF
		}
		for filled := 1; filled < len(mem.virt); filled *= 2 {
			copy(mem.virt[filled:], mem.virt[:filled])
		}
		setReg32(&dev.addr[0], IXGBE_TDBAL(i), uint32(mem.phy&0xFFFFFFFF))
		setReg32(&dev.addr[0], IXGBE_TDBAH(i), uint32(mem.phy>>32))
		setReg32(&dev.addr[0], IXGBE_TDLEN(i), ringSizeBytes)
		fmt.Printf("tx ring %v phy addr: %+#v", i, mem.phy)
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&mem.virt))
		fmt.Printf("tx ring %v virt addr: %+#v", i, hdr.Data)

		//descriptor writeback magic values, important to get good performance and low PCIe overhead
		//see 7.2.3.4.1 and 7.2.3.5 for an explanation of these values and how to find good ones
		//we just use the defaults from DPDK here, but this is a potentially interesting point for optimizations
		txdctl := getReg32(&dev.addr[0], IXGBE_TXDCTL(i))
		//there are no defines for this in ixgbe_h for some reason
		//pthresh: 6:0, hthresh: 14:8, wthresh: 22:16
		txdctl &= ^(uint32(0x3F | (0x3F << 8) | (0x3F << 16))) //clear bits
		txdctl |= (36 | (8 << 8) | (4 << 16))                  //from DPDK
		setReg32(&dev.addr[0], IXGBE_TXDCTL(i), txdctl)

		//private data for the driver, 0-initialized
		queue := dev.txQueues[i]
		queue.numEntries = numTxQueueEntries
		//see rxInit
		desc := make([]IxgbeAdvTxDesc, len(mem.virt)/16)
		for i := 0; i < len(mem.virt)/16; i++ {
			desc[i] = IxgbeAdvTxDesc{raw: mem.virt[i*16 : (i+1)*16]} //creating subslices should be not too inefficient since they all use the same underlying array (what we want ayways)
		}
		queue.descriptors = desc
	}
	//final step: enable dma
	setReg32(&dev.addr[0], IXGBE_DMATXCTL, IXGBE_DMATXCTL_TE)
}

func (dev *ixgbeDevice) waitForLink() {
	fmt.Printf("Waiting for link...")
	maxWait := time.Second * 10
	pollInterval := time.Millisecond * 10
	for speed := dev.getLinkSpeed(); speed == 0 && maxWait > 0; speed = dev.getLinkSpeed() {
		time.Sleep(pollInterval)
		maxWait -= pollInterval
	}
	fmt.Printf("Link speed is %v Mbit/s", dev.getLinkSpeed())
}

//see section 4.6.3
func (dev *ixgbeDevice) resetAndInit() {
	fmt.Printf("Resetting device %v\n", dev.ixy.PciAddr)
	//section 4.6.3.1
	setReg32(&dev.addr[0], IXGBE_EIMC, 0x7FFFFFFF)

	//section 4.6.3.2
	setReg32(&dev.addr[0], IXGBE_CTRL, IXGBE_CTRL_RST_MASK)
	waitClearReg32(&dev.addr[0], IXGBE_CTRL, IXGBE_CTRL_RST_MASK)
	time.Sleep(time.Millisecond)

	//section 4.6.3.1 - disable interrupts again after reset
	setReg32(&dev.addr[0], IXGBE_EIMC, 0x7FFFFFFF)

	fmt.Printf("Initializing device %v\n", dev.ixy.PciAddr)

	//section 4.6.3 - Wait for EEPROM auto read completion
	waitSetReg32(&dev.addr[0], IXGBE_EEC, IXGBE_EEC_ARD)

	//section 4.6.3 - Wait for DMA initialization done (RDRXCTL.DMAIDONE)
	waitSetReg32(&dev.addr[0], IXGBE_RDRXCTL, IXGBE_RDRXCTL_DMAIDONE)

	//section 4.6.4 - initialize link (auto negotiation)
	dev.initLink()

	//section 4.6.7 - init rx
	dev.initRx()

	//section 4.6.8 - init tx
	dev.initTx()

	//enables queue after initializing everything
	for i := uint16(0); i < dev.ixy.NumRxQueues; i++ {
		dev.startRxQueue(int(i))
	}
	for i := uint16(0); i < dev.ixy.NumTxQueues; i++ {
		dev.startTxQueue(int(i))
	}

	//skip last step from 4.6.3 - don't want interrupts
	//finally, enable promisc mode by default
	dev.setPromisc(true)

	//wait for some time for the link to come up
	dev.waitForLink()
}

//create struct and initialize it -> return pointer to struct (return an ixgbeDevice which conforms to the IxyInterface Interface)
func ixgbeInit(pciAddr string, rxQueues, txQueues uint16) IxyInterface {
	if syscall.Getuid() != 0 {
		fmt.Println("Not running as root, this will probably fail")
	}
	if rxQueues > maxQueues {
		log.Fatalf("cannot configure %v rx queues: limit is %v", rxQueues, maxQueues)
	}
	if txQueues > maxQueues {
		log.Fatalf("cannot configure %v tx queues: limit is %v", txQueues, maxQueues)
	}
	dev := new(ixgbeDevice)
	dev.ixy.PciAddr = pciAddr
	dev.ixy.DriverName = driverName
	dev.ixy.NumRxQueues = rxQueues
	dev.ixy.NumTxQueues = txQueues
	dev.addr = pciMapResource(pciAddr)             //todo: .addr -> [] bc a pointer ro a slice is not what we want -> change .addr uses to &.addr[x]
	dev.rxQueues = make([]*ixgbeRxQueue, rxQueues) //todo: see if that's correct
	dev.txQueues = make([]*ixgbeTxQueue, txQueues)
	dev.resetAndInit()
	return dev
}

func (dev *ixgbeDevice) getLinkSpeed() uint32 {
	links := getReg32(&dev.addr[0], IXGBE_LINKS)
	if links&IXGBE_LINKS_UP == 0 {
		return 0
	}
	switch links & IXGBE_LINKS_SPEED_82599 {
	case IXGBE_LINKS_SPEED_100_82599:
		return 100
	case IXGBE_LINKS_SPEED_1G_82599:
		return 1000
	case IXGBE_LINKS_SPEED_10G_82599:
		return 10000
	default:
		return 0
	}
}

func (dev *ixgbeDevice) setPromisc(enabled bool) {
	if enabled {
		fmt.Println("enabling promisc mode")
		setFlags32(&dev.addr[0], IXGBE_FCTRL, IXGBE_FCTRL_MPE|IXGBE_FCTRL_UPE)
	} else {
		fmt.Println("disabling promisc mode")
		clearFlags32(&dev.addr[0], IXGBE_FCTRL, IXGBE_FCTRL_MPE|IXGBE_FCTRL_UPE)
	}
}

//read stat counters and accumulate in stats
//stats may be nil to just reset the counters
func (dev *ixgbeDevice) ReadStats(stats *DeviceStats) {
	rxPkts := getReg32(&dev.addr[0], IXGBE_GPRC)
	txPkts := getReg32(&dev.addr[0], IXGBE_GPTC)
	rxBytes := uint64(getReg32(&dev.addr[0], IXGBE_GORCL)) | uint64(getReg32(&dev.addr[0], IXGBE_GORCH))<<32
	txBytes := uint64(getReg32(&dev.addr[0], IXGBE_GOTCL)) | uint64(getReg32(&dev.addr[0], IXGBE_GOTCH))<<32
	if stats != nil {
		stats.rxPackets += uint64(rxPkts)
		stats.txPackets += uint64(txPkts)
		stats.rxBytes += rxBytes
		stats.txBytes += txBytes
	}
}

func wrapRing(index, ringSize uint16) uint16 {
	return (index + 1) & (ringSize - 1)
}

//section 1.8.2 and 7.1
//try to receive a single packet if one is available, non-blocking
//see datasheet section 7.1.9 for an explanation of the rx ring structure
//tl;dr: we control the tail of the queue, the hardware the head
func (dev *ixgbeDevice) RxBatch(queueID uint16, bufs [][]byte, numBufs uint32) uint32 {
	queue := dev.rxQueues[queueID]
	rxIndex := queue.rxIndex
	lastRxIndex := rxIndex
	bufIndex := uint32(0)
	for ; bufIndex < numBufs; bufIndex++ {
		//rx descriptors are explained in 7.1.5, advances rx descriptors in 7.1.6
		//since go doesn't support unions, we just use the raw byte array (slice)
		descPtr := queue.descriptors[rxIndex]
		var status uint32
		if isBig {
			status = binary.BigEndian.Uint32(descPtr.raw[8:12]) //bit 64 - 95 of the advanced rx recieve descriptor are the status/error
		} else {
			status = binary.LittleEndian.Uint32(descPtr.raw[8:12])
		}
		if status&IXGBE_RXDADV_STAT_DD != 0 {
			if status&IXGBE_RXDADV_STAT_EOP == 0 {
				log.Fatal("multi-segment packets are not supported - increase buffer size or decrease MTU")
			}
			//got a packet, read and copy the whole descriptor
			desc := descPtr.raw
			buf := queue.virtualAddresses[rxIndex]
			if isBig {
				//buf.Size = uint32(binary.BigEndian.Uint16(desc[12:14])) //bit 96 - 111 of the advanced rx recieve descriptor are the length
				binary.BigEndian.PutUint32(buf[20:24], uint32(binary.BigEndian.Uint16(desc[12:14])))
			} else {
				//buf.Size = uint32(binary.LittleEndian.Uint16(desc[12:14]))
				binary.LittleEndian.PutUint32(buf[20:24], uint32(binary.LittleEndian.Uint16(desc[12:14])))
			}
			//this would be the place to implement RX offloading by translating the device-specific flags
			//to an independent representation in the buf (similiar to how DPDK works)
			//need a new mbuf for the descriptor
			newBuf := PktBufAlloc(queue.mempool)[0] //slice has only a single element
			if newBuf == nil {
				//we could handle empty mempools more gracefully here, but it would be quite messy...
				//make your mempools large enough
				log.Fatal("failed to allocate new mbuf for rx, you are either leaking memory or your mempool is too small")
			}
			//reset descriptor
			if isBig {
				binary.BigEndian.PutUint64(descPtr.raw[:8], uint64(virtToPhys(uintptr(unsafe.Pointer(&newBuf[64])))))
				binary.BigEndian.PutUint64(descPtr.raw[8:], uint64(0)) //resets the flags
			} else {
				binary.LittleEndian.PutUint64(descPtr.raw[:8], uint64(virtToPhys(uintptr(unsafe.Pointer(&newBuf[64])))))
				binary.LittleEndian.PutUint64(descPtr.raw[8:], uint64(0))
			}
			queue.virtualAddresses[rxIndex] = newBuf
			bufs[bufIndex] = buf
			//want to read the next one in the next iteration, but we still need the last/curretn to update RDT later
			lastRxIndex = rxIndex
			rxIndex = wrapRing(rxIndex, queue.numEntries)
		} else {
			break
		}
	}
	if rxIndex != lastRxIndex {
		//tell hardware that we are done
		//this is intentionally off by one, otherwise we'd set RDT=RDH if we are receiving faster than packets are coming in
		//RDT=RDH means queue is full
		setReg32(&dev.addr[0], IXGBE_RDT(int(queueID)), uint32(lastRxIndex))
		queue.rxIndex = rxIndex
	}
	return bufIndex // number of packets stored in bufs; bufIndex "points" to the next index
}

//section 1.8.1 and 7.2
//we control the tail, hardware the head
//huge performance gains possible here by sending packets in batches - writing to TDT for every packet is not efficient
//returns the number of packets transmitted, will not block when the queue is full
func (dev *ixgbeDevice) TxBatch(queueID uint16, bufs [][]byte, numBufs uint32) uint32 {
	queue := dev.txQueues[queueID]
	//the descriptor is explained in section 7.2.3.2.4
	//we just use a struct copy & pasted from intel (not), but it basically has two formats (hence a union):
	//1. the write-back format which is written by the NIC once sending it is finished this is used in step 1
	//2. the read format which is read by the NIC and written by us, this is used in step 2
	//in go we use a plain byte slice large enough to hold the union since unions don't exist in go

	cleanIndex := queue.cleanIndex //next descriptor to clean up
	curIndex := queue.txIndex      //next descriptor to use for tx

	//step 1: clean up descriptors that were sent out by the hardware and return them to the mempool
	//start by reading step 2 which is done first for each packet
	//cleaning up must be done in batches for performance reasons, so this is unfortunately somewhat complicated
	for {
		cleanable := curIndex - cleanIndex
		if cleanable < 0 { //wrap around
			cleanable = queue.numEntries + cleanable
		}
		if cleanable < txCleanBatch {
			break
		}
		//calculcate the index of the last transcriptor in the clean batch
		//we can't check all descriptors for performance reasons
		cleanupTo := cleanIndex + txCleanBatch - 1
		if cleanupTo >= queue.numEntries {
			cleanupTo -= queue.numEntries
		}
		txd := queue.descriptors[cleanupTo]
		var status uint32
		if isBig {
			status = binary.BigEndian.Uint32(txd.raw[12:]) //last 32 bit
		} else {
			status = binary.LittleEndian.Uint32(txd.raw[12:]) //last 32 bit
		}
		//hardware sets this flag as soon as it's sent out, we can give back all bufs in the batch back to the mempool
		if status&IXGBE_ADVTXD_STAT_DD != 0 {
			i := cleanIndex
			for {
				buf := queue.virtualAddresses[i]
				PktBufFree(buf)
				if i == cleanupTo {
					break
				}
				i = wrapRing(cleanupTo, queue.numEntries)
			}
			//next descriptor to be cleaned up is one after the one we just cleaned
			cleanIndex = wrapRing(cleanupTo, queue.numEntries)
		} else {
			//clean the whole batch or nothing; yes, this leaves some packets in
			//the queue forever if you stop transmitting, but that's not a real concern
			break
		}
	}
	queue.cleanIndex = cleanIndex

	//step 2: send out as many of our packets as possible
	sent := uint32(0)
	for ; sent < numBufs; sent++ {
		nextIndex := wrapRing(curIndex, queue.numEntries)
		//we are full if the next index is the one we are trying to reclaim
		if cleanIndex == nextIndex {
			break
		}
		buf := bufs[sent]
		//remember virtual address to clean it up later
		queue.virtualAddresses[curIndex] = buf
		queue.txIndex = wrapRing(queue.txIndex, queue.numEntries)
		txd := queue.descriptors[curIndex]
		//NIC reads from here
		if isBig {
			size := binary.BigEndian.Uint32(buf[20:24])
			binary.BigEndian.PutUint64(txd.raw[:8], uint64(virtToPhys(uintptr((unsafe.Pointer(&buf[64]))))))
			// always the same flags: one buffer (EOP), advanced data descriptor, CRC offload, data length
			binary.BigEndian.PutUint32(txd.raw[8:12], IXGBE_ADVTXD_DCMD_EOP|IXGBE_ADVTXD_DCMD_RS|IXGBE_ADVTXD_DCMD_IFCS|IXGBE_ADVTXD_DCMD_DEXT|IXGBE_ADVTXD_DTYP_DATA|size)
			//no fancy offloading stuff - only the total payload length
			//implement offloading flags here:
			// 	* ip checksum offloading is trivial: just set the offset
			// 	* tcp/udp checksum offloading is more annoying, you have to precalculate the pseudo-header checksum
			binary.BigEndian.PutUint32(txd.raw[12:16], size<<IXGBE_ADVTXD_PAYLEN_SHIFT)
		} else {
			size := binary.LittleEndian.Uint32(buf[20:24])
			binary.LittleEndian.PutUint64(txd.raw[:8], uint64(uintptr((unsafe.Pointer(&buf[64])))))
			binary.LittleEndian.PutUint32(txd.raw[8:12], IXGBE_ADVTXD_DCMD_EOP|IXGBE_ADVTXD_DCMD_RS|IXGBE_ADVTXD_DCMD_IFCS|IXGBE_ADVTXD_DCMD_DEXT|IXGBE_ADVTXD_DTYP_DATA|size)
			binary.LittleEndian.PutUint32(txd.raw[12:16], size<<IXGBE_ADVTXD_PAYLEN_SHIFT)
		}
		curIndex = nextIndex
	}
	//send out by advancing tail, i.e., pass control of the bufs to the nic
	//this seems like a textbook case for a release memory order, but Intel's driver doesn't even use a compiler barrier here
	setReg32(&dev.addr[0], IXGBE_TDT(int(queueID)), uint32(queue.txIndex))
	return sent
}

func test() {
	bufs := make([]*PktBuf, 5)
	bufs[0].Data[0]++
}
