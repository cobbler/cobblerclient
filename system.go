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
	"github.com/mitchellh/mapstructure"
)

// System is a created system.
type System struct {
	// These are internal fields and cannot be modified.
	Ctime                 float64                `mapstructure:"ctime"                 cobbler:"noupdate"` // TODO: convert to time
	Depth                 int                    `mapstructure:"depth"                 cobbler:"noupdate"`
	ID                    string                 `mapstructure:"uid"                   cobbler:"noupdate"`
	IPv6Autoconfiguration bool                   `mapstructure:"ipv6_autoconfiguration" cobbler:"noupdate"`
	Mtime                 float64                `mapstructure:"mtime"                 cobbler:"noupdate"` // TODO: convert to time
	ReposEnabled          bool                   `mapstructure:"repos_enabled"          cobbler:"noupdate"`
	Autoinstall           string                 `mapstructure:"autoinstall"`
	AutoinstallMeta       []string               `mapstructure:"autoinstall_meta"`
	BootFiles             string                 `mapstructure:"boot_files"`
	BootLoaders           []string               `mapstructure:"boot_loaders"`
	Comment               string                 `mapstructure:"comment"`
	EnableGPXE            bool                   `mapstructure:"enable_gpxe"`
	FetchableFiles        []string               `mapstructure:"fetchable_files"`
	Gateway               string                 `mapstructure:"gateway"`
	Hostname              string                 `mapstructure:"hostname"`
	Image                 string                 `mapstructure:"image"`
	Interfaces            map[string]interface{} `mapstructure:"interfaces" cobbler:"noupdate"`
	IPv6DefaultDevice     string                 `mapstructure:"ipv6_default_device"`
	KernelOptions         []string               `mapstructure:"kernel_options"`
	KernelOptionsPost     []string               `mapstructure:"kernel_options_post"`
	MGMTClasses           []string               `mapstructure:"mgmt_classes"`
	MGMTParameters        string                 `mapstructure:"mgmt_parameters"`
	Name                  string                 `mapstructure:"name"`
	NameServers           []string               `mapstructure:"name_servers"`
	NameServersSearch     []string               `mapstructure:"name_servers_search"`
	NetbootEnabled        bool                   `mapstructure:"netboot_enabled"`
	NextServerv4          string                 `mapstructure:"next_server_v4"`
	NextServerv6          string                 `mapstructure:"next_server_v6"`
	Owners                []string               `mapstructure:"owners"`
	PowerAddress          string                 `mapstructure:"power_address"`
	PowerID               string                 `mapstructure:"power_id"`
	PowerPass             string                 `mapstructure:"power_pass"`
	PowerType             string                 `mapstructure:"power_type"`
	PowerUser             string                 `mapstructure:"power_user"`
	Profile               string                 `mapstructure:"profile"`
	Proxy                 string                 `mapstructure:"proxy"`
	Status                string                 `mapstructure:"status"`
	TemplateFiles         []string               `mapstructure:"template_files"`
	VirtAutoBoot          string                 `mapstructure:"virt_auto_boot"`
	VirtCPUs              string                 `mapstructure:"virt_cpus"`
	VirtDiskDriver        string                 `mapstructure:"virt_disk_driver"`
	VirtFileSize          string                 `mapstructure:"virt_file_size"`
	VirtPath              string                 `mapstructure:"virt_path"`
	VirtPXEBoot           int                    `mapstructure:"virt_pxe_boot"`
	VirtRAM               string                 `mapstructure:"virt_ram"`
	VirtType              string                 `mapstructure:"virt_type"`

	Client
}

// Interface is an interface in a system.
type Interface struct {
	CNAMEs             []string `mapstructure:"cnames" structs:"cnames"`
	DHCPTag            string   `mapstructure:"dhcp_tag" structs:"dhcp_tag"`
	DNSName            string   `mapstructure:"dns_name" structs:"dns_name"`
	BondingOpts        string   `mapstructure:"bonding_opts" structs:"bonding_opts"`
	BridgeOpts         string   `mapstructure:"bridge_opts" structs:"bridge_opts"`
	Gateway            string   `mapstructure:"if_gateway" structs:"if_gateway"`
	InterfaceType      string   `mapstructure:"interface_type" structs:"interface_type"`
	InterfaceMaster    string   `mapstructure:"interface_master" structs:"interface_master"`
	IPAddress          string   `mapstructure:"ip_address" structs:"ip_address"`
	IPv6Address        string   `mapstructure:"ipv6_address" structs:"ipv6_address"`
	IPv6Secondaries    []string `mapstructure:"ipv6_secondaries" structs:"ipv6_secondaries"`
	IPv6MTU            string   `mapstructure:"ipv6_mtu" structs:"ipv6_mtu"`
	IPv6StaticRoutes   []string `mapstructure:"ipv6_static_routes" structs:"ipv6_static_routes"`
	IPv6DefaultGateway string   `mapstructure:"ipv6_default_gateway" structs:"ipv6_default_gateway"`
	MACAddress         string   `mapstructure:"mac_address" structs:"mac_address"`
	Management         bool     `mapstructure:"management" structs:"management"`
	Netmask            string   `mapstructure:"netmask" structs:"netmask"`
	Static             bool     `mapstructure:"static" structs:"static"`
	StaticRoutes       []string `mapstructure:"static_routes" structs:"static_routes"`
	VirtBridge         string   `mapstructure:"virt_bridge" structs:"virt_bridge"`
}

// Interfaces is a collection of interfaces in a system.
type Interfaces map[string]Interface

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

	return s, nil
}

func (c *Client) convertRawSystemsList(xmlrpcResult interface{}) ([]*System, error) {
	var systems []*System

	for _, s := range xmlrpcResult.([]interface{}) {
		system, err := c.convertRawSystem("unkown", s)
		if err != nil {
			return nil, err
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
func (c *Client) GetSystem(name string) (*System, error) {
	var system System

	result, err := c.Call("get_system", name, c.Token)
	if err != nil {
		return &system, err
	}

	return c.convertRawSystem(name, result)
}

// CreateSystem creates a system.
// It ensures that either a Profile or Image are set and then sets other default values.
func (c *Client) CreateSystem(system System) (*System, error) {
	// Check if a system with the same name already exists
	if _, err := c.GetSystem(system.Name); err == nil {
		return nil, fmt.Errorf("a system with the name %s already exists", system.Name)
	}

	if system.Profile == "" && system.Image == "" {
		return nil, fmt.Errorf("a system must have a profile or image set")
	}

	// Set default values. I guess these aren't taken care of by Cobbler?
	if system.BootFiles == "" {
		system.BootFiles = "<<inherit>>"
	}

	if len(system.BootLoaders) == 0 {
		system.BootLoaders = []string{"<<inherit>>"}
	}

	if len(system.BootLoaders) == 0 {
		system.FetchableFiles = []string{"<<inherit>>"}
	}

	if system.MGMTParameters == "" {
		system.MGMTParameters = "<<inherit>>"
	}

	if system.PowerType == "" {
		system.PowerType = "ipmilan"
	}

	if system.Status == "" {
		system.Status = "production"
	}

	if system.VirtAutoBoot == "" {
		system.VirtAutoBoot = "0"
	}

	if system.VirtCPUs == "" {
		system.VirtCPUs = "<<inherit>>"
	}

	if system.VirtDiskDriver == "" {
		system.VirtDiskDriver = "<<inherit>>"
	}

	if system.VirtFileSize == "" {
		system.VirtFileSize = "<<inherit>>"
	}

	if system.VirtPath == "" {
		system.VirtPath = "<<inherit>>"
	}

	if system.VirtRAM == "" {
		system.VirtRAM = "<<inherit>>"
	}

	if system.VirtType == "" {
		system.VirtType = "<<inherit>>"
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
	if _, err := c.Call("save_system", newID, c.Token); err != nil {
		return nil, err
	}

	// Return a clean copy of the system
	return c.GetSystem(system.Name)
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

// DeleteSystem deletes a single system by its name.
func (c *Client) DeleteSystem(name string) error {
	_, err := c.Call("remove_system", name, c.Token)
	return err
}

// CreateInterface controls network interfaces in Cobbler
func (s *System) CreateInterface(name string, iface Interface) error {
	i := structs.Map(iface)
	nic := make(map[string]interface{})
	for key, value := range i {
		attrName := fmt.Sprintf("%s-%s", key, name)
		log.Printf("[DEBUG] Cobblerclient: setting interface attr %s to %s", attrName, value)
		nic[attrName] = value
	}

	systemID, err := s.Client.GetItemHandle("system", s.Name)
	if err != nil {
		return err
	}

	_, err = s.Client.Call("modify_system", systemID, "modify_interface", nic, s.Client.Token)
	if err != nil {
		return err
	}

	// Save the final system
	_, err = s.Client.Call("save_system", systemID, s.Client.Token)
	if err != nil {
		return err
	}

	return nil
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

// DeleteInterface deletes a single interface in a System.
func (s *System) DeleteInterface(name string) error {
	if _, err := s.GetInterface(name); err != nil {
		return err
	}

	systemID, err := s.Client.GetItemHandle("system", s.Name)
	if err != nil {
		return err
	}

	if _, err := s.Client.Call("modify_system", systemID, "delete_interface", name, s.Client.Token); err != nil {
		return err
	}

	// Save the final system
	if _, err := s.Client.Call("save_system", systemID, s.Client.Token); err != nil {
		return err
	}

	return nil
}

// ListSystemNames is returning a list of all system names currently available in Cobbler.
func (c *Client) ListSystemNames() ([]string, error) {
	return c.GetItemNames("system")
}

// FindSystem is searching for one or more systems by any of its attributes.
func (c *Client) FindSystem(criteria map[string]interface{}) ([]*System, error) {
	result, err := c.Call("find_system", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return c.convertRawSystemsList(result)
}

// FindSystemNames is searching for one or more systems by any of its attributes.
func (c *Client) FindSystemNames(criteria map[string]interface{}) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("find_system", criteria, false, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetSystemsSince is returning all systems which were created after the specified date.
func (c *Client) GetSystemsSince(mtime time.Time) ([]*System, error) {
	result, err := c.Call("get_systems_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return c.convertRawSystemsList(result)
}
