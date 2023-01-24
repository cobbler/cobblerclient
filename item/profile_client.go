package item

import (
	"fmt"
	"github.com/cobbler/cobblerclient/client"
	client_internals "github.com/cobbler/cobblerclient/internal/client"
)

// GetProfiles returns all systems in Cobbler.
func GetProfiles(c client.Client) ([]*Profile, error) {
	var profiles []*Profile

	result, err := c.Call("get_profiles", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, p := range result.([]interface{}) {
		var profile Profile
		decodedResult, err := client_internals.DecodeCobblerItem(p, &profile)
		if err != nil {
			return nil, err
		}
		decodedProfile := decodedResult.(*Profile)
		profiles = append(profiles, decodedProfile)
	}

	return profiles, nil
}

func ListProfileNames(c client.Client) ([]string, error) {
	return c.GetItemNames("profile")
}

// GetProfile returns a single profile obtained by its name.
func GetProfile(c client.Client, name string) (*Profile, error) {
	var profile Profile

	result, err := c.Call("get_profile", name, c.Token)
	if err != nil {
		return &profile, err
	}

	if result == "~" {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	decodeResult, err := client_internals.DecodeCobblerItem(result, &profile)
	if err != nil {
		return nil, err
	}

	s := decodeResult.(*Profile)

	return s, nil
}

// FindProfile is ...
func FindProfile(c client.Client) error {
	_, err := c.Call("find_profile")
	return err
}

// GetProfilesSince is ...
func GetProfilesSince(c client.Client) error {
	_, err := c.Call("get_profiles_since")
	return err
}
