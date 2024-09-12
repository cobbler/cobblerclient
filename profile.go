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
	Item `mapstructure:",squash"`

	// These are internal fields and cannot be modified.
	ReposEnabled bool `mapstructure:"repos_enabled"          cobbler:"noupdate"`

	Autoinstall         string          `mapstructure:"autoinstall"`
	BootLoaders         Value[[]string] `mapstructure:"boot_loaders"`
	DHCPTag             string          `mapstructure:"dhcp_tag"`
	Distro              string          `mapstructure:"distro"`
	EnableIPXE          Value[bool]     `mapstructure:"enable_ipxe"`
	EnableMenu          Value[bool]     `mapstructure:"enable_menu"`
	Filename            string          `mapstructure:"filename"`
	Menu                string          `mapstructure:"menu"`
	NameServers         Value[[]string] `mapstructure:"name_servers"`
	NameServersSearch   Value[[]string] `mapstructure:"name_servers_search"`
	NextServerv4        string          `mapstructure:"next_server_v4"`
	NextServerv6        string          `mapstructure:"next_server_v6"`
	Proxy               string          `mapstructure:"proxy"`
	RedhatManagementKey string          `mapstructure:"redhat_management_key"`
	Repos               []string        `mapstructure:"repos"`
	Server              string          `mapstructure:"server"`
	VirtAutoBoot        Value[bool]     `mapstructure:"virt_auto_boot"`
	VirtBridge          string          `mapstructure:"virt_bridge"`
	VirtCPUs            int             `mapstructure:"virt_cpus"`
	VirtDiskDriver      string          `mapstructure:"virt_disk_driver"`
	VirtFileSize        Value[float64]  `mapstructure:"virt_file_size"`
	VirtPath            string          `mapstructure:"virt_path"`
	VirtRAM             Value[int]      `mapstructure:"virt_ram"`
	VirtType            string          `mapstructure:"virt_type"`

	Client
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
		NameServers: Value[[]string]{
			Data:        make([]string, 0),
			IsInherited: true,
		},
		NameServersSearch: Value[[]string]{
			Data:        make([]string, 0),
			IsInherited: true,
		},
		NextServerv4:        inherit,
		NextServerv6:        inherit,
		Proxy:               inherit,
		RedhatManagementKey: inherit,
		Repos:               make([]string, 0),
		Server:              inherit,
		VirtAutoBoot: Value[bool]{
			IsInherited: true,
		},
		VirtBridge:     inherit,
		VirtCPUs:       1,
		VirtDiskDriver: inherit,
		VirtFileSize: Value[float64]{
			IsInherited: true,
		},
		VirtRAM: Value[int]{
			IsInherited: true,
		},
		VirtType: inherit,
	}
	// Overwrite Item defaults
	profile.BootFiles = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.FetchableFiles = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.AutoinstallMeta = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.KernelOptions = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.KernelOptionsPost = Value[map[string]interface{}]{
		IsInherited: true,
	}
	profile.MgmtClasses = Value[[]string]{
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
	err = sanitizeValueMapStruct(&decodedProfile.FetchableFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedProfile.BootFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedProfile.TemplateFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedProfile.MgmtParameters)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.MgmtClasses)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.BootLoaders)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.NameServers)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedProfile.NameServersSearch)
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

	if profile.VirtType == "" {
		profile.VirtType = inherit
	}
	if profile.VirtDiskDriver == "" {
		profile.VirtDiskDriver = inherit
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
	err = c.SaveProfile(newID, "bypass")
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

	return c.SaveProfile(id, "bypass")
}

// SaveProfile saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveProfile(objectId, editmode string) error {
	_, err := c.Call("save_profile", objectId, c.Token, editmode)
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
	result, err := c.Call("find_profile", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawProfilesList(result)
}

// FindProfileNames searches for one or more profiles by any of its attributes.
func (c *Client) FindProfileNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_profile", criteria, false, c.Token)
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
	res, err := c.Call("get_profile_handle", name, c.Token)
	return returnString(res, err)
}
