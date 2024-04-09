package network

import (
	"testing"
	"time"

	"github.com/furon-kuina/microps-go/internet"
)

func TestLoopback(t *testing.T) {
	testData := []byte("Hello, world!")
	ns := NewNetworkStack()
	dev := ns.NewLoopbackDevice()
	ns.RegisterDevice(dev)
	ns.RegisterIrq("LOOPBACK", LoopbackIsr, false)
	ns.RegisterProtocol(internet.IpProtocol, IpInputHandler)

	if err := ns.Run(); err != nil {
		t.Errorf("net.Run: %v", err)
	}

	ticker := time.NewTicker(1 * time.Second)
	complete := time.After(20 * time.Second)
	done := make(chan interface{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := ns.Output(dev, internet.IpProtocol, testData, nil); err != nil {
					t.Errorf("Output(%q): %v", dev.Name, err)
				}
			case <-complete:
				done <- struct{}{}
				return
			}
		}
	}()
	<-done
	ticker.Stop()
	ns.Shutdown()
}
