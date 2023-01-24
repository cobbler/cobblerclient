package item

import (
	"fmt"
	"reflect"
)

// Create creates a Cobbler package.
func (linuxpackage *Package) Create() (*Package, error) {
	// Make sure a distro with the same name does not already exist
	if _, err := GetPackage(linuxpackage.Client, linuxpackage.Name.Get()); err == nil {
		return nil, fmt.Errorf("a Package with the name %s already exists", linuxpackage.Name.Get())
	}

	result, err := linuxpackage.Client.Call("new_package", linuxpackage.Client.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&linuxpackage).Elem()
	if err := linuxpackage.UpdateCobblerFields("package", item, newID); err != nil {
		return nil, err
	}

	if _, err := linuxpackage.Client.Call("save_package", newID, linuxpackage.Client.Token); err != nil {
		return nil, err
	}

	return GetPackage(linuxpackage.Client, linuxpackage.Name.Get())
}
