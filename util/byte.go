package util

import "slices"

func Concat2Uint8(x, y uint8) uint16 {
	return (uint16(x) << 8) + uint16(y)
}

func Concat4Uint8(a, b, c, d uint8) uint32 {
	return uint32(Concat2Uint8(a, b))<<16 + uint32(Concat2Uint8(c, d))
}

func OnesComplementAdd(x, y uint16) uint16 {
	s := uint32(x) + uint32(y)
	s += (s >> 16)
	s &= 0xFFFFFFFF ^ (1 << 16)
	return uint16(s)
}

// []byteをuint16の配列とみなしチェックサムを計算する
// len(data)が奇数の場合は末尾に0を追加して計算する
func Cksum16(data []byte) uint16 {
	dataCopy := slices.Clone(data)
	if len(data)%2 != 0 {
		dataCopy = append(dataCopy, 0)
	}
	var sum uint16 = 0
	for i := 0; i < len(dataCopy); i += 2 {
		tmp := Concat2Uint8(data[i], data[i+1])
		sum = OnesComplementAdd(sum, tmp)
	}
	return sum
}
