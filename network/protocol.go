package network

import (
	"fmt"

	"github.com/furon-kuina/microps-go/internet"
	"github.com/furon-kuina/microps-go/util"
)

type ProtocolManager struct {
	protocols map[internet.InternetProtocolType]protocol
}

func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{
		protocols: make(map[internet.InternetProtocolType]protocol),
	}
}

type internetHandler func(dev Device, data []byte)

func (pm *ProtocolManager) Register(ptype internet.InternetProtocolType, handler internetHandler) error {
	_, ok := pm.protocols[ptype]
	if ok {
		return fmt.Errorf("already registered")
	}
	pm.protocols[ptype] = protocol{
		handler: handler,
		queue:   util.NewConcurrentQueue[ProtocolQueueEntry](),
	}
	return nil
}

type protocol struct {
	handler internetHandler
	queue   *util.ConcurrentQueue[ProtocolQueueEntry]
}

type ProtocolQueueEntry struct {
	dev  Device
	data []byte
}

func IpInputHandler(dev Device, data []byte) {
	util.Debugf("IpInputHandler: %s", string(data))
}
