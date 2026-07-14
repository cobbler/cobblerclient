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

func TestNewSystemGroup(t *testing.T) {
	g := NewSystemGroup()
	if g.Members == nil {
		t.Error("Items should not be nil")
	}
}

func TestGetSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-group")
	g, err := c.GetSystemGroup("webservers", false, false)
	FailOnError(t, err)
	if g.Name != "webservers" {
		t.Errorf("wrong name %q", g.Name)
	}
}

func TestGetSystemGroupHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-group-handle")
	handle, err := c.GetSystemGroupHandle("webservers")
	FailOnError(t, err)
	if handle != "000000000000000000000000000000dd" {
		t.Errorf("wrong handle: %q", handle)
	}
}

func TestListSystemGroupNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-system-group")
	names, err := c.ListSystemGroupNames()
	FailOnError(t, err)
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestDeleteSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-system-group")
	err := c.DeleteSystemGroup("webservers")
	FailOnError(t, err)
}

func TestDeleteSystemGroupRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-system-group")
	err := c.DeleteSystemGroupRecursive("webservers", false)
	FailOnError(t, err)
}

func TestCopySystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-system-group")
	err := c.CopySystemGroup("000000000000000000000000000000dd", "webservers2")
	FailOnError(t, err)
}

func TestRenameSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-system-group")
	err := c.RenameSystemGroup("000000000000000000000000000000de", "webservers-new")
	FailOnError(t, err)
}

func TestSaveSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-system-group")
	err := c.SaveSystemGroup("000000000000000000000000000000dd", true, true, "bypass")
	FailOnError(t, err)
}

func TestGetSystemGroupAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-group-as-rendered")
	res, err := c.GetSystemGroupAsRendered("webservers")
	FailOnError(t, err)

	if res["name"] != "webservers" {
		t.Errorf("Wrong system group name returned: %v", res["name"])
	}
}

func TestFindSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-system-group")
	criteria := map[string]interface{}{"name": "webservers"}
	res, err := c.FindSystemGroup(criteria, false)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 system group, got %d.", len(res))
	}
}

func TestGetSystemGroups(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-groups")
	groups, err := c.GetSystemGroups()
	FailOnError(t, err)
	if len(groups) != 2 {
		t.Errorf("Expected 2 system groups, got %d.", len(groups))
	}
}

func TestGetSystemGroupsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-groups-since")
	groups, err := c.GetSystemGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(groups) != 2 {
		t.Errorf("Expected 2 system groups, got %d.", len(groups))
	}
}

func TestCreateSystemGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"new-system-group",
		"new-system-group-modify-name",
		"new-system-group-modify-comment",
		"new-system-group-modify-kernel-options",
		"new-system-group-modify-kernel-options-post",
		"new-system-group-modify-autoinstall-meta",
		"new-system-group-modify-template-files",
		"new-system-group-modify-owners",
		"new-system-group-modify-items",
		"new-system-group-save",
		"new-system-group-get",
	})
	g := NewSystemGroup()
	g.Name = "webservers"
	g.Members = []string{"member-a"}
	result, err := c.CreateSystemGroup(g)
	FailOnError(t, err)
	if result.Name != "webservers" {
		t.Errorf("Wrong group name: %v", result.Name)
	}
}

func TestUpdateSystemGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-system-group-handle",
		"update-system-group-modify-name",
		"update-system-group-modify-comment",
		"update-system-group-modify-kernel-options",
		"update-system-group-modify-kernel-options-post",
		"update-system-group-modify-autoinstall-meta",
		"update-system-group-modify-template-files",
		"update-system-group-modify-owners",
		"update-system-group-modify-items",
		"update-system-group-save",
	})
	g := NewSystemGroup()
	g.Name = "webservers"
	g.Members = []string{"member-a"}
	err := c.UpdateSystemGroup(&g)
	FailOnError(t, err)
}
