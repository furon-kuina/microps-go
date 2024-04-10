package network

import (
	"fmt"
	"sync"

	"github.com/furon-kuina/microps-go/internet"
	"github.com/furon-kuina/microps-go/util"
)

type NetworkStack struct {
	devices []Device
	mu      sync.Mutex
	irqm    *IrqManager
	pm      *ProtocolManager
	Enter   func(ptype internet.InternetProtocolType, data []byte)
}

func NewNetworkStack() *NetworkStack {
	irqm := NewIrqManager()
	pm := NewProtocolManager()
	ns := &NetworkStack{
		irqm: irqm,
		pm:   pm,
	}
	ns.Enter = func(ptype internet.InternetProtocolType, data []byte) {
		util.Debugf("enter: %s", ptype)
		ns.pm.protocols[ptype].queue.Enqueue(ProtocolQueueEntry{
			data: data,
		})
		ns.irqm.IrqChan <- Irq{
			name: "SOFT",
			dev:  nil,
		}
	}

	SoftIrqHandler := func(irq Irq) error {
		util.Debugf("SoftIrq")
		for _, proto := range ns.pm.protocols {
			for proto.queue.Len() > 0 {
				entry := proto.queue.Dequeue()
				err := proto.handler(entry.dev, entry.data)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	irqm.Register("SOFT", SoftIrqHandler, false)
	return ns
}

func (ns *NetworkStack) Run() error {
	ns.irqm.Run()
	for _, dev := range ns.devices {
		err := dev.Open()
		if err != nil {
			return err
		}
		util.Debugf("State: %s", dev.Info().State())
	}
	util.Debugf("running...")
	return nil
}

func (ns *NetworkStack) Shutdown() error {
	for _, dev := range ns.devices {
		err := dev.Close()
		if err != nil {
			return err
		}
	}
	util.Debugf("shutting down...")
	return nil
}

func (ns *NetworkStack) RegisterDevice(dev Device) {
	info := dev.Info()
	info.Index = len(ns.devices)
	info.Name = fmt.Sprintf("net%d", info.Index)
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.devices = append(ns.devices, dev)
}

func (ns *NetworkStack) SelectDevice(name string) Device {
	for _, dev := range ns.devices {
		if dev.Info().Name == name {
			return dev
		}
	}
	return nil
}

func (ns *NetworkStack) Transmit(srcName, dstName string, proto internet.InternetProtocolType, data []byte) (err error) {
	src := ns.SelectDevice(srcName)
	dst := ns.SelectDevice(dstName)
	info := src.Info()
	defer util.Wrap(&err, "Transmit(%q)", info.Name)
	if info.Flags&UpFlag == 0 {
		return fmt.Errorf("not open")
	}
	if len(data) > info.Mtu {
		return fmt.Errorf("too long, mtu=%d, len=%d", info.Mtu, len(data))
	}
	util.Infof("transmitting data %q", string(data))
	if err = src.Transmit(dst, proto, data); err != nil {
		return fmt.Errorf("transmit failed: %v", err)
	}
	return nil
}

func (ns *NetworkStack) Output(dev Device, ptype internet.InternetProtocolType, data []byte, dst Device) error {
	devInfo := dev.Info()
	if devInfo.Flags&UpFlag == 0 {
		return fmt.Errorf("%s not open", devInfo.Name)
	}
	if len(data) > devInfo.Mtu {
		return fmt.Errorf("too long")
	}
	if err := dev.Transmit(dst, ptype, data); err != nil {
		return err
	}
	return nil
}

func (ns *NetworkStack) RegisterIrq(irqName string, handler IrqHandler, shared bool) error {
	err := ns.irqm.Register(irqName, handler, shared)
	if err != nil {
		return err
	}
	return nil
}

func (ns *NetworkStack) RegisterProtocol(ptype internet.InternetProtocolType, handler internetHandler) error {
	return ns.pm.Register(ptype, handler)
}

func (ns *NetworkStack) NewDummyDevice(irqChan chan Irq) *DummyDevice {
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

func (ns *NetworkStack) NewLoopbackDevice() *LoopbackDevice {
	devInfo := DeviceInfo{
		Type:          Loopback,
		Mtu:           loopbackMtu,
		HeaderLength:  0,
		AddressLength: 0,
	}
	im := NewNetworkInterfaceManager()
	queue := util.NewConcurrentQueue[LoopbackQueueEntry]()
	return &LoopbackDevice{
		&devInfo,
		im,
		queue,
		ns.irqm.IrqChan,
		ns.Enter,
	}
}
