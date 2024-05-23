package cobblerclient

// Item general fields
type Item struct {
	Parent            string            `mapstructure:"parent"`
	Depth             int               `mapstructure:"depth"          cobbler:"noupdate"`
	Children          []string          `mapstructure:"children"`
	CTime             float64           `mapstructure:"ctime"          cobbler:"noupdate"`
	MTime             float64           `mapstructure:"mtime"          cobbler:"noupdate"`
	Uid               string            `mapstructure:"uid"            cobbler:"noupdate"`
	Name              string            `mapstructure:"name"`
	Comment           string            `mapstructure:"comment"`
	KernelOptions     map[string]string `mapstructure:"kernel_options"`
	KernelOptionsPost map[string]string `mapstructure:"kernel_options_post"`
	AutoinstallMeta   map[string]string `mapstructure:"autoinstall_meta"`
	FetchableFiles    map[string]string `mapstructure:"fetchable_files"`
	BootFiles         map[string]string `mapstructure:"boot_files"`
	TemplateFiles     map[string]string `mapstructure:"template_files"`
	Owners            []string          `mapstructure:"owners"`
	MgmtClasses       []string          `mapstructure:"mgmt_classes"`
	MgmtParameters    interface{}       `mapstructure:"mgmt_parameters"` // FIXME: This is not a str but a dict
}

// GetItemNames is the method which returns the list of names for a certain object type present inside Cobbler.
func (c *Client) GetItemNames(what string) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("get_item_names", what)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetItemResolvedValue is ...
func (c *Client) GetItemResolvedValue(itemUuid string, attribute string) error {
	_, err := c.Call("get_item_resolved_value", itemUuid, attribute)
	return err
}

// GetItem is ...
func (c *Client) GetItem(what string, name string, flatten bool) error {
	_, err := c.Call("get_item", what, name, flatten)
	return err
}

// FindItems is ...
func (c *Client) FindItems(what string, criteria map[string]interface{}, sortField string, expand bool) error {
	_, err := c.Call("find_items", what, criteria, sortField, expand)
	// TODO: Parse result
	return err
}

// FindItemsPaged is ...
func (c *Client) FindItemsPaged(what string, criteria map[string]interface{}, sortField string, page, itemsPerPage int32) error {
	_, err := c.Call("find_items_paged", what, criteria, sortField, page, itemsPerPage, c.Token)
	// TODO: Parse result
	return err
}

// HasItem is ...
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

// RenameItem is ...
func (c *Client) RenameItem(what, objectId, newName string) error {
	_, err := c.Call("rename_item", what, objectId, newName, c.Token)
	return err
}

// NewItem is ...
func (c *Client) NewItem(what string, isSubobject bool) error {
	_, err := c.Call("new_item", what, c.Token, isSubobject)
	return err
}

// SaveItem is ...
func (c *Client) SaveItem(what, objectId, token, editmode string) error {
	_, err := c.Call("save_item", what, objectId, token, editmode)
	return err
}

// RemoveItem is ...
func (c *Client) RemoveItem(what, name string, recursive bool) error {
	_, err := c.Call("remove_item", what, name, c.Token, recursive)
	return err
}

// CopyItem is ...
func (c *Client) CopyItem(what, objectId, newName string) error {
	_, err := c.Call("copy_item", what, objectId, newName, c.Token)
	return err
}