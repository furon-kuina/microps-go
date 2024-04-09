package net

import (
	"fmt"
	"sync"

	"github.com/furon-kuina/microps-go/util"
)

// デバイス関連

type NetDeviceType string

const (
	DummyDevice    = NetDeviceType("DUMMY")
	LoopbackDevice = NetDeviceType("LOOPBACK")
	EthernetDevice = NetDeviceType("ETHERNET")
)

type Device interface {
	Info() *DeviceInfo
	Interfaces() *NetworkInterfaceManager
	Open() error
	Close() error
	Transmit(ProtocolType, []byte, int, Device) error
}

type DeviceInfo struct {
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

func RegisterDevice(dev Device) (*DeviceInfo, error) {
	info := dev.Info()
	newInfo := info
	newInfo.index = deviceCounter.GetValue()
	deviceCounter.Increment()
	newInfo.Name = fmt.Sprintf("net%d", newInfo.index)
	devices = append(devices, dev)
	return newInfo, nil
}

func (device DeviceInfo) State() string {
	if device.Flags&UpFlag != 0 {
		return "up"
	} else {
		return "down"
	}
}

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

type NetworkDeviceManager struct {
	devices []Device
}

func NewNetworkDeviceManager() NetworkDeviceManager {
	return NetworkDeviceManager{}
}

var (
	deviceCounter = NewCounter()
	devices       []Device
	c             = sync.NewCond(&sync.Mutex{})
	irqReady      = false
)

// dev:    どのNICからデータが届いたか
// nptype: 届いたデータのプロトコル
func InputHandler(dev Device, nptype ProtocolType, data []byte, len int) error {
	info := dev.Info()
	util.Debugf("dev=%s, type=%s, len=%d", info.Name, nptype, len)
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

func Open(nd Device) (err error) {
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

func Close(nd Device) (err error) {
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

func Output(nd Device, nptype ProtocolType, data []byte, dst Device) (err error) {
	info := nd.Info()
	defer util.Wrap(&err, "Output(%q)", info.Name)
	if info.Flags&UpFlag == 0 {
		return fmt.Errorf("not open")
	}
	if len(data) > info.Mtu {
		return fmt.Errorf("too long, mtu=%d, len=%d", info.Mtu, len(data))
	}
	util.Infof("transmitting data %q", string(data))
	if err = nd.Transmit(nptype, data, len(data), dst); err != nil {
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
	// if err != nil {
	// 	return fmt.Errorf("IpInit: %v", err)
	// }
	// util.Infof("initialized")
	// return nil
	return nil
}

type NetProtocol struct {
	nptype  ProtocolType
	q       *util.ConcurrentQueue[NetProtocolQueueEntry]
	handler protocolHandler
}

// プロトコル関連

type ProtocolType string
type protocolHandler func(data []byte, len int, dev Device)

var protocols = make([]*NetProtocol, 0)

const (
	DummyProtocol = ProtocolType("DUMMY")
	IpProtocol    = ProtocolType("IP")
	ArpProtocol   = ProtocolType("ARP")
	IpV6Protocol  = ProtocolType("IPV6")
)

func RegisterNetProtocol(nptype ProtocolType, handler protocolHandler) error {
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
	dev  Device
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
