package network

import (
	"github.com/furon-kuina/microps-go/internet"
)

type DeviceType string

const (
	Dummy    = DeviceType("DUMMY")
	Loopback = DeviceType("LOOPBACK")
	Ethernet = DeviceType("ETHERNET")
)

type RxHandler func(internet.InternetProtocolType, []byte)

type Device interface {
	Info() *DeviceInfo
	Interfaces() *NetworkInterfaceManager
	Open() error
	Close() error
	// dstに向けてdataを送信する
	Transmit(dst Device, ptype internet.InternetProtocolType, data []byte) error
	RxHandler(ptype internet.InternetProtocolType, data []byte)
}

const (
	UpFlag = 1 << iota
	LoopbackFlag
	BroadcastFlag
	P2PFlag
	NeedArpFlag
)

type DeviceInfo struct {
	Index         int
	Name          string
	Type          DeviceType
	Mtu           int
	Flags         uint32
	HeaderLength  uint16
	AddressLength uint16
	Addr          []uint16
	Peer          []uint8
	Broadcast     []uint8
	Priv          *any
}

func (devI DeviceInfo) State() string {
	if devI.Flags&UpFlag != 0 {
		return "UP"
	} else {
		return "DOWN"
	}
}
