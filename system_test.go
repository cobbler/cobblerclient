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

func TestGetSystems(t *testing.T) {
	c := createStubHTTPClient(t, "get-systems-req.xml", "get-systems-res.xml")
	systems, err := c.GetSystems()
	utils.FailOnError(t, err)

	if len(systems) != 1 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystem(t *testing.T) {
	c := createStubHTTPClient(t, "get-system-req.xml", "get-system-res.xml")
	system, err := c.GetSystem("test")
	utils.FailOnError(t, err)

	if system.Name != "test" {
		t.Errorf("Wrong system returned.")
	}
}

func TestNewSystem(t *testing.T) {
	c := createStubHTTPClient(t, "new-system-req.xml", "new-system-res.xml")
	result, err := c.Call("new_system", c.Token)
	utils.FailOnError(t, err)
	newID := result.(string)

	if newID != "___NEW___system::abc123==" {
		t.Errorf("Wrong ID returned.")
	}

	c = createStubHTTPClient(t, "set-system-hostname-req.xml", "set-system-hostname-res.xml")
	result, err = c.Call("modify_system", newID, "hostname", "blahhost", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting hostname failed.")
	}

	c = createStubHTTPClient(t, "set-system-name-req.xml", "set-system-name-res.xml")
	result, err = c.Call("modify_system", newID, "name", "mytestsystem", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name failed.")
	}

	c = createStubHTTPClient(t, "set-system-nameservers-req.xml", "set-system-nameservers-res.xml")
	result, err = c.Call("modify_system", newID, "name_servers", "8.8.8.8 8.8.4.4", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name servers failed.")
	}

	c = createStubHTTPClient(t, "set-system-profile-req.xml", "set-system-profile-res.xml")
	result, err = c.Call("modify_system", newID, "profile", "centos7-x86_64", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name servers failed.")
	}

	/* I'm not sure how to get this test to pass with unordered maps
	nicInfo := map[string]interface{}{
		"macaddress-eth0":  "01:02:03:04:05:06",
		"ipaddress-eth0":   "1.2.3.4",
		"dnsname-eth0":     "deathstar",
		"subnetsmask-eth0": "255.255.255.0",
		"if-gateway-eth0":  "4.3.2.1",
	}

	c = createStubHTTPClient(t, "set-system-network-req.xml", "set-system-network-res.xml")
	result, err = c.Call("modify_system", newID, "modify_interface", nicInfo, c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting interface failed.")
	}
	*/

	c = createStubHTTPClient(t, "save-system-req.xml", "save-system-res.xml")
	err = c.SaveSystem(newID, "bypass")
	utils.FailOnError(t, err)
}

func TestDeleteSystem(t *testing.T) {
	c := createStubHTTPClient(t, "delete-system-req.xml", "delete-system-res.xml")
	err := c.DeleteSystem("test")
	utils.FailOnError(t, err)
}

func TestDeleteSystemRecursive(t *testing.T) {
	c := createStubHTTPClient(t, "delete-system-req.xml", "delete-system-res.xml")
	err := c.DeleteSystemRecursive("test", false)
	utils.FailOnError(t, err)
}

func TestListSystemNames(t *testing.T) {
	c := createStubHTTPClient(t, "get-item-names-system-req.xml", "get-item-names-system-res.xml")
	sytems, err := c.ListSystemNames()
	utils.FailOnError(t, err)

	if len(sytems) != 1 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystemsSince(t *testing.T) {
	c := createStubHTTPClient(t, "get-system-since-req.xml", "get-system-since-res.xml")
	systems, err := c.GetSystemsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	utils.FailOnError(t, err)

	if len(systems) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindSystem(t *testing.T) {
	c := createStubHTTPClient(t, "find-system-req.xml", "find-system-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	_, err := c.FindSystem(criteria)
	utils.FailOnError(t, err)
}

func TestFindSystemNames(t *testing.T) {
	c := createStubHTTPClient(t, "find-system-names-req.xml", "find-system-names-res.xml")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	_, err := c.FindSystem(criteria)
	utils.FailOnError(t, err)
}

func TestSaveSystem(t *testing.T) {
	c := createStubHTTPClient(t, "save-system-req.xml", "save-system-res.xml")
	err := c.SaveSystem("___NEW___system::abc123==", "bypass")
	utils.FailOnError(t, err)
}

func TestCopySystem(t *testing.T) {
	c := createStubHTTPClient(t, "copy-system-req.xml", "copy-system-res.xml")
	err := c.CopySystem("system::testsys", "testsys2")
	utils.FailOnError(t, err)
}

func TestRenameSystem(t *testing.T) {
	c := createStubHTTPClient(t, "rename-system-req.xml", "rename-system-res.xml")
	err := c.RenameSystem("system::testsys", "testsys1")
	utils.FailOnError(t, err)
}

func TestGetSystemHandle(t *testing.T) {
	c := createStubHTTPClient(t, "get-system-handle-req.xml", "get-system-handle-res.xml")
	res, err := c.GetSystemHandle("testsys")
	utils.FailOnError(t, err)

	if res != "system::testsys" {
		t.Error("Wrong object id returned.")
	}
}
