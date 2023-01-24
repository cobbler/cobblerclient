package client

// GetSignatures is ...
func (c *Client) GetSignatures() error {
	// TODO: Create object for signatures
	_, err := c.Call("get_signatures", c.Token)
	return err
}

// GetValidBreeds is ...
func (c *Client) GetValidBreeds() ([]string, error) {
	var result []string
	resultUnmarshalled, err := c.Call("get_valid_breeds", c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetValidOsVersionsForBreed is ...
func (c *Client) GetValidOsVersionsForBreed(breed string) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("get_valid_os_versions_for_breed", breed, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetValidOsVersions is ...
func (c *Client) GetValidOsVersions() ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("get_valid_os_versions", c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetValidArchs is ...
func (c *Client) GetValidArchs() ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("get_valid_archs", c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

func (c *Client) BackgroundSignatureUpdate() (string, error) {
	result, err := c.Call("background_signature_update", map[string]string{}, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}
