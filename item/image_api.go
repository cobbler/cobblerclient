package item

import (
	"reflect"
)

// CreateImage creates a profile.
func (i *Image) Create() (*Image, error) {
	// To create an image via the Cobbler API, first call new_image to obtain an ID
	result, err := i.Client.Call("new_image", i.Client.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)
	// Set the value of all fields
	item := reflect.ValueOf(i).Elem()
	if err := i.UpdateCobblerFields("profile", item, newID); err != nil {
		return nil, err
	}

	// Save the final profile
	_, err = i.Client.Call("save_profile", newID, i.Client.Token)
	if err != nil {
		return nil, err
	}

	// Return a clean copy of the profile
	return GetImage(i.Client, i.Name.Get())
}

// UpdateImage updates a single image.
func (i *Image) Update(image *Image) error {
	item := reflect.ValueOf(image).Elem()
	id, err := i.Client.GetItemHandle("image", image.Name.Get())
	if err != nil {
		return err
	}

	if err := i.UpdateCobblerFields("image", item, id); err != nil {
		return err
	}

	// Save the final image
	if _, err := i.Client.Call("save_image", id, i.Client.Token); err != nil {
		return err
	}

	return nil
}

// Delete deletes a single image by its name.
func (i *Image) Delete() error {
	_, err := i.Client.Call("remove_image", i.Name, i.Client.Token)
	return err
}

// GetImageHandle gets the internal ID of a Cobbler item.
func (i *Image) GetImageHandle(name string) (string, error) {
	result, err := i.Client.Call("get_image_handle", name, i.Client.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// CopyImage is ...
func (i *Image) CopyImage(objectId, newName string) error {
	_, err := i.Client.Call("copy_image", objectId, newName, i.Client.Token)
	return err
}

// Rename is ...
func (i *Image) Rename(newName string) error {
	_, err := i.Client.Call("rename_image", i.Handle, newName, i.Client.Token)
	return err
}

// GetAsRendered is ...
func (i *Image) GetAsRendered() error {
	_, err := i.Client.Call("get_image_as_rendered", i.Name)
	return err
}

// Save is ...
func (i *Image) Save(editmode string) error {
	_, err := i.Client.Call("save_image", i.Handle, i.Client.Token, editmode)
	return err
}

// GetValidBootLoaders is ...
func (i *Image) GetValidBootLoaders() error {
	_, err := i.Client.Call("get_valid_image_boot_loaders", i.Name, i.Client.Token)
	return err
}
