package net

import (
	"fmt"
	"sync"

	"github.com/furon-kuina/microps-go/util"
)

type irqEntry struct {
	irq     Irq
	handler IrqHandler
	shared  bool
	name    string
	dev     Device
}

type Irq int
type IrqHandler func(Irq, Device) error
type safeIrqEntries struct {
	entries []irqEntry
	mu      sync.Mutex
}

var (
	irqEntries = &safeIrqEntries{
		entries: []irqEntry{},
		mu:      sync.Mutex{},
	}
	irqChan       = make(chan Irq, 100)
	terminateChan = make(chan interface{})
	wg            sync.WaitGroup
)

func (irqs *safeIrqEntries) append(irq irqEntry) {
	irqs.mu.Lock()
	defer irqs.mu.Unlock()
	irqs.entries = append(irqs.entries, irq)
}

func (irqs *safeIrqEntries) getEntries() []irqEntry {
	irqs.mu.Lock()
	defer irqs.mu.Unlock()
	return irqs.entries
}

func IntrRun() {
	wg.Add(1)
	go runIrqHandler()
	wg.Wait()
}

func IntrShutdown() {
	terminateChan <- struct{}{}
}

func IntrRaiseIrq(irq Irq) {
	util.Debugf("raise irq %d", irq)
	irqChan <- irq
}

func runIrqHandler() {
	wg.Done()
loop:
	for {
		select {
		case <-terminateChan:
			break loop
		case irq := <-irqChan:
			util.Debugf("captured irq %d", irq)
			if irq == SoftIrq {
				SoftIrqHandler()
				break
			}
			entries := irqEntries.getEntries()
			for _, entry := range entries {
				if entry.irq == irq {
					util.Debugf("irq=%d, name=%s", entry.irq, entry.name)
					entry.handler(entry.irq, entry.dev)
				}
			}
		}
	}
	util.Debugf("terminated")
}

func RegisterIrqHandler(irq Irq, handler IrqHandler, shared bool, name string, dev Device) error {
	entries := irqEntries.getEntries()
	for _, entry := range entries {
		if entry.irq != irq {
			continue
		}
		if entry.shared || shared {
			return fmt.Errorf("IRQ conflict")
		}
	}

	entry := irqEntry{
		irq:     irq,
		handler: handler,
		shared:  shared,
		name:    name,
		dev:     dev,
	}

	irqEntries.append(entry)
	util.Debugf("registered: irq=%d, name=%s", irq, name)

	return nil
}
