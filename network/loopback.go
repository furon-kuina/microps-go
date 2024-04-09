package network

import (
	"fmt"
	"math"

	"github.com/furon-kuina/microps-go/internet"
	"github.com/furon-kuina/microps-go/util"
)

const (
	loopbackMtu = math.MaxInt
)

var (
	LoopbackIrq = "LOOPBACK"
)

type LoopbackDevice struct {
	*DeviceInfo
	*NetworkInterfaceManager
	queue   *util.ConcurrentQueue[LoopbackQueueEntry]
	irqChan chan Irq
	Rx      RxHandler
}

var _ Device = &LoopbackDevice{}

func NewLoopbackDevice(rxhandler RxHandler, irqChan chan Irq) *LoopbackDevice {
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
		irqChan,
		rxhandler,
	}
}

func (dev *LoopbackDevice) Open() error {
	if dev.Flags&UpFlag != 0 {
		return fmt.Errorf("already open")
	}
	dev.Flags |= UpFlag
	util.Debugf("Open: %s", dev.Name)
	return nil
}

func (dev *LoopbackDevice) Close() error {
	if dev.Flags&UpFlag == 0 {
		return fmt.Errorf("already closed")
	}
	dev.Flags ^= UpFlag
	util.Debugf("Close: %s", dev.Name)
	return nil
}

func (dev *LoopbackDevice) Info() *DeviceInfo {
	return dev.DeviceInfo
}

func (dev *LoopbackDevice) Transmit(dst Device, ptype internet.InternetProtocolType, data []byte) error {
	entry := LoopbackQueueEntry{
		ptype: ptype,
		data:  data,
	}
	dev.queue.Enqueue(entry)
	util.Debugf("enqueued: %d on queue", dev.queue.Len())
	dev.irqChan <- Irq{
		name: LoopbackIrq,
		dev:  dev,
	}
	return nil
}

func (dev *LoopbackDevice) RxHandler(ptype internet.InternetProtocolType, data []byte) {
	dev.Rx(ptype, data)
}

type LoopbackQueueEntry struct {
	ptype internet.InternetProtocolType
	data  []byte
}

// デバイスドライバ
func LoopbackIsr(irq Irq) error {
	dev, ok := irq.dev.(*LoopbackDevice)
	if !ok {
		return fmt.Errorf("invalid device")
	}
	for dev.queue.Len() > 0 {
		entry := dev.queue.Dequeue()
		util.Debugf("dequeued: %d on queue", dev.queue.Len())
		dev.RxHandler(entry.ptype, entry.data)
	}
	return nil
}
