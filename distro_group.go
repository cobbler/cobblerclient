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

type DistroGroup struct {
	Group `mapstructure:",squash" yaml:",inline"`
}

// NewDistroGroup returns a zero-valued DistroGroup with sensible defaults.
func NewDistroGroup() DistroGroup { return DistroGroup{Group: newGroup()} }

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

// GetDistroGroup returns a single distro group obtained by its uid.
func (c *Client) GetDistroGroup(uid string, flattened, resolved bool) (*DistroGroup, error) {
	result, err := c.getConcreteItem("get_distro_group", uid, flattened, resolved)
	if err != nil {
		return nil, err
	}
	var g DistroGroup
	if err := convertRawGroup("distro_group", uid, result, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (c *Client) CreateDistroGroup(g DistroGroup) (*DistroGroup, error) {
	// Make sure a distro group with the same name does not already exist
	if exists, err := c.HasItem("distro_group", g.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("a DistroGroup with the name %s already exists", g.Name)
	}

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
	return c.GetDistroGroup(objectID, false, false)
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

// DeleteDistroGroup deletes a single DistroGroup by its uid.
func (c *Client) DeleteDistroGroup(uid string) error {
	_, err := c.Call("remove_distro_group", uid, c.Token, false)
	return err
}

// DeleteDistroGroupRecursive deletes a single DistroGroup by its uid with the option to do so recursively.
func (c *Client) DeleteDistroGroupRecursive(uid string, recursive bool) error {
	_, err := c.Call("remove_distro_group", uid, c.Token, recursive)
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

// GetDistroGroupAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
//
// uid must be a Cobbler UID, not a name (Cobbler >=4.0.0b6's
// get_distro_group_as_rendered does a strict uid-keyed lookup server-side;
// an unresolvable uid silently returns an empty map, not an error).
func (c *Client) GetDistroGroupAsRendered(uid string) (map[string]interface{}, error) {
	result, err := c.Call("get_distro_group_as_rendered", uid, c.Token)
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
