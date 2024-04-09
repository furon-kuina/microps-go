package internet

// func Test_Ip(t *testing.T) {
// 	dev := driver.NewLoopbackDevice()
// 	iface, err := NewIpInterface(LoopbackIpAddress, LoopbackNetmask)
// 	if err != nil {
// 		t.Fatalf("NewIpInterface: %v", err)
// 	}
// 	RegisterIpInterface(dev, iface)
// }

// func Test_ParseIpHeader(t *testing.T) {
// 	var _ IpHeader

// }

// func Test_ToIpAddress(t *testing.T) {
// 	cases := map[string]struct {
// 		input string
// 		res   IpAddress
// 		err   error
// 	}{
// 		"loopback":       {"127.0.0.1", 0x7F000001, nil},
// 		"netmask":        {"255.0.0.0", 0xFF000000, nil},
// 		"invalid format": {"1270.0.1", 0, ErrInvalidFormat},
// 		"invalid value":  {"256.0.0.1", 0, ErrInvalidValue},
// 	}
// 	for name, tc := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			res, err := ToIpAddress(tc.input)
// 			if res != tc.res {
// 				t.Errorf("IpAddressStringToUint32(%q): expected %d, got %d", tc.input, tc.res, res)
// 			}
// 			if !errors.Is(err, tc.err) {
// 				t.Errorf("IpAddressStringToUint32(%q): expected %v, got %v", tc.input, tc.err, err)
// 			}
// 		})
// 	}
// }

// func Test_ToIpString(t *testing.T) {
// 	cases := []struct {
// 		name string
// 		give IpAddress
// 		want string
// 	}{
// 		{"broadcast", 0xFFFFFFFF, "255.255.255.255"},
// 		{"loopback", 0x7F000001, "127.0.0.1"},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			res := ToIpString(tc.give)
// 			if res != tc.want {
// 				t.Errorf("ToIpString(%d): expected %s, got %s", tc.give, tc.want, res)
// 			}
// 		})
// 	}
// }
