package cobblerclient

// CheckAccessNoFail validates if a certain resource can be accessed with the current token. "arg1" and "arg2" have
// different meanings depending on the authorization provider configured server side.
func (c *Client) CheckAccessNoFail(resource, arg1, arg2 string) (bool, error) {
	result, err := c.Call("check_access_no_fail", c.Token, resource, arg1, arg2)
	if err != nil {
		return false, err
	} else {
		convertedInteger, err := convertToInt(result)
		if err != nil {
			return convertIntBool(convertedInteger)
		}
		return false, err
	}
}

// CheckAccess performs the same check as [Client.CheckAccessNoFail] but returning the error message with the
// reason instead of a boolean.
func (c *Client) CheckAccess(resource, arg1, arg2 string) (int, error) {
	result, err := c.Call("check_access", c.Token, resource, arg1, arg2)
	if err != nil {
		return -1, err
	} else {
		convertedInteger, err := convertToInt(result)
		if err != nil {
			return -1, err
		}
		return convertedInteger, err
	}
}

// GetAuthnModuleName retrieves the currently configured authentication module name.
func (c *Client) GetAuthnModuleName() (string, error) {
	res, err := c.Call("get_authn_module_name", c.Token)
	return returnString(res, err)
}

// Login performs a login request to Cobbler using the credentials provided in the configuration in the initializer.
func (c *Client) Login() (bool, error) {
	result, err := c.Call("login", c.config.Username, c.config.Password)
	if err != nil {
		return false, err
	}

	c.Token = result.(string)
	return true, nil
}

// Logout performs a logout from the Cobbler server.
func (c *Client) Logout() (bool, error) {
	res, err := c.Call("logout", c.Token)
	return returnBool(res, err)
}

// TokenCheck returns if a given token is still valid or not.
func (c *Client) TokenCheck(token string) (bool, error) {
	res, err := c.Call("token_check", token)
	return returnBool(res, err)
}

// GetUserFromToken checks what user a given token is belonging to.
func (c *Client) GetUserFromToken(token string) (string, error) {
	res, err := c.Call("get_user_from_token", token)
	return returnString(res, err)
}
