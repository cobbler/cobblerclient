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

func TestGetPackages(t *testing.T) {
	c := createStubHTTPClient(t, "get-packages-req.xml", "get-packages-res.xml")
	linuxpackages, err := c.GetPackages()
	utils.FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestGetPackage(t *testing.T) {
	c := createStubHTTPClient(t, "get-package-req.xml", "get-package-res.xml")
	distro, err := c.GetPackage("testpackage")
	utils.FailOnError(t, err)

	if distro.Name != "testpackage" {
		t.Errorf("Wrong package returned.")
	}
}

func TestListPackageNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-package-req.xml", "get-item-names-package-res.xml")
	linuxpackages, err := c.ListPackageNames()
	utils.FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestGetPackagesSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-packages-since-req.xml", "get-packages-since-res.xml")
	linuxpackages, err := c.GetPackagesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestFindPackage(t *testing.T) {
	c := createStubHTTPClient(t, "find-package-req.xml", "find-package-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testpackage"
	linuxpackages, err := c.FindPackage(criteria)
	utils.FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestFindPackageNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-package-names-req.xml", "find-package-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testpackage"
	linuxpackages, err := c.FindPackageNames(criteria)
	utils.FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Error("Wrong number of packages returned.")
	}
}

func TestSavePackage(t *testing.T) {
	c := createStubHTTPClient(t, "save-package-req.xml", "save-package-res.xml")
	err := c.SavePackage("package::testpackage", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyPackage(t *testing.T) {
	c := createStubHTTPClient(t, "copy-package-req.xml", "copy-package-res.xml")
	err := c.CopyPackage("package::testpackage", "testpackage2")
	utils.FailOnError(t, err)
}

func TestRenamePackage(t *testing.T) {
	c := createStubHTTPClient(t, "rename-package-req.xml", "rename-package-res.xml")
	err := c.RenamePackage("package::testpackage2", "testpackage1")
	utils.FailOnError(t, err)
}

func TestGetPackageHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-package-handle-req.xml", "get-package-handle-res.xml")
	res, err := c.GetPackageHandle("testpackage")
	utils.FailOnError(t, err)

	if res != "package::testpackage" {
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
