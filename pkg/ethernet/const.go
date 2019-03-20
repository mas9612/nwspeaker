package ethernet

const (
	// EtherLen is the length of Ethernet address
	EtherLen = 6

	// HeaderLen is the length of Ethernet header
	HeaderLen = EtherLen*2 + 2

	// TypeIPv4 is the type number of IPv4.
	TypeIPv4 = 0x0800
	// TypeARP is the type number of ARP
	TypeARP = 0x0806
)
