package net

import "testing"

func TestChecksum16bit(t *testing.T) {
	cases := []struct {
		give []byte
		want uint16
	}{
		{give: []byte{0x00, 0x00, 0x00, 0x00}, want: 0xFFFF},
	}

	for _, tc := range cases {
		res := Cksum16(tc.give)
		if res != tc.want {
			t.Errorf("Checksum16bit(%v): got %d, want %d", tc.give, res, tc.want)
		}
	}
}
