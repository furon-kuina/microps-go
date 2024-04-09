package net

type InternetProtocol int

type TransportProtocol byte

const (
	IcmpProtoNum TransportProtocol = 0x01
	TcpProtoNum  TransportProtocol = 0x06
	UdpProtoNum  TransportProtocol = 0x11
)

func (pn TransportProtocol) String() string {
	switch pn {
	case IcmpProtoNum:
		return "ICMP"
	case TcpProtoNum:
		return "TCP"
	case UdpProtoNum:
		return "UDP"
	default:
		return "Unknown"
	}
}
