package endian

import (
	"encoding/binary"
	"unsafe"
)

var (
	nativeEndian binary.ByteOrder
)

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0x0001)
	switch buf[0] {
	case 0x00:
		nativeEndian = binary.BigEndian
	case 0x01:
		nativeEndian = binary.LittleEndian
	default:
		panic("could not determine the host byteorder")
	}
}

// HostEndian returns the binary.ByteOrder instance of host endianess.
func HostEndian() binary.ByteOrder {
	return nativeEndian
}

// Htons converts byteorder from host byteorder to network byteorder.
func Htons(n uint16) uint16 {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, n)
	return nativeEndian.Uint16(buf)
}
