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

// Distro is a created distro.
// Get the fields from cobbler/items/distro.py
type Distro struct {
	Item `mapstructure:",squash" yaml:",inline"`

	// These are internal fields and cannot be modified.
	SourceRepos              []string        `mapstructure:"source_repos"   cobbler:"noupdate" json:"source_repos" yaml:"source_repos"`
	TreeBuildTime            string          `mapstructure:"tree_build_time" cobbler:"noupdate" json:"tree_build_time" yaml:"tree_build_time"`
	Arch                     string          `mapstructure:"arch" json:"arch" yaml:"arch"`
	BootLoaders              Value[[]string] `mapstructure:"boot_loaders" json:"boot_loaders" yaml:"boot_loaders"`
	Breed                    string          `mapstructure:"breed" json:"breed" yaml:"breed"`
	Initrd                   string          `mapstructure:"initrd" json:"initrd" yaml:"initrd"`
	RemoteBootInitrd         string          `mapstructure:"remote_boot_initrd" json:"remote_boot_initrd" yaml:"remote_boot_initrd"`
	Kernel                   string          `mapstructure:"kernel" json:"kernel" yaml:"kernel"`
	RemoteBootKernel         string          `mapstructure:"remote_boot_kernel" json:"remote_boot_kernel" yaml:"remote_boot_kernel"`
	RedhatManagementKey      string          `mapstructure:"redhat_management_key" json:"redhat_management_key" yaml:"redhat_management_key"`
	RedhatManagementOrg      string          `mapstructure:"redhat_management_org" json:"redhat_management_org" yaml:"redhat_management_org"`
	RedhatManagementUser     string          `mapstructure:"redhat_management_user" json:"redhat_management_user" yaml:"redhat_management_user"`
	RedhatManagementPassword string          `mapstructure:"redhat_management_password" json:"redhat_management_password" yaml:"redhat_management_password"`
	OSVersion                string          `mapstructure:"os_version" json:"os_version" yaml:"os_version"`
	// SourceTreePath is new in Cobbler 4.0.0a5.
	SourceTreePath string `mapstructure:"source_tree_path" json:"source_tree_path" yaml:"source_tree_path"`
}

func NewDistro() Distro {
	return Distro{
		Item: NewItem(),
		Arch: "x86_64",
		BootLoaders: Value[[]string]{
			Data:        make([]string, 0),
			IsInherited: true,
		},
		RedhatManagementKey:      inherit,
		RedhatManagementOrg:      inherit,
		RedhatManagementUser:     inherit,
		RedhatManagementPassword: inherit,
	}
}

// convertRawDistro ...
func convertRawDistro(name string, xmlrpcResult interface{}) (*Distro, error) {
	var distro Distro

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("distro %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &distro)
	if err != nil {
		return nil, err
	}

	// Now clean the Value structs
	decodedDistro := decodeResult.(*Distro)
	err = sanitizeValueMapStruct(&decodedDistro.KernelOptions)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedDistro.KernelOptionsPost)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedDistro.AutoinstallMeta)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedDistro.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedDistro.BootLoaders)
	if err != nil {
		return nil, err
	}
	return &distro, nil
}

func sanitizeValueSliceStruct(value *Value[[]string]) error {
	if value.IsInherited {
		value.Data = make([]string, 0)
		value.FlattenedValue = ""
		value.RawData = make([]string, 0)
	} else {
		kopts, err := returnStringSlice(value.RawData, nil)
		if err == nil {
			value.Data = kopts
		} else {
			kopts, ok := value.RawData.(string)
			if ok {
				value.Data = make([]string, 0)
				value.FlattenedValue = kopts
			} else {
				if value.RawData == nil {
					value.Data = make([]string, 0)
					value.FlattenedValue = ""
					value.RawData = make([]string, 0)
				} else {
					return fmt.Errorf("error converting raw list value")
				}
			}
		}
	}
	return nil
}

func sanitizeValueMapStruct(value *Value[map[string]interface{}]) error {
	if value.IsInherited {
		value.Data = make(map[string]interface{})
		value.FlattenedValue = ""
		value.RawData = make(map[string]interface{})
	} else {
		kopts, ok := value.RawData.(map[string]interface{})
		if ok {
			value.Data = kopts
		} else {
			kopts, ok := value.RawData.(string)
			if ok {
				value.Data = make(map[string]interface{})
				value.FlattenedValue = kopts
			} else {
				if value.RawData == nil {
					value.Data = make(map[string]interface{})
					value.FlattenedValue = ""
					value.RawData = make(map[string]interface{})
				} else {
					return fmt.Errorf("error converting raw map value")
				}
			}
		}
	}
	return nil
}

// convertRawDistrosList...
func convertRawDistrosList(xmlrpcResult interface{}) ([]*Distro, error) {
	var distros []*Distro

	for _, d := range xmlrpcResult.([]interface{}) {
		distro, err := convertRawDistro("unknown", d)
		if err != nil {
			return nil, err
		}
		distro.Meta = ItemMeta{
			IsFlattened: false,
			IsResolved:  false,
		}
		distros = append(distros, distro)
	}

	return distros, nil
}

// GetDistros returns all distros in Cobbler.
func (c *Client) GetDistros() ([]*Distro, error) {
	result, err := c.Call("get_distros", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawDistrosList(result)
}

// GetDistro returns a single distro obtained by its name.
func (c *Client) GetDistro(name string, flattened, resolved bool) (*Distro, error) {
	result, err := c.getConcreteItem("get_distro", name, flattened, resolved)
	if err != nil {
		return nil, err
	}

	distro, err := convertRawDistro(name, result)
	if err != nil {
		return nil, err
	}
	distro.Meta = ItemMeta{
		IsFlattened: flattened,
		IsResolved:  resolved,
	}
	return distro, nil
}

// CreateDistro creates a distro.
func (c *Client) CreateDistro(distro Distro) (*Distro, error) {
	// Make sure a distro with the same name does not already exist
	if exists, err := c.HasItem("distro", distro.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("a Distro with the name %s already exists", distro.Name)
	}

	result, err := c.Call("new_distro", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&distro).Elem()
	if err := c.updateCobblerFields("distro", item, newID); err != nil {
		return nil, err
	}

	if err := c.SaveDistro(newID, true, true, "new"); err != nil {
		return nil, err
	}

	return c.GetDistro(distro.Name, false, false)
}

// UpdateDistro updates a single distro.
func (c *Client) UpdateDistro(distro *Distro) error {
	item := reflect.ValueOf(distro).Elem()
	id, err := c.GetItemHandle("distro", distro.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("distro", item, id); err != nil {
		return err
	}

	if err := c.SaveDistro(id, true, true, "bypass"); err != nil {
		return err
	}

	return nil
}

// SaveDistro saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveDistro(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_distro", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

// CopyDistro duplicates a distro on the server with a new name.
func (c *Client) CopyDistro(objectId, newName string) error {
	_, err := c.Call("copy_distro", objectId, newName, c.Token)
	return err
}

// DeleteDistro deletes a single Distro by its name.
func (c *Client) DeleteDistro(name string) error {
	return c.DeleteDistroRecursive(name, false)
}

// DeleteDistroRecursive deletes a single Distro by its name with the option to do so recursively.
func (c *Client) DeleteDistroRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_distro", name, c.Token, recursive)
	return err
}

// ListDistroNames returns a list of all distro names currently available in Cobbler.
func (c *Client) ListDistroNames() ([]string, error) {
	return c.GetItemNames("distro")
}

// GetDistrosSince returns all distros which were created after the specified date.
func (c *Client) GetDistrosSince(mtime time.Time) ([]*Distro, error) {
	var distros []*Distro

	result, err := c.Call("get_distros_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	for _, d := range result.([]interface{}) {
		var distro Distro
		decodedResult, err := decodeCobblerItem(d, &distro)
		if err != nil {
			return nil, err
		}

		distros = append(distros, decodedResult.(*Distro))
	}

	return distros, nil
}

// FindDistro searches for one or more distros by any of its attributes.
func (c *Client) FindDistro(criteria map[string]interface{}, resolved bool) ([]*Distro, error) {
	var distros []*Distro

	result, err := c.Call("find_distro", criteria, true, resolved, c.Token)
	if err != nil {
		return nil, err
	}

	for _, d := range result.([]interface{}) {
		var distro Distro
		decodedResult, err := decodeCobblerItem(d, &distro)
		if err != nil {
			return nil, err
		}

		distros = append(distros, decodedResult.(*Distro))
	}

	return distros, nil
}

// FindDistroNames searches for one or more distros by any of its attributes.
func (c *Client) FindDistroNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_distro", criteria, false, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// RenameDistro renames a distro with a given object id.
func (c *Client) RenameDistro(objectId, newName string) error {
	_, err := c.Call("rename_distro", objectId, newName, c.Token)
	return err
}

// GetDistroHandle gets the internal ID of a Cobbler item.
func (c *Client) GetDistroHandle(name string) (string, error) {
	res, err := c.Call("get_distro_handle", name)
	return returnString(res, err)
}

// GetValidDistroBootLoaders retrieves the list of bootloaders that can be assigned to a distro.
func (c *Client) GetValidDistroBootLoaders(distroName string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_distro_boot_loaders", distroName, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetDistroAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetDistroAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_distro_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}
