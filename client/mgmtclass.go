package client

// FindMgmtClass is ...
func (c *Client) FindMgmtClass() error {
	_, err := c.Call("find_mgmtclass")
	return err
}

// GetMgmtClassHandle gets the internal ID of a Cobbler item.
func (c *Client) GetMgmtClassHandle(name string) (string, error) {
	result, err := c.Call("get_mgmtclass_handle", name, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
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
func (c *Client) GetMgmtClassesSince() error {
	_, err := c.Call("get_mgmtclasses_since")
	return err
}

// GetMgmtClassAsRendered is ...
func (c *Client) GetMgmtClassAsRendered() error {
	_, err := c.Call("get_mgmtclass_as_rendered")
	return err
}

// SaveMgmtClass is ...
func (c *Client) SaveMgmtClass(objectId, token, editmode string) error {
	_, err := c.Call("save_mgmtclass", objectId, token, editmode)
	return err
}
