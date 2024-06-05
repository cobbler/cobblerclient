package cobblerclient

// GetSignatures retrieves the complete signatures that are loaded by Cobbler.
func (c *Client) GetSignatures() error {
	// TODO: Create object for signatures
	_, err := c.Call("get_signatures", c.Token)
	return err
}

// GetValidBreeds retrieves all valid OS breeds that a distro can have.
func (c *Client) GetValidBreeds() ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_breeds", c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetValidOsVersionsForBreed retrieves all valid OS versions for a given breed.
func (c *Client) GetValidOsVersionsForBreed(breed string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_os_versions_for_breed", breed, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetValidOsVersions retrieves all valid OS versions that a distro can have.
func (c *Client) GetValidOsVersions() ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_os_versions", c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetValidArchs retrieves all valid architectures that Cobbler is offering.
func (c *Client) GetValidArchs() ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_archs", c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// BackgroundSignatureUpdate runs a signatures update in the background on the server.
func (c *Client) BackgroundSignatureUpdate() (string, error) {
	res, err := c.Call("background_signature_update", map[string]string{}, c.Token)
	return returnString(res, err)
}
