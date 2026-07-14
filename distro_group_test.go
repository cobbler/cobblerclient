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

func TestNewDistroGroup(t *testing.T) {
	g := NewDistroGroup()
	if g.Members == nil {
		t.Error("Items should not be nil")
	}
}

func TestGetDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-group")
	g, err := c.GetDistroGroup("webservers", false, false)
	FailOnError(t, err)
	if g.Name != "webservers" {
		t.Errorf("wrong name %q", g.Name)
	}
	if len(g.Members) != 1 {
		t.Errorf("expected 1 member, got %d", len(g.Members))
	}
}

func TestGetDistroGroupHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-group-handle")
	handle, err := c.GetDistroGroupHandle("webservers")
	FailOnError(t, err)
	if handle != "000000000000000000000000000000d9" {
		t.Errorf("wrong handle: %q", handle)
	}
}

func TestListDistroGroupNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-distro-group")
	names, err := c.ListDistroGroupNames()
	FailOnError(t, err)
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestDeleteDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-distro-group")
	err := c.DeleteDistroGroup("webservers")
	FailOnError(t, err)
}

func TestDeleteDistroGroupRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-distro-group")
	err := c.DeleteDistroGroupRecursive("webservers", false)
	FailOnError(t, err)
}

func TestCopyDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-distro-group")
	err := c.CopyDistroGroup("000000000000000000000000000000d9", "webservers2")
	FailOnError(t, err)
}

func TestRenameDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-distro-group")
	err := c.RenameDistroGroup("000000000000000000000000000000da", "webservers-new")
	FailOnError(t, err)
}

func TestSaveDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-distro-group")
	err := c.SaveDistroGroup("000000000000000000000000000000d9", true, true, "bypass")
	FailOnError(t, err)
}

func TestGetDistroGroupAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-group-as-rendered")
	res, err := c.GetDistroGroupAsRendered("webservers")
	FailOnError(t, err)

	if res["name"] != "webservers" {
		t.Errorf("Wrong distro group name returned: %v", res["name"])
	}
}

func TestFindDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-distro-group")
	criteria := map[string]interface{}{"name": "webservers"}
	res, err := c.FindDistroGroup(criteria, false)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 distro group, got %d.", len(res))
	}
}

func TestGetDistroGroups(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-groups")
	groups, err := c.GetDistroGroups()
	FailOnError(t, err)
	if len(groups) != 2 {
		t.Errorf("Expected 2 distro groups, got %d.", len(groups))
	}
}

func TestGetDistroGroupsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-groups-since")
	groups, err := c.GetDistroGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(groups) != 2 {
		t.Errorf("Expected 2 distro groups, got %d.", len(groups))
	}
}

func TestCreateDistroGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"new-distro-group",
		"new-distro-group-modify-name",
		"new-distro-group-modify-comment",
		"new-distro-group-modify-kernel-options",
		"new-distro-group-modify-kernel-options-post",
		"new-distro-group-modify-autoinstall-meta",
		"new-distro-group-modify-template-files",
		"new-distro-group-modify-owners",
		"new-distro-group-modify-items",
		"new-distro-group-save",
		"new-distro-group-get",
	})
	g := NewDistroGroup()
	g.Name = "webservers"
	g.Members = []string{"member-a"}
	result, err := c.CreateDistroGroup(g)
	FailOnError(t, err)
	if result.Name != "webservers" {
		t.Errorf("Wrong group name: %v", result.Name)
	}
}

func TestUpdateDistroGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-distro-group-handle",
		"update-distro-group-modify-name",
		"update-distro-group-modify-comment",
		"update-distro-group-modify-kernel-options",
		"update-distro-group-modify-kernel-options-post",
		"update-distro-group-modify-autoinstall-meta",
		"update-distro-group-modify-template-files",
		"update-distro-group-modify-owners",
		"update-distro-group-modify-items",
		"update-distro-group-save",
	})
	g := NewDistroGroup()
	g.Name = "webservers"
	g.Members = []string{"member-a"}
	err := c.UpdateDistroGroup(&g)
	FailOnError(t, err)
}
