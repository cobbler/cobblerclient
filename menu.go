package cobblerclient

import (
	"fmt"
	"reflect"
	"time"
)

// Menu is a created menu.
// Get the fields from cobbler/items/menu.py
type Menu struct {
	Item `mapstructure:",squash"`

	// Menu specific fields
	DisplayName string `mapstructure:"display_name"`
}

func convertRawMenu(name string, xmlrpcResult interface{}) (*Menu, error) {
	var menu Menu

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("menu %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &menu)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Menu), nil
}

func convertRawMenusList(xmlrpcResult interface{}) ([]*Menu, error) {
	var menus []*Menu

	for _, d := range xmlrpcResult.([]interface{}) {
		menu, err := convertRawMenu("unknown", d)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}

	return menus, nil
}

// GetMenus returns all menus in Cobbler.
func (c *Client) GetMenus() ([]*Distro, error) {
	result, err := c.Call("get_menus", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawDistrosList(result)
}

// GetMenu returns a single menu obtained by its name.
func (c *Client) GetMenu(name string) (*Menu, error) {
	result, err := c.Call("get_menu", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawMenu(name, result)
}

// CreateMenu creates a menu.
func (c *Client) CreateMenu(menu Menu) (*Menu, error) {
	// Make sure a menu with the same name does not already exist
	if _, err := c.GetMenu(menu.Name); err == nil {
		return nil, fmt.Errorf("a Menu with the name %s already exists", menu.Name)
	}

	result, err := c.Call("new_menu", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&menu).Elem()
	if err := c.updateCobblerFields("menu", item, newID); err != nil {
		return nil, err
	}

	if err = c.SaveMenu(newID, "new"); err != nil {
		return nil, err
	}

	return c.GetMenu(menu.Name)
}

// UpdateMenu updates a single menu.
func (c *Client) UpdateMenu(menu *Menu) error {
	item := reflect.ValueOf(menu).Elem()
	id, err := c.GetItemHandle("menu", menu.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("menu", item, id); err != nil {
		return err
	}

	if err = c.SaveMenu(id, "bypass"); err != nil {
		return err
	}

	return nil
}

// ListMenuNames is returning a list of all menu names currently available in Cobbler.
func (c *Client) ListMenuNames() ([]string, error) {
	return c.GetItemNames("menu")
}

// FindMenu is searching for one or more menus by any of its attributes.
func (c *Client) FindMenu(criteria map[string]interface{}) ([]*Menu, error) {
	var menus []*Menu

	result, err := c.Call("find_menu", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	for _, d := range result.([]interface{}) {
		var menu Menu
		decodedResult, err := decodeCobblerItem(d, &menu)
		if err != nil {
			return nil, err
		}

		menus = append(menus, decodedResult.(*Menu))
	}

	return menus, nil
}

// FindMenuNames is searching for one or more menus by any of its attributes.
func (c *Client) FindMenuNames(criteria map[string]interface{}) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("find_menu", criteria, false, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetMenuHandle gets the internal ID of a Cobbler item.
func (c *Client) GetMenuHandle(name string) (string, error) {
	result, err := c.Call("get_menu_handle", name, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}

// CopyMenu is copying a given menu server side with a new name.
func (c *Client) CopyMenu(objectId, newName string) error {
	_, err := c.Call("copy_menu", objectId, newName, c.Token)
	return err
}

// GetMenusSince is returning all menus which were created after the specified date.
func (c *Client) GetMenusSince(mtime time.Time) ([]*Menu, error) {
	result, err := c.Call("get_menus_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawMenusList(result)
}

// GetMenuAsRendered is returning the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetMenuAsRendered() error {
	_, err := c.Call("get_menu_as_rendered")
	return err
}

// SaveMenu is persisting all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveMenu(objectId, editmode string) error {
	_, err := c.Call("save_menu", objectId, c.Token, editmode)
	return err
}

// RenameMenu is renaming a menu with a given object id.
func (c *Client) RenameMenu(objectId, newName string) error {
	_, err := c.Call("rename_menu", objectId, newName, c.Token)
	return err
}
