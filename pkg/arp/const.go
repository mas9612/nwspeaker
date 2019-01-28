package arp

const (
	// HardwareTypeEthernet represents Ethernet hardware type
	HardwareTypeEthernet = 1

	// ProtocolTypeIPv4 represents IPv4 protocol type
	ProtocolTypeIPv4 = 0x0800

	// OpRequest represents the packet is ARP request
	OpRequest = 1
	// OpReply represents the packet is ARP reply
	OpReply = 2
)

const (
	payloadLen = 4 * 7
)
