package item_test

import (
	"github.com/cobbler/cobblerclient/item"
	"testing"

	"github.com/ContainerSolutions/go-utils"

	cobblerTesting "github.com/cobbler/cobblerclient/internal/testing"
)

func TestGetDistros(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-distros-req.xml", "get-distros-res.xml")
	distros, err := item.GetDistros(c)
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestListDistroNames(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-item-names-distro.xml", "get-item-names-distro.xml")
	distros, err := item.ListDistroNames(&c)
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistro(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-distro-req.xml", "get-distro-res.xml")
	resDistro, err := item.GetDistro(c, "Ubuntu-20.04-x86_64")
	utils.FailOnError(t, err)

	if resDistro.Name.Get() != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong distro returned.")
	}
}

func TestGetDistrosSince(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "get-distros-since-req.xml", "get-distros-since-res.xml")
	err := item.GetDistrosSince(c)
	utils.FailOnError(t, err)
}

func TestFindDistro(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "find-distro-req.xml", "find-distro-res.xml")
	err := item.FindDistro(c)
	utils.FailOnError(t, err)
}
