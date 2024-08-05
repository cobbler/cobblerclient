/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

import (
	"testing"
	"time"

	"github.com/ContainerSolutions/go-utils"
)

func TestGetMenus(t *testing.T) {
	c := createStubHTTPClient(t, "get-menus-req.xml", "get-menus-res.xml")
	menus, err := c.GetMenus()
	utils.FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestGetMenu(t *testing.T) {
	c := createStubHTTPClient(t, "get-menu-req.xml", "get-menu-res.xml")
	menu, err := c.GetMenu("testmenu")
	utils.FailOnError(t, err)

	if menu.Name != "testmenu" {
		t.Errorf("Wrong menu returned.")
	}
}

func TestDeleteMenu(t *testing.T) {
	c := createStubHTTPClient(t, "delete-menu-req.xml", "delete-menu-res.xml")
	err := c.DeleteMenu("test")
	utils.FailOnError(t, err)
}

func TestDeleteMenuRecursive(t *testing.T) {
	c := createStubHTTPClient(t, "delete-menu-req.xml", "delete-menu-res.xml")
	err := c.DeleteMenuRecursive("test", false)
	utils.FailOnError(t, err)
}

func TestListMenuNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-menu-req.xml", "get-item-names-menu-res.xml")
	menus, err := c.ListMenuNames()
	utils.FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestGetMenusSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-menus-since-req.xml", "get-menus-since-res.xml")
	menus, err := c.GetMenusSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestFindMenu(t *testing.T) {
	c := createStubHTTPClient(t, "find-menu-req.xml", "find-menu-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmenu"
	menus, err := c.FindMenu(criteria)
	utils.FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestFindMenuNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-menu-names-req.xml", "find-menu-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmenu"
	menus, err := c.FindMenuNames(criteria)
	utils.FailOnError(t, err)

	if len(menus) != 1 {
		t.Error("Wrong number of menus returned.")
	}
}

func TestSaveMenu(t *testing.T) {
	c := createStubHTTPClient(t, "save-menu-req.xml", "save-menu-res.xml")
	err := c.SaveMenu("menu::testmenu", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyMenu(t *testing.T) {
	c := createStubHTTPClient(t, "copy-menu-req.xml", "copy-menu-res.xml")
	err := c.CopyMenu("menu::testmenu", "testmenu2")
	utils.FailOnError(t, err)
}

func TestRenameMenu(t *testing.T) {
	c := createStubHTTPClient(t, "rename-menu-req.xml", "rename-menu-res.xml")
	err := c.RenameMenu("menu::testmenu2", "testmenu1")
	utils.FailOnError(t, err)
}

func TestGetMenuHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-menu-handle-req.xml", "get-menu-handle-res.xml")
	res, err := c.GetMenuHandle("testmenu")
	utils.FailOnError(t, err)

	if res != "menu::testmenu" {
		t.Error("Wrong object id returned.")
	}
}

/*
 * NOTE: We're skipping the testing of CREATE, UPDATE, DELETE methods for now because
 *       the current implementation of the StubHTTPClient does not allow
 *       buffered mock responses so as soon as the method makes the second
 *       call to Cobbler it'll fail.
 *       This is a system test, so perhaps we can run Cobbler in a Docker container
 *       and take it from there.
 */
