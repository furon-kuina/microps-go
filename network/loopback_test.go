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

func TestLoopbackWithIp(t *testing.T) {
	testData := []byte{
		0x45, 0x00, 0x00, 0x30,
		0x00, 0x80, 0x00, 0x00,
		0xff, 0x01, 0x69, 0x9A,
		0x7f, 0x00, 0x00, 0x01,
		0x7f, 0x00, 0x00, 0x01,
		0x61, 0x62, 0x63, 0x64,
		0x65, 0x66, 0x67, 0x68,
		0x31, 0x32, 0x33, 0x34,
		0x35, 0x36, 0x37, 0x38,
		0x39, 0x30, 0x21, 0x40,
		0x23, 0x24, 0x25, 0x5e,
		0x26, 0x2a, 0x28, 0x29,
	}
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
