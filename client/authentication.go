package client

import "errors"

func convertIntBool(integer int) (bool, error) {
	if integer == 0 {
		return false, nil
	}
	if integer == 1 {
		return true, nil
	}
	return false, errors.New("integer was neither 0 nor 1")
}

func convertToInt(integer interface{}) (int, error) {
	switch integer.(type) {
	case int8:
		return int(integer.(int8)), nil
	case int16:
		return int(integer.(int16)), nil
	case int32:
		return int(integer.(int32)), nil
	case int64:
		return int(integer.(int64)), nil
	default:
		return -1, errors.New("integer could not be converted")
	}
}

// CheckAccessNoFail is ...
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

// CheckAccess is ...
func (c *Client) CheckAccess(resource, arg1, arg2 string) (int, error) {
	result, err := c.Call("check_access", c.Token, resource, arg1, arg2)
	if err != nil {
		return -1, err
	} else {
		switch result.(type) {
		case int8:
			return int(result.(int8)), err
		case int16:
			return int(result.(int16)), err
		case int32:
			return int(result.(int32)), err
		case int64:
			return int(result.(int64)), err
		default:
			return -1, errors.New("integer could not be converted")
		}
	}
}

// GetAuthnModuleName is ...
func (c *Client) GetAuthnModuleName() (string, error) {
	result, err := c.Call("get_authn_module_name", c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// Login will perform a login request to Cobbler using the credentials provided
// in the configuration in the initializer.
func (c *Client) Login() (bool, error) {
	result, err := c.Call("login", c.config.Username, c.config.Password)
	if err != nil {
		return false, err
	}

	c.Token = result.(string)
	return true, nil
}

// Logout is ...
func (c *Client) Logout() (bool, error) {
	result, err := c.Call("logout", c.Token)
	if err != nil {
		return false, err
	} else {
		return result.(bool), err
	}
}

// TokenCheck is ...
func (c *Client) TokenCheck(token string) (bool, error) {
	result, err := c.Call("token_check", token)
	if err != nil {
		return false, err
	} else {
		return result.(bool), err
	}
}

// GetUserFromToken is ...
func (c *Client) GetUserFromToken(token string) (string, error) {
	result, err := c.Call("get_user_from_token", token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}
