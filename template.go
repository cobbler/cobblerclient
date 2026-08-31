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

import (
	"fmt"
	"reflect"
	"time"
)

// Template is a Cobbler 4.0.0 first-class item type representing an
// autoinstall template (Jinja2, Cheetah, etc.). See cobbler/items/template.py.
type Template struct {
	Item `mapstructure:",squash" yaml:",inline"`

	TemplateType string    `mapstructure:"template_type" json:"template_type" yaml:"template_type"`
	URI          URIOption `mapstructure:"uri" json:"uri" yaml:"uri"`
	BuiltIn      bool      `mapstructure:"built_in" cobbler:"noupdate" json:"built_in" yaml:"built_in"`
	Tags         []string  `mapstructure:"tags" json:"tags" yaml:"tags"`
	Content      string    `mapstructure:"content" json:"content" yaml:"content"`
}

// NewTemplate returns a zero-valued Template with sensible defaults.
func NewTemplate() Template {
	return Template{
		Item:         NewItem(),
		TemplateType: "jinja",
		URI:          URIOption{Schema: TemplateSchemaFile},
		Tags:         []string{},
	}
}

func convertRawTemplate(name string, xmlrpcResult interface{}) (*Template, error) {
	var t Template
	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("template %s not found", name)
	}
	decoded, err := decodeCobblerItem(xmlrpcResult, &t)
	if err != nil {
		return nil, err
	}
	result, ok := decoded.(*Template)
	if !ok {
		return nil, fmt.Errorf("unexpected decoder result type %T", decoded)
	}
	return result, nil
}

func convertRawTemplatesList(xmlrpcResult interface{}) ([]*Template, error) {
	var out []*Template
	for _, raw := range xmlrpcResult.([]interface{}) {
		t, err := convertRawTemplate("unknown", raw)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (c *Client) GetTemplates() ([]*Template, error) {
	result, err := c.Call("get_templates", "-1", c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawTemplatesList(result)
}

// GetTemplate returns a single template obtained by its uid.
func (c *Client) GetTemplate(uid string, flattened, resolved bool) (*Template, error) {
	result, err := c.getConcreteItem("get_template", uid, flattened, resolved)
	if err != nil {
		return nil, err
	}
	return convertRawTemplate(uid, result)
}

func (c *Client) GetTemplateHandle(name string) (string, error) {
	return c.GetItemHandle("template", name)
}

func (c *Client) GetTemplatesSince(mtime time.Time) ([]*Template, error) {
	result, err := c.Call("get_templates_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}
	return convertRawTemplatesList(result)
}

func (c *Client) FindTemplate(criteria map[string]interface{}, resolved bool) ([]*Template, error) {
	result, err := c.Call("find_template", criteria, true, resolved, c.Token)
	if err != nil {
		return nil, err
	}
	return convertRawTemplatesList(result)
}

func (c *Client) FindTemplateNames(criteria map[string]interface{}) ([]string, error) {
	result, err := c.Call("find_template", criteria, false, false, c.Token)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, n := range result.([]interface{}) {
		names = append(names, n.(string))
	}
	return names, nil
}

func (c *Client) CreateTemplate(tpl Template) (*Template, error) {
	// Make sure a template with the same name does not already exist
	if exists, err := c.HasItem("template", tpl.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("a Template with the name %s already exists", tpl.Name)
	}

	id, err := c.Call("new_template", c.Token)
	if err != nil {
		return nil, err
	}
	objectID, ok := id.(string)
	if !ok {
		return nil, fmt.Errorf("new_template returned %T, want string", id)
	}
	if err := c.updateCobblerFields("template", reflect.ValueOf(&tpl).Elem(), objectID); err != nil {
		return nil, err
	}
	if err := c.SaveTemplate(objectID, true, true, "new"); err != nil {
		return nil, err
	}
	return c.GetTemplate(objectID, false, false)
}

func (c *Client) UpdateTemplate(tpl *Template) error {
	objectID, err := c.GetTemplateHandle(tpl.Name)
	if err != nil {
		return err
	}
	if err := c.updateCobblerFields("template", reflect.ValueOf(tpl).Elem(), objectID); err != nil {
		return err
	}
	return c.SaveTemplate(objectID, true, true, "bypass")
}

// SaveTemplate saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveTemplate(objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_template", objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

func (c *Client) CopyTemplate(objectId, newName string) error {
	_, err := c.Call("copy_template", objectId, newName, c.Token)
	return err
}

// DeleteTemplate deletes a single Template by its uid.
func (c *Client) DeleteTemplate(uid string) error {
	return c.DeleteTemplateRecursive(uid, false)
}

// DeleteTemplateRecursive deletes a single Template by its uid with the option to do so recursively.
func (c *Client) DeleteTemplateRecursive(uid string, recursive bool) error {
	_, err := c.Call("remove_template", uid, c.Token, recursive)
	return err
}

func (c *Client) RenameTemplate(objectId, newName string) error {
	_, err := c.Call("rename_template", objectId, newName, c.Token)
	return err
}

func (c *Client) ListTemplateNames() ([]string, error) {
	return c.GetItemNames("template")
}

// GetTemplateContent returns the in-memory rendered content of a Template by UID.
func (c *Client) GetTemplateContent(uid string) (string, error) {
	result, err := c.Call("get_template_content", uid, c.Token)
	return returnString(result, err)
}

// TemplatesRefreshContent forces a synchronous reload of the listed templates'
// content. Pass nil or an empty slice to refresh all known templates.
func (c *Client) TemplatesRefreshContent(objects []string) error {
	if objects == nil {
		objects = []string{}
	}
	_, err := c.Call("templates_refresh_content", objects, c.Token)
	return err
}

// BackgroundTemplatesRefreshContent kicks off an async content refresh and
// returns the resulting event id.
func (c *Client) BackgroundTemplatesRefreshContent(objects []string) (string, error) {
	if objects == nil {
		objects = []string{}
	}
	options := map[string]interface{}{"objects": objects}
	result, err := c.Call("background_templates_refresh_content", options, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// GetTemplateFileForProfile returns the rendered template file content for a profile.
//
// profileUid must be a Cobbler UID, not a name (Cobbler >=4.0.0b6's
// get_template_file_for_profile does a strict uid-keyed lookup server-side;
// an unresolvable uid returns a "# object not found" comment string, not an
// error). Normally reached anonymously via the /cblr/svc/op/template HTTP
// endpoint, which resolves a name to a uid server-side before calling this
// RPC — direct callers of this wrapper must resolve the uid themselves.
func (c *Client) GetTemplateFileForProfile(profileUid, path string) (string, error) {
	result, err := c.Call("get_template_file_for_profile", profileUid, path)
	return returnString(result, err)
}

// GetTemplateFileForSystem returns the rendered template file content for a system.
//
// systemUid must be a Cobbler UID, not a name (Cobbler >=4.0.0b6's
// get_template_file_for_system does a strict uid-keyed lookup server-side;
// an unresolvable uid returns a "# object not found" comment string, not an
// error). Normally reached anonymously via the /cblr/svc/op/template HTTP
// endpoint, which resolves a name to a uid server-side before calling this
// RPC — direct callers of this wrapper must resolve the uid themselves.
func (c *Client) GetTemplateFileForSystem(systemUid, path string) (string, error) {
	result, err := c.Call("get_template_file_for_system", systemUid, path)
	return returnString(result, err)
}
