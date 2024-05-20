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

// Getmgmtclasses returns all mgmtclasses in Cobbler.
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

// ListMgmtClassNames is returning a list of all distro names currently available in Cobbler.
func (c *Client) ListMgmtClassNames() ([]string, error) {
	return c.GetItemNames("mgmtclass")
}

// FindMgmtClass is ...
func (c *Client) FindMgmtClass(criteria map[string]interface{}) ([]*MgmtClass, error) {
	result, err := c.Call("find_mgmtclass", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawMgmtClassList(result)
}

// FindMgmtClassNames is searching for one or more mgmtclasses by any of its attributes.
func (c *Client) FindMgmtClassNames(criteria map[string]interface{}) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("find_mgmtclass", criteria, false, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetMgmtClassHandle gets the internal ID of a Cobbler item.
func (c *Client) GetMgmtClassHandle(name string) (string, error) {
	result, err := c.Call("get_mgmtclass_handle", name, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}

// CopyMgmtClass is ...
func (c *Client) CopyMgmtClass(objectId, newName string) error {
	_, err := c.Call("copy_mgmtclass", objectId, newName, c.Token)
	return err
}

// RenameMgmtClass is ...
func (c *Client) RenameMgmtClass(objectId, newName string) error {
	_, err := c.Call("rename_mgmtclass", objectId, newName, c.Token)
	return err
}

// GetMgmtClassesSince is ...
func (c *Client) GetMgmtClassSince(mtime time.Time) ([]*MgmtClass, error) {
	result, err := c.Call("get_mgmtclasses_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawMgmtClassList(result)
}

// GetMgmtClassAsRendered is ...
func (c *Client) GetMgmtClassAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_mgmtclass_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// SaveMgmtClass is ...
func (c *Client) SaveMgmtClass(objectId, editmode string) error {
	_, err := c.Call("save_mgmtclass", objectId, c.Token, editmode)
	return err
}
