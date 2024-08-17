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

func TestGetSystems(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-systems")
	systems, err := c.GetSystems()
	FailOnError(t, err)

	if len(systems) != 1 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system")
	system, err := c.GetSystem("test")
	FailOnError(t, err)

	if system.Name != "test" {
		t.Errorf("Wrong system returned.")
	}
}

func TestNewSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "new-system")
	result, err := c.Call("new_system", c.Token)
	FailOnError(t, err)
	newID := result.(string)

	if newID != "___NEW___system::abc123==" {
		t.Errorf("Wrong ID returned.")
	}

	c = createStubHTTPClientSingle(t, "set-system-hostname")
	result, err = c.Call("modify_system", newID, "hostname", "blahhost", c.Token)
	FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting hostname failed.")
	}

	c = createStubHTTPClientSingle(t, "set-system-name")
	result, err = c.Call("modify_system", newID, "name", "mytestsystem", c.Token)
	FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name failed.")
	}

	c = createStubHTTPClientSingle(t, "set-system-nameservers")
	result, err = c.Call("modify_system", newID, "name_servers", "8.8.8.8 8.8.4.4", c.Token)
	FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name servers failed.")
	}

	c = createStubHTTPClientSingle(t, "set-system-profile")
	result, err = c.Call("modify_system", newID, "profile", "centos7-x86_64", c.Token)
	FailOnError(t, err)

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

	c = createStubHTTPClientSingle(t, "set-system-network-req.xml", "set-system-network-res.xml")
	result, err = c.Call("modify_system", newID, "modify_interface", nicInfo, c.Token)
	FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting interface failed.")
	}
	*/

	c = createStubHTTPClientSingle(t, "save-system")
	err = c.SaveSystem(newID, "bypass")
	FailOnError(t, err)
}

func TestDeleteSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-system")
	err := c.DeleteSystem("test")
	FailOnError(t, err)
}

func TestDeleteSystemRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-system")
	err := c.DeleteSystemRecursive("test", false)
	FailOnError(t, err)
}

func TestListSystemNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-system")
	sytems, err := c.ListSystemNames()
	FailOnError(t, err)

	if len(sytems) != 1 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystemsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-since")
	systems, err := c.GetSystemsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(systems) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-system")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	_, err := c.FindSystem(criteria)
	FailOnError(t, err)
}

func TestFindSystemNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-system-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	_, err := c.FindSystem(criteria)
	FailOnError(t, err)
}

func TestSaveSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-system")
	err := c.SaveSystem("___NEW___system::abc123==", "bypass")
	FailOnError(t, err)
}

func TestCopySystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-system")
	err := c.CopySystem("system::testsys", "testsys2")
	FailOnError(t, err)
}

func TestRenameSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-system")
	err := c.RenameSystem("system::testsys", "testsys1")
	FailOnError(t, err)
}

func TestGetSystemHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-handle")
	res, err := c.GetSystemHandle("testsys")
	FailOnError(t, err)

	if res != "system::testsys" {
		t.Error("Wrong object id returned.")
	}
}
