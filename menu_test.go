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
)

func TestNewMenu(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewMenu()
}

func TestGetMenus(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-menus")
	menus, err := c.GetMenus()
	FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestGetMenu(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-menu")
	c.CachedVersion = CobblerVersion{3, 3, 2}

	// Act
	menu, err := c.GetMenu("testmenu", false, false)

	// Assert
	FailOnError(t, err)
	if menu.Name != "testmenu" {
		t.Errorf("Wrong menu returned.")
	}
}

func TestDeleteMenu(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-menu")
	err := c.DeleteMenu("test")
	FailOnError(t, err)
}

func TestDeleteMenuRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-menu")
	err := c.DeleteMenuRecursive("test", false)
	FailOnError(t, err)
}

func TestListMenuNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-menu")
	menus, err := c.ListMenuNames()
	FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestGetMenusSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-menus-since")
	menus, err := c.GetMenusSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestFindMenu(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-menu")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmenu"
	menus, err := c.FindMenu(criteria)
	FailOnError(t, err)

	if len(menus) != 1 {
		t.Errorf("Wrong number of menus returned.")
	}
}

func TestFindMenuNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-menu-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmenu"
	menus, err := c.FindMenuNames(criteria)
	FailOnError(t, err)

	if len(menus) != 1 {
		t.Error("Wrong number of menus returned.")
	}
}

func TestSaveMenu(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-menu")
	err := c.SaveMenu("menu::testmenu", "bypass")
	FailOnError(t, err)
}

func TestCopyMenu(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-menu")
	err := c.CopyMenu("menu::testmenu", "testmenu2")
	FailOnError(t, err)
}

func TestRenameMenu(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-menu")
	err := c.RenameMenu("menu::testmenu2", "testmenu1")
	FailOnError(t, err)
}

func TestGetMenuHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-menu-handle")
	res, err := c.GetMenuHandle("testmenu")
	FailOnError(t, err)

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
