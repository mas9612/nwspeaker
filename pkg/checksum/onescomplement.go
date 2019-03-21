package checksum

import (
	"encoding/binary"
)

func onesComplement(b byte) byte {
	var complement byte
	for i := 7; i >= 0; i-- {
		if (b >> uint(i) & 0x1) == 0x0 {
			complement |= 0x1
		}
		if i > 0 {
			complement <<= 1
		}
	}
	return complement
}

func bytesOnesComplement(data []byte) []byte {
	complement := make([]byte, len(data))
	for i, d := range data {
		complement[i] = onesComplement(d)
	}
	return complement
}

// SumOfOnesComplement16 calculate one's complement sum of one's complement for each 16bit.
func SumOfOnesComplement16(b []byte) []byte {
	var sum uint32
	for i := 0; i < len(b); i += 2 {
		if i+2 > len(b) { // in the last word, if only one byte remain, add 0x00 as padding
			sum += uint32(binary.BigEndian.Uint16(append(b[i:], 0x00)))
		} else {
			sum += uint32(binary.BigEndian.Uint16(b[i : i+2]))
		}
		if (sum >> 16) > 0x0 {
			sum += sum >> 16
			sum &= 0xffff // clear carry bit
		}
	}
	sumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sumBytes, uint16(sum))
	return bytesOnesComplement(sumBytes)
}
