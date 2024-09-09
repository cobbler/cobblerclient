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

func TestNewDistro(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewDistro()
}

func TestGetDistros(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distros")
	distros, err := c.GetDistros()
	FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro")
	c.CachedVersion = CobblerVersion{3, 3, 2}
	distro, err := c.GetDistro("Ubuntu-20.04-x86_64", false, false)
	FailOnError(t, err)

	if distro.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong distro returned.")
	}
}

func TestDeleteDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-distro")
	err := c.DeleteDistro("test")
	FailOnError(t, err)
}

func TestDeleteDistroRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-distro")
	err := c.DeleteDistroRecursive("test", false)
	FailOnError(t, err)
}

func TestListDistroNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-distro")
	distros, err := c.ListDistroNames()
	FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistrosSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distros-since")
	distros, err := c.GetDistrosSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestFindDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-distro")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	distros, err := c.FindDistro(criteria)
	FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestFindDistroNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-distro-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	distros, err := c.FindDistroNames(criteria)
	FailOnError(t, err)

	if len(distros) != 1 {
		t.Error("Wrong number of distros returned.")
	}
}

func TestSaveDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-distro")
	err := c.SaveDistro("distro::test", "bypass")
	FailOnError(t, err)
}

func TestCopyDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-distro")
	err := c.CopyDistro("distro::test", "test2")
	FailOnError(t, err)
}

func TestRenameDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-distro")
	err := c.RenameDistro("distro::test2", "test1")
	FailOnError(t, err)
}

func TestGetDistroHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-handle")
	res, err := c.GetDistroHandle("test")
	FailOnError(t, err)

	if res != "distro::test" {
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
