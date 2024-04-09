package network

import (
	"fmt"
	"math"

	"github.com/furon-kuina/microps-go/internet"
	"github.com/furon-kuina/microps-go/util"
)

const (
	dummyMtu = math.MaxInt
)

type DummyDevice struct {
	*DeviceInfo
	*NetworkInterfaceManager
	irqChan chan Irq
}

var _ Device = &DummyDevice{}

func NewDummyDevice(irqChan chan Irq) *DummyDevice {
	devInfo := DeviceInfo{
		Type:          Dummy,
		Mtu:           dummyMtu,
		HeaderLength:  0,
		AddressLength: 0,
	}
	im := NewNetworkInterfaceManager()

	return &DummyDevice{
		&devInfo,
		im,
		irqChan,
	}
}

func (dev *DummyDevice) Open() error {
	if dev.Flags&UpFlag != 0 {
		return fmt.Errorf("already open")
	}
	dev.Flags |= UpFlag
	util.Debugf("Open: %s", dev.Name)
	return nil
}

func (dev *DummyDevice) Close() error {
	if dev.Flags&UpFlag == 0 {
		return fmt.Errorf("already closed")
	}
	dev.Flags ^= UpFlag
	util.Debugf("Close: %s", dev.Name)
	return nil
}

func (dev *DummyDevice) Info() *DeviceInfo {
	return dev.DeviceInfo
}

func (dev *DummyDevice) Transmit(dst Device, ptype internet.InternetProtocolType, data []byte) error {
	info := dev.Info()
	if dst != nil {
		util.Debugf("Transmit: %s -> %s", info.Name, dst.Info().Name)
	}
	dev.irqChan <- Irq{
		name: DummyIrq,
		dev:  dev,
	}
	return nil
}

func (dev *DummyDevice) RxHandler(ptype internet.InternetProtocolType, data []byte) {}

func DummyIsr(irq Irq) error {
	util.Debugf("Dummy IRQ")
	return nil
}
