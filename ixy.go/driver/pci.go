package driver

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"syscall"
)

func removeDriver(pciAddr string) {
	path := fmt.Sprintf("/sys/bus/pci/devices/%v/driver/unbind", pciAddr)
	fd, err := os.OpenFile(path, os.O_WRONLY, 0700)
	defer fd.Close()
	if err != nil {
		fmt.Printf("no driver loaded: %v\n", err)
		return
	}
	_, err = fd.WriteAt([]byte(pciAddr), 0)
	if err != nil {
		log.Fatalf("failed to unload driver for device %v: %v\n", pciAddr, err)
	}
}

func enableDma(pciAddr string) {
	path := fmt.Sprintf("/sys/bus/pci/devices/%v/config", pciAddr)
	fd, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer fd.Close()
	if err != nil {
		log.Fatalf("Error opening pci config: %v", err)
	}
	//write to the command register (offset 4) in the PCIe config space
	//bit 2 is "bus master enable", --see PCIe 3.0 specification section 7.5.1.1-- see datasheet section 9.3.3.3
	dma := make([]byte, 2)
	_, err = fd.ReadAt(dma, 4)
	if err != nil {
		log.Fatalf("Error reading from config: %v", err)
	}
	var dmaInt uint16
	if isBig {
		dmaInt = binary.BigEndian.Uint16(dma)
	} else {
		dmaInt = binary.LittleEndian.Uint16(dma)
	}
	dmaInt |= 1 << 2
	if isBig {
		binary.BigEndian.PutUint16(dma, dmaInt)
	} else {
		binary.LittleEndian.PutUint16(dma, dmaInt)
	}
	_, err = fd.WriteAt(dma, 4)
	if err != nil {
		log.Fatalf("Error writing dma flag to config: %v\n", err)
	}
}

func pciMapResource(pciAddr string) []byte {
	path := fmt.Sprintf("/sys/bus/pci/devices/%v/resource0", pciAddr)
	fmt.Printf("Mapping PCI resource at %v\n", path)
	removeDriver(pciAddr)
	enableDma(pciAddr)
	fd, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("Error opening pci resource: %v", err)
	}
	stat, _ := fd.Stat()
	mmap, err := syscall.Mmap(int(fd.Fd()), 0, int(stat.Size()), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("Error trying to mmap: %v", err)
	}
	return mmap
}

func pciOpenResource(pciAddr string, resource string) *os.File {
	path := fmt.Sprintf("/sys/bus/pci/devices/%v/%v", pciAddr, resource)
	fmt.Printf("Opening PCI resource at %v \n", path)
	fd, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening pci resource: %v", err)
	}
	return fd
}
