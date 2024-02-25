package test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	net "github.com/furon-kuina/microps-go"
	"github.com/furon-kuina/microps-go/driver"
)

var LoopbackData = []byte{
	0x45, 0x00, 0x00, 0x30,
	0x00, 0x80, 0x00, 0x00,
	0xff, 0x01, 0xbd, 0x4a,
	0x7f, 0x00, 0x00, 0x01,
	0x7f, 0x00, 0x00, 0x01,
	0x08, 0x00, 0x35, 0x64,
	0x00, 0x80, 0x00, 0x01,
	0x31, 0x32, 0x33, 0x34,
	0x35, 0x36, 0x37, 0x38,
	0x39, 0x30, 0x21, 0x40,
	0x23, 0x24, 0x25, 0x5e,
	0x26, 0x2a, 0x28, 0x29,
}

func TestLoopback(t *testing.T) {
	err := net.Init()
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	dev := driver.NewLoopbackDevice()
	if err = net.Run(); err != nil {
		t.Errorf("run failed: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-sigs:
				fmt.Println("gracefully shutting down...")
				done <- true
				return
			case <-ticker.C:
				err := net.Output(dev, net.DummyProtocol, LoopbackData, len(LoopbackData), dev)
				if err != nil {
					t.Errorf("transmit failed: %v", err)
				}
			}
		}
	}()
	<-done
	ticker.Stop()
	net.Shutdown()
}
