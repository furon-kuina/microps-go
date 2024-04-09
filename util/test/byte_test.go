package test

import (
	"testing"

	"github.com/furon-kuina/microps-go/util"
)

func TestConcat2Uint8(t *testing.T) {
	var a uint8 = 0x1F
	var b uint8 = 0x0A
	var expected uint16 = 0x1F0A
	res := util.Concat2Uint8(a, b)
	if res != expected {
		t.Errorf("didn't match: expected %d, got %d", expected, res)
	}
}

func TestConcat4Uint8(t *testing.T) {
	var a uint8 = 0x0F
	var b uint8 = 0x0A
	var c uint8 = 0xCD
	var d uint8 = 0x23
	var expected uint32 = 0x0F0ACD23
	res := util.Concat4Uint8(a, b, c, d)
	if res != expected {
		t.Errorf("didn't match: expected %d, got %d", expected, res)
	}
}
