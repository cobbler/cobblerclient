package item

import (
	"fmt"

	"github.com/cobbler/cobblerclient/client"
	client_internals "github.com/cobbler/cobblerclient/internal/client"
)

// GetPackages returns all systems in Cobbler.
func GetPackages(c client.Client) ([]*Package, error) {
	var packages []*Package

	result, err := c.Call("get_packages", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, d := range result.([]interface{}) {
		var linuxpackage Package
		decodedResult, err := client_internals.DecodeCobblerItem(d, &linuxpackage)
		if err != nil {
			return nil, err
		}

		packages = append(packages, decodedResult.(*Package))
	}

	return packages, nil
}

// ListPackageNames is returning a list of all package names currently available in Cobbler.
func ListPackageNames(c *client.Client) ([]string, error) {
	return c.GetItemNames("package")
}

// GetPackage returns a single package obtained by its name.
func GetPackage(c client.Client, name string) (*Package, error) {
	var linuxpackage Package

	result, err := c.Call("get_package", name, c.Token)
	if result == "~" {
		return nil, fmt.Errorf("package %s not found", name)
	}

	if err != nil {
		return nil, err
	}

	decodeResult, err := client_internals.DecodeCobblerItem(result, &linuxpackage)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Package), nil
}

// GetPackagesSince is returning all packages which were created after the specified date.
func GetPackagesSince(c client.Client) error {
	_, err := c.Call("get_packages_since")
	return err
}

// FindPackage is searching for one or more packages by any of its attributes.
func FindPackage(c client.Client) error {
	_, err := c.Call("find_package")
	return err
}
