package driver

import (
	"fmt"
	"math"

	net "github.com/furon-kuina/microps-go"
	"github.com/furon-kuina/microps-go/util"
)

const (
	loopbackMtu        = math.MaxInt
	loopbackQueueLimit = 10
)

type LoopbackDevice struct {
	irq net.Irq
	q   *util.ConcurrentQueue[LoopbackEntry]
	*net.DeviceInfo
	net.NetworkInterfaceManager
}

func (ld *LoopbackDevice) Info() *net.DeviceInfo {
	return ld.DeviceInfo
}

func (ld *LoopbackDevice) Open() error {
	return nil
}

func (ld *LoopbackDevice) Close() error {
	return nil
}

func (ld *LoopbackDevice) Transmit(nptype net.ProtocolType, data []byte, len int, dst net.Device) (err error) {
	defer util.Wrap(&err, "Loopback transmit: dev=%s, dst=%s", ld.Info().Name, dst.Info().Name)
	if ld.q.Len() >= loopbackQueueLimit {
		return fmt.Errorf("queue is full")
	}
	entry := LoopbackEntry{
		nptype: nptype,
		data:   data,
	}
	ld.q.Enqueue(entry)
	util.Debugf("enqueued: (num:%d), dev=%s, type=0x%04x, len=%d", ld.q.Len(), ld.Info().Name, nptype, len)
	util.Infof("sending data: %v", data)
	net.IntrRaiseIrq(ld.irq)
	return nil
}

type LoopbackEntry struct {
	nptype net.ProtocolType
	data   []byte
}

func NewLoopbackDevice() *LoopbackDevice {
	info := net.DeviceInfo{
		Type:          net.LoopbackDevice,
		Mtu:           loopbackMtu,
		Flags:         net.LoopbackFlag,
		HeaderLength:  0,
		AddressLength: 0,
	}

	q := util.NewConcurrentQueue[LoopbackEntry]()

	ld := LoopbackDevice{
		net.LoopbackIrq,
		q,
		&info,
		net.NetworkInterfaceManager{},
	}

	newInfo, err := net.RegisterDevice(&ld)
	if err != nil {
		fmt.Printf("net.Register(): %v", err)
		return nil
	}
	ld.DeviceInfo = newInfo
	net.RegisterIrqHandler(ld.irq, ld.LoopbackIsr, true, ld.Info().Name, &ld)
	util.Infof("initialized, dev=%s", ld.Info().Name)
	return &ld
}

func (ld *LoopbackDevice) LoopbackIsr(irq net.Irq, dev net.Device) error {
	util.Debugf("irq=%d, dev=%s", irq, ld.Info().Name)
	for {
		if ld.q.IsEmpty() {
			break
		}
		entry := ld.q.Dequeue()
		util.Debugf("dequeued: (num: %d), dev=%s, type=0x%04x", ld.q.Len(), ld.Info().Name, entry.nptype)
		net.InputHandler(dev, entry.nptype, entry.data, len(entry.data))
	}

	return nil
}
