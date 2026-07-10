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

// System is a created system.
// Get the fields from cobbler/items/system.py
type System struct {
	Item `mapstructure:",squash" yaml:",inline"`

	// These are internal fields and cannot be modified.
	IPv6Autoconfiguration    bool                         `mapstructure:"ipv6_autoconfiguration" cobbler:"noupdate" json:"ipv6_autoconfiguration" yaml:"ipv6_autoconfiguration"`
	ReposEnabled             bool                         `mapstructure:"repos_enabled"          cobbler:"noupdate" json:"repos_enabled" yaml:"repos_enabled"`
	Autoinstall              string                       `mapstructure:"autoinstall" json:"autoinstall" yaml:"autoinstall"`
	BootLoaders              Value[[]string]              `mapstructure:"boot_loaders" json:"boot_loaders" yaml:"boot_loaders"`
	DNS                      DNSOptions                   `mapstructure:"dns" json:"dns" yaml:"dns"`
	EnableIPXE               Value[bool]                  `mapstructure:"enable_ipxe" json:"enable_ipxe" yaml:"enable_ipxe"`
	Filename                 string                       `mapstructure:"filename" json:"filename" yaml:"filename"`
	Gateway                  string                       `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	Hostname                 string                       `mapstructure:"hostname" json:"hostname" yaml:"hostname"`
	IPv6DefaultDevice        string                       `mapstructure:"ipv6_default_device" json:"ipv6_default_device" yaml:"ipv6_default_device"`
	Image                    string                       `mapstructure:"image" json:"image" yaml:"image"`
	Interfaces               map[string]*NetworkInterface `mapstructure:"interfaces" cobbler:"noupdate" json:"interfaces" yaml:"interfaces"`
	NetbootEnabled           bool                         `mapstructure:"netboot_enabled" json:"netboot_enabled" yaml:"netboot_enabled"`
	Power                    PowerOptions                 `mapstructure:"power" json:"power" yaml:"power"`
	Profile                  string                       `mapstructure:"profile" json:"profile" yaml:"profile"`
	Proxy                    string                       `mapstructure:"proxy" json:"proxy" yaml:"proxy"`
	RedhatManagementKey      string                       `mapstructure:"redhat_management_key" json:"redhat_management_key" yaml:"redhat_management_key"`
	RedhatManagementOrg      string                       `mapstructure:"redhat_management_org" json:"redhat_management_org" yaml:"redhat_management_org"`
	RedhatManagementUser     string                       `mapstructure:"redhat_management_user" json:"redhat_management_user" yaml:"redhat_management_user"`
	RedhatManagementPassword string                       `mapstructure:"redhat_management_password" json:"redhat_management_password" yaml:"redhat_management_password"`
	SerialBaudRate           int                          `mapstructure:"serial_baud_rate" json:"serial_baud_rate" yaml:"serial_baud_rate"`
	SerialDevice             int                          `mapstructure:"serial_device" json:"serial_device" yaml:"serial_device"`
	Server                   string                       `mapstructure:"server" json:"server" yaml:"server"`
	Status                   string                       `mapstructure:"status" json:"status" yaml:"status"`
	TFTP                     TFTPOptions                  `mapstructure:"tftp" json:"tftp" yaml:"tftp"`
	Virt                     VirtOptions                  `mapstructure:"virt" json:"virt" yaml:"virt"`
	VirtPXEBoot              bool                         `mapstructure:"virt_pxe_boot" json:"virt_pxe_boot" yaml:"virt_pxe_boot"`
}

// Interface type removed in 4.0.0 - network interfaces are now first-class items.
// Use NetworkInterface type and dedicated CRUD methods instead.
// System.Interfaces field now maps to map[string]*NetworkInterface.

func NewSystem() System {
	system := System{
		Item:        NewItem(),
		Autoinstall: inherit,
		BootLoaders: Value[[]string]{
			IsInherited: true,
		},
		EnableIPXE: Value[bool]{
			IsInherited: true,
		},
		DNS:                      newDNSOptions(),
		Filename:                 inherit,
		Interfaces:               make(map[string]*NetworkInterface),
		NetbootEnabled:           false,
		Proxy:                    inherit,
		RedhatManagementKey:      inherit,
		RedhatManagementOrg:      inherit,
		RedhatManagementUser:     inherit,
		RedhatManagementPassword: inherit,
		SerialBaudRate:           -1,
		SerialDevice:             -1,
		Server:                   inherit,
		TFTP:                     newTFTPOptions(),
		Virt:                     newVirtOptions(),
	}
	// Overwrite defaults from Item
	system.Owners = Value[[]string]{
		IsInherited: true,
	}
	system.AutoinstallMeta = Value[map[string]interface{}]{
		IsInherited: true,
	}
	system.KernelOptions = Value[map[string]interface{}]{
		IsInherited: true,
	}
	system.KernelOptionsPost = Value[map[string]interface{}]{
		IsInherited: true,
	}
	return system
}

// NewInterface removed in 4.0.0 - use NewNetworkInterface() instead.

func (c *Client) convertRawSystem(name string, xmlrpcResult interface{}) (*System, error) {
	var system System

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("system %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &system)
	if err != nil {
		return nil, err
	}

	s := decodeResult.(*System)

	// Now clean the network interface struct
	if s.Interfaces == nil {
		s.Interfaces = make(map[string]*NetworkInterface)
	}
	// Now clean the Value structs
	err = sanitizeValueMapStruct(&s.KernelOptions)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&s.KernelOptionsPost)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&s.AutoinstallMeta)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&s.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&s.BootLoaders)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&s.DNS.NameServers)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Client) convertRawSystemsList(xmlrpcResult interface{}) ([]*System, error) {
	var systems []*System

	for _, s := range xmlrpcResult.([]interface{}) {
		system, err := c.convertRawSystem("unkown", s)
		if err != nil {
			return nil, err
		}
		system.Meta = ItemMeta{
			IsFlattened: false,
			IsResolved:  false,
		}
		systems = append(systems, system)
	}

	return systems, nil
}

// GetSystems returns all systems in Cobbler.
func (c *Client) GetSystems() ([]*System, error) {

	result, err := c.Call("get_systems", "", c.Token)
	if err != nil {
		return nil, err
	}

	return c.convertRawSystemsList(result)
}

// GetSystem returns a single system obtained by its name.
func (c *Client) GetSystem(name string, flattened, resolved bool) (*System, error) {
	result, err := c.getConcreteItem("get_system", name, flattened, resolved)
	if err != nil {
		return nil, err
	}

	system, err := c.convertRawSystem(name, result)
	if err != nil {
		return nil, err
	}
	system.Meta = ItemMeta{
		IsFlattened: flattened,
		IsResolved:  resolved,
	}
	return system, nil
}

// CreateSystem creates a system.
// It ensures that either a Profile or Image are set and then sets other default values.
func (c *Client) CreateSystem(system System) (*System, error) {
	// Check if a system with the same name already exists
	if _, err := c.GetSystem(system.Name, false, false); err == nil {
		return nil, fmt.Errorf("a system with the name %s already exists", system.Name)
	}

	if system.Profile == "" && system.Image == "" {
		return nil, fmt.Errorf("a system must have a profile or image set")
	}

	// Set default values. I guess these aren't taken care of by Cobbler?
	if len(system.BootLoaders.Data) == 0 {
		system.BootLoaders.IsInherited = true
	}

	if system.Power.Type == "" {
		system.Power.Type = "ipmilanplus"
	}

	if system.Status == "" {
		system.Status = "production"
	}

	if system.Virt.DiskDriver == "" {
		system.Virt.DiskDriver = inherit
	}

	if system.Virt.Path == "" {
		system.Virt.Path = inherit
	}

	if system.Virt.Type == "" {
		system.Virt.Type = inherit
	}

	// To create a system via the Cobbler API, first call new_system to obtain an ID
	result, err := c.Call("new_system", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	// Set the value of all fields
	item := reflect.ValueOf(&system).Elem()
	if err := c.updateCobblerFields("system", item, newID); err != nil {
		return nil, err
	}

	// Save the final system
	if err := c.SaveSystem(newID, true, true, "new"); err != nil {
		return nil, err
	}

	// Return a clean copy of the system
	return c.GetSystem(system.Name, false, false)
}

// UpdateSystem updates a single system.
func (c *Client) UpdateSystem(system *System) error {
	item := reflect.ValueOf(system).Elem()
	id, err := c.GetItemHandle("system", system.Name)
	if err != nil {
		return err
	}
	return c.updateCobblerFields("system", item, id)
}

// SaveSystem saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveSystem(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_system", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

// CopySystem duplicates a given system on the server with a new name.
func (c *Client) CopySystem(objectId, newName string) error {
	_, err := c.Call("copy_system", objectId, newName, c.Token)
	return err
}

// DeleteSystem deletes a single System by its name.
func (c *Client) DeleteSystem(name string) error {
	return c.DeleteSystemRecursive(name, false)
}

// DeleteSystemRecursive deletes a single System by its name with the option to do so recursively.
func (c *Client) DeleteSystemRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_system", name, c.Token, recursive)
	return err
}

// GetInterfaces and GetInterface removed in 4.0.0.
// Access System.Interfaces directly (now map[string]*NetworkInterface).
// For CRUD operations, use Client.{Create,Update,Delete,Rename}NetworkInterface methods.

// ListSystemNames returns a list of all system names currently available in Cobbler.
func (c *Client) ListSystemNames() ([]string, error) {
	return c.GetItemNames("system")
}

// FindSystem searches for one or more systems by any of its attributes.
func (c *Client) FindSystem(criteria map[string]interface{}) ([]*System, error) {
	result, err := c.Call("find_system", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return c.convertRawSystemsList(result)
}

// FindSystemNames searches for one or more systems by any of its attributes.
func (c *Client) FindSystemNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_system", criteria, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetSystemsSince returns all systems which were created after the specified date.
func (c *Client) GetSystemsSince(mtime time.Time) ([]*System, error) {
	result, err := c.Call("get_systems_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return c.convertRawSystemsList(result)
}

// RenameSystem renames a System with a given object id.
func (c *Client) RenameSystem(objectId, newName string) error {
	_, err := c.Call("rename_system", objectId, newName, c.Token)
	return err
}

// GetSystemHandle gets the internal ID of a Cobbler item.
func (c *Client) GetSystemHandle(name string) (string, error) {
	res, err := c.Call("get_system_handle", name)
	return returnString(res, err)
}

// DisableNetboot disables PXE booting for the named system (pxe_just_once feature).
func (c *Client) DisableNetboot(name string) error {
	_, err := c.Call("disable_netboot", name, c.Token)
	return err
}

// UploadLogData uploads Anaconda log data for the named system (anamon logging).
func (c *Client) UploadLogData(sysName, file string, size, offset int, data string) (bool, error) {
	result, err := c.Call("upload_log_data", sysName, file, size, offset, data, c.Token)
	return returnBool(result, err)
}

// ClearSystemLogs clears the console logs for the system identified by its object ID.
func (c *Client) ClearSystemLogs(objectId string) (bool, error) {
	result, err := c.Call("clear_system_logs", objectId, c.Token)
	return returnBool(result, err)
}

// GetValidSystemBootLoaders retrieves the list of bootloaders that can be assigned to a system.
func (c *Client) GetValidSystemBootLoaders(systemName string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_system_boot_loaders", systemName, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetSystemAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetSystemAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_system_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}
