package net

import (
	"fmt"
	"sync"

	"github.com/furon-kuina/microps-go/util"
)

type Counter struct {
	mu sync.Mutex
	v  int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v++
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

func NewCounter() *Counter {
	return &Counter{
		mu: sync.Mutex{},
		v:  0,
	}
}

var (
	deviceCounter *Counter = NewCounter()
	devices       []NetDevice
)

type NetDeviceType int

const (
	Dummy NetDeviceType = iota + 1
	Loopback
	Ethernet
)

type NetDevice interface {
	Info() *NetDeviceInfo
	Open() error
	Close() error
	Transmit(NetDeviceType, []byte, uint, *any) error
}

type NetDeviceInfo struct {
	index         int
	Name          string
	Type          NetDeviceType
	Mtu           uint
	isUp          bool
	isLoopback    bool
	isBroadcast   bool
	isP2P         bool
	needsArp      bool
	HeaderLength  uint16
	AddressLength uint16
	addr          []uint16
	peer          []uint8
	broadcast     []uint8
	priv          *any
}

func (ndi *NetDeviceInfo) Info() *NetDeviceInfo {
	return ndi
}

func (ndc *NetDeviceInfo) Open() error {
	return nil
}

func (ndc *NetDeviceInfo) Close() error {
	return nil
}

func (ndc *NetDeviceInfo) Transmit(ndType NetDeviceType, data []byte, len uint, dst *any) error {
	return nil
}

func Register(nd NetDevice) error {
	info := nd.Info()
	info.index = deviceCounter.GetValue()
	deviceCounter.Increment()
	info.Name = fmt.Sprintf("net%d", info.index)
	devices = append(devices, nd)
	return nil
}

func (device NetDeviceInfo) State() string {
	if device.isUp {
		return "up"
	} else {
		return "down"
	}
}

func Wrap(errp *error, format string, args ...any) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
	}
}

func HandleInput(nd NetDevice, ndtype NetDeviceType, data *uint8, len uint64) {
	info := nd.Info()
	util.Debugf("dev=%s, type=%d, len=%d", info.Name, ndtype, len)
}

func Open(nd NetDevice) (err error) {
	info := nd.Info()
	defer Wrap(&err, "Open(%q)", info.Name)
	if info.isUp {
		return fmt.Errorf("already open")
	}
	if err = nd.Open(); err != nil {
		return err
	}
	info.isUp = true
	util.Debugf("dev=%s, state=%s", info.Name, info.State())
	return nil
}

func Close(nd NetDevice) (err error) {
	info := nd.Info()
	defer Wrap(&err, "Close(%q)", info.Name)
	if !info.isUp {
		return fmt.Errorf("not open")
	}
	if err = nd.Close(); err != nil {
		return err
	}
	info.isUp = false
	util.Debugf("dev=%s, state=%s", info.Name, info.State())
	return nil
}

func Output(nd NetDevice, ndtype NetDeviceType, data []byte, len uint, dst *any) (err error) {
	info := nd.Info()
	defer Wrap(&err, "Output(%q)", info.Name)
	if !info.isUp {
		return fmt.Errorf("not open")
	}
	if len > info.Mtu {
		return fmt.Errorf("too long, mtu=%d, len=%d", info.Mtu, len)
	}
	util.Infof("transmitting data %v", data)
	if err = nd.Transmit(ndtype, data, len, dst); err != nil {
		return fmt.Errorf("transmit failed: %v", err)
	}
	return nil
}

func Run() error {
	for _, dev := range devices {
		Open(dev)
	}
	util.Infof("running...")
	return nil
}

func Shutdown() {
	for _, dev := range devices {
		Close(dev)
	}
	util.Infof("shutting down")
}

func Init() error {
	util.Infof("initialized")
	return nil
}
