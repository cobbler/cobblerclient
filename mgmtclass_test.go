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

func TestGetMgmtclasses(t *testing.T) {
	c := createStubHTTPClient(t, "get-mgmtclasses-req.xml", "get-mgmtclasses-res.xml")
	mgmtclasses, err := c.GetMgmtClasses()
	utils.FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclass returned.")
	}
}

func TestGetMgmtclass(t *testing.T) {
	c := createStubHTTPClient(t, "get-mgmtclass-req.xml", "get-mgmtclass-res.xml")
	mgmtclass, err := c.GetMgmtClass("testmgmtclass")
	utils.FailOnError(t, err)

	if mgmtclass.Name != "testmgmtclass" {
		t.Errorf("Wrong mgmtclass returned.")
	}
}

func TestDeleteMgmtClass(t *testing.T) {
	c := createStubHTTPClient(t, "delete-mgmtclass-req.xml", "delete-mgmtclass-res.xml")
	err := c.DeleteMgmtClass("test")
	utils.FailOnError(t, err)
}

func TestDeleteMgmtClassRecursive(t *testing.T) {
	c := createStubHTTPClient(t, "delete-mgmtclass-req.xml", "delete-mgmtclass-res.xml")
	err := c.DeleteMgmtClassRecursive("test", false)
	utils.FailOnError(t, err)
}

func TestListMgmtClassNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-mgmtclass-req.xml", "get-item-names-mgmtclass-res.xml")
	mgmtclasses, err := c.ListMgmtClassNames()
	utils.FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclasses returned.")
	}
}

func TestGetMgmtclassSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-mgmtclasses-since-req.xml", "get-mgmtclasses-since-res.xml")
	mgmtclasses, err := c.GetMgmtClassesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclasses returned.")
	}
}

func TestFindMgmtclass(t *testing.T) {
	c := createStubHTTPClient(t, "find-mgmtclass-req.xml", "find-mgmtclass-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmgmtclass"
	mgmtclasses, err := c.FindMgmtClass(criteria)
	utils.FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Errorf("Wrong number of mgmtclasses returned.")
	}
}

func TestFindMgmtClassNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-mgmtclass-names-req.xml", "find-mgmtclass-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testmgmtclass"
	mgmtclasses, err := c.FindMgmtClassNames(criteria)
	utils.FailOnError(t, err)

	if len(mgmtclasses) != 1 {
		t.Error("Wrong number of mgmtclasses returned.")
	}
}

func TestSaveMgmtClass(t *testing.T) {
	c := createStubHTTPClient(t, "save-mgmtclass-req.xml", "save-mgmtclass-res.xml")
	err := c.SaveMgmtClass("mgmtclass::testmgmtclass", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyMgmtClas(t *testing.T) {
	c := createStubHTTPClient(t, "copy-mgmtclass-req.xml", "copy-mgmtclass-res.xml")
	err := c.CopyMgmtClass("mgmtclass::testmgmtclass", "testmgmtclass2")
	utils.FailOnError(t, err)
}

func TestRenameMgmtClass(t *testing.T) {
	c := createStubHTTPClient(t, "rename-mgmtclass-req.xml", "rename-mgmtclass-res.xml")
	err := c.RenameMgmtClass("mgmtclass::testmgmtclass2", "mgmtclass1")
	utils.FailOnError(t, err)
}

func TestGetMgmtClassHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-mgmtclass-handle-req.xml", "get-mgmtclass-handle-res.xml")
	res, err := c.GetMgmtClassHandle("testmgmtclass")
	utils.FailOnError(t, err)

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
