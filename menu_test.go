package cobblerclient

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"
)

func TestDeleteMenu(t *testing.T) {}

func TestFindMenu(t *testing.T) {
	c := createStubHTTPClient(t, "find-menu-req.xml", "find-menu-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	menus, err := c.FindMenu(criteria)
	utils.FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestGetMenuHandle(t *testing.T) {}

func TestCopyMenu(t *testing.T) {}

func TestGetMenusSince(t *testing.T) {}

func TestGetMenuAsRendered(t *testing.T) {}

func TestSaveMenu(t *testing.T) {}

func TestRenameMenu(t *testing.T) {}
