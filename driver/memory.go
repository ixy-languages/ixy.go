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

//PktBuf stuct that describes a packet buffer
type PktBuf struct {
	BufAddrPhy       uintptr
	Mempool          *Mempool
	MempoolIdx, Size uint32
	HeadRoom         [sizePktBufHeadroom]uint8
	Data             []byte //probably the biggest problem: Data has to be directly after HeadRoom
}

//Mempool struct that describes a mempool
type Mempool struct {
	buf []byte
	//bufSize != len(buf)
	bufSize, numEntries, freeStackTop uint32
	freeStack                         []uint32
}

type dmaMemory struct {
	virt []byte
	phy  uintptr
}

//todo: update driver to use this
//readPktBuf interprets and translates the arg as a PktBuf
func readPktBuf(mem []byte) *PktBuf {
	lenWoData := 24 + sizePktBufHeadroom
	var hdr [sizePktBufHeadroom]uint8
	copy(hdr[:], mem[24:lenWoData])
	var ret PktBuf
	if isBig {
		ret = PktBuf{
			BufAddrPhy: uintptr(binary.BigEndian.Uint64(mem[:8])),
			Mempool:    (*Mempool)(unsafe.Pointer(uintptr(binary.BigEndian.Uint64(mem[8:16])))),
			MempoolIdx: binary.BigEndian.Uint32(mem[16:20]),
			Size:       binary.BigEndian.Uint32(mem[20:24]),
		}
	} else {
		ret = PktBuf{
			BufAddrPhy: uintptr(binary.LittleEndian.Uint64(mem[:8])),
			Mempool:    (*Mempool)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(mem[8:16])))),
			MempoolIdx: binary.LittleEndian.Uint32(mem[16:20]),
			Size:       binary.LittleEndian.Uint32(mem[20:24]),
		}
	}
	ret.HeadRoom = hdr
	ret.Data = mem[lenWoData : len(mem)-lenWoData]
	return &ret
}

//setPktBuf writes the content of the PktBuf into the slice mem
func setPktBuf(mem []byte, pbuf *PktBuf) {
	lenWoData := 24 + sizePktBufHeadroom
	if isBig {
		binary.BigEndian.PutUint64(mem[:8], uint64(pbuf.BufAddrPhy))
		binary.BigEndian.PutUint64(mem[8:16], uint64(uintptr(unsafe.Pointer(pbuf.Mempool))))
		binary.BigEndian.PutUint32(mem[16:20], pbuf.MempoolIdx)
		binary.BigEndian.PutUint32(mem[20:24], pbuf.Size)
	} else {
		binary.LittleEndian.PutUint64(mem[:8], uint64(pbuf.BufAddrPhy))
		binary.LittleEndian.PutUint64(mem[8:16], uint64(uintptr(unsafe.Pointer(pbuf.Mempool))))
		binary.LittleEndian.PutUint32(mem[16:20], pbuf.MempoolIdx)
		binary.LittleEndian.PutUint32(mem[20:24], pbuf.Size)
	}
	copy(mem[24:lenWoData], pbuf.HeadRoom[:])
	copy(mem[lenWoData:lenWoData+len(pbuf.Data)], pbuf.Data)
}

//todo: check datasheet
//pktBufToByteSlice takes a PacketBuffer and returns a []byte that conforms to the needs of the NIC
/*func pktBufToByteSlice(pbuf PktBuf) []byte {
	lenWoData := 24 + sizePktBufHeadroom
	ret := make([]byte, lenWoData+len(pbuf.Data))
	if isBig {
		binary.BigEndian.PutUint64(ret[:8], uint64(pbuf.BufAddrPhy))
		binary.BigEndian.PutUint64(ret[8:16], uint64(uintptr(unsafe.Pointer(pbuf.Mempool))))
		binary.BigEndian.PutUint32(ret[16:20], pbuf.MempoolIdx)
		binary.BigEndian.PutUint32(ret[20:24], pbuf.Size)
	} else {
		binary.LittleEndian.PutUint64(ret[:8], uint64(pbuf.BufAddrPhy))
		binary.LittleEndian.PutUint64(ret[8:16], uint64(uintptr(unsafe.Pointer(pbuf.Mempool))))
		binary.LittleEndian.PutUint32(ret[16:20], pbuf.MempoolIdx)
		binary.LittleEndian.PutUint32(ret[20:24], pbuf.Size)
	}
	copy(ret[24:lenWoData], pbuf.HeadRoom[:])
	copy(ret[lenWoData:lenWoData+len(pbuf.Data)], pbuf.Data)
	return ret
}

//if raw is not a valid PktBuf the result will be undefined garbage values
func byteSliceToPktBuf(raw []byte) *PktBuf {
	lenWoData := 24 + sizePktBufHeadroom
	var hdr [sizePktBufHeadroom]uint8
	data := make([]byte, len(raw)-lenWoData) //allocate new array so it doesn't point so the underlying values, though I'm not sure if that would be the correct choice in this case -> will be bad on the performance
	copy(hdr[:], raw[24:lenWoData])
	copy(data[:], raw[lenWoData:len(raw)])
	var ret PktBuf
	if isBig {
		ret = PktBuf{
			BufAddrPhy: uintptr(binary.BigEndian.Uint64(raw[:8])),
			Mempool:    (*Mempool)(unsafe.Pointer(uintptr(binary.BigEndian.Uint64(raw[8:16])))),
			MempoolIdx: binary.BigEndian.Uint32(raw[16:20]),
			Size:       binary.BigEndian.Uint32(raw[20:24]),
		}
	} else {
		ret = PktBuf{
			BufAddrPhy: uintptr(binary.LittleEndian.Uint64(raw[:8])),
			Mempool:    (*Mempool)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint64(raw[8:16])))),
			MempoolIdx: binary.LittleEndian.Uint32(raw[16:20]),
			Size:       binary.LittleEndian.Uint32(raw[20:24]),
		}
	}
	ret.HeadRoom = hdr
	ret.Data = data
	return &ret
}
*/

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
	//hdr := (*reflect.SliceHeader)(unsafe.Pointer(&mmap))
	return dmaMemory{mmap, virtToPhys(uintptr(unsafe.Pointer(&mmap[0])))}
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
	mem := memoryAllocateDma(numEntries*entrySize, false)
	mempool := Mempool{
		buf:          mem.virt,
		bufSize:      entrySize,
		numEntries:   numEntries,
		freeStackTop: numEntries,
	}
	for i := uint32(0); i < numEntries; i++ {
		mempool.freeStack[i] = i
		//idea: get pointer of the byte array and interpret as PktBuf
		buf := new(PktBuf) //(*PktBuf)(unsafe.Pointer(&mempool.buf[(i * entrySize) /*:((i + 1) * entrySize)*/]))
		buf.BufAddrPhy = virtToPhys(uintptr(unsafe.Pointer(buf)))
		buf.MempoolIdx = i
		buf.Mempool = &mempool
		buf.Size = 0
		setPktBuf(mempool.buf[(i*entrySize):((i+1)*entrySize)], buf)
	}
	return &mempool
}

//PktBufAllocBatch allocates a batch of packets in the mempool
func PktBufAllocBatch(mempool *Mempool, bufs [][]byte) uint32 {
	numBufs := uint32(len(bufs))
	if mempool.freeStackTop < numBufs {
		fmt.Printf("memory pool %v only has %v free bufs, requested %v\n", mempool, mempool.freeStackTop, numBufs)
		numBufs = mempool.freeStackTop
	}
	for i := uint32(0); i < numBufs; i++ {
		entryID := mempool.freeStack[mempool.freeStackTop-1]
		mempool.freeStackTop--
		bufs[i] = mempool.buf[entryID*mempool.bufSize : (entryID+1)*mempool.bufSize]
		/*bufs[i] = (*PktBuf)(unsafe.Pointer(uintptr(unsafe.Pointer(&mempool.buf[0])) + uintptr(entryID)*uintptr(mempool.bufSize)))*/
	}
	return numBufs
}

//PktBufAlloc allocates a single packet in the mempool
func PktBufAlloc(mempool *Mempool) [][]byte {
	buf := make([][]byte, 1)
	PktBufAllocBatch(mempool, buf)
	return buf
}

//PktBufFree frees PktBuf
func PktBufFree(buf []byte) {
	pbuf := readPktBuf(buf)
	mempool := pbuf.Mempool
	mempool.freeStack[mempool.freeStackTop] = pbuf.MempoolIdx
	mempool.freeStackTop++
}
