package cobblerclient

// Menu is ...
type Menu struct {
	Item `mapstructure:",squash"`

	// Menu specific fields
	DisplayName string `mapstructure:"display_name"`
}

// FindMenu is ...
func (c *Client) FindMenu(criteria map[string]interface{}) ([]*Menu, error) {
	_, err := c.Call("find_menu")
	return nil, err
}

// GetMenuHandle gets the internal ID of a Cobbler item.
func (c *Client) GetMenuHandle(name string) (string, error) {
	result, err := c.Call("get_menu_handle", name, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// CopyMenu is ...
func (c *Client) CopyMenu(objectId, newName string) error {
	_, err := c.Call("copy_menu", objectId, newName, c.Token)
	return err
}

// GetMenusSince is ...
func (c *Client) GetMenusSince() error {
	_, err := c.Call("get_menus_since")
	return err
}

// GetMenuAsRendered is ...
func (c *Client) GetMenuAsRendered() error {
	_, err := c.Call("get_menu_as_rendered")
	return err
}

// SaveMenu is ...
func (c *Client) SaveMenu(objectId, token, editmode string) error {
	_, err := c.Call("save_menu", objectId, token, editmode)
	return err
}

// RenameMenu is ...
func (c *Client) RenameMenu(objectId, newName string) error {
	_, err := c.Call("rename_menu", objectId, newName, c.Token)
	return err
}
