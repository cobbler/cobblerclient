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

// Group is the shared structure for the three 4.0.0 group item types. Each
// concrete group (DistroGroup, ProfileGroup, SystemGroup) embeds Group with
// `mapstructure:",squash"` so callers see typed signatures while the wire
// representation stays uniform.
type Group struct {
	Item    `mapstructure:",squash" yaml:",inline"`
	Members []string `mapstructure:"members" json:"members" yaml:"members"`
}

type DistroGroup struct {
	Group `mapstructure:",squash" yaml:",inline"`
}
type ProfileGroup struct {
	Group `mapstructure:",squash" yaml:",inline"`
}
type SystemGroup struct {
	Group `mapstructure:",squash" yaml:",inline"`
}

func newGroup() Group {
	return Group{Item: NewItem(), Members: []string{}}
}

// NewDistroGroup returns a zero-valued DistroGroup with sensible defaults.
func NewDistroGroup() DistroGroup { return DistroGroup{Group: newGroup()} }

// NewProfileGroup returns a zero-valued ProfileGroup with sensible defaults.
func NewProfileGroup() ProfileGroup { return ProfileGroup{Group: newGroup()} }

// NewSystemGroup returns a zero-valued SystemGroup with sensible defaults.
func NewSystemGroup() SystemGroup { return SystemGroup{Group: newGroup()} }

func convertRawGroup(what, name string, xmlrpcResult interface{}, dest interface{}) error {
	if xmlrpcResult == "~" {
		return fmt.Errorf("%s %s not found", what, name)
	}
	_, err := decodeCobblerItem(xmlrpcResult, dest)
	return err
}

// --- DistroGroup ---

func (c *Client) GetDistroGroups() ([]*DistroGroup, error) {
	result, err := c.Call("get_distro_groups", "-1", c.Token)
	if err != nil {
		return nil, err
	}
	var out []*DistroGroup
	for _, raw := range result.([]interface{}) {
		var g DistroGroup
		if err := convertRawGroup("distro_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

func (c *Client) GetDistroGroup(name string, flattened, resolved bool) (*DistroGroup, error) {
	result, err := c.getConcreteItem("get_distro_group", name, flattened, resolved)
	if err != nil {
		return nil, err
	}
	var g DistroGroup
	if err := convertRawGroup("distro_group", name, result, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (c *Client) CreateDistroGroup(g DistroGroup) (*DistroGroup, error) {
	id, err := c.Call("new_distro_group", c.Token)
	if err != nil {
		return nil, err
	}
	objectID := id.(string)
	if err := c.updateCobblerFields("distro_group", reflect.ValueOf(&g).Elem(), objectID); err != nil {
		return nil, err
	}
	if err := c.SaveDistroGroup(objectID, true, true, "new"); err != nil {
		return nil, err
	}
	return c.GetDistroGroup(g.Name, false, false)
}

func (c *Client) UpdateDistroGroup(g *DistroGroup) error {
	result, err := c.Call("get_distro_group_handle", g.Name)
	if err != nil {
		return err
	}
	objectID := result.(string)
	if err := c.updateCobblerFields("distro_group", reflect.ValueOf(g).Elem(), objectID); err != nil {
		return err
	}
	return c.SaveDistroGroup(objectID, true, true, "bypass")
}

// SaveDistroGroup saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveDistroGroup(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_distro_group", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

func (c *Client) DeleteDistroGroup(name string) error {
	_, err := c.Call("remove_distro_group", name, c.Token, false)
	return err
}
func (c *Client) DeleteDistroGroupRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_distro_group", name, c.Token, recursive)
	return err
}
func (c *Client) RenameDistroGroup(objectId, newName string) error {
	_, err := c.Call("rename_distro_group", objectId, newName, c.Token)
	return err
}
func (c *Client) CopyDistroGroup(objectId, newName string) error {
	_, err := c.Call("copy_distro_group", objectId, newName, c.Token)
	return err
}
func (c *Client) ListDistroGroupNames() ([]string, error) {
	return c.GetItemNames("distro_group")
}
func (c *Client) GetDistroGroupHandle(name string) (string, error) {
	return c.GetItemHandle("distro_group", name)
}

func (c *Client) GetDistroGroupsSince(mtime time.Time) ([]*DistroGroup, error) {
	result, err := c.Call("get_distro_groups_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}
	var out []*DistroGroup
	for _, raw := range result.([]interface{}) {
		var g DistroGroup
		if err := convertRawGroup("distro_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

// --- ProfileGroup ---

func (c *Client) GetProfileGroups() ([]*ProfileGroup, error) {
	result, err := c.Call("get_profile_groups", "-1", c.Token)
	if err != nil {
		return nil, err
	}
	var out []*ProfileGroup
	for _, raw := range result.([]interface{}) {
		var g ProfileGroup
		if err := convertRawGroup("profile_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

func (c *Client) GetProfileGroup(name string, flattened, resolved bool) (*ProfileGroup, error) {
	result, err := c.getConcreteItem("get_profile_group", name, flattened, resolved)
	if err != nil {
		return nil, err
	}
	var g ProfileGroup
	if err := convertRawGroup("profile_group", name, result, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (c *Client) CreateProfileGroup(g ProfileGroup) (*ProfileGroup, error) {
	id, err := c.Call("new_profile_group", c.Token)
	if err != nil {
		return nil, err
	}
	objectID := id.(string)
	if err := c.updateCobblerFields("profile_group", reflect.ValueOf(&g).Elem(), objectID); err != nil {
		return nil, err
	}
	if err := c.SaveProfileGroup(objectID, true, true, "new"); err != nil {
		return nil, err
	}
	return c.GetProfileGroup(g.Name, false, false)
}

func (c *Client) UpdateProfileGroup(g *ProfileGroup) error {
	objectID, err := c.GetItemHandle("profile_group", g.Name)
	if err != nil {
		return err
	}
	if err := c.updateCobblerFields("profile_group", reflect.ValueOf(g).Elem(), objectID); err != nil {
		return err
	}
	return c.SaveProfileGroup(objectID, true, true, "bypass")
}

// SaveProfileGroup saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveProfileGroup(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_profile_group", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

func (c *Client) DeleteProfileGroup(name string) error {
	_, err := c.Call("remove_profile_group", name, c.Token, false)
	return err
}
func (c *Client) DeleteProfileGroupRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_profile_group", name, c.Token, recursive)
	return err
}
func (c *Client) RenameProfileGroup(objectId, newName string) error {
	_, err := c.Call("rename_profile_group", objectId, newName, c.Token)
	return err
}
func (c *Client) CopyProfileGroup(objectId, newName string) error {
	_, err := c.Call("copy_profile_group", objectId, newName, c.Token)
	return err
}
func (c *Client) ListProfileGroupNames() ([]string, error) {
	return c.GetItemNames("profile_group")
}
func (c *Client) GetProfileGroupHandle(name string) (string, error) {
	result, err := c.Call("get_profile_group_handle", name)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (c *Client) GetProfileGroupsSince(mtime time.Time) ([]*ProfileGroup, error) {
	result, err := c.Call("get_profile_groups_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}
	var out []*ProfileGroup
	for _, raw := range result.([]interface{}) {
		var g ProfileGroup
		if err := convertRawGroup("profile_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

// --- SystemGroup ---

func (c *Client) GetSystemGroups() ([]*SystemGroup, error) {
	result, err := c.Call("get_system_groups", "-1", c.Token)
	if err != nil {
		return nil, err
	}
	var out []*SystemGroup
	for _, raw := range result.([]interface{}) {
		var g SystemGroup
		if err := convertRawGroup("system_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

func (c *Client) GetSystemGroup(name string, flattened, resolved bool) (*SystemGroup, error) {
	result, err := c.getConcreteItem("get_system_group", name, flattened, resolved)
	if err != nil {
		return nil, err
	}
	var g SystemGroup
	if err := convertRawGroup("system_group", name, result, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (c *Client) CreateSystemGroup(g SystemGroup) (*SystemGroup, error) {
	id, err := c.Call("new_system_group", c.Token)
	if err != nil {
		return nil, err
	}
	objectID := id.(string)
	if err := c.updateCobblerFields("system_group", reflect.ValueOf(&g).Elem(), objectID); err != nil {
		return nil, err
	}
	if err := c.SaveSystemGroup(objectID, true, true, "new"); err != nil {
		return nil, err
	}
	return c.GetSystemGroup(g.Name, false, false)
}

func (c *Client) UpdateSystemGroup(g *SystemGroup) error {
	objectID, err := c.GetItemHandle("system_group", g.Name)
	if err != nil {
		return err
	}
	if err := c.updateCobblerFields("system_group", reflect.ValueOf(g).Elem(), objectID); err != nil {
		return err
	}
	return c.SaveSystemGroup(objectID, true, true, "bypass")
}

// SaveSystemGroup saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveSystemGroup(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_system_group", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

func (c *Client) DeleteSystemGroup(name string) error {
	_, err := c.Call("remove_system_group", name, c.Token, false)
	return err
}
func (c *Client) DeleteSystemGroupRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_system_group", name, c.Token, recursive)
	return err
}
func (c *Client) RenameSystemGroup(objectId, newName string) error {
	return c.RenameItem("system_group", objectId, newName)
}
func (c *Client) CopySystemGroup(objectId, newName string) error {
	_, err := c.Call("copy_system_group", objectId, newName, c.Token)
	return err
}
func (c *Client) ListSystemGroupNames() ([]string, error) {
	return c.GetItemNames("system_group")
}
func (c *Client) GetSystemGroupHandle(name string) (string, error) {
	result, err := c.Call("get_system_group_handle", name)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (c *Client) GetSystemGroupsSince(mtime time.Time) ([]*SystemGroup, error) {
	result, err := c.Call("get_system_groups_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}
	var out []*SystemGroup
	for _, raw := range result.([]interface{}) {
		var g SystemGroup
		if err := convertRawGroup("system_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

// GetDistroGroupAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetDistroGroupAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_distro_group_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

// GetProfileGroupAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetProfileGroupAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_profile_group_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

// GetSystemGroupAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetSystemGroupAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_system_group_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

// FindDistroGroup searches for distro groups matching the given criteria.
func (c *Client) FindDistroGroup(criteria map[string]interface{}, resolved bool) ([]*DistroGroup, error) {
	result, err := c.Call("find_distro_group", criteria, true, resolved, c.Token)
	if err != nil {
		return nil, err
	}
	var out []*DistroGroup
	for _, raw := range result.([]interface{}) {
		var g DistroGroup
		if err := convertRawGroup("distro_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

// FindProfileGroup searches for profile groups matching the given criteria.
func (c *Client) FindProfileGroup(criteria map[string]interface{}, resolved bool) ([]*ProfileGroup, error) {
	result, err := c.Call("find_profile_group", criteria, true, resolved, c.Token)
	if err != nil {
		return nil, err
	}
	var out []*ProfileGroup
	for _, raw := range result.([]interface{}) {
		var g ProfileGroup
		if err := convertRawGroup("profile_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}

// FindSystemGroup searches for system groups matching the given criteria.
func (c *Client) FindSystemGroup(criteria map[string]interface{}, resolved bool) ([]*SystemGroup, error) {
	result, err := c.Call("find_system_group", criteria, true, resolved, c.Token)
	if err != nil {
		return nil, err
	}
	var out []*SystemGroup
	for _, raw := range result.([]interface{}) {
		var g SystemGroup
		if err := convertRawGroup("system_group", "unknown", raw, &g); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, nil
}
