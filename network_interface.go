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

import (
	"fmt"
	"reflect"
	"time"
)

// NetworkInterface is a Cobbler 4.0.0 first-class item type representing a
// network interface attached to a System. Get the fields from
// cobbler/items/network_interface.py.
type NetworkInterface struct {
	Item `mapstructure:",squash" yaml:",inline"`

	// Parent system reference. Set at creation time via new_network_interface(system_uid, token).
	SystemUid  string `mapstructure:"system_uid" cobbler:"noupdate" json:"system_uid" yaml:"system_uid"`
	SystemName string `mapstructure:"system_name" cobbler:"noupdate" json:"system_name" yaml:"system_name"`

	// Layer 2 / link
	MacAddress      string               `mapstructure:"mac_address" json:"mac_address" yaml:"mac_address"`
	InterfaceType   NetworkInterfaceType `mapstructure:"interface_type" json:"interface_type" yaml:"interface_type"`
	InterfaceMaster string               `mapstructure:"interface_master" json:"interface_master" yaml:"interface_master"`
	BondingOpts     string               `mapstructure:"bonding_opts" json:"bonding_opts" yaml:"bonding_opts"`
	BridgeOpts      string               `mapstructure:"bridge_opts" json:"bridge_opts" yaml:"bridge_opts"`
	ConnectedMode   bool                 `mapstructure:"connected_mode" structs:"connected_mode" json:"connected_mode" yaml:"connected_mode"`
	Management      bool                 `mapstructure:"management" structs:"management" json:"management" yaml:"management"`
	Static          bool                 `mapstructure:"static" structs:"static" json:"static" yaml:"static"`
	DHCPTag         string               `mapstructure:"dhcp_tag" json:"dhcp_tag" yaml:"dhcp_tag"`
	IfGateway       string               `mapstructure:"if_gateway" json:"if_gateway" yaml:"if_gateway"`
	MTU             string               `mapstructure:"mtu" json:"mtu" yaml:"mtu"`

	// These live top-level on NetworkInterface in Python (not nested under
	// IPv6), backed by their own private attributes independent of
	// IPv6Option.StaticRoutes. See cobbler/items/network_interface.py:
	// ipv6_default_gateway, ipv6_static_routes.
	Ipv6DefaultGateway string   `mapstructure:"ipv6_default_gateway" json:"ipv6_default_gateway" yaml:"ipv6_default_gateway"`
	Ipv6StaticRoutes   []string `mapstructure:"ipv6_static_routes" json:"ipv6_static_routes" yaml:"ipv6_static_routes"`

	// Inheritable
	VirtBridge Value[string] `mapstructure:"virt_bridge" json:"virt_bridge" yaml:"virt_bridge"`

	// Layer 3 / addressing — nested option objects
	IPv4 IPv4Option         `mapstructure:"ipv4" json:"ipv4" yaml:"ipv4"`
	IPv6 IPv6Option         `mapstructure:"ipv6" json:"ipv6" yaml:"ipv6"`
	DNS  DNSInterfaceOption `mapstructure:"dns" json:"dns" yaml:"dns"`
}

// NewNetworkInterface returns a zero-valued NetworkInterface with sensible defaults.
func NewNetworkInterface() NetworkInterface {
	return NetworkInterface{
		Item:             NewItem(),
		InterfaceType:    NetworkInterfaceTypeNA,
		IPv4:             IPv4Option{StaticRoutes: []string{}},
		IPv6:             IPv6Option{Secondaries: []string{}, StaticRoutes: []string{}},
		DNS:              DNSInterfaceOption{CNames: []string{}},
		Ipv6StaticRoutes: []string{},
		VirtBridge:       Value[string]{IsInherited: true},
	}
}

// convertRawNetworkInterface decodes a raw XML-RPC result into a NetworkInterface struct.
func convertRawNetworkInterface(name string, xmlrpcResult interface{}) (*NetworkInterface, error) {
	var ni NetworkInterface
	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("network interface %s not found", name)
	}
	decoded, err := decodeCobblerItem(xmlrpcResult, &ni)
	if err != nil {
		return nil, err
	}
	result, ok := decoded.(*NetworkInterface)
	if !ok {
		return nil, fmt.Errorf("unexpected decoder result type %T", decoded)
	}
	return result, nil
}

// convertRawNetworkInterfacesList converts a list of raw XML-RPC results.
func convertRawNetworkInterfacesList(xmlrpcResult interface{}) ([]*NetworkInterface, error) {
	var ifaces []*NetworkInterface
	for _, raw := range xmlrpcResult.([]interface{}) {
		ni, err := convertRawNetworkInterface("unknown", raw)
		if err != nil {
			return nil, err
		}
		ifaces = append(ifaces, ni)
	}
	return ifaces, nil
}

// GetNetworkInterfaces returns all network interfaces known to the backend.
func (c *Client) GetNetworkInterfaces() ([]*NetworkInterface, error) {
	result, err := c.Call("get_network_interfaces", "-1", c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawNetworkInterfacesList(result)
}

// GetNetworkInterface returns a single network interface by its uid.
func (c *Client) GetNetworkInterface(uid string, flattened, resolved bool) (*NetworkInterface, error) {
	result, err := c.getConcreteItem("get_network_interface", uid, flattened, resolved)
	if err != nil {
		return nil, err
	}
	return convertRawNetworkInterface(uid, result)
}

// GetNetworkInterfaceHandle returns the in-memory handle for a network interface.
func (c *Client) GetNetworkInterfaceHandle(name string) (string, error) {
	return c.GetItemHandle("network_interface", name)
}

// GetNetworkInterfacesSince returns network interfaces modified at or after mtime.
func (c *Client) GetNetworkInterfacesSince(mtime time.Time) ([]*NetworkInterface, error) {
	result, err := c.Call("get_network_interfaces_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}
	return convertRawNetworkInterfacesList(result)
}

// FindNetworkInterface searches for network interfaces matching the criteria.
func (c *Client) FindNetworkInterface(criteria map[string]interface{}, resolved bool) ([]*NetworkInterface, error) {
	result, err := c.Call("find_network_interface", criteria, true, resolved, c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawNetworkInterfacesList(result)
}

// FindNetworkInterfaceNames returns the names of network interfaces matching the criteria.
func (c *Client) FindNetworkInterfaceNames(criteria map[string]interface{}) ([]string, error) {
	result, err := c.Call("find_network_interface", criteria, false, false, c.Token)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, n := range result.([]interface{}) {
		names = append(names, n.(string))
	}
	return names, nil
}

// CreateNetworkInterface creates a new network interface attached to systemUid
// and persists the field values supplied via iface.
func (c *Client) CreateNetworkInterface(systemUid string, iface NetworkInterface) (*NetworkInterface, error) {
	// Make sure a network interface with the same name does not already exist
	if exists, err := c.HasItem("network_interface", iface.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("a NetworkInterface with the name %s already exists", iface.Name)
	}

	id, err := c.Call("new_network_interface", systemUid, c.Token)
	if err != nil {
		return nil, err
	}
	objectID, ok := id.(string)
	if !ok {
		return nil, fmt.Errorf("new_network_interface returned %T, want string", id)
	}
	if err := c.updateCobblerFields("network_interface", reflect.ValueOf(&iface).Elem(), objectID); err != nil {
		return nil, err
	}
	if err := c.SaveNetworkInterface(objectID, true, true, "new"); err != nil {
		return nil, err
	}
	return c.GetNetworkInterface(objectID, false, false)
}

// UpdateNetworkInterface persists changes to an existing network interface.
func (c *Client) UpdateNetworkInterface(iface *NetworkInterface) error {
	objectID, err := c.GetNetworkInterfaceHandle(iface.Name)
	if err != nil {
		return err
	}
	if err := c.updateCobblerFields("network_interface", reflect.ValueOf(iface).Elem(), objectID); err != nil {
		return err
	}
	return c.SaveNetworkInterface(objectID, true, true, "bypass")
}

// SaveNetworkInterface flushes in-memory changes to the backend.
func (c *Client) SaveNetworkInterface(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_network_interface", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

// CopyNetworkInterface copies an interface to a new name.
func (c *Client) CopyNetworkInterface(objectId, newName string) error {
	_, err := c.Call("copy_network_interface", objectId, newName, c.Token)
	return err
}

// DeleteNetworkInterface removes a network interface by its uid.
func (c *Client) DeleteNetworkInterface(uid string) error {
	_, err := c.Call("remove_network_interface", uid, c.Token, false)
	return err
}

// DeleteNetworkInterfaceRecursive removes a network interface by its uid, optionally cascading.
func (c *Client) DeleteNetworkInterfaceRecursive(uid string, recursive bool) error {
	_, err := c.Call("remove_network_interface", uid, c.Token, recursive)
	return err
}

// RenameNetworkInterface renames an interface.
func (c *Client) RenameNetworkInterface(objectId, newName string) error {
	_, err := c.Call("rename_network_interface", objectId, newName, c.Token)
	return err
}

// ListNetworkInterfaceNames returns all interface names known to the backend.
func (c *Client) ListNetworkInterfaceNames() ([]string, error) {
	return c.GetItemNames("network_interface")
}
