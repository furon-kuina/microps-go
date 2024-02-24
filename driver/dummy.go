package driver

import (
	"fmt"
	"math"

	net "github.com/furon-kuina/microps-go"
	"github.com/furon-kuina/microps-go/util"
)

const (
	dummyMtu uint = math.MaxUint
)

type DummyDevice struct {
	net.NetDevice
}

func (dd *DummyDevice) Transmit(ndType net.NetDeviceType, data []byte, len uint, dst *any) error {
	info := dd.Info()
	util.Infof("dev=%s, type=0x%04x, len=%d", info.Name, ndType, len)
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
	util.Infof("initialized, dev=%s", dd.Info().Name)
	return dd
}
