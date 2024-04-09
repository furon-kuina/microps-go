package driver

import (
	"fmt"
	"math"

	net "github.com/furon-kuina/microps-go"
	"github.com/furon-kuina/microps-go/util"
)

const (
	dummyMtu = math.MaxInt
)

type DummyDevice struct {
	*net.DeviceInfo
	*net.NetworkInterfaceManager
}

func (dd *DummyDevice) Open() error {
	return nil
}

func (dd *DummyDevice) Close() error {
	return nil
}

func (dd *DummyDevice) Info() *net.DeviceInfo {
	return dd.DeviceInfo
}

func (dd *DummyDevice) Transmit(nptype net.ProtocolType, data []byte, len int, dst net.Device) error {
	info := dd.Info()
	util.Infof("dev=%s, type=0x%04x, len=%d", info.Name, nptype, len)
	net.IntrRaiseIrq(net.DummyIrq)
	return nil
}

func NewDummyDevice() *DummyDevice {
	info := net.DeviceInfo{
		Type:          net.DummyDevice,
		Mtu:           dummyMtu,
		HeaderLength:  0,
		AddressLength: 0,
	}
	dd := DummyDevice{
		&info,
		net.NewNetworkInterfaceManager(),
	}
	newInfo, err := net.RegisterDevice(&dd)
	if err != nil {
		fmt.Printf("net.Register(): %v", err)
		return nil
	}
	dd.DeviceInfo = newInfo

	util.Infof("registered device %s", dd.Info().Name)
	net.RegisterIrqHandler(net.DummyIrq, dd.DummyIsr, true, dd.Info().Name, &dd)
	util.Infof("initialized, dev=%s", dd.Info().Name)
	return &dd
}

func (dd *DummyDevice) DummyIsr(irq net.Irq, dev net.Device) error {
	util.Debugf("irq=%d, dev=%s", irq, dd.Info().Name)
	return nil
}
