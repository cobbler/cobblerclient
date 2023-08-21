package item_test

import (
	"github.com/cobbler/cobblerclient/item"
	"testing"
	"time"

	"github.com/ContainerSolutions/go-utils"

	cobblerTesting "github.com/cobbler/cobblerclient/internal/testing"
)

func TestGetProfiles(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-profiles-req.xml", "get-profiles-res.xml")
	profiles, err := item.GetProfiles(c)
	utils.FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestListProfileNames(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-item-names-profile-req.xml", "get-item-names-profile-res.xml")
	profiles, err := item.ListProfileNames(c)
	utils.FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfile(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-profile-req.xml", "get-profile-res.xml")
	profile, err := item.GetProfile(c, "Ubuntu-20.04-x86_64")
	utils.FailOnError(t, err)

	if profile.Name.Get() != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile returned.")
	}
}

func TestGetProfilesSince(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-profiles-since-req.xml", "get-profiles-since-res.xml")
	profiles, err := item.GetProfilesSince(c, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfile(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-distro-req.xml", "get-distro-res.xml")
	err := item.FindProfile(c)
	utils.FailOnError(t, err)
}
