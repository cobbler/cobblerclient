package item

import (
	"fmt"
	"github.com/cobbler/cobblerclient/client"
	clientInternals "github.com/cobbler/cobblerclient/internal/client"
	"github.com/cobbler/cobblerclient/internal/item/raw"
)

// GetDistros returns all systems in Cobbler.
func GetDistros(c client.Client) ([]*Distro, error) {
	var distros []*Distro

	result, err := c.Call("get_distros", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, d := range result.([]interface{}) {
		var distro = BuildDistro(c)
		var rawDistro raw.Distro
		decodedResult, err := clientInternals.DecodeCobblerItem(d, &rawDistro)
		if err != nil {
			return nil, err
		}
		distro.raw = decodedResult.(*raw.Distro)

		distros = append(distros, &distro)
	}

	return distros, nil
}

// ListDistroNames is returning a list of all distro names currently available in Cobbler.
func ListDistroNames(c *client.Client) ([]string, error) {
	return c.GetItemNames("distro")
}

// GetDistro returns a single distro obtained by its name.
func GetDistro(c client.Client, name string) (*Distro, error) {
	var distro = BuildDistro(c)
	var rawDistro raw.Distro

	result, err := c.Call("get_distro", name, c.Token)
	if result == "~" {
		return nil, fmt.Errorf("distro %s not found", name)
	}

	if err != nil {
		return nil, err
	}

	decodeResult, err := clientInternals.DecodeCobblerItem(result, &rawDistro)
	if err != nil {
		return nil, err
	}

	distro.raw = decodeResult.(*raw.Distro)
	refreshDistroPointers(&distro)
	return &distro, nil
}

// GetDistrosSince is returning all distros which were created after the specified date.
func GetDistrosSince(c client.Client) error {
	_, err := c.Call("get_distros_since")
	return err
}

// FindDistro is searching for one or more distros by any of its attributes.
func FindDistro(c client.Client) error {
	_, err := c.Call("find_distro")
	return err
}
