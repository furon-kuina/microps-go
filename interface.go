package net

import "fmt"

// Network Interface 関連

type NetworkInterfaceFamily string

const (
	Ip   = NetworkInterfaceFamily("IP")
	IpV6 = NetworkInterfaceFamily("IPV6")
)

type NetworkInterface interface {
	Dev() Device
	Family() NetworkInterfaceFamily
}

type NetworkInterfaceManager struct {
	interfaces []NetworkInterface
}

func NewNetworkInterfaceManager() *NetworkInterfaceManager {
	return &NetworkInterfaceManager{}
}

func (nim *NetworkInterfaceManager) Interfaces() *NetworkInterfaceManager {
	return nim
}

func (m *NetworkInterfaceManager) AddInterface(iface NetworkInterface) error {
	for _, i := range m.interfaces {
		if i.Family() == iface.Family() {
			return fmt.Errorf("already registered: %s", i.Family())
		}
	}
	m.interfaces = append(m.interfaces, iface)
	return nil
}

func (m *NetworkInterfaceManager) SelectInterface(family NetworkInterfaceFamily) (NetworkInterface, error) {
	for _, i := range m.interfaces {
		if i.Family() == family {
			return i, nil
		}
	}
	return nil, fmt.Errorf("interface for family %s not found", family)
}
