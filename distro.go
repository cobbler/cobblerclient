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
	Item `mapstructure:",squash"`

	// These are internal fields and cannot be modified.
	SourceRepos         []string        `mapstructure:"source_repos"   cobbler:"noupdate"`
	TreeBuildTime       string          `mapstructure:"tree_build_time" cobbler:"noupdate"`
	Arch                string          `mapstructure:"arch"`
	BootLoaders         Value[[]string] `mapstructure:"boot_loaders"`
	Breed               string          `mapstructure:"breed"`
	Initrd              string          `mapstructure:"initrd"`
	RemoteBootInitrd    string          `mapstructure:"remote_boot_initrd"`
	Kernel              string          `mapstructure:"kernel"`
	RemoteBootKernel    string          `mapstructure:"remote_boot_kernel"`
	RedhatManagementKey string          `mapstructure:"redhat_management_key"`
	OSVersion           string          `mapstructure:"os_version"`
}

func NewDistro() Distro {
	return Distro{
		Item: NewItem(),
		Arch: "x86_64",
		BootLoaders: Value[[]string]{
			Data:        make([]string, 0),
			IsInherited: true,
		},
		RedhatManagementKey: inherit,
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
	err = sanitizeValueMapStruct(&decodedDistro.FetchableFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedDistro.BootFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedDistro.TemplateFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedDistro.MgmtParameters)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedDistro.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedDistro.MgmtClasses)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedDistro.BootLoaders)
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
	if _, err := c.GetDistro(distro.Name, false, false); err == nil {
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

	if _, err := c.Call("save_distro", newID, c.Token); err != nil {
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

	if _, err := c.Call("save_distro", id, c.Token); err != nil {
		return err
	}

	return nil
}

// SaveDistro saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveDistro(objectId, editmode string) error {
	_, err := c.Call("save_distro", objectId, c.Token, editmode)
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
func (c *Client) FindDistro(criteria map[string]interface{}) ([]*Distro, error) {
	var distros []*Distro

	result, err := c.Call("find_distro", criteria, true, c.Token)
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
	resultUnmarshalled, err := c.Call("find_distro", criteria, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// RenameDistro renames a distro with a given object id.
func (c *Client) RenameDistro(objectId, newName string) error {
	_, err := c.Call("rename_distro", objectId, newName, c.Token)
	return err
}

// GetDistroHandle gets the internal ID of a Cobbler item.
func (c *Client) GetDistroHandle(name string) (string, error) {
	res, err := c.Call("get_distro_handle", name, c.Token)
	return returnString(res, err)
}
