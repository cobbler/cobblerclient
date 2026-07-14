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

type SystemGroup struct {
	Group `mapstructure:",squash" yaml:",inline"`
}

// NewSystemGroup returns a zero-valued SystemGroup with sensible defaults.
func NewSystemGroup() SystemGroup { return SystemGroup{Group: newGroup()} }

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
	// Make sure a system group with the same name does not already exist
	if exists, err := c.HasItem("system_group", g.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("a SystemGroup with the name %s already exists", g.Name)
	}

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

// GetSystemGroupAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetSystemGroupAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_system_group_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
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
