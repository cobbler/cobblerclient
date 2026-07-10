/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

import "fmt"

// InputBoolean parses the given value (string, int, bool) and returns
// whether Cobbler considers it truthy. If Cobbler resolves the value to the
// inheritance sentinel ("<<inherit>>"), isInherited is true and result is the
// zero value — this is not an error.
func (c *Client) InputBoolean(value interface{}) (result bool, isInherited bool, err error) {
	res, err := c.Call("input_boolean", value)
	if err != nil {
		return false, false, err
	}
	if c.IsValueInherit(res) {
		return false, true, nil
	}
	b, ok := res.(bool)
	if !ok {
		return false, false, fmt.Errorf("input_boolean returned %T, want bool", res)
	}
	return b, false, nil
}

// InputInt parses the given value into an integer. If Cobbler resolves the
// value to the inheritance sentinel ("<<inherit>>"), isInherited is true and
// result is the zero value — this is not an error.
func (c *Client) InputInt(value interface{}) (result int, isInherited bool, err error) {
	res, err := c.Call("input_int", value)
	if err != nil {
		return 0, false, err
	}
	if c.IsValueInherit(res) {
		return 0, true, nil
	}
	i, err := convertToInt(res)
	if err != nil {
		return 0, false, err
	}
	return i, false, nil
}

// InputStringOrList parses a string ("a,b,c") or list into a Go []string. If
// Cobbler resolves the value to the inheritance sentinel ("<<inherit>>"),
// isInherited is true and result is nil — this is not an error.
func (c *Client) InputStringOrList(value interface{}) (result []string, isInherited bool, err error) {
	return c.inputStringOrList("input_string_or_list", value)
}

// InputStringOrListNoInherit parses like InputStringOrList but instructs
// Cobbler to reject the inheritance sentinel as input. Added in 4.0.0. If the
// server nonetheless returns the sentinel anyway, isInherited is reported the
// same way — this is not an error.
func (c *Client) InputStringOrListNoInherit(value interface{}) (result []string, isInherited bool, err error) {
	return c.inputStringOrList("input_string_or_list_no_inherit", value)
}

func (c *Client) inputStringOrList(method string, value interface{}) (result []string, isInherited bool, err error) {
	res, err := c.Call(method, value)
	if err != nil {
		return nil, false, err
	}
	if c.IsValueInherit(res) {
		return nil, true, nil
	}
	raw, ok := res.([]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%s returned %T, want list", method, res)
	}
	out := make([]string, len(raw))
	for i, v := range raw {
		s, ok := v.(string)
		if !ok {
			return nil, false, fmt.Errorf("%s element %d is %T, want string", method, i, v)
		}
		out[i] = s
	}
	return out, false, nil
}

// InputStringOrDict parses a string ("k=v k2=v2") or dict into a Go map. If
// Cobbler resolves the value to the inheritance sentinel ("<<inherit>>"),
// isInherited is true and result is nil — this is not an error.
func (c *Client) InputStringOrDict(value interface{}) (result map[string]interface{}, isInherited bool, err error) {
	return c.inputStringOrDict("input_string_or_dict", value)
}

// InputStringOrDictNoInherit parses like InputStringOrDict but instructs
// Cobbler to reject the inheritance sentinel as input. Added in 4.0.0. If the
// server nonetheless returns the sentinel anyway, isInherited is reported the
// same way — this is not an error.
func (c *Client) InputStringOrDictNoInherit(value interface{}) (result map[string]interface{}, isInherited bool, err error) {
	return c.inputStringOrDict("input_string_or_dict_no_inherit", value)
}

func (c *Client) inputStringOrDict(method string, value interface{}) (result map[string]interface{}, isInherited bool, err error) {
	res, err := c.Call(method, value)
	if err != nil {
		return nil, false, err
	}
	if c.IsValueInherit(res) {
		return nil, true, nil
	}
	m, ok := res.(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%s returned %T, want map", method, res)
	}
	return m, false, nil
}
