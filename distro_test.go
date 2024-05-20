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

func TestGetDistros(t *testing.T) {
	c := createStubHTTPClient(t, "get-distros-req.xml", "get-distros-res.xml")
	distros, err := c.GetDistros()
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistro(t *testing.T) {
	c := createStubHTTPClient(t, "get-distro-req.xml", "get-distro-res.xml")
	distro, err := c.GetDistro("Ubuntu-20.04-x86_64")
	utils.FailOnError(t, err)

	if distro.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong distro returned.")
	}
}

func TestListDistroNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-distro-req.xml", "get-item-names-distro-res.xml")
	distros, err := c.ListDistroNames()
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistrosSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-distros-since-req.xml", "get-distros-since-res.xml")
	distros, err := c.GetDistrosSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestFindDistro(t *testing.T) {
	c := createStubHTTPClient(t, "find-distro-req.xml", "find-distro-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	distros, err := c.FindDistro(criteria)
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestFindDistroNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-distro-names-req.xml", "find-distro-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	distros, err := c.FindDistroNames(criteria)
	utils.FailOnError(t, err)

	if len(distros) != 1 {
		t.Error("Wrong number of distros returned.")
	}
}

func TestSaveDistro(t *testing.T) {
	c := createStubHTTPClient(t, "save-distro-req.xml", "save-distro-res.xml")
	err := c.SaveDistro("distro::test", "bypass")
	utils.FailOnError(t, err)
}

func TestCopyDistro(t *testing.T) {
	c := createStubHTTPClient(t, "copy-distro-req.xml", "copy-distro-res.xml")
	err := c.CopyDistro("distro::test", "test2")
	utils.FailOnError(t, err)
}

func TestRenameDistro(t *testing.T) {
	c := createStubHTTPClient(t, "rename-distro-req.xml", "rename-distro-res.xml")
	err := c.RenameDistro("distro::test2", "test1")
	utils.FailOnError(t, err)
}

func TestGetDistroHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-distro-handle-req.xml", "get-distro-handle-res.xml")
	res, err := c.GetDistroHandle("test")
	utils.FailOnError(t, err)

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
