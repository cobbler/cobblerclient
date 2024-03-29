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
)

// Profile is a created profile.
type Profile struct {
	// These are internal fields and cannot be modified.
	Ctime        float64 `mapstructure:"ctime"                  cobbler:"noupdate"` // TODO: convert to time
	Depth        int     `mapstructure:"depth"                  cobbler:"noupdate"`
	ID           string  `mapstructure:"uid"                    cobbler:"noupdate"`
	Mtime        float64 `mapstructure:"mtime"                  cobbler:"noupdate"` // TODO: convert to time
	ReposEnabled bool    `mapstructure:"repos_enabled"          cobbler:"noupdate"`

	Autoinstall       string   `mapstructure:"autoinstall"`
	AutoinstallMeta   []string `mapstructure:"autoinstall_meta"`
	BootFiles         []string `mapstructure:"boot_files"`
	Comment           string   `mapstructure:"comment"`
	DHCPTag           string   `mapstructure:"dhcp_tag"`
	Distro            string   `mapstructure:"distro"`
	EnableGPXE        bool     `mapstructure:"enable_gpxe"`
	EnableMenu        bool     `mapstructure:"enable_menu"`
	FetchableFiles    []string `mapstructure:"fetchable_files"`
	KernelOptions     []string `mapstructure:"kernel_options"`
	KernelOptionsPost []string `mapstructure:"kernel_options_post"`
	MGMTClasses       []string `mapstructure:"mgmt_classes"`
	MGMTParameters    string   `mapstructure:"mgmt_parameters"`
	Name              string   `mapstructure:"name"`
	NameServers       []string `mapstructure:"name_servers"`
	NameServersSearch []string `mapstructure:"name_servers_search"`
	NextServerv4      string   `mapstructure:"next_server_v4"`
	NextServerv6      string   `mapstructure:"next_server_v6"`
	Owners            []string `mapstructure:"owners"`
	Parent            string   `mapstructure:"parent"`
	Proxy             string   `mapstructure:"proxy"`
	Repos             []string `mapstructure:"repos"`
	Server            string   `mapstructure:"server"`
	TemplateFiles     []string `mapstructure:"template_files"`
	VirtAutoBoot      string   `mapstructure:"virt_auto_boot"`
	VirtBridge        string   `mapstructure:"virt_bridge"`
	VirtCPUs          string   `mapstructure:"virt_cpus"`
	VirtDiskDriver    string   `mapstructure:"virt_disk_driver"`
	VirtFileSize      string   `mapstructure:"virt_file_size"`
	VirtPath          string   `mapstructure:"virt_path"`
	VirtRAM           string   `mapstructure:"virt_ram"`
	VirtType          string   `mapstructure:"virt_type"`

	Client
}

// GetProfiles returns all systems in Cobbler.
func (c *Client) GetProfiles() ([]*Profile, error) {
	var profiles []*Profile

	result, err := c.Call("get_profiles", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, p := range result.([]interface{}) {
		var profile Profile
		decodedResult, err := decodeCobblerItem(p, &profile)
		if err != nil {
			return nil, err
		}
		decodedProfile := decodedResult.(*Profile)
		decodedProfile.Client = *c
		profiles = append(profiles, decodedProfile)
	}

	return profiles, nil
}

// GetProfile returns a single profile obtained by its name.
func (c *Client) GetProfile(name string) (*Profile, error) {
	var profile Profile

	result, err := c.Call("get_profile", name, c.Token)
	if err != nil {
		return &profile, err
	}

	if result == "~" {
		return nil, fmt.Errorf("Profile %s not found.", name)
	}

	decodeResult, err := decodeCobblerItem(result, &profile)
	if err != nil {
		return nil, err
	}

	s := decodeResult.(*Profile)
	s.Client = *c

	return s, nil
}

// CreateProfile creates a profile.
// It ensures that a Distro is set and then sets other default values.
func (c *Client) CreateProfile(profile Profile) (*Profile, error) {
	// Check if a profile with the same name already exists
	if _, err := c.GetProfile(profile.Name); err == nil {
		return nil, fmt.Errorf("A profile with the name %s already exists.", profile.Name)
	}

	if profile.Distro == "" {
		return nil, fmt.Errorf("A profile must have a distro set.")
	}

	if profile.MGMTParameters == "" {
		profile.MGMTParameters = "<<inherit>>"
	}
	if profile.VirtType == "" {
		profile.VirtType = "<<inherit>>"
	}
	if profile.VirtDiskDriver == "" {
		profile.VirtDiskDriver = "<<inherit>>"
	}

	// To create a profile via the Cobbler API, first call new_profile to obtain an ID
	result, err := c.Call("new_profile", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)
	// Set the value of all fields
	item := reflect.ValueOf(&profile).Elem()
	if err := c.updateCobblerFields("profile", item, newID); err != nil {
		return nil, err
	}

	// Save the final profile
	result, err = c.Call("save_profile", newID, c.Token)
	if err != nil {
		return nil, err
	}

	// Return a clean copy of the profile
	return c.GetProfile(profile.Name)
}

// UpdateProfile updates a single profile.
func (c *Client) UpdateProfile(profile *Profile) error {
	item := reflect.ValueOf(profile).Elem()
	id, err := c.GetItemHandle("profile", profile.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("profile", item, id); err != nil {
		return err
	}

	// Save the final profile
	if _, err := c.Call("save_profile", id, c.Token); err != nil {
		return err
	}

	return nil
}

// DeleteProfile deletes a single profile by its name.
func (c *Client) DeleteProfile(name string) error {
	_, err := c.Call("remove_profile", name, c.Token)
	return err
}
