package cobblerclient

import (
	"fmt"
	"reflect"
	"time"
)

// File is a created file.
// Get the fields from cobbler/items/file.py
type File struct {
	Resource `mapstructure:",squash"`

	// File specific fields
	IsDir bool `mapstructure:"is_dir"`
}

func convertRawFile(name string, xmlrpcResult interface{}) (*File, error) {
	var file File

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &file)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*File), nil
}

func convertRawFilesList(xmlrpcResult interface{}) ([]*File, error) {
	var files []*File

	for _, d := range xmlrpcResult.([]interface{}) {
		file, err := convertRawFile("unknown", d)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// GetFiles returns a list of all files.
func (c *Client) GetFiles() ([]*File, error) {
	result, err := c.Call("get_files", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawFilesList(result)
}

// GetFile returns a single file obtained by its name.
func (c *Client) GetFile(name string) (*File, error) {
	result, err := c.Call("get_file", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawFile(name, result)
}

// CreateFile creates a single file.
func (c *Client) CreateFile(file File) (*File, error) {
	// Make sure a file with the same name does not already exist
	if _, err := c.GetFile(file.Name); err == nil {
		return nil, fmt.Errorf("a File with the name %s already exists", file.Name)
	}

	result, err := c.Call("new_file", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&file).Elem()
	if err := c.updateCobblerFields("file", item, newID); err != nil {
		return nil, err
	}

	if err = c.SaveFile(newID, "new"); err != nil {
		return nil, err
	}

	return c.GetFile(file.Name)
}

// UpdateFile updates a single file.
func (c *Client) UpdateFile(file *File) error {
	item := reflect.ValueOf(file).Elem()
	id, err := c.GetItemHandle("file", file.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("file", item, id); err != nil {
		return err
	}

	if err = c.SaveFile(id, "bypass"); err != nil {
		return err
	}

	return nil
}

// DeleteFile deletes a single file by its name.
func (c *Client) DeleteFile(name string) error {
	_, err := c.Call("remove_file", name, c.Token)
	return err
}

// ListDistroNames returns a list of all files names currently available in Cobbler.
func (c *Client) ListFileNames() ([]string, error) {
	return c.GetItemNames("file")
}

// FindFile searches for one or more files by any of its attributes.
func (c *Client) FindFile(criteria map[string]interface{}) ([]*File, error) {
	result, err := c.Call("find_file", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawFilesList(result)
}

// FindFileNames searches for one or more files by any of its attributes.
func (c *Client) FindFileNames(criteria map[string]interface{}) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("find_file", criteria, false, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetFileHandle gets the internal ID of a Cobbler item.
func (c *Client) GetFileHandle(name string) (string, error) {
	result, err := c.Call("get_file_handle", name, c.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// CopyFile duplicates a file on the server with a new name.
func (c *Client) CopyFile(objectId, newName string) error {
	_, err := c.Call("copy_file", objectId, newName, c.Token)
	return err
}

// GetFilesSince returns all files which were created after the specified date.
func (c *Client) GetFilesSince(mtime time.Time) ([]*File, error) {
	result, err := c.Call("get_files_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawFilesList(result)
}

// GetFileAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetFileAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_file_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// SaveFile saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveFile(objectId, editmode string) error {
	_, err := c.Call("save_file", objectId, c.Token, editmode)
	return err
}

// RenameFile renames a file with a given object id.
func (c *Client) RenameFile(objectId, newName string) error {
	_, err := c.Call("rename_file", objectId, newName, c.Token)
	return err
}
