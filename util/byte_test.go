package util

import "testing"

func TestOnesComplementAdd(t *testing.T) {
	cases := []struct {
		giveX uint16
		giveY uint16
		want  uint16
	}{
		{0xFFFF, 0x1011, 0x1011},
		{0x0005, 0xFFFE, 0x0004},
		{0x0500, 0xFEFF, 0x0400},
	}
	for _, tc := range cases {
		res := OnesComplementAdd(tc.giveX, tc.giveY)
		if res != tc.want {
			t.Errorf("OnesComplementAdd(%x, %x): expected %x, got %x", tc.giveX, tc.giveY, tc.want, res)
		}
	}
}
