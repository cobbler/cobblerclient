package cobblerclient

import (
	"fmt"
	"reflect"
	"time"
)

// File is ...
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

func (c *Client) GetFiles() ([]*File, error) {
	result, err := c.Call("get_files", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawFilesList(result)
}

// GetFile returns a single distro obtained by its name.
func (c *Client) GetFile(name string) (*File, error) {
	result, err := c.Call("get_file", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawFile(name, result)
}

// CreateFile creates a distro.
func (c *Client) CreateFile(file File) (*File, error) {
	// Make sure a distro with the same name does not already exist
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

// DeleteDistro deletes a single distro by its name.
func (c *Client) DeleteFile(name string) error {
	_, err := c.Call("remove_file", name, c.Token)
	return err
}

// ListDistroNames is returning a list of all distro names currently available in Cobbler.
func (c *Client) ListFileNames() ([]string, error) {
	return c.GetItemNames("file")
}

// FindFile is ...
func (c *Client) FindFile(criteria map[string]interface{}) ([]*File, error) {
	result, err := c.Call("find_file", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawFilesList(result)
}

// FindFileNames is searching for one or more files by any of its attributes.
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

// CopyFile is ...
func (c *Client) CopyFile(objectId, newName string) error {
	_, err := c.Call("copy_file", objectId, newName, c.Token)
	return err
}

// GetFilesSince is ...
func (c *Client) GetFilesSince(mtime time.Time) ([]*File, error) {
	result, err := c.Call("get_files_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawFilesList(result)
}

// GetFileAsRendered is ...
func (c *Client) GetFileAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_file_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// SaveFile is ...
func (c *Client) SaveFile(objectId, editmode string) error {
	_, err := c.Call("save_file", objectId, c.Token, editmode)
	return err
}

// RenameFile is ...
func (c *Client) RenameFile(objectId, newName string) error {
	_, err := c.Call("rename_file", objectId, newName, c.Token)
	return err
}
