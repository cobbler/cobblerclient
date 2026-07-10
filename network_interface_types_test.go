/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

import "testing"

func TestNetworkInterfaceTypeString(t *testing.T) {
	cases := []struct {
		input    NetworkInterfaceType
		expected string
	}{
		{NetworkInterfaceTypeNA, "na"},
		{NetworkInterfaceTypeBond, "bond"},
		{NetworkInterfaceTypeBondSlave, "bond_slave"},
		{NetworkInterfaceTypeBridge, "bridge"},
		{NetworkInterfaceTypeBridgeSlave, "bridge_slave"},
		{NetworkInterfaceTypeBondedBridgeSlave, "bonded_bridge_slave"},
		{NetworkInterfaceTypeInfiniband, "infiniband"},
	}
	for _, tc := range cases {
		if got := tc.input.String(); got != tc.expected {
			t.Errorf("NetworkInterfaceType(%d).String() = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestIPv4OptionRoundtrip(t *testing.T) {
	o := IPv4Option{
		Address:      "192.168.1.10",
		Netmask:      "255.255.255.0",
		Gateway:      "192.168.1.1",
		StaticRoutes: []string{"10.0.0.0/8 via 192.168.1.1"},
	}
	if len(o.StaticRoutes) != 1 {
		t.Errorf("StaticRoutes length = %d, want 1", len(o.StaticRoutes))
	}
}
