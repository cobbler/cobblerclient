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

func TestNewProfileGroup(t *testing.T) {
	g := NewProfileGroup()
	if g.Members == nil {
		t.Error("Items should not be nil")
	}
}

func TestGetProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-group")
	g, err := c.GetProfileGroup("000000000000000000000000000000db", false, false)
	FailOnError(t, err)
	if g.Name != "webservers" {
		t.Errorf("wrong name %q", g.Name)
	}
}

func TestGetProfileGroupHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-group-handle")
	handle, err := c.GetProfileGroupHandle("webservers")
	FailOnError(t, err)
	if handle != "000000000000000000000000000000db" {
		t.Errorf("wrong handle: %q", handle)
	}
}

func TestListProfileGroupNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-profile-group")
	names, err := c.ListProfileGroupNames()
	FailOnError(t, err)
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestDeleteProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile-group")
	err := c.DeleteProfileGroup("000000000000000000000000000000db")
	FailOnError(t, err)
}

func TestDeleteProfileGroupRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile-group")
	err := c.DeleteProfileGroupRecursive("000000000000000000000000000000db", false)
	FailOnError(t, err)
}

func TestCopyProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-profile-group")
	err := c.CopyProfileGroup("000000000000000000000000000000db", "webservers2")
	FailOnError(t, err)
}

func TestRenameProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-profile-group")
	err := c.RenameProfileGroup("000000000000000000000000000000dc", "webservers-new")
	FailOnError(t, err)
}

func TestSaveProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-profile-group")
	err := c.SaveProfileGroup("000000000000000000000000000000db", true, true, "bypass")
	FailOnError(t, err)
}

func TestGetProfileGroupAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-group-as-rendered")
	res, err := c.GetProfileGroupAsRendered("webservers")
	FailOnError(t, err)

	if res["name"] != "webservers" {
		t.Errorf("Wrong profile group name returned: %v", res["name"])
	}
}

func TestFindProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-profile-group")
	criteria := map[string]interface{}{"name": "webservers"}
	res, err := c.FindProfileGroup(criteria, false)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 profile group, got %d.", len(res))
	}
}

func TestGetProfileGroups(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-groups")
	groups, err := c.GetProfileGroups()
	FailOnError(t, err)
	if len(groups) != 2 {
		t.Errorf("Expected 2 profile groups, got %d.", len(groups))
	}
}

func TestGetProfileGroupsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-groups-since")
	groups, err := c.GetProfileGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(groups) != 2 {
		t.Errorf("Expected 2 profile groups, got %d.", len(groups))
	}
}

func TestCreateProfileGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"create-profile-group-name-check",
		"new-profile-group",
		"new-profile-group-modify-name",
		"new-profile-group-modify-comment",
		"new-profile-group-modify-kernel-options",
		"new-profile-group-modify-kernel-options-post",
		"new-profile-group-modify-autoinstall-meta",
		"new-profile-group-modify-template-files",
		"new-profile-group-modify-owners",
		"new-profile-group-modify-items",
		"new-profile-group-save",
		"new-profile-group-get",
	})
	g := NewProfileGroup()
	g.Name = "webservers"
	g.Members = []string{"member-a"}
	result, err := c.CreateProfileGroup(g)
	FailOnError(t, err)
	if result.Name != "webservers" {
		t.Errorf("Wrong group name: %v", result.Name)
	}
}

func TestUpdateProfileGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-profile-group-handle",
		"update-profile-group-modify-name",
		"update-profile-group-modify-comment",
		"update-profile-group-modify-kernel-options",
		"update-profile-group-modify-kernel-options-post",
		"update-profile-group-modify-autoinstall-meta",
		"update-profile-group-modify-template-files",
		"update-profile-group-modify-owners",
		"update-profile-group-modify-items",
		"update-profile-group-save",
	})
	g := NewProfileGroup()
	g.Name = "webservers"
	g.Members = []string{"member-a"}
	err := c.UpdateProfileGroup(&g)
	FailOnError(t, err)
}
