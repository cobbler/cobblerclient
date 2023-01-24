package item

import (
	"fmt"
	"github.com/cobbler/cobblerclient/client"
	client_internals "github.com/cobbler/cobblerclient/internal/client"
)

// GetSystems returns all systems in Cobbler.
func GetSystems(c client.Client) ([]*System, error) {
	var systems []*System

	result, err := c.Call("get_systems", "", c.Token)
	if err != nil {
		return nil, err
	}

	for _, s := range result.([]interface{}) {
		var system System
		decodedResult, err := client_internals.DecodeCobblerItem(s, &system)
		if err != nil {
			return nil, err
		}
		decodedSystem := decodedResult.(*System)
		systems = append(systems, decodedSystem)
	}

	return systems, nil
}

func ListSystemNames(c client.Client) ([]string, error) {
	return c.GetItemNames("system")
}

// GetSystem returns a single system obtained by its name.
func GetSystem(c client.Client, name string) (*System, error) {
	var system System

	result, err := c.Call("get_system", name, c.Token)
	if err != nil {
		return &system, err
	}

	if result == "~" {
		return nil, fmt.Errorf("system %s not found", name)
	}

	decodeResult, err := client_internals.DecodeCobblerItem(result, &system)
	if err != nil {
		return nil, err
	}

	s := decodeResult.(*System)

	return s, nil
}

// FindSystem is ...
func FindSystem(c client.Client) error {
	_, err := c.Call("find_system")
	return err
}

// GetSystemsSince is ...
func GetSystemsSince(c client.Client) error {
	_, err := c.Call("get_systems_since")
	return err
}
