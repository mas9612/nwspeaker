package icmp

import (
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/mas9612/nwspeaker/pkg/checksum"
)

var (
	supported = []uint8{
		TypeEcho,
		TypeEchoReply,
	}
)

// Message represents the ICMP message.
type Message struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Data     Payload
}

// Encode returns byte-encoded data of ICMP message.
func (m *Message) Encode() []byte {
	payload := m.Data.Encode()
	buffer := make([]byte, HeaderLen+len(payload))
	buffer[0] = m.Type
	buffer[1] = m.Code
	copy(buffer[4:], payload)
	checksum := checksum.SumOfOnesComplement16(buffer)
	copy(buffer[2:], checksum)

	return buffer
}

// Payload represents the ICMP data.
type Payload interface {
	Encode() []byte
}

// Option is the option for ICMP packet crafting.
type Option func(*config)

type config struct {
	srcMac string
	srcIP  string
	data   []byte
}

// Echo represents the data of ICMP Echo and Echo Reply message.
type Echo struct {
	Identifier     uint16
	SequenceNumber uint16
	Data           []byte
}

// SetData sets data for ICMP echo message.
func SetData(data []byte) Option {
	return func(c *config) {
		c.data = data
	}
}

// NewEcho creates ICMP Echo message and return it.
func NewEcho(outIfname, dstIP, dstMac string, opts ...Option) (*Message, error) {
	c := config{}
	for _, o := range opts {
		o(&c)
	}

	rand.Seed(time.Now().Unix())

	msg := &Message{
		Type: TypeEcho,
	}
	echoMsg := &Echo{
		Identifier: uint16(rand.Uint32()),
	}
	if len(c.data) > 0 {
		echoMsg.Data = c.data
	} else {
		echoMsg.Data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	msg.Data = echoMsg

	return msg, nil
}

// Encode returns byte-encoded data of Echo message.
func (e *Echo) Encode() []byte {
	buffer := make([]byte, 4+len(e.Data))
	binary.BigEndian.PutUint16(buffer[0:], e.Identifier)
	binary.BigEndian.PutUint16(buffer[2:], e.SequenceNumber)
	copy(buffer[4:], e.Data)
	return buffer
}
