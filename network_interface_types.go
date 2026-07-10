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

// NetworkInterfaceType enumerates the kinds of interface configurations
// supported by Cobbler. Mirrors cobbler.enums.NetworkInterfaceType.
type NetworkInterfaceType int

const (
	NetworkInterfaceTypeNA NetworkInterfaceType = iota
	NetworkInterfaceTypeBond
	NetworkInterfaceTypeBondSlave
	NetworkInterfaceTypeBridge
	NetworkInterfaceTypeBridgeSlave
	NetworkInterfaceTypeBondedBridgeSlave
	NetworkInterfaceTypeInfiniband
)

func (t NetworkInterfaceType) String() string {
	switch t {
	case NetworkInterfaceTypeNA:
		return "na"
	case NetworkInterfaceTypeBond:
		return "bond"
	case NetworkInterfaceTypeBondSlave:
		return "bond_slave"
	case NetworkInterfaceTypeBridge:
		return "bridge"
	case NetworkInterfaceTypeBridgeSlave:
		return "bridge_slave"
	case NetworkInterfaceTypeBondedBridgeSlave:
		return "bonded_bridge_slave"
	case NetworkInterfaceTypeInfiniband:
		return "infiniband"
	}
	return "na"
}

// IPv4Option models the per-interface IPv4 configuration.
type IPv4Option struct {
	Address      string   `mapstructure:"address" json:"address" yaml:"address"`
	Netmask      string   `mapstructure:"netmask" json:"netmask" yaml:"netmask"`
	Gateway      string   `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	StaticRoutes []string `mapstructure:"static_routes" json:"static_routes" yaml:"static_routes"`
}

// IPv6Option models the per-interface IPv6 configuration.
type IPv6Option struct {
	Address        string   `mapstructure:"address" json:"address" yaml:"address"`
	Prefix         string   `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Secondaries    []string `mapstructure:"secondaries" json:"secondaries" yaml:"secondaries"`
	MTU            string   `mapstructure:"mtu" json:"mtu" yaml:"mtu"`
	StaticRoutes   []string `mapstructure:"static_routes" json:"static_routes" yaml:"static_routes"`
	DefaultGateway string   `mapstructure:"default_gateway" json:"default_gateway" yaml:"default_gateway"`
}

// DNSInterfaceOption models the per-interface DNS configuration.
type DNSInterfaceOption struct {
	Name   string   `mapstructure:"name" json:"name" yaml:"name"`
	CNames []string `mapstructure:"cnames" json:"cnames" yaml:"cnames"`
}
