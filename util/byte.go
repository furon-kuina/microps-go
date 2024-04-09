package util

func Concat2Uint8(x, y uint8) uint16 {
	return (uint16(x) << 8) + uint16(y)
}

func Concat4Uint8(a, b, c, d uint8) uint32 {
	return uint32(Concat2Uint8(a, b))<<16 + uint32(Concat2Uint8(c, d))
}
