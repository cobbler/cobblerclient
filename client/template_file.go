package client

// TemplateFile is the former TemplateFile
type TemplateFile struct {
	Name string // The name the template file will be saved in Cobbler
	Body string // The contents of the template file
}

// CreateTemplateFile to create a template file in Cobbler.
// Takes a TemplateFile struct as input.
// Requires 3 arguments: file, data and token
// Returns true/false and error if creation failed.
func (c *Client) CreateTemplateFile(f TemplateFile) error {
	_, err := c.Call("write_autoinstall_template", f.Name, f.Body, c.Token) // TODO: check name
	return err
}

// GetTemplateFile to get a template file in Cobbler.
// Takes a template file name as input.
// Requires 2 arguments: short filename and token
// Returns *TemplateFile and error if read failed.
func (c *Client) GetTemplateFile(ksName string) (*TemplateFile, error) {
	result, err := c.Call("read_autoinstall_template", ksName, c.Token) // TODO: check name

	if err != nil {
		return nil, err
	}

	ks := TemplateFile{
		Name: ksName,
		Body: result.(string),
	}

	return &ks, nil
}

// DeleteTemplateFile to delete a template file in Cobbler.
// Takes a template file name as input.
// Requires 2 arguments: short filename and token
// Returns error if delete failed.
func (c *Client) DeleteTemplateFile(name string) error {
	_, err := c.Call("remove_autoinstall_template", name, c.Token) // TODO: check name
	return err
}
