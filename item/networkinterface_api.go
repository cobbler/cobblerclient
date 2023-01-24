package item

import (
	"fmt"
	"log"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// CreateInterface controls network interfaces in Cobbler
func (iface *Interface) CreateInterface(system System, name string) error {
	i := structs.Map(iface)
	nic := make(map[string]interface{})
	for key, value := range i {
		attrName := fmt.Sprintf("%s-%s", key, name)
		log.Printf("[DEBUG] Cobblerclient: setting interface attr %s to %s", attrName, value)
		nic[attrName] = value
	}

	systemID, err := system.Client.GetItemHandle("system", system.Name.Get())
	if err != nil {
		return err
	}

	_, err = system.Client.Call("modify_system", systemID, "modify_interface", nic, system.Client.Token)
	if err != nil {
		return err
	}

	// Save the final system
	_, err = system.Client.Call("save_system", systemID, system.Client.Token)
	if err != nil {
		return err
	}

	return nil
}

// GetInterfaces returns all interfaces in a System.
func GetInterfaces(system System) (Interfaces, error) {
	nics := make(Interfaces)
	for nicName, nicData := range system.Interfaces {
		var nic Interface
		if err := mapstructure.Decode(nicData, &nic); err != nil {
			return nil, err
		}
		nics[nicName] = nic
	}

	return nics, nil
}

// GetInterface returns a single interface in a System.
func GetInterface(system System, name string) (Interface, error) {
	nics := make(Interfaces)
	var iface Interface
	for nicName, nicData := range system.Interfaces {
		var nic Interface
		if err := mapstructure.Decode(nicData, &nic); err != nil {
			return iface, err
		}
		nics[nicName] = nic
	}

	if iface, ok := nics[name]; ok {
		return iface, nil
	} else {
		return iface, fmt.Errorf("interface %s not found", name)
	}
}

// DeleteInterface deletes a single interface in a System.
func DeleteInterface(system System, name string) error {
	if _, err := GetInterface(system, name); err != nil {
		return err
	}

	systemID, err := system.Client.GetItemHandle("system", system.Name.Get())
	if err != nil {
		return err
	}

	if _, err := system.Client.Call("modify_system", systemID, "delete_interface", name, system.Client.Token); err != nil {
		return err
	}

	// Save the final system
	if _, err := system.Client.Call("save_system", systemID, system.Client.Token); err != nil {
		return err
	}

	return nil
}
