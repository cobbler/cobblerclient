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
	// These are internal fields and cannot be modified.
	Ctime             float64  `mapstructure:"ctime"           cobbler:"noupdate"` // TODO: convert to time
	Depth             int      `mapstructure:"depth"          cobbler:"noupdate"`
	Mtime             float64  `mapstructure:"mtime"           cobbler:"noupdate"` // TODO: convert to time
	SourceRepos       []string `mapstructure:"source_repos"   cobbler:"noupdate"`
	TreeBuildTime     string   `mapstructure:"tree_build_time" cobbler:"noupdate"`
	UID               string   `mapstructure:"uid"            cobbler:"noupdate"`
	Arch              string   `mapstructure:"arch"`
	BootFiles         []string `mapstructure:"boot_files"`
	BootLoaders       []string `mapstructure:"boot_loaders"`
	Breed             string   `mapstructure:"breed"`
	Comment           string   `mapstructure:"comment"`
	FetchableFiles    []string `mapstructure:"fetchable_files"`
	Initrd            string   `mapstructure:"initrd"`
	Kernel            string   `mapstructure:"kernel"`
	KernelOptions     []string `mapstructure:"kernel_options"`
	KernelOptionsPost []string `mapstructure:"kernel_options_post"`
	MGMTClasses       []string `mapstructure:"mgmt_classes"`
	Name              string   `mapstructure:"name"`
	OSVersion         string   `mapstructure:"os_version"`
	Owners            []string `mapstructure:"owners"`
	TemplateFiles     []string `mapstructure:"template_files"`
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

	return decodeResult.(*Distro), nil
}

//convertRawDistrosList...
func convertRawDistrosList(xmlrpcResult interface{}) ([]*Distro, error) {
	var distros []*Distro

	for _, d := range xmlrpcResult.([]interface{}) {
		distro, err := convertRawDistro("unknown", d)
		if err != nil {
			return nil, err
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
func (c *Client) GetDistro(name string) (*Distro, error) {
	result, err := c.Call("get_distro", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawDistro(name, result)
}

// CreateDistro creates a distro.
func (c *Client) CreateDistro(distro Distro) (*Distro, error) {
	// Make sure a distro with the same name does not already exist
	if _, err := c.GetDistro(distro.Name); err == nil {
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

	return c.GetDistro(distro.Name)
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

// DeleteDistro deletes a single distro by its name.
func (c *Client) DeleteDistro(name string) error {
	_, err := c.Call("remove_distro", name, c.Token)
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
	var result []string

	resultUnmarshalled, err := c.Call("find_distro", criteria, false, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// RenameDistro renames a distro with a given object id.
func (c *Client) RenameDistro(objectId, newName string) error {
	_, err := c.Call("rename_distro", objectId, newName, c.Token)
	return err
}

// GetDistroHandle gets the internal ID of a Cobbler item.
func (c *Client) GetDistroHandle(name string) (string, error) {
	result, err := c.Call("get_distro_handle", name, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}
