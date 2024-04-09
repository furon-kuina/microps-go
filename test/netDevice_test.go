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

var dummyData = []byte("Hello, world!")

func TestNetDevice(t *testing.T) {
	err := net.Init()
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	dev := driver.NewDummyDevice()
	if err = net.Run(); err != nil {
		t.Errorf("run failed: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)
	complete := time.After(20 * time.Second)

	go func() {
		for {
			select {
			case <-sigs:
				fmt.Println("gracefully shutting down...")
				done <- true
				return
			case <-ticker.C:
				err := net.Output(dev, net.DummyProtocol, dummyData, nil)
				if err != nil {
					t.Errorf("transmit failed: %v", err)
				}
			case <-complete:
				done <- true
				return
			}
		}
	}()
	<-done
	ticker.Stop()
	net.Shutdown()
}
