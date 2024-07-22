package cobblerclient

// Item general fields
type Item struct {
	Parent            string      `mapstructure:"parent"`
	Depth             int         `mapstructure:"depth"          cobbler:"noupdate"`
	Children          []string    `mapstructure:"children"`
	CTime             float64     `mapstructure:"ctime"          cobbler:"noupdate"`
	MTime             float64     `mapstructure:"mtime"          cobbler:"noupdate"`
	Uid               string      `mapstructure:"uid"            cobbler:"noupdate"`
	Name              string      `mapstructure:"name"`
	Comment           string      `mapstructure:"comment"`
	KernelOptions     interface{} `mapstructure:"kernel_options"`
	KernelOptionsPost interface{} `mapstructure:"kernel_options_post"`
	AutoinstallMeta   interface{} `mapstructure:"autoinstall_meta"`
	FetchableFiles    interface{} `mapstructure:"fetchable_files"`
	BootFiles         interface{} `mapstructure:"boot_files"`
	TemplateFiles     interface{} `mapstructure:"template_files"`
	Owners            []string    `mapstructure:"owners"`
	MgmtClasses       []string    `mapstructure:"mgmt_classes"`
	MgmtParameters    interface{} `mapstructure:"mgmt_parameters"` // FIXME: This is not a str but a dict
}

// GetItemNames returns the list of names for a specified object type present inside Cobbler.
func (c *Client) GetItemNames(what string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_item_names", what)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetItemResolvedValue retrieves the value of a single attribute of a single item which was passed through the
// inheritance chain of Cobbler.
func (c *Client) GetItemResolvedValue(itemUuid string, attribute string) error {
	_, err := c.Call("get_item_resolved_value", itemUuid, attribute)
	return err
}

// GetItem retrieves a single item from the database.
func (c *Client) GetItem(what string, name string, flatten bool) error {
	_, err := c.Call("get_item", what, name, flatten)
	return err
}

// FindItems searches for one or more items by any of its attributes.
func (c *Client) FindItems(what string, criteria map[string]interface{}, sortField string, expand bool) error {
	_, err := c.Call("find_items", what, criteria, sortField, expand)
	// TODO: Parse result
	return err
}

// FindItemsPaged searches for items with the given criteria and returning
func (c *Client) FindItemsPaged(what string, criteria map[string]interface{}, sortField string, page, itemsPerPage int32) error {
	_, err := c.Call("find_items_paged", what, criteria, sortField, page, itemsPerPage, c.Token)
	// TODO: Parse result
	return err
}

// HasItem checks if an item with the given name exists.
func (c *Client) HasItem(what string, name string) (bool, error) {
	result, err := c.Call("has_item", what, name, c.Token)
	return result.(bool), err
}

// GetItemHandle gets the internal ID of a Cobbler item.
func (c *Client) GetItemHandle(what, name string) (string, error) {
	result, err := c.Call("get_item_handle", what, name, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}

// RenameItem renames an item.
func (c *Client) RenameItem(what, objectId, newName string) error {
	_, err := c.Call("rename_item", what, objectId, newName, c.Token)
	return err
}

// NewItem creates a new empty item that has to be filled with data. The item does not exist in the database
// before [Client.SaveItem] was called.
func (c *Client) NewItem(what string, isSubobject bool) error {
	_, err := c.Call("new_item", what, c.Token, isSubobject)
	return err
}

// SaveItem saves the changes done via XML-RPC.
func (c *Client) SaveItem(what, objectId, token, editmode string) error {
	_, err := c.Call("save_item", what, objectId, token, editmode)
	return err
}

// RemoveItem deletes an item from the Cobbler database.
func (c *Client) RemoveItem(what, name string, recursive bool) error {
	_, err := c.Call("remove_item", what, name, c.Token, recursive)
	return err
}

// CopyItem duplicates an item on the server with a new name.
func (c *Client) CopyItem(what, objectId, newName string) error {
	_, err := c.Call("copy_item", what, objectId, newName, c.Token)
	return err
}
