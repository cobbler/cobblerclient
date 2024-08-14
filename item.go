package cobblerclient

import (
	"errors"
	"strings"
)

// Item general fields
type Item struct {
	Parent            string      `mapstructure:"parent"`
	Depth             int         `mapstructure:"depth"          cobbler:"noupdate"`
	Children          []string    `mapstructure:"children"`
	CTime             float64     `mapstructure:"ctime"          cobbler:"noupdate"`
	MTime             float64     `mapstructure:"mtime"          cobbler:"noupdate"`
	Uid               string      `mapstructure:"uid"            cobbler:"noupdate"`
	Name              string      `mapstructure:"name"`
	Comment           string      `mapstructure:"comment"`
	KernelOptions     interface{} `mapstructure:"kernel_options"`
	KernelOptionsPost interface{} `mapstructure:"kernel_options_post"`
	AutoinstallMeta   interface{} `mapstructure:"autoinstall_meta"`
	FetchableFiles    interface{} `mapstructure:"fetchable_files"`
	BootFiles         interface{} `mapstructure:"boot_files"`
	TemplateFiles     interface{} `mapstructure:"template_files"`
	Owners            []string    `mapstructure:"owners"`
	MgmtClasses       []string    `mapstructure:"mgmt_classes"`
	MgmtParameters    interface{} `mapstructure:"mgmt_parameters"`
}

// ModifyItem is a generic method to modify items. Changes made with this method are not persisted until a call to
// SaveItem or one of its other concrete methods.
func (c *Client) ModifyItem(what, objectId, attribute string, arg interface{}) error {
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
		"boot_files",
		"fetchable_files",
		"params",
	}
	if !stringInSlice(attribute, itemKey) {
		return errors.New("invalid attribute")
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
	err = c.ModifyItem(what, itemHandle, attribute, newMap)
	if err != nil {
		return err
	}
	return c.SaveItem(what, itemHandle, c.Token, "bypass")
}

// GetItemNames returns the list of names for a specified object type present inside Cobbler.
func (c *Client) GetItemNames(what string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_item_names", what)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetItemResolvedValue retrieves the value of a single attribute of a single item which was passed through the
// inheritance chain of Cobbler.
func (c *Client) GetItemResolvedValue(itemUuid string, attribute string) error {
	_, err := c.Call("get_item_resolved_value", itemUuid, attribute)
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

// FindItems searches for one or more items by any of its attributes.
func (c *Client) FindItems(what string, criteria map[string]interface{}, sortField string, expand bool) ([]interface{}, error) {
	unmarshalledResult, err := c.Call("find_items", what, criteria, sortField, expand)
	return unmarshalledResult.([]interface{}), err
}

func (c *Client) FindItemNames(what string, criteria map[string]interface{}, sortField string) ([]string, error) {
	unmarshalledResult, err := c.Call("find_items", what, criteria, sortField, false)
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
func (c *Client) SaveItem(what, objectId, token, editmode string) error {
	_, err := c.Call("save_item", what, objectId, token, editmode)
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
