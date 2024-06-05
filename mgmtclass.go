package cobblerclient

import (
	"fmt"
	"reflect"
	"time"
)

type MgmtClass struct {
	Item `mapstructure:",squash"`

	// Mgmtclass specific fields
	IsDefiniton bool                   `mapstructure:"is_definition"`
	Params      map[string]interface{} `mapstructure:"params"`
	ClassName   string                 `mapstructure:"class_name"`
	Files       []string               `mapstructure:"files"`
	Packages    []string               `mapstructure:"packages"`
}

func convertRawMgmtClass(name string, xmlrpcResult interface{}) (*MgmtClass, error) {
	var mgmtclass MgmtClass

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("mgmtclass %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &mgmtclass)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*MgmtClass), nil
}

func convertRawMgmtClassList(xmlrpcResult interface{}) ([]*MgmtClass, error) {
	var mgmtclasses []*MgmtClass

	for _, d := range xmlrpcResult.([]interface{}) {
		mgmtclass, err := convertRawMgmtClass("unknown", d)
		if err != nil {
			return nil, err
		}
		mgmtclasses = append(mgmtclasses, mgmtclass)
	}

	return mgmtclasses, nil
}

// GetMgmtClasses returns all mgmtclasses in Cobbler.
func (c *Client) GetMgmtClasses() ([]*MgmtClass, error) {
	result, err := c.Call("get_mgmtclasses", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawMgmtClassList(result)
}

// GetMgmtClass returns a single mgmtclass obtained by its name.
func (c *Client) GetMgmtClass(name string) (*MgmtClass, error) {
	result, err := c.Call("get_mgmtclass", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawMgmtClass(name, result)
}

// CreateMgmtClass creates a mgmtclass.
func (c *Client) CreateMgmtClass(mgmtclass MgmtClass) (*MgmtClass, error) {
	// Make sure a distro with the same name does not already exist
	if _, err := c.GetMgmtClass(mgmtclass.Name); err == nil {
		return nil, fmt.Errorf("a MgmtClass with the name %s already exists", mgmtclass.Name)
	}

	result, err := c.Call("new_mgmtclass", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&mgmtclass).Elem()
	if err := c.updateCobblerFields("mgmtclass", item, newID); err != nil {
		return nil, err
	}

	if err := c.SaveMgmtClass(newID, "new"); err != nil {
		return nil, err
	}

	return c.GetMgmtClass(mgmtclass.Name)
}

// UpdateMgmtClass updates a single MgmtClass.
func (c *Client) UpdateMgmtClass(distro *Distro) error {
	item := reflect.ValueOf(distro).Elem()
	id, err := c.GetItemHandle("mgmtclass", distro.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("mgmtclass", item, id); err != nil {
		return err
	}

	if err := c.SaveMgmtClass(id, "bypass"); err != nil {
		return err
	}

	return nil
}

// ListMgmtClassNames returns a list of all managementclass names currently available in Cobbler.
func (c *Client) ListMgmtClassNames() ([]string, error) {
	return c.GetItemNames("mgmtclass")
}

// FindMgmtClass searches for one or more managementclasses by any of its attributes.
func (c *Client) FindMgmtClass(criteria map[string]interface{}) ([]*MgmtClass, error) {
	result, err := c.Call("find_mgmtclass", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawMgmtClassList(result)
}

// FindMgmtClassNames searches for one or more managementclasses by any of its attributes.
func (c *Client) FindMgmtClassNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_mgmtclass", criteria, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetMgmtClassHandle gets the internal ID of a Cobbler item.
func (c *Client) GetMgmtClassHandle(name string) (string, error) {
	res, err := c.Call("get_mgmtclass_handle", name, c.Token)
	return returnString(res, err)
}

// CopyMgmtClass copies a given managementclass server side with a new name.
func (c *Client) CopyMgmtClass(objectId, newName string) error {
	_, err := c.Call("copy_mgmtclass", objectId, newName, c.Token)
	return err
}

// RenameMgmtClass renames a managementclass with a given object id.
func (c *Client) RenameMgmtClass(objectId, newName string) error {
	_, err := c.Call("rename_mgmtclass", objectId, newName, c.Token)
	return err
}

// GetMgmtClassesSince returns all managementclasses which were created after the specified date.
func (c *Client) GetMgmtClassesSince(mtime time.Time) ([]*MgmtClass, error) {
	result, err := c.Call("get_mgmtclasses_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawMgmtClassList(result)
}

// GetMgmtClassAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetMgmtClassAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_mgmtclass_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// SaveMgmtClass saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveMgmtClass(objectId, editmode string) error {
	_, err := c.Call("save_mgmtclass", objectId, c.Token, editmode)
	return err
}
