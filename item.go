package cobblerclient

import (
	"errors"
	"strings"
)

const inherit string = "<<inherit>>"
const none string = "none"

// Value is a helper struct that wraps the multi-typed values being returned from the Cobbler API.
type Value[T any] struct {
	// Data contains the unresolved or resolved data of the attribute. This is not set in case the value is flattened.
	Data T
	// FlattenedValue contains the unresolved or resolved flattened data of the Attribute. This is not set in case the
	// value is not flattened.
	FlattenedValue string
	// IsInherited is a flag that signals if the attribute is inherited or not. If this flag is true then both Data and
	// FlattenedValue are not set.
	IsInherited bool
	// RawData contains the data as received by the API. This field is not evaluated when updating an Item via the
	// API.
	RawData interface{}
}

type ItemMeta struct {
	IsFlattened bool `json:"is_flattened" yaml:"is_flattened"`
	IsResolved  bool `json:"is_resolved" yaml:"is_resolved"`
	// This flag signals if the item was modified by a called method server-side.
	IsDirty bool `json:"is_dirty" yaml:"is_dirty"`
}

// Item general fields
type Item struct {
	// Meta information about an item
	Meta ItemMeta `cobbler:"noupdate" json:"meta" yaml:"meta"`

	// Item fields
	Parent            string                        `mapstructure:"parent" json:"parent" yaml:"parent"`
	Depth             int                           `mapstructure:"depth" cobbler:"noupdate" json:"depth" yaml:"depth"`
	Children          []string                      `mapstructure:"children"       cobbler:"noupdate" json:"children" yaml:"children"`
	CTime             float64                       `mapstructure:"ctime"          cobbler:"noupdate" json:"ctime" yaml:"ctime"`
	MTime             float64                       `mapstructure:"mtime"          cobbler:"noupdate" json:"mtime" yaml:"mtime"`
	Uid               string                        `mapstructure:"uid"            cobbler:"noupdate" json:"uid" yaml:"uid"`
	Name              string                        `mapstructure:"name" json:"name" yaml:"name"`
	Comment           string                        `mapstructure:"comment" json:"comment" yaml:"comment"`
	KernelOptions     Value[map[string]interface{}] `mapstructure:"kernel_options" json:"kernel_options" yaml:"kernel_options"`
	KernelOptionsPost Value[map[string]interface{}] `mapstructure:"kernel_options_post" json:"kernel_options_post" yaml:"kernel_options_post"`
	AutoinstallMeta   Value[map[string]interface{}] `mapstructure:"autoinstall_meta" json:"autoinstall_meta" yaml:"autoinstall_meta"`
	// TemplateFiles is not inheritable (cobbler/items/abstract/bootable_item.py: template_files is a plain
	// @LazyProperty, unlike kernel_options/kernel_options_post/autoinstall_meta).
	TemplateFiles map[string]string `mapstructure:"template_files" json:"template_files" yaml:"template_files"`
	Owners        Value[[]string]   `mapstructure:"owners" json:"owners" yaml:"owners"`
}

// NewItem is a method to initialize the struct with the values that the server-side would internally use. Using this is
// important since the client overwrites all fields with those chosen locally inside the item.
func NewItem() Item {
	return Item{
		AutoinstallMeta: Value[map[string]interface{}]{
			Data: make(map[string]interface{}),
		},
		Children: make([]string, 0),
		KernelOptions: Value[map[string]interface{}]{
			Data: make(map[string]interface{}),
		},
		KernelOptionsPost: Value[map[string]interface{}]{
			Data: make(map[string]interface{}),
		},
		Owners: Value[[]string]{
			Data:        make([]string, 0),
			IsInherited: true,
		},
		TemplateFiles: make(map[string]string),
	}
}

// ModifyItem is a generic method to modify items. Changes made with this method are not persisted until a call to
// SaveItem or one of its other concrete methods.
func (c *Client) ModifyItem(what, objectId string, attribute []string, arg interface{}) error {
	_, err := c.Call("modify_item", what, objectId, attribute, arg, c.Token)
	return err
}

// ModifyItemInPlace attempts to recreate the functionality of the "in_place" parameter for the "xapi_object_edit"
// XML-RPC method.
func (c *Client) ModifyItemInPlace(what, name, attribute string, value map[string]interface{}) error {
	itemKey := []string{
		"autoinstall_meta",
		"kernel_options",
		"kernel_options_post",
		"template_files",
		"params",
	}
	if !stringInSlice(attribute, itemKey) {
		return errors.New("invalid attribute for in-place modification")
	}
	rawItem, err := c.GetItem(what, name, false, false)
	if err != nil {
		return err
	}
	newMapInterface, keyExists := rawItem[attribute]
	if !keyExists {
		return errors.New("attribute not found in ")
	}
	newMap, castSuccessful := newMapInterface.(map[string]interface{})
	if !castSuccessful {
		return errors.New("failed to cast to map[string]interface{}")
	}
	if newMap == nil {
		// An empty XML-RPC struct decodes to a nil map rather than an empty one.
		newMap = map[string]interface{}{}
	}
	for key, mapValue := range value {
		if strings.HasPrefix(key, "~") && len(key) > 1 {
			delete(newMap, key[1:])
		} else {
			newMap[key] = mapValue
		}
	}
	itemHandle, err := c.GetItemHandle(what, name)
	if err != nil {
		return err
	}
	err = c.ModifyItem(what, itemHandle, []string{attribute}, newMap)
	if err != nil {
		return err
	}
	return c.SaveItem(what, itemHandle, true, true, "bypass")
}

// GetItemNames returns the list of names for a specified object type present inside Cobbler.
func (c *Client) GetItemNames(what string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_item_names", what)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetItemResolvedValue returns the value of a single attribute of a single
// item with the inheritance chain applied. attribute is the dotted path to
// the property (e.g. []string{"ipv4", "address"} for a NetworkInterface).
// Cobbler 4.0.0 returns Union[str, int, float, List, dict].
func (c *Client) GetItemResolvedValue(itemUuid string, attribute []string) (interface{}, error) {
	return c.Call("get_item_resolved_value", itemUuid, attribute)
}

// SetItemResolvedValue sets a single attribute on an item with inheritance
// rules applied. attribute is the dotted path; value is the desired value.
// Added in Cobbler 4.0.0.
func (c *Client) SetItemResolvedValue(itemUuid string, attribute []string, value interface{}) error {
	_, err := c.Call("set_item_resolved_value", itemUuid, attribute, value, c.Token)
	return err
}

// GetItem retrieves a single item from the database. An empty map means that the item could not be found.
func (c *Client) GetItem(what string, name string, flatten, resolved bool) (map[string]interface{}, error) {
	unmarshalledResult, err := c.Call("get_item", what, name, flatten, resolved)
	if err != nil {
		return nil, err
	}
	marshalledResult, marshallSuccessful := unmarshalledResult.(map[string]interface{})
	if !marshallSuccessful {
		notFoundMarker, marshallSuccessful := unmarshalledResult.(string)
		if !marshallSuccessful {
			return nil, errors.New("marshall to map unsuccessful and not-found marker not detected")
		}
		if notFoundMarker == "~" {
			return make(map[string]interface{}), nil
		}
	}
	return marshalledResult, nil
}

// getConcreteItem dispatches the get_<item> XML-RPC call. Cobbler 4.0.0 always
// supports the `resolved` parameter (added in 3.3.3); the v1.x client targets
// 4.0.0+ only and so always sends it. The CachedVersion field on Client is
// kept around for compatibility with [Client.ExtendedVersion] consumers but
// no longer gates parameter shape.
func (c *Client) getConcreteItem(method, name string, flattened, resolved bool) (interface{}, error) {
	return c.Call(method, name, flattened, resolved, c.Token)
}

// FindItems searches for one or more items by any of its attributes.
func (c *Client) FindItems(what string, criteria map[string]interface{}, sortField string, expand bool) ([]interface{}, error) {
	unmarshalledResult, err := c.Call("find_items", what, criteria, sortField, expand, false, c.Token)
	return unmarshalledResult.([]interface{}), err
}

func (c *Client) FindItemNames(what string, criteria map[string]interface{}, sortField string) ([]string, error) {
	unmarshalledResult, err := c.Call("find_items", what, criteria, sortField, false, false, c.Token)
	return returnStringSlice(unmarshalledResult, err)
}

type PageInfo struct {
	Page             int   `mapstructure:"page"`
	PrevPage         int   `mapstructure:"prev_page"`
	NextPage         int   `mapstructure:"next_page"`
	Pages            []int `mapstructure:"pages"`
	NumPages         int   `mapstructure:"num_pages"`
	NumItems         int   `mapstructure:"num_items"`
	StartItem        int   `mapstructure:"start_item"`
	EndItem          int   `mapstructure:"end_item"`
	ItemsPerPage     int   `mapstructure:"items_per_page"`
	ItemsPerPageList []int `mapstructure:"items_per_page_list"`
}

type PagedSearchResult struct {
	FoundItems []interface{} `mapstructure:"items"`
	PageInfo   PageInfo      `mapstructure:"pageinfo"`
}

// FindItemsPaged searches for items with the given criteria and returning
func (c *Client) FindItemsPaged(what string, criteria map[string]interface{}, sortField string, page, itemsPerPage int32) (*PagedSearchResult, error) {
	var pagedSearchResult PagedSearchResult
	unmarshalledResult, err := c.Call("find_items_paged", what, criteria, sortField, page, itemsPerPage, c.Token)
	if err != nil {
		return nil, err
	}
	parsedResult, err := decodeCobblerItem(unmarshalledResult, &pagedSearchResult)
	if err != nil {
		return nil, err
	}
	return parsedResult.(*PagedSearchResult), err
}

// HasItem checks if an item with the given name exists.
func (c *Client) HasItem(what string, name string) (bool, error) {
	result, err := c.Call("has_item", what, name, c.Token)
	return result.(bool), err
}

// GetItemHandle gets the internal ID of a Cobbler item.
func (c *Client) GetItemHandle(what, name string) (string, error) {
	result, err := c.Call("get_item_handle", what, name, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}

// RenameItem renames an item.
func (c *Client) RenameItem(what, objectId, newName string) error {
	_, err := c.Call("rename_item", what, objectId, newName, c.Token)
	return err
}

// NewItem creates a new empty item that has to be filled with data. The item does not exist in the database
// before [Client.SaveItem] was called.
func (c *Client) NewItem(what string, isSubobject bool) error {
	_, err := c.Call("new_item", what, c.Token, isSubobject)
	return err
}

// SaveItem saves the changes done via XML-RPC.
func (c *Client) SaveItem(what, objectId string, withTriggers, withSync bool, editmode string) error {
	_, err := c.Call("save_item", what, objectId, withTriggers, withSync, editmode, c.Token)
	return err
}

// RemoveItem deletes an item from the Cobbler database.
func (c *Client) RemoveItem(what, name string, recursive bool) error {
	_, err := c.Call("remove_item", what, name, c.Token, recursive)
	return err
}

// CopyItem duplicates an item on the server with a new name.
func (c *Client) CopyItem(what, objectId, newName string) error {
	_, err := c.Call("copy_item", what, objectId, newName, c.Token)
	return err
}
