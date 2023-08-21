package item

import (
	"fmt"
	"github.com/cobbler/cobblerclient/client"
	clientInternals "github.com/cobbler/cobblerclient/internal/client"
	"time"
)

func convertRawProfile(name string, xmlrpcResult interface{}) (*Profile, error) {
	var profile Profile

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	decodeResult, err := clientInternals.DecodeCobblerItem(xmlrpcResult, &profile)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Profile), nil
}

func convertRawProfilesList(xmlrpcResult interface{}) ([]*Profile, error) {
	var profiles []*Profile

	for _, p := range xmlrpcResult.([]interface{}) {
		profile, err := convertRawProfile("unknown", p)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetProfiles returns all systems in Cobbler.
func GetProfiles(c client.Client) ([]*Profile, error) {
	result, err := c.Call("get_profiles", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawProfilesList(result)
}

func ListProfileNames(c client.Client) ([]string, error) {
	return c.GetItemNames("profile")
}

// GetProfile returns a single profile obtained by its name.
func GetProfile(c client.Client, name string) (*Profile, error) {
	var profile = BuildProfile(c)

	result, err := c.Call("get_profile", name, c.Token)
	if err != nil {
		return &profile, err
	}

	return convertRawProfile(name, result)
}

// FindProfile is ...
func FindProfile(c client.Client) error {
	_, err := c.Call("find_profile")
	return err
}

// GetProfilesSince is ...
func GetProfilesSince(c client.Client, time time.Time) ([]*Profile, error) {
	result, err := c.Call("get_profiles_since", float64(time.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawProfilesList(result)
}
