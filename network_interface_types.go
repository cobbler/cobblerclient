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

// NetworkInterfaceType enumerates the kinds of interface configurations supported by Cobbler. Mirrors
// cobbler.enums.NetworkInterfaceType. This is a plain string (not a Go int-backed enum) because the XML-RPC codec
// (kolo/xmlrpc) picks the wire representation purely from reflect.Kind — an int-kind type is always encoded as
// <int>, which cobbler/items/network_interface.py's interface_type setter rejects (it only accepts a string name
// or an enum instance, never a raw int). Every other Cobbler enum-like field in this client (VirtType, PowerType,
// ...) is a plain string for the same reason.
type NetworkInterfaceType = string

const (
	NetworkInterfaceTypeNA                = "na"
	NetworkInterfaceTypeBond              = "bond"
	NetworkInterfaceTypeBondSlave         = "bond_slave"
	NetworkInterfaceTypeBridge            = "bridge"
	NetworkInterfaceTypeBridgeSlave       = "bridge_slave"
	NetworkInterfaceTypeBondedBridgeSlave = "bonded_bridge_slave"
	NetworkInterfaceTypeBmc               = "bmc"
	NetworkInterfaceTypeInfiniband        = "infiniband"
)

// IPv4Option models the per-interface IPv4 configuration. There is no
// "gateway" field here — the real field is NetworkInterface.IfGateway
// (cobbler/items/network_interface.py: if_gateway).
type IPv4Option struct {
	Address      string   `mapstructure:"address" json:"address" yaml:"address"`
	Netmask      string   `mapstructure:"netmask" json:"netmask" yaml:"netmask"`
	StaticRoutes []string `mapstructure:"static_routes" json:"static_routes" yaml:"static_routes"`
}

// IPv6Option models the per-interface IPv6 configuration. There is no
// "default_gateway" field here — the real field is
// NetworkInterface.Ipv6DefaultGateway (cobbler/items/network_interface.py:
// ipv6_default_gateway), which is a separate, independently backed property.
type IPv6Option struct {
	Address      string   `mapstructure:"address" json:"address" yaml:"address"`
	Prefix       string   `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Secondaries  []string `mapstructure:"secondaries" json:"secondaries" yaml:"secondaries"`
	MTU          string   `mapstructure:"mtu" json:"mtu" yaml:"mtu"`
	StaticRoutes []string `mapstructure:"static_routes" json:"static_routes" yaml:"static_routes"`
}

// DNSInterfaceOption models the per-interface DNS configuration.
type DNSInterfaceOption struct {
	Name   string   `mapstructure:"name" json:"name" yaml:"name"`
	CNames []string `mapstructure:"cnames" json:"cnames" yaml:"cnames"`
}
