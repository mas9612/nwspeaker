package arp

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/mas9612/nwspeaker/pkg/ethernet"
)

var newRequestTests = []struct {
	in  string
	out *Packet
}{
	{
		in: "192.168.0.1",
		out: &Packet{
			HType:    HardwareTypeEthernet,
			PType:    ProtocolTypeIPv4,
			HLen:     ethernet.EtherLen,
			PLen:     net.IPv4len,
			Op:       OpRequest,
			DstHAddr: ethernet.Zero,
			DstPAddr: net.IPv4(192, 168, 0, 1),
		},
	},
}

func TestNewRequest(t *testing.T) {
	for _, tt := range newRequestTests {
		b, err := NewRequest(tt.in)
		if err != nil {
			t.Errorf("NewRequest(%s) should not return error, but got %v\n", tt.in, err)
		}
		if !reflect.DeepEqual(b, tt.out) {
			t.Errorf("NewRequest(%s) = %x, but got %x\n", tt.in, tt.out, b)
		}
	}
}

var encodeTests = []struct {
	in  *Packet
	out []byte
}{
	{
		in: &Packet{
			HType:    HardwareTypeEthernet,
			PType:    ProtocolTypeIPv4,
			HLen:     ethernet.EtherLen,
			PLen:     net.IPv4len,
			Op:       OpRequest,
			SrcHAddr: net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			SrcPAddr: net.IPv4(192, 168, 1, 0),
			DstHAddr: ethernet.Zero,
			DstPAddr: net.IPv4(192, 168, 0, 1),
		},
		out: []byte{
			0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66,
			0xc0, 0xa8, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0xa8, 0x00, 0x01,
		},
	},
}

func TestEncode(t *testing.T) {
	for _, tt := range encodeTests {
		b := tt.in.Encode()
		if !bytes.Equal(b, tt.out) {
			t.Errorf("Encode() = %x, but got %x\n", tt.out, b)
		}
	}
}

var parseTests = []struct {
	in  []byte
	out *Packet
}{
	{
		in: []byte{
			0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x08, 0x06,
			0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x01, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0xc0, 0xa8,
			0x01, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0xc0, 0xa8, 0x00, 0x01,
		},
		out: &Packet{
			HType:    HardwareTypeEthernet,
			PType:    ProtocolTypeIPv4,
			HLen:     ethernet.EtherLen,
			PLen:     net.IPv4len,
			Op:       OpRequest,
			SrcHAddr: net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			SrcPAddr: net.IPv4(192, 168, 1, 0),
			DstHAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			DstPAddr: net.IPv4(192, 168, 0, 1),
		},
	},
	{
		in: []byte{
			0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x08, 0x06,
			0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x02, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0xc0, 0xa8,
			0x01, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0xc0, 0xa8, 0x00, 0x01,
		},
		out: &Packet{
			HType:    HardwareTypeEthernet,
			PType:    ProtocolTypeIPv4,
			HLen:     ethernet.EtherLen,
			PLen:     net.IPv4len,
			Op:       OpReply,
			SrcHAddr: net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			SrcPAddr: net.IPv4(192, 168, 1, 0),
			DstHAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			DstPAddr: net.IPv4(192, 168, 0, 1),
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		p := Parse(tt.in)
		if !reflect.DeepEqual(p, tt.out) {
			t.Errorf("Parse(%x) = %v, but got %v\n", tt.in, tt.out, p)
		}
	}
}
