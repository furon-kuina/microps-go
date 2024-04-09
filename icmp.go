package net

// import (
// 	"encoding/binary"
// 	"fmt"
// )

// type IcmpHeader struct {
// 	Type            byte
// 	Code            byte
// 	Checksum        uint16
// 	MessageSpecific uint32
// }

// type IcmpPacket struct {
// 	Header IcmpHeader
// 	Data   []byte
// }

// const (
// 	IcmpHeaderLenByte = 8
// )

// func (iface IpInterface) IcmpInput(data []byte, src, dst IpAddress) error {
// 	icmpPacket, err := ParseIcmpData(data)
// 	if err != nil {
// 		return err
// 	}
// 	PrintIcmp(icmpPacket)
// 	return nil
// }

// // ms はネットワークバイトオーダー
// func IcmpOutput(ty, code byte, ms uint32, data []byte, src, dst IpAddress) {
// }

// func ParseIcmpData(data []byte) (IcmpPacket, error) {
// 	if len(data) < IcmpHeaderLenByte {
// 		return IcmpPacket{}, ErrTooShort
// 	}
// 	tmpData := data
// 	tmpData[2] = 0
// 	tmpData[3] = 0
// 	if Cksum16(tmpData) != binary.BigEndian.Uint16(data[2:4]) {
// 		return IcmpPacket{}, ErrWrongChecksum
// 	}
// 	return IcmpPacket{
// 		Header: IcmpHeader{
// 			Type:            data[0],
// 			Code:            data[1],
// 			Checksum:        binary.BigEndian.Uint16(data[2:4]),
// 			MessageSpecific: binary.BigEndian.Uint32(data[4:8]),
// 		},
// 		Data: data[8:],
// 	}, nil
// }

// func PrintIcmp(icmp IcmpPacket) {
// 	fmt.Printf("type: %x\n", icmp.Header.Type)
// 	fmt.Printf("code: %x\n", icmp.Header.Code)
// 	fmt.Printf("sum : %x\n", icmp.Header.Checksum)
// }
