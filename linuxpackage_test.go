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

func TestNewPackage(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewPackage()
}

func TestGetPackages(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-packages")
	linuxpackages, err := c.GetPackages()
	FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestGetPackage(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-package")
	c.CachedVersion = CobblerVersion{3, 3, 2}

	// Act
	linuxpackage, err := c.GetPackage("testpackage", false, false)

	// Assert
	FailOnError(t, err)
	if linuxpackage.Name != "testpackage" {
		t.Errorf("Wrong package returned.")
	}
}

func TestDeletePackage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-package")
	err := c.DeletePackage("test")
	FailOnError(t, err)
}

func TestDeletePackageRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-package")
	err := c.DeletePackageRecursive("test", false)
	FailOnError(t, err)
}

func TestListPackageNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-package")
	linuxpackages, err := c.ListPackageNames()
	FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestGetPackagesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-packages-since")
	linuxpackages, err := c.GetPackagesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestFindPackage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-package")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testpackage"
	linuxpackages, err := c.FindPackage(criteria)
	FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Errorf("Wrong number of packages returned.")
	}
}

func TestFindPackageNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-package-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "testpackage"
	linuxpackages, err := c.FindPackageNames(criteria)
	FailOnError(t, err)

	if len(linuxpackages) != 1 {
		t.Error("Wrong number of packages returned.")
	}
}

func TestSavePackage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-package")
	err := c.SavePackage("package::testpackage", "bypass")
	FailOnError(t, err)
}

func TestCopyPackage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-package")
	err := c.CopyPackage("package::testpackage", "testpackage2")
	FailOnError(t, err)
}

func TestRenamePackage(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-package")
	err := c.RenamePackage("package::testpackage2", "testpackage1")
	FailOnError(t, err)
}

func TestGetPackageHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-package-handle")
	res, err := c.GetPackageHandle("testpackage")
	FailOnError(t, err)

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
