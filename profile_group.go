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

type ProfileGroup struct {
	Group `mapstructure:",squash" yaml:",inline"`
}

// NewProfileGroup returns a zero-valued ProfileGroup with sensible defaults.
func NewProfileGroup() ProfileGroup { return ProfileGroup{Group: newGroup()} }

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

// GetProfileGroup returns a single profile group obtained by its uid.
func (c *Client) GetProfileGroup(uid string, flattened, resolved bool) (*ProfileGroup, error) {
	result, err := c.getConcreteItem("get_profile_group", uid, flattened, resolved)
	if err != nil {
		return nil, err
	}
	var g ProfileGroup
	if err := convertRawGroup("profile_group", uid, result, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (c *Client) CreateProfileGroup(g ProfileGroup) (*ProfileGroup, error) {
	// Make sure a profile group with the same name does not already exist
	if exists, err := c.HasItem("profile_group", g.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("a ProfileGroup with the name %s already exists", g.Name)
	}

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
	return c.GetProfileGroup(objectID, false, false)
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

// DeleteProfileGroup deletes a single ProfileGroup by its uid.
func (c *Client) DeleteProfileGroup(uid string) error {
	_, err := c.Call("remove_profile_group", uid, c.Token, false)
	return err
}

// DeleteProfileGroupRecursive deletes a single ProfileGroup by its uid with the option to do so recursively.
func (c *Client) DeleteProfileGroupRecursive(uid string, recursive bool) error {
	_, err := c.Call("remove_profile_group", uid, c.Token, recursive)
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

// GetProfileGroupAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetProfileGroupAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_profile_group_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
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
