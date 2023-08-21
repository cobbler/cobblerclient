package item

import (
	"fmt"
	"reflect"
)

// Create creates a profile.
// It ensures that a Distro is set and then sets other default values.
func (p *Profile) Create() (*Profile, error) {
	// Check if a profile with the same name already exists
	if _, err := GetProfile(p.Client, p.Name.Get()); err == nil {
		return nil, fmt.Errorf("a profile with the name %s already exists", p.Name)
	}

	if p.Distro.Get() == "" {
		return nil, fmt.Errorf("a profile must have a distro set")
	}

	if len(p.MgmtParameters.Get()) == 0 {
		p.MgmtParameters.SetRaw("<<inherit>>")
	}
	if p.VirtType.Get() == "" {
		p.VirtType.Set("<<inherit>>")
	}
	if p.VirtDiskDriver.Get() == "" {
		p.VirtDiskDriver.Set("<<inherit>>")
	}

	// To create a profile via the Cobbler API, first call new_profile to obtain an ID
	result, err := p.Client.Call("new_profile", p.Client.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)
	// Set the value of all fields
	item := reflect.ValueOf(&p).Elem()
	if err := p.UpdateCobblerFields("profile", item, newID); err != nil {
		return nil, err
	}

	// Save the final profile
	_, err = p.Client.Call("save_profile", newID, p.Client.Token)
	if err != nil {
		return nil, err
	}

	// Return a clean copy of the profile
	return GetProfile(p.Client, p.Name.Get())
}

// Update updates a single profile.
func (p *Profile) Update(profile *Profile) error {
	item := reflect.ValueOf(profile).Elem()
	id, err := p.Client.GetItemHandle("profile", profile.Name.Get())
	if err != nil {
		return err
	}

	if err := p.UpdateCobblerFields("profile", item, id); err != nil {
		return err
	}

	// Save the final profile
	if _, err := p.Client.Call("save_profile", id, p.Client.Token); err != nil {
		return err
	}

	return nil
}

// Delete deletes a single profile by its name.
func (p *Profile) Delete() error {
	_, err := p.Client.Call("remove_profile", p.Name, p.Client.Token)
	return err
}

// GenerateAutoinstall is ...
func (p *Profile) GenerateAutoinstall() (string, error) {
	result, err := p.Client.Call("generate_profile_autoinstall", p.Name)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// GetHandle gets the internal ID of a Cobbler item.
func (p *Profile) GetHandle() (string, error) {
	result, err := p.Client.Call("get_profile_handle", p.Name, p.Client.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// Copy is ...
func (p *Profile) Copy(newName string) error {
	_, err := p.Client.Call("copy_profile", p.Handle, newName, p.Client.Token)
	return err
}

// Rename is ...
func (p *Profile) Rename(newName string) error {
	_, err := p.Client.Call("rename_profile", p.Handle, newName, p.Client.Token)
	return err
}

// GetAsRendered is ...
func (p *Profile) GetAsRendered() error {
	_, err := p.Client.Call("get_profile_as_rendered")
	return err
}

// Save is ...
func (p *Profile) Save(editmode string) error {
	_, err := p.Client.Call("save_profile", p.Handle, p.Client.Token, editmode)
	return err
}

// GetRepoConfig is ...
func (p *Profile) GetRepoConfig(profileName string) error {
	_, err := p.Client.Call("get_repo_config_for_profile", profileName)
	return err
}

// GetValidBootLoaders is ...
func (p *Profile) GetValidBootLoaders(profileName string) error {
	_, err := p.Client.Call("get_valid_profile_boot_loaders", profileName, p.Client.Token)
	return err
}

// GetTemplateFile is ...
func (p *Profile) GetTemplateFile(path string) error {
	_, err := p.Client.Call("get_template_file_for_profile", p.Name, path)
	return err
}
