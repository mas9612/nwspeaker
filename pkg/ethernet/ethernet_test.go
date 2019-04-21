package ethernet

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

var headerEncodeTests = []struct {
	in  *Header
	out []byte
}{
	{
		in: &Header{
			DstAddr:   net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			SrcAddr:   net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			EtherType: 0x0806,
		},
		out: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x08, 0x06},
	},
}

func TestHeaderEncode(t *testing.T) {
	for _, tt := range headerEncodeTests {
		b := tt.in.Encode()
		if !bytes.Equal(b, tt.out) {
			t.Errorf("Encode() = %x, but got %x\n", tt.out, b)
		}
	}
}

var headerParseTests = []struct {
	in  []byte
	out *Header
}{
	{
		in: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x08, 0x06},
		out: &Header{
			DstAddr:   net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			SrcAddr:   net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			EtherType: 0x0806,
		},
	},
}

func TestHeaderParse(t *testing.T) {
	for _, tt := range headerParseTests {
		h := Parse(tt.in)
		if !reflect.DeepEqual(h, tt.out) {
			t.Errorf("Parse(%x) = %v, but got %v\n", tt.in, tt.out, h)
		}
	}
}
