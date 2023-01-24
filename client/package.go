package client

// FindPackage is ...
func (c *Client) FindPackage() error {
	_, err := c.Call("find_package")
	return err
}

// GetPackageHandle gets the internal ID of a Cobbler item.
func (c *Client) GetPackageHandle(name string) (string, error) {
	result, err := c.Call("get_package_handle", name, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// CopyPackage is ...
func (c *Client) CopyPackage(objectId, newName string) error {
	_, err := c.Call("copy_package", objectId, newName, c.Token)
	return err
}

// RenamePackage is ...
func (c *Client) RenamePackage(objectId, newName string) error {
	_, err := c.Call("rename_package", objectId, newName, c.Token)
	return err
}

// GetPackagesSince is ...
func (c *Client) GetPackagesSince() error {
	_, err := c.Call("get_packages_since")
	return err
}

// GetPackageAsRendered is ...
func (c *Client) GetPackageAsRendered() error {
	_, err := c.Call("get_package_as_rendered")
	return err
}

// SavePackage is ...
func (c *Client) SavePackage(objectId, token, editmode string) error {
	_, err := c.Call("save_package", objectId, token, editmode)
	return err
}
