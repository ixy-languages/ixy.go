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
	Pkt     []byte
	mempool *Mempool
}

//Mempool struct that describes a mempool
type Mempool struct {
	buf                               []byte
	bufSize, numEntries, freeStackTop uint32
	freeStack                         []uint32
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
func memoryAllocateDma(size uint32, requireContiguous bool) ([]byte, uintptr) { //maybe change so we don't have to return a struct butinstead return 2 values
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
	// require entries that neatly fit into the page size, this makes the memory pool much easier
	// otherwise our base_addr + index * size formula would be wrong because we can't cross a page-boundary
	if hugePageSize%entrySize != 0 {
		log.Fatalf("entry size must be a divisor of the huge page size (%v)", hugePageSize)
	}
	memvirt, _ := memoryAllocateDma(numEntries*entrySize, false)
	fStack := make([]uint32, numEntries)
	mempool := &Mempool{
		buf:          memvirt,
		bufSize:      entrySize,
		numEntries:   numEntries,
		freeStackTop: numEntries,
		freeStack:    fStack,
	}
	for i := uint32(0); i < numEntries; i++ {
		mempool.freeStack[i] = i
		pBufStart := i * entrySize
		if isBig {
			binary.BigEndian.PutUint64(mempool.buf[pBufStart:pBufStart+8], uint64(virtToPhys(uintptr(unsafe.Pointer(&mempool.buf[i*entrySize])))))
			binary.BigEndian.PutUint64(mempool.buf[pBufStart+8:pBufStart+16], uint64(uintptr(unsafe.Pointer(mempool))))
			binary.BigEndian.PutUint32(mempool.buf[pBufStart+16:pBufStart+20], i)
			binary.BigEndian.PutUint32(mempool.buf[pBufStart+20:pBufStart+24], 0)
		} else {
			binary.LittleEndian.PutUint64(mempool.buf[pBufStart:pBufStart+8], uint64(virtToPhys(uintptr(unsafe.Pointer(&mempool.buf[i*entrySize])))))
			binary.LittleEndian.PutUint64(mempool.buf[pBufStart+8:pBufStart+16], uint64(uintptr(unsafe.Pointer(mempool))))
			binary.LittleEndian.PutUint32(mempool.buf[pBufStart+16:pBufStart+20], i)
			binary.LittleEndian.PutUint32(mempool.buf[pBufStart+20:pBufStart+24], 0)
		}
	}
	return mempool
}

//PktBufAllocBatch allocates a batch of packets in the mempool
func PktBufAllocBatch(mempool *Mempool, bufs []*PktBuf) uint32 { //might try to fill from the front
	numBufs := uint32(len(bufs))
	if mempool.freeStackTop < numBufs {
		fmt.Printf("memory pool %v only has %v free bufs, requested %v\n", mempool, mempool.freeStackTop, numBufs)
		numBufs = mempool.freeStackTop
	}
	for i := uint32(0); i < numBufs; i++ {
		if bufs[i] == nil {
			bufs[i] = new(PktBuf)
		}
		entryID := mempool.freeStack[mempool.freeStackTop-1]
		mempool.freeStackTop--
		bufs[i].Pkt = mempool.buf[entryID*mempool.bufSize : (entryID+1)*mempool.bufSize]
		bufs[i].mempool = mempool
	}
	return numBufs
}

//PktBufAlloc allocates a single packet in the mempool
func PktBufAlloc(mempool *Mempool) *PktBuf {
	buf := make([]*PktBuf, 1)
	if PktBufAllocBatch(mempool, buf) == uint32(0) {
		return nil
	}
	return buf[0]
}

//PktBufFree frees PktBuf
func PktBufFree(buf *PktBuf) {
	if isBig {
		buf.mempool.freeStack[buf.mempool.freeStackTop] = binary.BigEndian.Uint32(buf.Pkt[16:20])
	} else {
		buf.mempool.freeStack[buf.mempool.freeStackTop] = binary.LittleEndian.Uint32(buf.Pkt[16:20])
	}
	buf.mempool.freeStackTop++
}
