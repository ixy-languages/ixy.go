package driver

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

//much to do... produced some pretty shitty translation -> use the features Go gives us!

//try to circumvent the use of pointer as much as possible

const (
	hugePageBits       = 21
	hugePageSize       = 1 << hugePageBits
	sizePktBufHeadroom = 40
)

//PktBuf stuct that describes a packet buffer
type PktBuf struct {
	BufAddrPhy       uintptr
	Mempool          *Mempool
	MempoolIdx, Size uint32
	HeadRoom         [sizePktBufHeadroom]uint8
	Data             []byte
}

//Mempool struct tthat describes a mempool
type Mempool struct {
	buf []byte
	//bufSize != len(buf)
	numEntries, entrySize, freeStackTop, bufSize uint32
	freeStack                                    []uint32
}

type dmaMemory struct {
	virt []byte
	phy  uintptr
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
	//since the uintptr we read is either 4 or 8 bytes, we need to reconstruct it from the byte slice
	var phy uintptr
	for i, v := range rbuf {
		phy += uintptr(v) << (8 * uint(len(rbuf)-i-1)) //sizeof(byte) = 8
	}
	if phy == 0 {
		log.Fatalf("failed to translate virtual address %#v to physical address", virt)
	}
	//bits 0-54 are the page number; 0x7fffffffffffffULL -> ULL not available in go, but & with uintptr should be large enough
	return (phy&0x7fffffffffffff)*uintptr(pagesize) + virt%uintptr(pagesize)
}

// allocate memory suitable for DMA access in huge pages
// this requires hugetlbfs to be mounted at /mnt/huge
func memoryAllocateDma(size uint32, requireContiguous bool) dmaMemory {
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
		log.Fatal("opening hugetlbfs file failed, check that /mnt/huge is mounted\n")
	}
	err = fd.Truncate(int64(size))
	if err != nil {
		log.Fatal("allocating huge page memory failed, check hugetlbfs configuration\n")
	}
	var mmap []byte
	mmap, err = syscall.Mmap(int(fd.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_HUGETLB)
	if err != nil {
		log.Fatal("mmaping hugepage failed\n")
	}
	//never swap out DMA memory
	err = syscall.Mlock(mmap)
	if err != nil {
		log.Fatal("disabling swap for DMA memory failed\n")
	}
	syscall.Unlink(path)
	//virt is a slice, phys a pointer -> extract Slice pointer for translation
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&mmap))
	return dmaMemory{mmap, virtToPhys(hdr.Data)}
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
		log.Fatal("entry size must be a divisor of the huge page size (%v)", hugePageSize)
	}
	mem := memoryAllocateDma(numEntries*entrySize, false)
	mempool := Mempool{
		buf:          mem.virt,
		numEntries:   numEntries,
		entrySize:    entrySize,
		bufSize:      entrySize,
		freeStackTop: numEntries,
	}
	for i := uint32(0); i < numEntries; i++ {
		mempool.freeStack[i] = i
		//idea: get pointer of the byte array and interpret as PktBuf
		buf := (*PktBuf)(unsafe.Pointer(&mempool.buf[(i * entrySize) /*:((i + 1) * entrySize)*/]))
		buf.BufAddrPhy = virtToPhys(uintptr(unsafe.Pointer(buf)))
		buf.MempoolIdx = i
		buf.Mempool = &mempool
		buf.Size = 0
	}
	return &mempool
}

//PktBufAllocBatch allocates a batch of packets in the mempool
func PktBufAllocBatch(mempool *Mempool, bufs []*PktBuf) uint32 {
	numBufs := uint32(len(bufs))
	if mempool.freeStackTop < numBufs {
		fmt.Println("memory pool %v only has %v free bufs, requested %v", mempool, mempool.freeStackTop, numBufs)
		numBufs = mempool.freeStackTop
	}
	for i := uint32(0); i < numBufs; i++ {
		entryID := mempool.freeStack[mempool.freeStackTop-1]
		mempool.freeStackTop--
		bufs[i] = (*PktBuf)(unsafe.Pointer(uintptr(unsafe.Pointer(&mempool.buf[0])) + uintptr(entryID)*uintptr(mempool.bufSize)))
	}
	return numBufs
}

//PktBufAlloc allocates a single packet in the mempool
func PktBufAlloc(mempool *Mempool) []*PktBuf {
	buf := make([]*PktBuf, 1)
	PktBufAllocBatch(mempool, buf)
	return buf
}

//PktBufFree frees PktBuf
func PktBufFree(buf *PktBuf) {
	mempool := buf.Mempool
	mempool.freeStack[mempool.freeStackTop] = buf.MempoolIdx
	mempool.freeStackTop++
}

/*func main() {
	//test funtionality
	//flags to switch between mutually exclusive test cases
	dmaFlag := flag.Bool("dmaOnly", false, "test dma only")
	flag.Parse()
	//address translation:
	asdf := [1]int{5}
	fmt.Printf("virt addr: %p\nphys addr: %#v\n", &asdf, VirtToPhys(uintptr(unsafe.Pointer(&asdf))))
	//dma allocation
	if *dmaFlag {
		dma := memoryAllocateDma((1 << hugePageBits), false)
		fmt.Printf("dma success: %v\n", dma)
		return
	}
	//mempool allocation (contains dma allocation -> use either this or dma allocation)
	mempool := MemoryAllocateMempool(20, 2048)
	fmt.Printf("allocated mempool: %v\n", mempool)
	//allocate PktBuf
	pbuf := pktBufAlloc(mempool)
	fmt.Printf("allocated a single PktBuf\nmempool: %v\nPktBuf: %v\n", mempool, pbuf)
	//free PktBuf
	PktBufFree(pbuf[0])
	fmt.Printf("freed PktBuf. mempool: %v", mempool)
}*/
