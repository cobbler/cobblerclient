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
	// NetworkInterfaceType is a plain string alias (see network_interface_types.go) so that the XML-RPC codec,
	// which picks its wire representation from reflect.Kind, sends these as <string> rather than <int> — this just
	// pins the wire values the constants must keep sending.
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
		{NetworkInterfaceTypeBmc, "bmc"},
		{NetworkInterfaceTypeInfiniband, "infiniband"},
	}
	for _, tc := range cases {
		if tc.input != tc.expected {
			t.Errorf("NetworkInterfaceType = %q, want %q", tc.input, tc.expected)
		}
	}
}

func TestIPv4OptionRoundtrip(t *testing.T) {
	o := IPv4Option{
		Address:      "192.168.1.10",
		Netmask:      "255.255.255.0",
		StaticRoutes: []string{"10.0.0.0/8 via 192.168.1.1"},
	}
	if len(o.StaticRoutes) != 1 {
		t.Errorf("StaticRoutes length = %d, want 1", len(o.StaticRoutes))
	}
}
