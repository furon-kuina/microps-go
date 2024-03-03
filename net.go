package net

import (
	"fmt"
	"sync"

	"github.com/furon-kuina/microps-go/util"
)

const (
	Dummy NetDeviceType = iota + 1
	Loopback
	Ethernet
)

const (
	DummyIrq Irq = iota + 1
	LoopbackIrq
	SoftIrq
)

const (
	UpFlag = 1 << iota
	LoopbackFlag
	BroadcastFlag
	P2PFlag
	needArpFlag
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
	deviceCounter = NewCounter()
	devices       []NetDevice
	c             = sync.NewCond(&sync.Mutex{})
	irqReady      = false
)

type NetDeviceType int

type NetDevice interface {
	Info() *NetDeviceInfo
	Open() error
	Close() error
	Transmit(NetProtocolType, []byte, int, NetDevice) error
}

type NetDeviceInfo struct {
	index         int
	Name          string
	Type          NetDeviceType
	Mtu           int
	Flags         uint32
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

func (ndc *NetDeviceInfo) Transmit(npType NetProtocolType, data []byte, len int, dst NetDevice) error {
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
	if device.Flags&UpFlag != 0 {
		return "up"
	} else {
		return "down"
	}
}

// dev:    どのNICからデータが届いたか
// nptype: 届いたデータのプロトコル
func InputHandler(dev NetDevice, nptype NetProtocolType, data []byte, len int) error {
	info := dev.Info()
	util.Debugf("dev=%s, type=%d, len=%d", info.Name, nptype, len)
	for _, proto := range protocols {
		if proto.nptype == nptype {
			entry := NetProtocolQueueEntry{
				dev:  dev,
				data: data,
				len:  len,
			}
			// nptype の受信キューにデータを積む
			proto.q.Enqueue(entry)
			util.Debugf("enqueued (num:%d), dev=%s, type=0x%04x, len=%d", proto.q.Len(), info.Name, nptype, len)
			// ソフトウェア割り込み
			IntrRaiseIrq(SoftIrq)
		}
	}
	return nil
}

func Open(nd NetDevice) (err error) {
	info := nd.Info()
	defer util.Wrap(&err, "Open(%q)", info.Name)
	if info.Flags&UpFlag != 0 {
		return fmt.Errorf("already open")
	}
	if err = nd.Open(); err != nil {
		return err
	}
	info.Flags |= UpFlag
	util.Debugf("dev=%s, state=%s", info.Name, info.State())
	return nil
}

func Close(nd NetDevice) (err error) {
	info := nd.Info()
	defer util.Wrap(&err, "Close(%q)", info.Name)
	if info.Flags&UpFlag == 0 {
		return fmt.Errorf("not open")
	}
	if err = nd.Close(); err != nil {
		return err
	}
	info.Flags ^= UpFlag
	util.Debugf("dev=%s, state=%s", info.Name, info.State())
	return nil
}

func Output(nd NetDevice, nptype NetProtocolType, data []byte, len int, dst NetDevice) (err error) {
	info := nd.Info()
	defer util.Wrap(&err, "Output(%q)", info.Name)
	if info.Flags&UpFlag == 0 {
		return fmt.Errorf("not open")
	}
	if len > info.Mtu {
		return fmt.Errorf("too long, mtu=%d, len=%d", info.Mtu, len)
	}
	util.Infof("transmitting data %v", data)
	if err = nd.Transmit(nptype, data, len, dst); err != nil {
		return fmt.Errorf("transmit failed: %v", err)
	}
	return nil
}

func Run() error {
	IntrRun()
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
	IntrShutdown()
	util.Infof("shutting down")
}

func Init() error {
	err := IpInit()
	if err != nil {
		return fmt.Errorf("IpInit: %v", err)
	}
	util.Infof("initialized")
	return nil
}

type NetProtocol struct {
	nptype  NetProtocolType
	q       *util.ConcurrentQueue[NetProtocolQueueEntry]
	handler protocolHandler
}

type NetProtocolType int
type protocolHandler func(data []byte, len int, dev NetDevice)

var protocols = make([]*NetProtocol, 0)

const (
	DummyProtocol NetProtocolType = iota + 1
	IpProtocol
	ArpProtocol
	IpV6Protocol
)

func RegisterNetProtocol(nptype NetProtocolType, handler protocolHandler) error {
	for _, proto := range protocols {
		if nptype == proto.nptype {
			return fmt.Errorf("protocol already registered, type=0x%04x", nptype)
		}
	}
	proto := NetProtocol{
		nptype:  nptype,
		q:       util.NewConcurrentQueue[NetProtocolQueueEntry](),
		handler: handler,
	}
	protocols = append(protocols, &proto)
	util.Infof("protocol registered, type=0x%04x", nptype)
	return nil
}

type NetProtocolQueueEntry struct {
	dev  NetDevice
	len  int
	data []byte
}

func SoftIrqHandler() error {
	for _, proto := range protocols {
		for {
			if proto.q.IsEmpty() {
				break
			}
			entry := proto.q.Dequeue()
			util.Debugf("dequeued (num:%d), dev=%s, protocol=0x%04x, len=%d", proto.q.Len(), entry.dev.Info().Name, proto.nptype, entry.len)
			util.Infof("received: %+v", entry.data)
			proto.handler(entry.data, entry.len, entry.dev)
		}
	}
	return nil
}
