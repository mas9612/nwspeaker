package ipv4

const (
	// HeaderLen is the length of IPv4 header which is not have any option.
	HeaderLen = 20

	// Version4 is the version number of IPv4.
	Version4 = 4
	// Version6 is the version number of IPv6.
	Version6 = 6

	// FlagUnused is the flag which is not used.
	FlagUnused = 0x1
	// FlagDontFragment is the flag which shows this packet must not be fragmented.
	FlagDontFragment = 0x1 << 1
	// FlagMoreFragment is the flag which shows this packet is the last one of fragmented packets.
	FlagMoreFragment = 0x1 << 2

	// ProtoICMP is the protocol nunber of ICMP.
	ProtoICMP = 1
	// ProtoTCP is the protocol number of TCP.
	ProtoTCP = 6
	// ProtoUDP is the protocol number of UDP.
	ProtoUDP = 17

	// DefaultTTL is the default Time To Live.
	DefaultTTL = 255
)
