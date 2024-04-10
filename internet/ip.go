package internet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/furon-kuina/microps-go/util"
	"github.com/m-mizutani/goerr"
)

const (
	IpVersionIpV4       = 4
	IpHeaderSizeMinByte = 20
	IpHeaderSizeMaxByte = 60
)

const (
	mfBit      = 0x2000
	offsetBits = 0x1FFF
)

type TransportProtocol uint8

const ()

// func IpInit() error {
// 	err := net.RegisterNetProtocol(net.IpProtocol, IpHandler)
// 	if err != nil {
// 		return fmt.Errorf("IpInit: %v", err)
// 	}
// 	return nil
// }

// func IpHandler(data []byte, len int, dev net.Device) {
// 	util.Debugf("dev=%s, len=%d, data=%+v", dev, len, data)
// }

type IpAddress uint32

type IpHeader struct {
	Vhl      uint8
	Tos      uint8
	Len      uint16
	Id       uint16
	Offset   uint16
	Ttl      uint8
	Protocol TransportProtocol
	Sum      uint16
	Src      IpAddress
	Dst      IpAddress
}

type IpDatagram struct {
	Header  IpHeader
	Payload []byte
}

func ParseIpHeader(data []uint8) (*IpDatagram, error) {
	if len(data) < IpHeaderSizeMinByte {
		return nil, goerr.New("too short")
	}
	hdr := IpHeader{}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	if hdr.Vhl>>4 != IpVersionIpV4 {
		return nil, fmt.Errorf("unsupported version")
	}
	hlen := int((hdr.Vhl & 0x0F) << 2)
	if len(data) < hlen {
		return nil, goerr.New("too short")
	}
	cksum := util.Cksum16(data)
	if ^cksum != 0 {
		return nil, fmt.Errorf("checksum error: got %x", cksum)
	}
	if len(data) != int(hdr.Len) {
		return nil, fmt.Errorf("wrong length")
	}
	if hdr.Ttl == 0 {
		return nil, fmt.Errorf("dead packet")
	}
	if (hdr.Offset&mfBit) != 0 && (hdr.Offset&offsetBits) != 0 {
		return nil, goerr.New("fragment not supported")
	}

	return &IpDatagram{
		Header:  hdr,
		Payload: data[hlen:hdr.Len],
	}, nil
}

func PrintIpDatagram(datagram IpDatagram) {
	hdr := datagram.Header
	fmt.Println("IP Header")
	fmt.Printf("version: %1x\n", (hdr.Vhl >> 4))
	fmt.Printf("header length: %d bytes\n", (hdr.Vhl&0x0F)<<2)
	fmt.Printf("type of service: %2x\n", hdr.Tos)
	fmt.Printf("total length: %4x\n", hdr.Len)
	fmt.Printf("ID: %4x\n", hdr.Id)
	fmt.Printf("offset: %b\n", hdr.Offset)
	fmt.Printf("time to live: %2x\n", hdr.Ttl)
	fmt.Printf("protocol: %2x\n", hdr.Protocol)
	fmt.Printf("checksum: %4x\n", hdr.Sum)
	fmt.Printf("source address: %s\n", ToIpString(hdr.Src))
	fmt.Printf("destination address: %s\n", ToIpString(hdr.Dst))
	fmt.Printf("data: %s\n", string(datagram.Payload))
}

// func PrintIpDatagram(datagram IpDatagram) {
// }

// // IpInput handles data input to dev
// func HandleIpInput(data []uint8, dev net.Device) error {
// 	if len(data) < IpHeaderSizeMinByte {
// 		return fmt.Errorf("len too short")
// 	}
// 	ipDatagram, err := ParseIpHeader(data)
// 	if err != nil {
// 		return err
// 	}
// 	for _, iface := range dev.Interfaces().interfaces {
// 		ipIface, ok := iface.(IpInterface)
// 		if !ok {
// 			continue
// 		}
// 		if ipIface.dev.Info().index == dev.Info().index {
// 			acceptIp := []IpAddress{broadcast, ipIface.unicast, ipIface.broadcast}
// 			// dropping packet not in acceptIp
// 			if slices.Contains(acceptIp, ipDatagram.Header.Dst) {
// 				return nil
// 			}
// 			util.Debugf("dev=%s, iface=%s", dev.Info().Name, ToIpString(ipIface.unicast))
// 		}
// 	}
// 	return nil
// }

// func IpOutput(protocol uint8, data []byte, src, dst IpAddress) {

// }

// type IpInterface struct {
// 	dev       net.Device
// 	family    net.NetworkInterfaceFamily
// 	unicast   IpAddress
// 	netmask   IpAddress
// 	broadcast IpAddress
// }

// func (ii IpInterface) Dev() net.Device {
// 	return ii.dev
// }

// func (ii IpInterface) Family() net.NetworkInterfaceFamily {
// 	return ii.family
// }

// func NewIpInterface(unicastStr, netmaskStr string) (*IpInterface, error) {
// 	unicast, err := ToIpAddress(unicastStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	netmask, err := ToIpAddress(netmaskStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &IpInterface{
// 		family:    net.Ip,
// 		unicast:   unicast,
// 		netmask:   netmask,
// 		broadcast: (unicast & netmask) | ^netmask,
// 	}, nil
// }

// const (
// 	LoopbackIpAddress = "127.0.0.1"
// 	LoopbackNetmask   = "255.0.0.0"
// )

// var (
// 	ErrInvalidFormat = errors.New("invalid format")
// 	ErrInvalidValue  = errors.New("invalid value")
// )

// func ToIpAddress(addrStr string) (IpAddress, error) {
// 	bytes := strings.Split(addrStr, ".")
// 	if len(bytes) != 4 {
// 		return 0, fmt.Errorf("AddressStringToUint32(%q): %w", addrStr, ErrInvalidFormat)
// 	}
// 	var res uint32 = 0
// 	for i, s := range bytes {
// 		num, err := strconv.Atoi(s)
// 		if err != nil {
// 			return 0, err
// 		}
// 		if num < 0 || 255 < num {
// 			return 0, fmt.Errorf("AddressStringToUint32(%q): %w", addrStr, ErrInvalidValue)
// 		}
// 		res += uint32(num)
// 		if i < 3 {
// 			res <<= 8
// 		}
// 	}
// 	return IpAddress(res), nil
// }

func ToIpString(addr IpAddress) string {
	var res string
	var mask IpAddress = 0xFF000000
	for i := range 4 {
		res += strconv.Itoa(int((mask & addr) >> ((3 - i) * 8)))
		if i == 3 {
			break
		}
		res += "."
		mask >>= 8
	}
	return res
}

// var ipInterfaces []IpInterface

// func RegisterIpInterface(dev net.Device, ipIface IpInterface) error {
// 	if err := dev.Interfaces().AddInterface(ipIface); err != nil {
// 		return err
// 	}
// 	ipInterfaces = append(ipInterfaces, ipIface)
// 	return nil
// }

// func SelectIpInterface(address IpAddress) (IpInterface, error) {
// 	for _, ipIface := range ipInterfaces {
// 		if ipIface.unicast == address {
// 			return ipIface, nil
// 		}
// 	}
// 	return IpInterface{}, fmt.Errorf("not found")
// }

// func ntoh32(data []byte) uint32 {
// 	return binary.BigEndian.Uint32(data)
// }

// func ntoh16(data []byte) uint16 {
// 	return binary.BigEndian.Uint16(data)
// }
