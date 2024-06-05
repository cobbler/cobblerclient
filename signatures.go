package cobblerclient

// GetSignatures retrieves the complete signatures that are loaded by Cobbler.
func (c *Client) GetSignatures() error {
	// TODO: Create object for signatures
	_, err := c.Call("get_signatures", c.Token)
	return err
}

// GetValidBreeds retrieves all valid OS breeds that a distro can have.
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

// GetValidOsVersionsForBreed retrieves all valid OS versions for a given breed.
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

// GetValidOsVersions retrieves all valid OS versions that a distro can have.
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

// GetValidArchs retrieves all valid architectures that Cobbler is offering.
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

// BackgroundSignatureUpdate runs a signatures update in the background on the server.
func (c *Client) BackgroundSignatureUpdate() (string, error) {
	result, err := c.Call("background_signature_update", map[string]string{}, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}
