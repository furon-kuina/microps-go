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
	net.NetDevice
}

func (dd *DummyDevice) Transmit(nptype net.NetProtocolType, data []byte, len int, dst net.NetDevice) error {
	info := dd.Info()
	util.Infof("dev=%s, type=0x%04x, len=%d", info.Name, nptype, len)
	net.IntrRaiseIrq(net.DummyIrq)
	return nil
}

func NewDummyDevice() *DummyDevice {
	info := net.NetDeviceInfo{
		Type:          net.Dummy,
		Mtu:           dummyMtu,
		HeaderLength:  0,
		AddressLength: 0,
	}
	dd := &DummyDevice{
		&info,
	}
	if err := net.Register(dd); err != nil {
		fmt.Printf("net.Register(): %v", err)
		return nil
	}
	net.RegisterIrqHandler(net.DummyIrq, dd.DummyIsr, true, dd.Info().Name, dd)
	util.Infof("initialized, dev=%s", dd.Info().Name)
	return dd
}

func (dd *DummyDevice) DummyIsr(irq net.Irq, dev net.NetDevice) error {
	util.Debugf("irq=%d, dev=%s", irq, dd.Info().Name)
	return nil
}
