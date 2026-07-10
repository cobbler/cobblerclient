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
		t.Error("Members should not be nil")
	}
}
func TestNewProfileGroup(t *testing.T) {
	g := NewProfileGroup()
	if g.Members == nil {
		t.Error("Members should not be nil")
	}
}
func TestNewSystemGroup(t *testing.T) {
	g := NewSystemGroup()
	if g.Members == nil {
		t.Error("Members should not be nil")
	}
}

func TestGetDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-group")
	g, err := c.GetDistroGroup("webservers", false, false)
	FailOnError(t, err)
	if g.Name != "webservers" {
		t.Errorf("wrong name %q", g.Name)
	}
	if len(g.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(g.Members))
	}
}

func TestGetProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-group")
	g, err := c.GetProfileGroup("webservers", false, false)
	FailOnError(t, err)
	if g.Name != "webservers" {
		t.Errorf("wrong name %q", g.Name)
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

func TestGetDistroGroupHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-group-handle")
	handle, err := c.GetDistroGroupHandle("webservers")
	FailOnError(t, err)
	if handle != "distro_group::webservers" {
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
	err := c.CopyDistroGroup("distro_group::webservers", "webservers2")
	FailOnError(t, err)
}

func TestRenameDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-distro-group")
	err := c.RenameDistroGroup("distro_group::webservers", "webservers-new")
	FailOnError(t, err)
}

func TestSaveDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-distro-group")
	err := c.SaveDistroGroup("distro_group::webservers", true, true, "bypass")
	FailOnError(t, err)
}

func TestGetProfileGroupHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-group-handle")
	handle, err := c.GetProfileGroupHandle("webservers")
	FailOnError(t, err)
	if handle != "profile_group::webservers" {
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
	err := c.DeleteProfileGroup("webservers")
	FailOnError(t, err)
}

func TestDeleteProfileGroupRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile-group")
	err := c.DeleteProfileGroupRecursive("webservers", false)
	FailOnError(t, err)
}

func TestCopyProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-profile-group")
	err := c.CopyProfileGroup("profile_group::webservers", "webservers2")
	FailOnError(t, err)
}

func TestRenameProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-profile-group")
	err := c.RenameProfileGroup("profile_group::webservers", "webservers-new")
	FailOnError(t, err)
}

func TestSaveProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-profile-group")
	err := c.SaveProfileGroup("profile_group::webservers", true, true, "bypass")
	FailOnError(t, err)
}

func TestGetSystemGroupHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-group-handle")
	handle, err := c.GetSystemGroupHandle("webservers")
	FailOnError(t, err)
	if handle != "system_group::webservers" {
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
	err := c.CopySystemGroup("system_group::webservers", "webservers2")
	FailOnError(t, err)
}

func TestRenameSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-system-group")
	err := c.RenameSystemGroup("system_group::webservers", "webservers-new")
	FailOnError(t, err)
}

func TestSaveSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-system-group")
	err := c.SaveSystemGroup("system_group::webservers", true, true, "bypass")
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

func TestGetProfileGroupAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-group-as-rendered")
	res, err := c.GetProfileGroupAsRendered("webservers")
	FailOnError(t, err)

	if res["name"] != "webservers" {
		t.Errorf("Wrong profile group name returned: %v", res["name"])
	}
}

func TestGetSystemGroupAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-group-as-rendered")
	res, err := c.GetSystemGroupAsRendered("webservers")
	FailOnError(t, err)

	if res["name"] != "webservers" {
		t.Errorf("Wrong system group name returned: %v", res["name"])
	}
}

func TestFindDistroGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-distro-group")
	criteria := map[string]interface{}{"name": "webservers"}
	res, err := c.FindDistroGroup(criteria)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 distro group, got %d.", len(res))
	}
}

func TestFindProfileGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-profile-group")
	criteria := map[string]interface{}{"name": "webservers"}
	res, err := c.FindProfileGroup(criteria)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 profile group, got %d.", len(res))
	}
}

func TestFindSystemGroup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-system-group")
	criteria := map[string]interface{}{"name": "webservers"}
	res, err := c.FindSystemGroup(criteria)
	FailOnError(t, err)
	if len(res) != 1 {
		t.Errorf("Expected 1 system group, got %d.", len(res))
	}
}

func TestGetDistroGroups(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-groups")
	groups, err := c.GetDistroGroups()
	FailOnError(t, err)
	if len(groups) != 1 {
		t.Errorf("Expected 1 distro group, got %d.", len(groups))
	}
}

func TestGetDistroGroupsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-groups-since")
	groups, err := c.GetDistroGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(groups) != 1 {
		t.Errorf("Expected 1 distro group, got %d.", len(groups))
	}
}

func TestGetProfileGroups(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-groups")
	groups, err := c.GetProfileGroups()
	FailOnError(t, err)
	if len(groups) != 1 {
		t.Errorf("Expected 1 profile group, got %d.", len(groups))
	}
}

func TestGetProfileGroupsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-groups-since")
	groups, err := c.GetProfileGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(groups) != 1 {
		t.Errorf("Expected 1 profile group, got %d.", len(groups))
	}
}

func TestGetSystemGroups(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-groups")
	groups, err := c.GetSystemGroups()
	FailOnError(t, err)
	if len(groups) != 1 {
		t.Errorf("Expected 1 system group, got %d.", len(groups))
	}
}

func TestGetSystemGroupsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-groups-since")
	groups, err := c.GetSystemGroupsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)
	if len(groups) != 1 {
		t.Errorf("Expected 1 system group, got %d.", len(groups))
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
		"new-distro-group-modify-fetchable-files",
		"new-distro-group-modify-boot-files",
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
		"update-distro-group-modify-fetchable-files",
		"update-distro-group-modify-boot-files",
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

func TestCreateProfileGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"new-profile-group",
		"new-profile-group-modify-name",
		"new-profile-group-modify-comment",
		"new-profile-group-modify-kernel-options",
		"new-profile-group-modify-kernel-options-post",
		"new-profile-group-modify-autoinstall-meta",
		"new-profile-group-modify-fetchable-files",
		"new-profile-group-modify-boot-files",
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
		"update-profile-group-modify-fetchable-files",
		"update-profile-group-modify-boot-files",
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

func TestCreateSystemGroup(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"new-system-group",
		"new-system-group-modify-name",
		"new-system-group-modify-comment",
		"new-system-group-modify-kernel-options",
		"new-system-group-modify-kernel-options-post",
		"new-system-group-modify-autoinstall-meta",
		"new-system-group-modify-fetchable-files",
		"new-system-group-modify-boot-files",
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
		"update-system-group-modify-fetchable-files",
		"update-system-group-modify-boot-files",
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
