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

func NewMenu() Menu {
	return Menu{
		Item: NewItem(),
	}
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

	// Now clean the Value structs
	decodedMenu := decodeResult.(*Menu)
	err = sanitizeValueMapStruct(&decodedMenu.KernelOptions)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedMenu.KernelOptionsPost)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedMenu.AutoinstallMeta)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedMenu.FetchableFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedMenu.BootFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedMenu.TemplateFiles)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueMapStruct(&decodedMenu.MgmtParameters)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedMenu.Owners)
	if err != nil {
		return nil, err
	}
	err = sanitizeValueSliceStruct(&decodedMenu.MgmtClasses)
	return decodedMenu, nil
}

func convertRawMenusList(xmlrpcResult interface{}) ([]*Menu, error) {
	var menus []*Menu

	for _, d := range xmlrpcResult.([]interface{}) {
		menu, err := convertRawMenu("unknown", d)
		if err != nil {
			return nil, err
		}
		menu.Meta = ItemMeta{
			IsFlattened: false,
			IsResolved:  false,
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
func (c *Client) GetMenu(name string, flattened, resolved bool) (*Menu, error) {
	result, err := c.getConcreteItem("get_menu", name, flattened, resolved)
	if err != nil {
		return nil, err
	}

	menu, err := convertRawMenu(name, result)
	if err != nil {
		return nil, err
	}
	menu.Meta = ItemMeta{
		IsFlattened: flattened,
		IsResolved:  resolved,
	}
	return menu, nil
}

// CreateMenu creates a menu.
func (c *Client) CreateMenu(menu Menu) (*Menu, error) {
	// Make sure a menu with the same name does not already exist
	if _, err := c.GetMenu(menu.Name, false, false); err == nil {
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

	return c.GetMenu(menu.Name, false, false)
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

// DeleteMenu deletes a single Menu by its name.
func (c *Client) DeleteMenu(name string) error {
	return c.DeleteMenuRecursive(name, false)
}

// DeleteMenuRecursive deletes a single Menu by its name with the option to do so recursively.
func (c *Client) DeleteMenuRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_menu", name, c.Token, recursive)
	return err
}

// ListMenuNames returns a list of all menu names currently available in Cobbler.
func (c *Client) ListMenuNames() ([]string, error) {
	return c.GetItemNames("menu")
}

// FindMenu searches for one or more menus by any of its attributes.
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

// FindMenuNames searches for one or more menus by any of its attributes.
func (c *Client) FindMenuNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_menu", criteria, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetMenuHandle gets the internal ID of a Cobbler item.
func (c *Client) GetMenuHandle(name string) (string, error) {
	result, err := c.Call("get_menu_handle", name, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}

// CopyMenu duplicates a menu on the server with a new name.
func (c *Client) CopyMenu(objectId, newName string) error {
	_, err := c.Call("copy_menu", objectId, newName, c.Token)
	return err
}

// GetMenusSince returns all menus which were created after the specified date.
func (c *Client) GetMenusSince(mtime time.Time) ([]*Menu, error) {
	result, err := c.Call("get_menus_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawMenusList(result)
}

// GetMenuAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetMenuAsRendered() error {
	_, err := c.Call("get_menu_as_rendered")
	return err
}

// SaveMenu saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveMenu(objectId, editmode string) error {
	_, err := c.Call("save_menu", objectId, c.Token, editmode)
	return err
}

// RenameMenu renames a menu with a given object id.
func (c *Client) RenameMenu(objectId, newName string) error {
	_, err := c.Call("rename_menu", objectId, newName, c.Token)
	return err
}
