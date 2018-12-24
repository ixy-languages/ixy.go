package driver

// #include <device.h>
import "C"
import (
	"os"
	"unsafe"
	"sync/atomic"
	"fmt"
	"time"
	"encoding/binary"
	"log"
)

//map C functions to Go

func setReg32(addr []byte, reg int, value uint32) {
	atomic.StoreUint32((*uint32)(unsafe.Pointer(&addr[reg])), value)
}

func getReg32(addr []byte, reg int) uint32 {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
}

func setFlags32(addr []byte, reg int, flags uint32) {
	setReg32(addr, reg, getReg32(addr, reg) | flags)
}

func clearFlags32(addr []byte, reg int, flags uint32) {
	setReg32(addr, reg, getReg32(addr, reg) &^ flags)
}

func waitClearReg32(addr []byte, reg int, mask uint32) {
	cur := atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	for (cur & mask) != 0 {
		fmt.Printf("waiting for flags %#x in register %#x to clear, current value %#x", mask, reg, cur)
		time.Sleep(10000*time.Microsecond)
		cur = atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	}
}

func waitSetReg32(addr []byte, reg int, mask uint32) {
	cur := atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	for (cur & mask) != mask {
		fmt.Printf("waiting for flags %#x in register %#x to clear, current value %#x", mask, reg, cur)
		time.Sleep(10000*time.Microsecond)
		cur = atomic.LoadUint32((*uint32)(unsafe.Pointer(&addr[reg])))
	}
}

//getter for pci io port resources
func readIo32(fd *os.File, offset uint) uint32 {
	fd.Sync()
	b := make([]byte, 4)
	n, err := fd.ReadAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci read wrong offset")
	}
	if isBig {
		return binary.BigEndian.Uint32(b)
	}
	return binary.LittleEndian.Uint32(b)
}

func readIo16(fd *os.File, offset uint) uint16 {
	fd.Sync()
	b := make([]byte, 2)
	n, err := fd.ReadAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci read wrong offset")
	}
	if isBig {
		return binary.BigEndian.Uint16(b)
	}
	return binary.LittleEndian.Uint16(b)
}

func readIo8(fd *os.File, offset uint) uint8 {
	fd.Sync()
	b := make([]byte, 1)
	n, err := fd.ReadAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci read wrong offset")
	}
	return uint8(b[0])
}

//setter for pci io port resources
func writeIo32(fd *os.File, value uint32, offset uint) {
	b := make([]byte, 4)
	if isBig {
		binary.BigEndian.PutUint32(b, value)
	} else {
		binary.LittleEndian.PutUint32(b, value)
	}
	n, err := fd.WriteAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci write wrong offset")
	}
	fd.Sync()
}

func writeIo16(fd *os.File, value uint16, offset uint) {
	b := make([]byte, 2)
	if isBig {
		binary.BigEndian.PutUint16(b, value)
	} else {
		binary.LittleEndian.PutUint16(b, value)
	}
	n, err := fd.WriteAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci write wrong offset")
	}
	fd.Sync()
}

func writeIo8(fd *os.File, value uint8, offset uint) {
	b := make([]byte, 1)
	b[0] = byte(value)
	n, err := fd.WriteAt(b, int64(offset))
	if err != nil || n < len(b) {
		log.Fatalf("Pci write wrong offset")
	}
	fd.Sync()
}


//getter/setter for PCIe memory mapped registers
/*func setCReg32(addr []byte, reg int, value uint32) {
	C.set_reg32((*C.uint8_t)(unsafe.Pointer(&addr[0])), C.int(reg), C.uint32_t(value))
}

func getCReg32(addr []byte, reg int) uint32 {
	return uint32(C.get_reg32((*C.uint8_t)(unsafe.Pointer(&addr[0])), C.int(reg)))
}

func setCFlags32(addr []byte, reg int, flags uint32) {
	C.set_flags32((*C.uint8_t)(unsafe.Pointer(&addr[0])), C.int(reg), C.uint32_t(flags))
}

func clearCFlags32(addr []byte, reg int, flags uint32) {
	C.clear_flags32((*C.uint8_t)(unsafe.Pointer(&addr[0])), C.int(reg), C.uint32_t(flags))
}

func waitClearCReg32(addr []byte, reg int, mask uint32) {
	C.wait_clear_reg32((*C.uint8_t)(unsafe.Pointer(&addr[0])), C.int(reg), C.uint32_t(mask))
}

func waitSetCReg32(addr []byte, reg int, mask uint32) {
	C.wait_set_reg32((*C.uint8_t)(unsafe.Pointer(&addr[0])), C.int(reg), C.uint32_t(mask))
}

//getter for pci io port resources
func readIo32C(fd *os.File, offset uint) uint32 {
	return uint32(C.read_io32(C.int(int(fd.Fd())), C.size_t(offset)))
}

func readIo16C(fd *os.File, offset uint) uint16 {
	return uint16(C.read_io16(C.int(int(fd.Fd())), C.size_t(offset)))
}

func readIo8C(fd *os.File, offset uint) uint8 {
	return uint8(C.read_io8(C.int(int(fd.Fd())), C.size_t(offset)))
}

//setter for pci io port resources
func writeIo32C(fd *os.File, value uint32, offset uint) {
	C.write_io32(C.int(int(fd.Fd())), C.uint32_t(value), C.size_t(offset))
}

func writeIo16C(fd *os.File, value uint16, offset uint) {
	C.write_io16(C.int(int(fd.Fd())), C.uint16_t(value), C.size_t(offset))
}

func writeIo8C(fd *os.File, value uint8, offset uint) {
	C.write_io8(C.int(int(fd.Fd())), C.uint8_t(value), C.size_t(offset))
}
*/
