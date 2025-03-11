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
	"log"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/go-viper/mapstructure/v2"
)

// System is a created system.
type System struct {
	Item `mapstructure:",squash" yaml:",inline"`

	// These are internal fields and cannot be modified.
	IPv6Autoconfiguration bool            `mapstructure:"ipv6_autoconfiguration" cobbler:"noupdate" json:"ipv6_autoconfiguration" yaml:"ipv6_autoconfiguration"`
	ReposEnabled          bool            `mapstructure:"repos_enabled"          cobbler:"noupdate" json:"repos_enabled" yaml:"repos_enabled"`
	Autoinstall           string          `mapstructure:"autoinstall" json:"autoinstall" yaml:"autoinstall"`
	BootLoaders           Value[[]string] `mapstructure:"boot_loaders" json:"boot_loaders" yaml:"boot_loaders"`
	EnableIPXE            Value[bool]     `mapstructure:"enable_ipxe" json:"enable_ipxe" yaml:"enable_ipxe"`
	Filename              string          `mapstructure:"filename" json:"filename" yaml:"filename"`
	Gateway               string          `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	Hostname              string          `mapstructure:"hostname" json:"hostname" yaml:"hostname"`
	IPv6DefaultDevice     string          `mapstructure:"ipv6_default_device" json:"ipv6_default_device" yaml:"ipv6_default_device"`
	Image                 string          `mapstructure:"image" json:"image" yaml:"image"`
	Interfaces            Interfaces      `mapstructure:"interfaces" cobbler:"noupdate" json:"interfaces" yaml:"interfaces"`
	NameServers           []string        `mapstructure:"name_servers" json:"name_servers" yaml:"name_servers"`
	NameServersSearch     []string        `mapstructure:"name_servers_search" json:"name_servers_search" yaml:"name_servers_search"`
	NetbootEnabled        bool            `mapstructure:"netboot_enabled" json:"netboot_enabled" yaml:"netboot_enabled"`
	NextServerv4          string          `mapstructure:"next_server_v4" json:"next_server_v4" yaml:"next_server_v4"`
	NextServerv6          string          `mapstructure:"next_server_v6" json:"next_server_v6" yaml:"next_server_v6"`
	PowerAddress          string          `mapstructure:"power_address" json:"power_address" yaml:"power_address"`
	PowerID               string          `mapstructure:"power_id" json:"power_id" yaml:"power_id"`
	PowerIdentityFile     string          `mapstructure:"power_identity_file" json:"power_identity_file" yaml:"power_identity_file"`
	PowerOptions          string          `mapstructure:"power_options" json:"power_options" yaml:"power_options"`
	PowerPass             string          `mapstructure:"power_pass" json:"power_pass" yaml:"power_pass"`
	PowerType             string          `mapstructure:"power_type" json:"power_type" yaml:"power_type"`
	PowerUser             string          `mapstructure:"power_user" json:"power_user" yaml:"power_user"`
	Profile               string          `mapstructure:"profile" json:"profile" yaml:"profile"`
	Proxy                 string          `mapstructure:"proxy" json:"proxy" yaml:"proxy"`
	RedhatManagementKey   string          `mapstructure:"redhat_management_key" json:"redhat_management_key" yaml:"redhat_management_key"`
	SerialBaudRate        int             `mapstructure:"serial_baud_rate" json:"serial_baud_rate" yaml:"serial_baud_rate"`
	SerialDevice          int             `mapstructure:"serial_device" json:"serial_device" yaml:"serial_device"`
	Server                string          `mapstructure:"server" json:"server" yaml:"server"`
	Status                string          `mapstructure:"status" json:"status" yaml:"status"`
	VirtAutoBoot          Value[bool]     `mapstructure:"virt_auto_boot" json:"virt_auto_boot" yaml:"virt_auto_boot"`
	VirtCPUs              Value[int]      `mapstructure:"virt_cpus" json:"virt_cpus" yaml:"virt_cpus"`
	VirtDiskDriver        string          `mapstructure:"virt_disk_driver" json:"virt_disk_driver" yaml:"virt_disk_driver"`
	VirtFileSize          Value[float64]  `mapstructure:"virt_file_size" json:"virt_file_size" yaml:"virt_file_size"`
	VirtPXEBoot           bool            `mapstructure:"virt_pxe_boot" json:"virt_pxe_boot" yaml:"virt_pxe_boot"`
	VirtPath              string          `mapstructure:"virt_path" json:"virt_path" yaml:"virt_path"`
	VirtRAM               Value[int]      `mapstructure:"virt_ram" json:"virt_ram" yaml:"virt_ram"`
	VirtType              string          `mapstructure:"virt_type" json:"virt_type" yaml:"virt_type"`

	Client
}

// Interface is an interface in a system.
type Interface struct {
	BondingOpts        string   `mapstructure:"bonding_opts" structs:"bonding_opts" json:"bonding_opts" yaml:"bonding_opts"`
	BridgeOpts         string   `mapstructure:"bridge_opts" structs:"bridge_opts" json:"bridge_opts" yaml:"bridge_opts"`
	CNAMEs             []string `mapstructure:"cnames" structs:"cnames" json:"cnames" yaml:"cnames"`
	ConnectedMode      bool     `mapstructure:"connected_mode" structs:"connected_mode" json:"connected_mode" yaml:"connected_mode"`
	DHCPTag            string   `mapstructure:"dhcp_tag" structs:"dhcp_tag" json:"dhcp_tag" yaml:"dhcp_tag"`
	DNSName            string   `mapstructure:"dns_name" structs:"dns_name" json:"dns_name" yaml:"dns_name"`
	Gateway            string   `mapstructure:"if_gateway" structs:"if_gateway" json:"gateway" yaml:"gateway"`
	IPAddress          string   `mapstructure:"ip_address" structs:"ip_address" json:"ip_address" yaml:"ip_address"`
	IPv6Address        string   `mapstructure:"ipv6_address" structs:"ipv6_address" json:"ipv6_address" yaml:"ipv6_address"`
	IPv6DefaultGateway string   `mapstructure:"ipv6_default_gateway" structs:"ipv6_default_gateway" json:"ipv6_default_gateway" yaml:"ipv6_default_gateway"`
	IPv6MTU            string   `mapstructure:"ipv6_mtu" structs:"ipv6_mtu" json:"ipv6_mtu" yaml:"ipv6_mtu"`
	IPv6Prefix         string   `mapstructure:"ipv6_prefix" structs:"ipv6_prefix" json:"ipv6_prefix" yaml:"ipv6_prefix"`
	IPv6Secondaries    []string `mapstructure:"ipv6_secondaries" structs:"ipv6_secondaries" json:"ipv6_secondaries" yaml:"ipv6_secondaries"`
	IPv6StaticRoutes   []string `mapstructure:"ipv6_static_routes" structs:"ipv6_static_routes" json:"ipv6_static_routes" yaml:"ipv6_static_routes"`
	InterfaceMaster    string   `mapstructure:"interface_master" structs:"interface_master" json:"interface_master" yaml:"interface_master"`
	InterfaceType      string   `mapstructure:"interface_type" structs:"interface_type" json:"interface_type" yaml:"interface_type"`
	MACAddress         string   `mapstructure:"mac_address" structs:"mac_address" json:"mac_address" yaml:"mac_address"`
	MTU                string   `mapstructure:"mtu" structs:"mtu" json:"mtu" yaml:"mtu"`
	Management         bool     `mapstructure:"management" structs:"management" json:"management" yaml:"management"`
	Netmask            string   `mapstructure:"netmask" structs:"netmask" json:"netmask" yaml:"netmask"`
	Static             bool     `mapstructure:"static" structs:"static" json:"static" yaml:"static"`
	StaticRoutes       []string `mapstructure:"static_routes" structs:"static_routes" json:"static_routes" yaml:"static_routes"`
	VirtBridge         string   `mapstructure:"virt_bridge" structs:"virt_bridge" json:"virt_bridge" yaml:"virt_bridge"`
}

// Interfaces is a collection of interfaces in a system.
type Interfaces map[string]Interface

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
		Filename:            inherit,
		Interfaces:          make(map[string]Interface),
		NameServers:         make([]string, 0),
		NameServersSearch:   make([]string, 0),
		NetbootEnabled:      false,
		NextServerv4:        inherit,
		NextServerv6:        inherit,
		Proxy:               inherit,
		RedhatManagementKey: inherit,
		SerialBaudRate:      -1,
		SerialDevice:        -1,
		Server:              inherit,
		VirtAutoBoot: Value[bool]{
			IsInherited: true,
		},
		VirtCPUs: Value[int]{
			IsInherited: true,
		},
		VirtDiskDriver: inherit,
		VirtFileSize: Value[float64]{
			IsInherited: true,
		},
		VirtPath: inherit,
		VirtRAM: Value[int]{
			IsInherited: true,
		},
		VirtType: inherit,
	}
	// Overwrite defaults from Item
	system.Owners = Value[[]string]{
		IsInherited: true,
	}
	system.BootFiles = Value[map[string]interface{}]{
		IsInherited: true,
	}
	system.FetchableFiles = Value[map[string]interface{}]{
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
	system.MgmtClasses = Value[[]string]{
		IsInherited: true,
	}
	return system
}

func NewInterface() Interface {
	return Interface{
		InterfaceType:    "na",
		CNAMEs:           make([]string, 0),
		IPv6Secondaries:  make([]string, 0),
		IPv6StaticRoutes: make([]string, 0),
		StaticRoutes:     make([]string, 0),
	}
}

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
	s.Client = *c

	// Now clean the network interface struct
	if s.Interfaces == nil {
		s.Interfaces = make(map[string]Interface)
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
	err = sanitizeValueMapStruct(&s.FetchableFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&s.BootFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&s.TemplateFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&s.MgmtParameters)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&s.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&s.MgmtClasses)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&s.BootLoaders)
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
	if len(system.BootFiles.Data) == 0 {
		system.BootFiles.IsInherited = true
	}

	if len(system.BootLoaders.Data) == 0 {
		system.BootLoaders.IsInherited = true
	}

	if len(system.FetchableFiles.Data) == 0 {
		system.FetchableFiles.IsInherited = true
	}

	if len(system.MgmtParameters.Data) == 0 {
		system.MgmtParameters.IsInherited = true
	}

	if system.PowerType == "" {
		system.PowerType = "ipmilanplus"
	}

	if system.Status == "" {
		system.Status = "production"
	}

	if system.VirtDiskDriver == "" {
		system.VirtDiskDriver = inherit
	}

	if system.VirtPath == "" {
		system.VirtPath = inherit
	}

	if system.VirtType == "" {
		system.VirtType = inherit
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
	if err := c.SaveSystem(newID, "new"); err != nil {
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
func (c *Client) SaveSystem(objectId, editmode string) error {
	_, err := c.Call("save_system", objectId, c.Token, editmode)
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

func makeInterfaceOptionsMap(name string, iface Interface) map[string]interface{} {
	i := structs.Map(iface)
	nic := make(map[string]interface{})
	for key, value := range i {
		attrName := fmt.Sprintf("%s-%s", key, name)
		log.Printf("[DEBUG] Cobblerclient: setting interface attr %s to %s", attrName, value)
		nic[attrName] = value
	}
	return nic
}

func (c *Client) ModifyInterface(systemID string, nic map[string]interface{}) error {
	editUncasted, err := c.Call("modify_system", systemID, "modify_interface", nic, c.Token)
	if err != nil {
		return err
	}
	ok, editSuccessful := editUncasted.(bool)
	if !ok {
		return fmt.Errorf("editing interfaces of system %s failed due to an invalid return value of the server", systemID)
	}
	if !editSuccessful {
		return fmt.Errorf("editing interface of system %s failed", systemID)
	}
	return nil
}

// CreateInterface creates network interfaces in Cobbler
func (s *System) CreateInterface(name string, iface Interface) error {
	nic := makeInterfaceOptionsMap(name, iface)

	systemID, err := s.Client.GetItemHandle("system", s.Name)
	if err != nil {
		return err
	}

	err = s.Client.ModifyInterface(systemID, nic)
	if err != nil {
		return err
	}

	// Save the final system
	err = s.Client.SaveSystem(systemID, "bypass")
	if err != nil {
		return err
	}

	return nil
}

// ModifyNetworkInterface updates the attributes of an existing interface
func (s *System) ModifyNetworkInterface(name string, iface Interface) error {
	return s.CreateInterface(name, iface)
}

// GetInterfaces returns all interfaces in a System.
func (s *System) GetInterfaces() (Interfaces, error) {
	nics := make(Interfaces)
	for nicName, nicData := range s.Interfaces {
		var nic Interface
		if err := mapstructure.Decode(nicData, &nic); err != nil {
			return nil, err
		}
		nics[nicName] = nic
	}

	return nics, nil
}

// GetInterface returns a single interface in a System.
func (s *System) GetInterface(name string) (Interface, error) {
	nics := make(Interfaces)
	var iface Interface
	for nicName, nicData := range s.Interfaces {
		var nic Interface
		if err := mapstructure.Decode(nicData, &nic); err != nil {
			return iface, err
		}
		nics[nicName] = nic
	}

	if iface, ok := nics[name]; ok {
		return iface, nil
	} else {
		return iface, fmt.Errorf("interface %s not found", name)
	}
}

func (c *Client) DeleteNetworkInterface(systemID, name string) error {
	editUncasted, err := c.Call("modify_system", systemID, "delete_interface", name, c.Token)
	if err != nil {
		return err
	}
	ok, editSuccessful := editUncasted.(bool)
	if !ok {
		return fmt.Errorf("deleting interfaces of system %s failed due to an invalid return value of the server", systemID)
	}
	if !editSuccessful {
		return fmt.Errorf("deleting interface of system %s failed", systemID)
	}
	return nil
}

// DeleteInterface deletes a single interface in a System.
func (s *System) DeleteInterface(name string) error {
	if _, err := s.GetInterface(name); err != nil {
		return err
	}

	systemID, err := s.Client.GetItemHandle("system", s.Name)
	if err != nil {
		return err
	}

	if err = s.Client.DeleteNetworkInterface(systemID, name); err != nil {
		return err
	}

	// Save the final system
	if err = s.Client.SaveSystem(systemID, "bypass"); err != nil {
		return err
	}

	return nil
}

func (c *Client) RenameNetworkInterface(systemName, oldName, newName string) error {
	systemID, err := c.GetItemHandle("system", systemName)
	if err != nil {
		return err
	}
	args := make(map[string]string)
	args["interface"] = oldName
	args["rename_interface"] = newName
	unparsedOk, err := c.Call("modify_system", systemID, "rename_interface", args, c.Token)
	if err != nil {
		return err
	}
	conversionSuccess, ok := unparsedOk.(bool)
	if !conversionSuccess {
		return fmt.Errorf("failed to convert return value of interface rename")
	}
	if !ok {
		return fmt.Errorf("failed to rename interface %s to %s", oldName, newName)
	}
	return nil
}

// RenameInterface renames an Interface.
func (s *System) RenameInterface(name string, newName string) error {
	err := s.Client.RenameNetworkInterface(s.Name, name, newName)
	if err != nil {
		return err
	}
	iface := s.Interfaces[name]
	delete(s.Interfaces, name)
	s.Interfaces[newName] = iface
	return nil
}

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
	res, err := c.Call("get_system_handle", name, c.Token)
	return returnString(res, err)
}
