package driver

// void mbarrier()
// {
// 		__asm__ volatile ("" : : : "memory");
// }
import "C" //use cgo to ensure read as volatile does not exist in go (at least I hope that it works, if not write the whole package in C and import that in ixgbe.go)
import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////
//TODO: memory barriers                                                                   //
//EDIT: not necessarily a memory barrier, just make sure, that it's read (no cached value)//
//-> only use with O_DIRECT flag                                                          //
//otherwise: try c call                                                                   //
//NO VOLATILE EQUIVALENT IN GO -> maybe use cgo                                           //
////////////////////////////////////////////////////////////////////////////////////////////

//https://stackoverflow.com/questions/18491032/does-go-support-volatile-non-volatile-variables

//DON'T use pointer, instead unse slices -> addr will always be []byte

//getter/setter for PCIe memory mapped registers
func setReg32(addr []byte, reg int, value uint32) {
	C.mbarrier()
	//*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) + uintptr(reg))) = value
	if isBig {
		binary.BigEndian.PutUint32(addr[reg:reg+4], value)
	} else {
		binary.LittleEndian.PutUint32(addr[reg:reg+4], value)
	}
}

func getReg32(addr []byte, reg int) uint32 {
	C.mbarrier()
	//return *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) + uintptr(reg)))
	if isBig {
		return binary.BigEndian.Uint32(addr[reg : reg+4])
	}
	return binary.LittleEndian.Uint32(addr[reg : reg+4])
}

func setFlags32(addr []byte, reg int, flags uint32) {
	setReg32(addr, reg, getReg32(addr, reg)|flags)
}

func clearFlags32(addr []byte, reg int, flags uint32) {
	setReg32(addr, reg, getReg32(addr, reg)&^flags)
}

func waitClearReg32(addr []byte, reg int, mask uint32) {
	C.mbarrier()
	var cur uint32
	if isBig {
		cur = binary.BigEndian.Uint32(addr[reg : reg+4])
	} else {
		cur = binary.LittleEndian.Uint32(addr[reg : reg+4])
	}
	for (cur & mask) != 0 {
		fmt.Printf("waiting for flags %+#v in register %+#v to clear, current value %+#v\n", mask, reg, cur)
		time.Sleep(10 * time.Millisecond)
		C.mbarrier()
		if isBig {
			cur = binary.BigEndian.Uint32(addr[reg : reg+4])
		} else {
			cur = binary.LittleEndian.Uint32(addr[reg : reg+4])
		}
	}
}

func waitSetReg32(addr []byte, reg int, mask uint32) {
	C.mbarrier()
	var cur uint32
	if isBig {
		cur = binary.BigEndian.Uint32(addr[reg : reg+4])
	} else {
		cur = binary.LittleEndian.Uint32(addr[reg : reg+4])
	}
	for (cur & mask) != mask {
		fmt.Printf("waiting for flags %+#v in register %+#v, current value %+#v\n", mask, reg, cur)
		time.Sleep(10 * time.Millisecond)
		C.mbarrier()
		if isBig {
			cur = binary.BigEndian.Uint32(addr[reg : reg+4])
		} else {
			cur = binary.LittleEndian.Uint32(addr[reg : reg+4])
		}
	}
}

//getter for pci io port resources
func readIo32(fd *os.File, offset uint) uint32 {
	C.mbarrier()
	buf := make([]byte, 4)
	_, err := fd.ReadAt(buf, int64(offset))
	if err != nil {
		log.Fatal("ReadAt io resource failed")
	}
	if isBig {
		return binary.BigEndian.Uint32(buf[0:])
	}
	return binary.LittleEndian.Uint32(buf[0:])
	//return *(*uint32)(unsafe.Pointer(&buf[0]))
}

func readIo16(fd *os.File, offset uint) uint16 {
	C.mbarrier()
	buf := make([]byte, 2)
	_, err := fd.ReadAt(buf, int64(offset))
	if err != nil {
		log.Fatal("ReadAt io resource failed")
	}
	if isBig {
		return binary.BigEndian.Uint16(buf[0:])
	}
	return binary.LittleEndian.Uint16(buf[0:])
}

func readIo8(fd *os.File, offset uint) uint8 {
	C.mbarrier()
	buf := make([]byte, 1)
	_, err := fd.ReadAt(buf, int64(offset))
	if err != nil {
		log.Fatal("ReadAt io resource failed")
	}
	return uint8(buf[0])
}

//setter for pci io port resources
func writeIo32(fd *os.File, value uint32, offset uint) {
	buf := make([]byte, 4)
	if isBig {
		binary.BigEndian.PutUint32(buf, value)
	} else {
		binary.LittleEndian.PutUint32(buf, value)
	}
	_, err := fd.WriteAt(buf, int64(offset))
	if err != nil {
		log.Fatal("WriteAt io resource failed")
	}
	C.mbarrier()
}

func writeIo16(fd *os.File, value uint16, offset uint) {
	buf := make([]byte, 2)
	if isBig {
		binary.BigEndian.PutUint16(buf, value)
	} else {
		binary.LittleEndian.PutUint16(buf, value)
	}
	_, err := fd.WriteAt(buf, int64(offset))
	if err != nil {
		log.Fatal("WriteAt io resource failed")
	}
	C.mbarrier()
}

func writeIo8(fd *os.File, value uint8, offset uint) {
	buf := make([]byte, 1)
	buf[0] = byte(value)
	_, err := fd.WriteAt(buf, int64(offset))
	if err != nil {
		log.Fatal("WriteAt io resource failed")
	}
	C.mbarrier()
}
