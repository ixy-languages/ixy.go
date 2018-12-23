package driver

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	hugePageBits       = 21
	hugePageSize       = 1 << hugePageBits
	sizePktBufHeadroom = 40
)

//PktBuf bundles the raw byte representation of a buffer with the corresponding mempool
type PktBuf struct {
	//test whether this makes more sense, should not slow the driver down but costs a minimal amount of additional memory
	Pkt        []byte
	PhyAddr    uint64
	mempool    *Mempool
	mempoolIdx uint32
	Size       uint32
}

//Mempool struct that describes a mempool
type Mempool struct {
	buf                               []byte
	bufSize, numEntries, freeStackTop uint32
	freeStack                         []uint32
	packetBuf                         []*PktBuf
}

func virtToPhys(virt uintptr) uintptr {
	pagesize := syscall.Getpagesize()
	fd, err := os.OpenFile("/proc/self/pagemap", os.O_RDONLY, 0700)
	defer fd.Close()
	if err != nil {
		log.Fatalf("Error opening pagemap: %v", err)
	}
	rbuf := make([]byte, unsafe.Sizeof(uintptr(0)))
	_, err = fd.ReadAt(rbuf, int64(virt)/int64(pagesize)*int64(unsafe.Sizeof(uintptr(0))))
	if err != nil {
		log.Fatalf("Error translating address: %v", err)
	}
	var phy uintptr
	if isBig {
		phy = uintptr(binary.BigEndian.Uint64(rbuf))
	} else {
		phy = uintptr(binary.LittleEndian.Uint64(rbuf))
	}
	if phy == 0 {
		log.Fatalf("failed to translate virtual address %#v to physical address", virt)
	}
	//bits 0-54 are the page number; 0x7fffffffffffffULL -> ULL not available in go, but & with uintptr should be large enough
	return (phy&0x7fffffffffffff)*uintptr(pagesize) + virt%uintptr(pagesize)
}

// allocate memory suitable for DMA access in huge pages
// this requires hugetlbfs to be mounted at /mnt/huge
func memoryAllocateDma(size uint32, requireContiguous bool) ([]byte, uintptr) {
	//round up to multiples of 2 MB if necessary, this is the wasteful part
	if size%hugePageSize != 0 {
		size = ((size >> hugePageBits) + 1) << hugePageBits
	}
	if requireContiguous && size > hugePageSize {
		//this is the place to implement larger contiguous physical mappings if that's ever needed
		log.Fatal("could not map physically contiguous memory\n")
	}
	//unique filename
	rand.Seed(time.Now().UTC().UnixNano())
	id := rand.Uint32()
	path := fmt.Sprintf("/mnt/huge/ixy-%v-%v", os.Getpid(), id)
	//temporary file, will be deleted to prevent leaks of persistent pages
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, syscall.S_IRWXU)
	defer fd.Close()
	if err != nil {
		log.Fatalf("opening hugetlbfs file failed, check that /mnt/huge is mounted. Error: %v\n", err)
	}
	err = fd.Truncate(int64(size))
	if err != nil {
		log.Fatalf("allocating huge page memory failed, check hugetlbfs configuration. Error: %v\n", err)
	}
	var mmap []byte
	mmap, err = syscall.Mmap(int(fd.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_HUGETLB)
	if err != nil {
		log.Fatalf("mmaping hugepage failed. Error: %v\n", err)
	}
	//never swap out DMA memory
	err = syscall.Mlock(mmap)
	if err != nil {
		log.Fatalf("disabling swap for DMA memory failed. Error: %v\n", err)
	}
	syscall.Unlink(path)
	return mmap, virtToPhys(uintptr(unsafe.Pointer(&mmap[0])))
}

//MemoryAllocateMempool allocate mempool with numEntries*entrySize
//allocate a memory pool from which DMA'able packet buffers can be allocated
//this is currently not yet thread-safe, i.e., a pool can only be used by one thread,
//this means a packet can only be sent/received by a single thread
//entry_size can be 0 to use the default
func MemoryAllocateMempool(numEntries, entrySize uint32) *Mempool {
	if entrySize == 0 {
		entrySize = 2048
	}
	//require entries that neatly fit into the page size, this makes the memory pool much easier
	//otherwise our base_addr + index * size formula would be wrong because we can't cross a page-boundary
	if hugePageSize%entrySize != 0 {
		log.Fatalf("entry size must be a divisor of the huge page size (%v)", hugePageSize)
	}
	memvirt, _ := memoryAllocateDma(numEntries*entrySize, false)
	fStack := make([]uint32, numEntries)
	pBuf := make([]*PktBuf, numEntries)
	mempool := &Mempool{
		buf:          memvirt,
		bufSize:      entrySize,
		numEntries:   numEntries,
		freeStackTop: numEntries,
		freeStack:    fStack,
		packetBuf:    pBuf,
	}
	for i := uint32(0); i < numEntries; i++ {
		mempool.freeStack[i] = i
		mempool.packetBuf[i] = &PktBuf{
			Pkt:        mempool.buf[i*entrySize : (i+1)*entrySize],
			PhyAddr:    uint64(virtToPhys(uintptr(unsafe.Pointer(&mempool.buf[i*entrySize])))),
			mempool:    mempool,
			mempoolIdx: i,
			Size:       0,
		}
	}
	return mempool
}

//PktBufAllocBatch allocates a batch of packets in the mempool, use PktBufAlloc for single packets
func PktBufAllocBatch(mempool *Mempool, numBufs uint32) []*PktBuf {
	if mempool.freeStackTop < numBufs {
		fmt.Printf("memory pool %v only has %v free bufs, requested %v\n", mempool, mempool.freeStackTop, numBufs)
		numBufs = mempool.freeStackTop
	}
	bufs := make([]*PktBuf, numBufs)
	for i := uint32(0); i < numBufs; i++ {
		mempool.freeStackTop--
		id := mempool.freeStack[mempool.freeStackTop]
		bufs[i] = mempool.packetBuf[id]
	}
	return bufs
}

//PktBufAlloc allocates a single packet in the mempool
func PktBufAlloc(mempool *Mempool) *PktBuf {
	//while it is a special case of PktBufAllocBatch, it is better no not make a slice we don't need
	if mempool.freeStackTop < 1 {
		fmt.Printf("memory pool %v is currently full, cannot allocate packet\n", mempool)
		return nil
	}
	mempool.freeStackTop--
	id := mempool.freeStack[mempool.freeStackTop]
	return mempool.packetBuf[id]
}

//PktBufFree frees PktBuf
func PktBufFree(buf *PktBuf) {
	buf.mempool.freeStack[buf.mempool.freeStackTop] = buf.mempoolIdx
	buf.mempool.freeStackTop++
}
