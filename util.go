package net

import "encoding/binary"

func Cksum16(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i:]))
		sum = (sum >> 16) + (sum & 0xFFFF)
	}
	if len(data)%2 != 0 {
		sum += uint32(data[len(data)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xFFFF)
	return ^uint16(sum)
}
