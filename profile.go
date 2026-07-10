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

// Profile is a created profile.
// Get the fields from cobbler/items/profile.py
type Profile struct {
	Item `mapstructure:",squash" yaml:",inline"`

	// These are internal fields and cannot be modified.
	ReposEnabled bool `mapstructure:"repos_enabled"          cobbler:"noupdate"`

	Autoinstall              string          `mapstructure:"autoinstall" json:"autoinstall" yaml:"autoinstall"`
	BootLoaders              Value[[]string] `mapstructure:"boot_loaders" json:"boot_loaders" yaml:"boot_loaders"`
	DHCPTag                  string          `mapstructure:"dhcp_tag" json:"dhcp_tag" yaml:"dhcp_tag"`
	Distro                   string          `mapstructure:"distro" json:"distro" yaml:"distro"`
	DNS                      DNSOptions      `mapstructure:"dns" json:"dns" yaml:"dns"`
	EnableIPXE               Value[bool]     `mapstructure:"enable_ipxe" json:"enable_ipxe" yaml:"enable_ipxe" cobbler_min_inherit:"3.3.5"`
	EnableMenu               Value[bool]     `mapstructure:"enable_menu" json:"enable_menu" yaml:"enable_menu" cobbler_min_inherit:"3.3.5"`
	Filename                 string          `mapstructure:"filename" json:"filename" yaml:"filename"`
	Menu                     string          `mapstructure:"menu" json:"menu" yaml:"menu"`
	Proxy                    string          `mapstructure:"proxy" json:"proxy" yaml:"proxy"`
	RedhatManagementKey      string          `mapstructure:"redhat_management_key" json:"redhat_management_key" yaml:"redhat_management_key"`
	RedhatManagementOrg      string          `mapstructure:"redhat_management_org" json:"redhat_management_org" yaml:"redhat_management_org"`
	RedhatManagementUser     string          `mapstructure:"redhat_management_user" json:"redhat_management_user" yaml:"redhat_management_user"`
	RedhatManagementPassword string          `mapstructure:"redhat_management_password" json:"redhat_management_password" yaml:"redhat_management_password"`
	Repos                    []string        `mapstructure:"repos" json:"repos" yaml:"repos"`
	Server                   string          `mapstructure:"server" json:"server" yaml:"server"`
	TFTP                     TFTPOptions     `mapstructure:"tftp" json:"tftp" yaml:"tftp"`
	Virt                     VirtOptions     `mapstructure:"virt" json:"virt" yaml:"virt"`
	VirtBridge               string          `mapstructure:"virt_bridge" json:"virt_bridge" yaml:"virt_bridge"`
}

func NewProfile() Profile {
	profile := Profile{
		Item:         NewItem(),
		ReposEnabled: false,
		Autoinstall:  inherit,
		BootLoaders: Value[[]string]{
			Data:        make([]string, 0),
			IsInherited: true,
		},
		EnableIPXE: Value[bool]{
			IsInherited: true,
		},
		EnableMenu: Value[bool]{
			IsInherited: true,
		},
		DNS:                      newDNSOptions(),
		Proxy:                    inherit,
		RedhatManagementKey:      inherit,
		RedhatManagementOrg:      inherit,
		RedhatManagementUser:     inherit,
		RedhatManagementPassword: inherit,
		Repos:                    make([]string, 0),
		Server:                   inherit,
		TFTP:                     newTFTPOptions(),
		Virt:                     newVirtOptions(),
		VirtBridge:               inherit,
	}
	// Overwrite Item defaults
	profile.AutoinstallMeta = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.KernelOptions = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.KernelOptionsPost = Value[map[string]interface{}]{
		IsInherited: true,
	}
	return profile
}

func convertRawProfile(name string, xmlrpcResult interface{}) (*Profile, error) {
	var profile Profile

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &profile)
	if err != nil {
		return nil, err
	}

	// Now clean the Value structs
	decodedProfile := decodeResult.(*Profile)
	err = sanitizeValueMapStruct(&decodedProfile.KernelOptions)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedProfile.KernelOptionsPost)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedProfile.AutoinstallMeta)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.BootLoaders)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.DNS.NameServers)
	if err != nil {
		return nil, err
	}
	return decodedProfile, nil
}

func convertRawProfilesList(xmlrpcResult interface{}) ([]*Profile, error) {
	var profiles []*Profile

	for _, p := range xmlrpcResult.([]interface{}) {
		profile, err := convertRawProfile("unknown", p)
		if err != nil {
			return nil, err
		}
		profile.Meta = ItemMeta{
			IsFlattened: false,
			IsResolved:  false,
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetProfiles returns all profiles in Cobbler.
func (c *Client) GetProfiles() ([]*Profile, error) {
	result, err := c.Call("get_profiles", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawProfilesList(result)
}

// GetProfile returns a single profile obtained by its name.
func (c *Client) GetProfile(name string, flattened, resolved bool) (*Profile, error) {
	result, err := c.getConcreteItem("get_profile", name, flattened, resolved)

	if err != nil {
		return nil, err
	}

	profile, err := convertRawProfile(name, result)
	if err != nil {
		return nil, err
	}
	profile.Meta = ItemMeta{
		IsFlattened: flattened,
		IsResolved:  resolved,
	}
	return profile, nil
}

// CreateProfile creates a profile.
// It ensures that a Distro is set and then sets other default values.
func (c *Client) CreateProfile(profile Profile) (*Profile, error) {
	// Check if a profile with the same name already exists
	if _, err := c.GetProfile(profile.Name, false, false); err == nil {
		return nil, fmt.Errorf("a profile with the name %s already exists", profile.Name)
	}

	if profile.Distro == "" {
		return nil, fmt.Errorf("a profile must have a distro set")
	}

	if profile.Virt.Type == "" {
		profile.Virt.Type = inherit
	}
	if profile.Virt.DiskDriver == "" {
		profile.Virt.DiskDriver = inherit
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
	err = c.SaveProfile(newID, true, true, "bypass")
	if err != nil {
		return nil, err
	}

	// Return a clean copy of the profile
	return c.GetProfile(profile.Name, false, false)
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

	return c.SaveProfile(id, true, true, "bypass")
}

// SaveProfile saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveProfile(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_profile", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

// CopyProfile duplicates a given profile on the server with a new name.
func (c *Client) CopyProfile(objectId, newName string) error {
	_, err := c.Call("copy_profile", objectId, newName, c.Token)
	return err
}

// DeleteProfile deletes a single profile by its name.
func (c *Client) DeleteProfile(name string) error {
	return c.DeleteProfileRecursive(name, false)
}

// DeleteProfileRecursive deletes a single profile by its name.
func (c *Client) DeleteProfileRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_profile", name, c.Token, recursive)
	return err
}

// ListProfileNames returns a list of all profile names currently available in Cobbler.
func (c *Client) ListProfileNames() ([]string, error) {
	return c.GetItemNames("profile")
}

// FindProfile searches for one or more profiles by any of its attributes.
func (c *Client) FindProfile(criteria map[string]interface{}) ([]*Profile, error) {
	result, err := c.Call("find_profile", criteria, true, false, c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawProfilesList(result)
}

// FindProfileNames searches for one or more profiles by any of its attributes.
func (c *Client) FindProfileNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_profile", criteria, false, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetProfilesSince returns all profiles which were created after the specified date.
func (c *Client) GetProfilesSince(mtime time.Time) ([]*Profile, error) {
	result, err := c.Call("get_profiles_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawProfilesList(result)
}

// RenameProfile renames a profile with a given object id.
func (c *Client) RenameProfile(objectId, newName string) error {
	_, err := c.Call("rename_profile", objectId, newName, c.Token)
	return err
}

// GetProfileHandle gets the internal ID of a Cobbler item.
func (c *Client) GetProfileHandle(name string) (string, error) {
	res, err := c.Call("get_profile_handle", name)
	return returnString(res, err)
}

// GetValidProfileBootLoaders retrieves the list of bootloaders that can be assigned to a profile.
func (c *Client) GetValidProfileBootLoaders(profileName string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_profile_boot_loaders", profileName, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetProfileAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetProfileAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_profile_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

// NewSubprofile creates a new blank sub-profile on the server and returns its object ID.
// A sub-profile inherits from another profile (rather than a distro).
func (c *Client) NewSubprofile() (string, error) {
	result, err := c.Call("new_subprofile", c.Token)
	return returnString(result, err)
}
