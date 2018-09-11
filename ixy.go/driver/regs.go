package driver

// #include <device.h>
import "C"
import (
	"os"
	"unsafe"
)

//map C functions to Go

//getter/setter for PCIe memory mapped registers
func setCReg32(addr []byte, reg int, value uint32) {
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
