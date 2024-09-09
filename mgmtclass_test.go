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

func TestNewMgmtClass(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewMgmtClass()
}

func TestGetMgmtclasses(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-mgmtclasses")
	mgmtclasses, err := c.GetMgmtClasses()
	FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclass returned.")
	}
}

func TestGetMgmtclass(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-mgmtclass")
	c.CachedVersion = CobblerVersion{3, 3, 2}

	// Act
	mgmtclass, err := c.GetMgmtClass("testmgmtclass", false, false)

	// Assert
	FailOnError(t, err)
	if mgmtclass.Name != "testmgmtclass" {
		t.Errorf("Wrong mgmtclass returned.")
	}
}

func TestDeleteMgmtClass(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-mgmtclass")
	err := c.DeleteMgmtClass("test")
	FailOnError(t, err)
}

func TestDeleteMgmtClassRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-mgmtclass")
	err := c.DeleteMgmtClassRecursive("test", false)
	FailOnError(t, err)
}

func TestListMgmtClassNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-mgmtclass")
	mgmtclasses, err := c.ListMgmtClassNames()
	FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclasses returned.")
	}
}

func TestGetMgmtclassSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-mgmtclasses-since")
	mgmtclasses, err := c.GetMgmtClassesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclasses returned.")
	}
}

func TestFindMgmtclass(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-mgmtclass")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmgmtclass"
	mgmtclasses, err := c.FindMgmtClass(criteria)
	FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclasses returned.")
	}
}

func TestFindMgmtClassNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-mgmtclass-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmgmtclass"
	mgmtclasses, err := c.FindMgmtClassNames(criteria)
	FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Error("Wrong number of mgmtclasses returned.")
	}
}

func TestSaveMgmtClass(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-mgmtclass")
	err := c.SaveMgmtClass("mgmtclass::testmgmtclass", "bypass")
	FailOnError(t, err)
}

func TestCopyMgmtClas(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-mgmtclass")
	err := c.CopyMgmtClass("mgmtclass::testmgmtclass", "testmgmtclass2")
	FailOnError(t, err)
}

func TestRenameMgmtClass(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-mgmtclass")
	err := c.RenameMgmtClass("mgmtclass::testmgmtclass2", "mgmtclass1")
	FailOnError(t, err)
}

func TestGetMgmtClassHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-mgmtclass-handle")
	res, err := c.GetMgmtClassHandle("testmgmtclass")
	FailOnError(t, err)

	if res != "mgmtclass::testmgmtclass" {
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
