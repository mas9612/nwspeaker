package endian

import (
	"encoding/binary"
	"testing"
)

var htonsTests = []struct {
	endian binary.ByteOrder
	in     uint16
	out    uint16
}{
	{
		endian: binary.BigEndian,
		in:     0x0806,
		out:    0x0806,
	},
	{
		endian: binary.LittleEndian,
		in:     0x0806,
		out:    0x0608,
	},
}

func TestHtons(t *testing.T) {
	for _, tt := range htonsTests {
		nativeEndian = tt.endian
		result := Htons(tt.in)
		if result != tt.out {
			t.Errorf("Htons(%d) = %d, but got %d\n", tt.in, tt.out, result)
		}
	}
}
