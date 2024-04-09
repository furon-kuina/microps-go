package network

import (
	"fmt"
	"sync"

	"github.com/furon-kuina/microps-go/util"
)

type Irq struct {
	name string
	dev  Device
}
type IrqHandler func(irq Irq) error

const (
	DummyIrq = "DUMMY"
	SoftIrq  = "SOFT"
)

// 1つのIRQ IDに紐づくIRQたち
type IrqEntry struct {
	handlers []IrqHandler
	shared   bool
}

type IrqManager struct {
	irqs map[string]*IrqEntry
	mu   sync.Mutex

	IrqChan       chan Irq
	TerminateChan chan interface{}
	wg            sync.WaitGroup
}

func NewIrqManager() *IrqManager {
	im := &IrqManager{
		mu:            sync.Mutex{},
		irqs:          make(map[string]*IrqEntry),
		IrqChan:       make(chan Irq),
		TerminateChan: make(chan interface{}),
	}
	return im
}

func (im *IrqManager) Run() {
	im.wg.Add(1)
	go im.run()
	util.Debugf("waiting...")
	im.wg.Wait()
}

func (im *IrqManager) Shutdown() {
	im.TerminateChan <- struct{}{}
}

func (im *IrqManager) Register(irqName string, handler IrqHandler, shared bool) error {
	im.mu.Lock()
	defer im.mu.Unlock()
	entries, ok := im.irqs[irqName]
	if ok {
		if !entries.shared || shared {
			return fmt.Errorf("sharing not allowed: %s", irqName)
		}
		entries.handlers = append(entries.handlers, handler)
		entries.shared = entries.shared && shared
	} else {
		im.irqs[irqName] = &IrqEntry{
			handlers: []IrqHandler{handler},
			shared:   shared,
		}
	}
	return nil
}

func (im *IrqManager) run() {
	im.wg.Done()
	util.Debugf("starting IRQ goroutine")
loop:
	for {
		util.Debugf("waiting interrupt...")
		select {
		case <-im.TerminateChan:
			break loop
		case irq := <-im.IrqChan:
			util.Debugf("captured IRQ %s", irq.name)
			entries, ok := im.irqs[irq.name]
			if !ok {
				util.Debugf("IRQ handler not found: %s", irq)
				break
			}
			for _, handler := range entries.handlers {
				go handler(irq)
			}
		}
	}
	util.Debugf("terminated")
}
