package client

// FindFile is ...
func (c *Client) FindFile() error {
	_, err := c.Call("find_file")
	return err
}

// GetFileHandle gets the internal ID of a Cobbler item.
func (c *Client) GetFileHandle(name string) (string, error) {
	result, err := c.Call("get_file_handle", name, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// CopyFile is ...
func (c *Client) CopyFile(objectId, newName string) error {
	_, err := c.Call("copy_file", objectId, newName, c.Token)
	return err
}

// GetFilesSince is ...
func (c *Client) GetFilesSince() error {
	_, err := c.Call("get_files_since")
	return err
}

// GetFileAsRendered is ...
func (c *Client) GetFileAsRendered() error {
	_, err := c.Call("get_file_as_rendered")
	return err
}

// SaveFile is ...
func (c *Client) SaveFile(objectId, token, editmode string) error {
	_, err := c.Call("save_file", objectId, token, editmode)
	return err
}

// RenameFile is ...
func (c *Client) RenameFile(objectId, newName string) error {
	_, err := c.Call("rename_file", objectId, newName, c.Token)
	return err
}
