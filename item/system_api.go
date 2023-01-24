package item

import (
	"fmt"
	"reflect"
)

// CreateSystem creates a system.
// It ensures that either a Profile or Image are set and then sets other default values.
func (s *System) Create() (*System, error) {
	// Check if a system with the same name already exists
	if _, err := GetSystem(s.Client, s.Name.Get()); err == nil {
		return nil, fmt.Errorf("a system with the name %s already exists", s.Name)
	}

	if s.Profile == "" && s.Image == "" {
		return nil, fmt.Errorf("a system must have a profile or image set")
	}

	// Set default values. I guess these aren't taken care of by Cobbler?
	if s.BootFiles == "" {
		s.BootFiles = "<<inherit>>"
	}

	if len(s.BootLoaders) == 0 {
		s.BootLoaders = []string{"<<inherit>>"}
	}

	if len(s.BootLoaders) == 0 {
		s.FetchableFiles = []string{"<<inherit>>"}
	}

	if s.MGMTParameters == "" {
		s.MGMTParameters = "<<inherit>>"
	}

	if s.PowerType == "" {
		s.PowerType = "ipmilan"
	}

	if s.Status == "" {
		s.Status = "production"
	}

	if s.VirtAutoBoot == "" {
		s.VirtAutoBoot = "0"
	}

	if s.VirtCPUs == "" {
		s.VirtCPUs = "<<inherit>>"
	}

	if s.VirtDiskDriver == "" {
		s.VirtDiskDriver = "<<inherit>>"
	}

	if s.VirtFileSize == "" {
		s.VirtFileSize = "<<inherit>>"
	}

	if s.VirtPath == "" {
		s.VirtPath = "<<inherit>>"
	}

	if s.VirtRAM == "" {
		s.VirtRAM = "<<inherit>>"
	}

	if s.VirtType == "" {
		s.VirtType = "<<inherit>>"
	}

	// To create a system via the Cobbler API, first call new_system to obtain an ID
	result, err := s.Client.Call("new_system", s.Client.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)

	// Set the value of all fields
	item := reflect.ValueOf(s).Elem()
	if err := s.UpdateCobblerFields("system", item, newID); err != nil {
		return nil, err
	}

	// Save the final system
	if _, err := s.Client.Call("save_system", newID, s.Client.Token); err != nil {
		return nil, err
	}

	// Return a clean copy of the system
	return GetSystem(s.Client, s.Name.Get())
}

// UpdateSystem updates a single system.
func (s *System) Update() error {
	item := reflect.ValueOf(s).Elem()
	id, err := s.Client.GetItemHandle("system", s.Name.Get())
	if err != nil {
		return err
	}
	return s.UpdateCobblerFields("system", item, id)
}

// DeleteSystem deletes a single system by its name.
func (s *System) Delete() error {
	_, err := s.Client.Call("remove_system", s.Name, s.Client.Token)
	return err
}

// ClearSystemLogs is ..
func (s *System) ClearLogs() error {
	_, err := s.Client.Call("clear_system_logs", s.Handle)
	return err
}

// GenerateAutoinstall is ..
func (s *System) GenerateAutoinstall() (string, error) {
	result, err := s.Client.Call("generate_system_autoinstall", s.Name)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// GetSystemHandle gets the internal ID of a Cobbler item.
func (s *System) GetHandle() (string, error) {
	result, err := s.Client.Call("get_system_handle", s.Name, s.Client.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// Copy is ...
func (s *System) Copy(newName string) error {
	_, err := s.Client.Call("copy_system", s.Handle, newName, s.Client.Token)
	return err
}

// Rename is ...
func (s *System) RenameSystem(newName string) error {
	_, err := s.Client.Call("rename_system", s.Handle, newName, s.Client.Token)
	return err
}

// GetAsRendered is ...
func (s *System) GetSystemAsRendered() error {
	_, err := s.Client.Call("get_system_as_rendered", s.Name)
	return err
}

// Save is ...
func (s *System) Save(editmode string) error {
	_, err := s.Client.Call("save_system", s.Handle, s.Client.Token, editmode)
	return err
}

// GetRepoConfigForSystem is ...
func (s *System) GetRepoConfigForSystem(systemName string) error {
	_, err := s.Client.Call("get_repo_config_for_system", systemName)
	return err
}

// GetValidSystemBootLoaders is ...
func (s *System) GetValidSystemBootLoaders(systemName string) error {
	_, err := s.Client.Call("get_valid_system_boot_loaders", systemName, s.Client.Token)
	return err
}

// GetTemplateFileForSystem is ...
func (s *System) GetTemplateFileForSystem() error {
	_, err := s.Client.Call("get_template_file_for_system")
	return err
}

func (s *System) BackgroundPowerSystem() (string, error) {
	result, err := s.Client.Call("background_power_system", s.Client.Token)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

func (s *System) PowerSystem(systemId string, power string) (bool, error) {
	result, err := s.Client.Call("power_system", systemId, power, s.Client.Token)
	if err != nil {
		return false, err
	} else {
		return result.(bool), err
	}
}

// DisableNetboot is ...
func (s *System) DisableNetboot() error {
	_, err := s.Client.Call("disable_netboot")
	return err
}

// UploadLogData is ...
func (s *System) UploadLogData() error {
	_, err := s.Client.Call("upload_log_data")
	return err
}
