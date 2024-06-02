package cobblerclient

import (
	"fmt"
	"reflect"
	"time"
)

// Package is a created package.
// Get the fields from cobbler/items/package.py
type Package struct {
	Resource `mapstructure:",squash"`

	// Package specific attributes
	Installer string `mapstructure:"installer"`
	Version   string `mapstructure:"version"`
}

func convertRawLinuxPackage(name string, xmlrpcResult interface{}) (*Package, error) {
	var linuxpackage Package

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("package %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &linuxpackage)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Package), nil
}

func convertRawLinuxPackageList(xmlrpcResult interface{}) ([]*Package, error) {
	var linuxpackages []*Package

	for _, d := range xmlrpcResult.([]interface{}) {
		linuxpackage, err := convertRawLinuxPackage("unknown", d)
		if err != nil {
			return nil, err
		}
		linuxpackages = append(linuxpackages, linuxpackage)
	}

	return linuxpackages, nil
}

// GetPackages returns all packages in Cobbler.
func (c *Client) GetPackages() ([]*Package, error) {
	result, err := c.Call("get_packages", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawLinuxPackageList(result)
}

// GetPackage returns a single package obtained by its name.
func (c *Client) GetPackage(name string) (*Package, error) {
	result, err := c.Call("get_package", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawLinuxPackage(name, result)
}

// CreatePackage creates a package.
func (c *Client) CreatePackage(linuxpackage Package) (*Package, error) {
	// Make sure a package with the same name does not already exist
	if _, err := c.GetPackage(linuxpackage.Name); err == nil {
		return nil, fmt.Errorf("a Package with the name %s already exists", linuxpackage.Name)
	}

	result, err := c.Call("new_package", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&linuxpackage).Elem()
	if err := c.updateCobblerFields("package", item, newID); err != nil {
		return nil, err
	}

	if err = c.SavePackage(newID, "new"); err != nil {
		return nil, err
	}

	return c.GetPackage(linuxpackage.Name)
}

// UpdatePackage updates a single package.
func (c *Client) UpdatePackage(linuxpackage *Package) error {
	item := reflect.ValueOf(linuxpackage).Elem()
	id, err := c.GetItemHandle("package", linuxpackage.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("package", item, id); err != nil {
		return err
	}

	if err = c.SavePackage(id, "bypass"); err != nil {
		return err
	}

	return nil
}

// ListPackageNames is returning a list of all packages names currently available in Cobbler.
func (c *Client) ListPackageNames() ([]string, error) {
	return c.GetItemNames("package")
}

// FindPackage is the search method that allows looking for an object by any of its attributes.
func (c *Client) FindPackage(criteria map[string]interface{}) ([]*Package, error) {
	result, err := c.Call("find_package", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawLinuxPackageList(result)
}

// FindPackageNames is searching for one or more distros by any of its attributes.
func (c *Client) FindPackageNames(criteria map[string]interface{}) ([]string, error) {
	var result []string

	resultUnmarshalled, err := c.Call("find_package", criteria, false, c.Token)

	if err != nil {
		return nil, err
	}

	for _, name := range resultUnmarshalled.([]interface{}) {
		result = append(result, name.(string))
	}

	return result, nil
}

// GetPackageHandle gets the internal ID of a Cobbler item.
func (c *Client) GetPackageHandle(name string) (string, error) {
	result, err := c.Call("get_package_handle", name, c.Token)
	if err != nil {
		return "", err
	}
	return result.(string), err
}

// CopyPackage is copying a given package server side with a new name.
func (c *Client) CopyPackage(objectId, newName string) error {
	_, err := c.Call("copy_package", objectId, newName, c.Token)
	return err
}

// DeletePackage deletes a single distro by its name.
func (c *Client) DeletePackage(name string) error {
	_, err := c.Call("remove_package", name, c.Token)
	return err
}

// RenamePackage is renaming a package with a given object id.
func (c *Client) RenamePackage(objectId, newName string) error {
	_, err := c.Call("rename_package", objectId, newName, c.Token)
	return err
}

// GetPackagesSince is returning all packages that have been edited since a given timestamp.
func (c *Client) GetPackagesSince(mtime time.Time) ([]*Package, error) {
	result, err := c.Call("get_packages_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawLinuxPackageList(result)
}

// GetPackageAsRendered is returning the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetPackageAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_package_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// SavePackage is persisting all changes performed via XML-RPC to disk on the server side.
func (c *Client) SavePackage(objectId, editmode string) error {
	_, err := c.Call("save_package", objectId, c.Token, editmode)
	return err
}
