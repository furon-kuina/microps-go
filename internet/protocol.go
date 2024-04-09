package internet

type InternetProtocolType string

const (
	DummyProtocol = InternetProtocolType("DUMMY")
	IpProtocol    = InternetProtocolType("IP")
	ArpProtocol   = InternetProtocolType("ARP")
	IpV6Protocol  = InternetProtocolType("IPV6")
)
