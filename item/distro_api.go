package item

import (
	"fmt"
	"reflect"
)

// Create creates a distro.
func (d *Distro) Create() (*Distro, error) {
	// Make sure a d with the same name does not already exist
	if _, err := GetDistro(d.Client, d.Name.Get()); err == nil {
		return nil, fmt.Errorf("a Distro with the name %s already exists", d.Name)
	}

	result, err := d.Client.Call("new_distro", d.Client.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	item := reflect.ValueOf(&d).Elem()
	if err := d.UpdateCobblerFields("d", item, newID); err != nil {
		return nil, err
	}

	if _, err := d.Client.Call("save_distro", newID, d.Client.Token); err != nil {
		return nil, err
	}

	return GetDistro(d.Client, d.Name.Get())
}

// Update updates a single distro.
func (d *Distro) Update() error {
	item := reflect.ValueOf(d).Elem()
	id, err := d.Client.GetItemHandle("d", d.Name.Get())
	if err != nil {
		return err
	}

	if err := d.UpdateCobblerFields("d", item, id); err != nil {
		return err
	}

	if _, err := d.Client.Call("save_distro", id, d.Client.Token); err != nil {
		return err
	}

	return nil
}

// Delete deletes a single distro by its name.
func (d *Distro) Delete() error {
	_, err := d.Client.Call("remove_distro", d.Name, d.Client.Token)
	return err
}

// GetHandle gets the internal ID of a Cobbler item.
func (d *Distro) GetHandle() (string, error) {
	result, err := d.Client.Call("get_distro_handle", d.Name, d.Client.Token)
	if err != nil {
		return "", err
	} else {
		d.Handle = result.(string)
		return result.(string), err
	}
}

// Copy is copying a distro.
// Note: The new object does not possess object equality with the old one because the names and other internal
// attributes are different.
func (d *Distro) Copy(newName string) (Distro, error) {
	_, err := d.Client.Call("copy_distro", d.Handle, newName, d.Client.Token)
	result, err := GetDistro(d.Client, newName)
	return *result, err
}

// Rename is renaming a distro.
func (d *Distro) Rename(newName string) error {
	_, err := d.Client.Call("rename_distro", d.Handle, newName, d.Client.Token)
	return err
}

// SaveDistro is saving a Distro to disk on the Cobbler server.
func (d *Distro) SaveDistro(editmode string) error {
	_, err := d.Client.Call("save_distro", d.Handle, d.Client.Token, editmode)
	return err
}

// GetAsRendered is resolving all values from the object inheritance and returns the final values the object has.
func (d *Distro) GetAsRendered() error {
	_, err := d.Client.Call("get_distro_as_rendered")
	return err
}

// GetValidBootLoaders is returning a list (which may be empty) of bootloaders which are valid for the distro in
// question.
func (d *Distro) GetValidBootLoaders() error {
	_, err := d.Client.Call("get_valid_distro_boot_loaders", d.Name, d.Client.Token)
	return err
}
