package net

import (
	"fmt"

	"github.com/furon-kuina/microps-go/util"
)

func IpInit() error {
	err := RegisterNetProtocol(IpProtocol, IpHandler)
	if err != nil {
		return fmt.Errorf("IpInit: %v", err)
	}
	return nil
}

func IpHandler(data []byte, len int, dev NetDevice) {
	util.Debugf("dev=%s, len=%d, data=%+v", dev, len, data)
}
