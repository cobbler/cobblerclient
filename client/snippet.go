package client

// Snippet is a snippet file
type Snippet struct {
	Name string // The name the snippet file will be saved in Cobbler
	Body string // The contents of the template file
}

// CreateSnippet creates a snippet in Cobbler.
// Takes a Snippet struct as input
// Returns true/false and error if creation failed.
func (c *Client) CreateSnippet(s Snippet) error {
	_, err := c.Call("write_autoinstall_snippet", s.Name, s.Body, c.Token)
	return err
}

// GetSnippet gets a snippet file in Cobbler.
// Takes a snippet file name as input.
// Returns *Snippet and error if read failed.
func (c *Client) GetSnippet(name string) (*Snippet, error) {
	result, err := c.Call("read_autoinstall_snippet", name, c.Token)

	if err != nil {
		return nil, err
	}

	snippet := Snippet{
		Name: name,
		Body: result.(string),
	}

	return &snippet, nil
}

// DeleteSnippet deletes a snippet file in Cobbler.
// Takes a snippet file name as input.
// Returns error if delete failed.
func (c *Client) DeleteSnippet(name string) error {
	_, err := c.Call("remove_autoinstall_snippet", name, c.Token)
	return err
}
