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

func TestNewProfile(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewProfile()
}

func TestGetProfiles(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profiles")
	profiles, err := c.GetProfiles()
	FailOnError(t, err)

	if len(profiles) != 5 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfile(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-profile")

	// Act
	profile, err := c.GetProfile("Ubuntu-20.04-x86_64", false, false)

	// Assert
	FailOnError(t, err)
	if profile.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile returned.")
	}
}

func TestDeleteProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile")
	err := c.DeleteProfile("test")
	FailOnError(t, err)
}

func TestDeleteProfileRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-profile")
	err := c.DeleteProfileRecursive("test", false)
	FailOnError(t, err)
}

func TestListProfileNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-profile")
	profiles, err := c.ListProfileNames()
	FailOnError(t, err)

	if len(profiles) != 5 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestGetProfilesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profiles-since")
	profiles, err := c.GetProfilesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(profiles) != 5 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-profile")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	profiles, err := c.FindProfile(criteria, false)
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindProfileNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-profile-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	profiles, err := c.FindProfileNames(criteria)
	FailOnError(t, err)

	if len(profiles) != 1 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestCreateProfile(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"create-profile-name-check",
		"new-profile",
		"new-profile-modify-name",
		"new-profile-modify-parent",
		"new-profile-modify-distro",
		"new-profile-modify-comment",
		"new-profile-modify-kernel-options",
		"new-profile-modify-kernel-options-post",
		"new-profile-modify-autoinstall-meta",
		"new-profile-modify-template-files",
		"new-profile-modify-owners",
		"new-profile-modify-autoinstall",
		"new-profile-modify-boot-loaders",
		"new-profile-modify-dhcp-tag",
		"new-profile-modify-name-servers",
		"new-profile-modify-name-servers-search",
		"new-profile-modify-enable-ipxe",
		"new-profile-modify-enable-menu",
		"new-profile-modify-filename",
		"new-profile-modify-menu",
		"new-profile-modify-proxy",
		"new-profile-modify-redhat-management-key",
		"new-profile-modify-redhat-management-org",
		"new-profile-modify-redhat-management-user",
		"new-profile-modify-redhat-management-password",
		"new-profile-modify-repos",
		"new-profile-modify-server",
		"new-profile-modify-next-server-v4",
		"new-profile-modify-next-server-v6",
		"new-profile-modify-virt-auto-boot",
		"new-profile-modify-virt-cpus",
		"new-profile-modify-virt-disk-driver",
		"new-profile-modify-virt-file-size",
		"new-profile-modify-virt-path",
		"new-profile-modify-virt-pxe-boot",
		"new-profile-modify-virt-ram",
		"new-profile-modify-virt-type",
		"new-profile-modify-virt-bridge",
		"new-profile-save",
		"new-profile-get",
	})
	p := NewProfile()
	p.Name = "Ubuntu-20.04-x86_64"
	p.Distro = "0000000000000000000000000000000e"
	p.Virt.Path = "/var/lib/libvirt/images"
	p.Virt.Cpus = Value[int]{Data: 2}

	newProfile, err := c.CreateProfile(p)
	FailOnError(t, err)

	if newProfile.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile name returned.")
	}
}

func TestUpdateProfile(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-profile-handle",
		"update-profile-modify-name",
		"update-profile-modify-parent",
		"update-profile-modify-distro",
		"update-profile-modify-comment",
		"update-profile-modify-kernel-options",
		"update-profile-modify-kernel-options-post",
		"update-profile-modify-autoinstall-meta",
		"update-profile-modify-template-files",
		"update-profile-modify-owners",
		"update-profile-modify-autoinstall",
		"update-profile-modify-boot-loaders",
		"update-profile-modify-dhcp-tag",
		"update-profile-modify-name-servers",
		"update-profile-modify-name-servers-search",
		"update-profile-modify-enable-ipxe",
		"update-profile-modify-enable-menu",
		"update-profile-modify-filename",
		"update-profile-modify-menu",
		"update-profile-modify-proxy",
		"update-profile-modify-redhat-management-key",
		"update-profile-modify-redhat-management-org",
		"update-profile-modify-redhat-management-user",
		"update-profile-modify-redhat-management-password",
		"update-profile-modify-repos",
		"update-profile-modify-server",
		"update-profile-modify-next-server-v4",
		"update-profile-modify-next-server-v6",
		"update-profile-modify-virt-auto-boot",
		"update-profile-modify-virt-cpus",
		"update-profile-modify-virt-disk-driver",
		"update-profile-modify-virt-file-size",
		"update-profile-modify-virt-path",
		"update-profile-modify-virt-pxe-boot",
		"update-profile-modify-virt-ram",
		"update-profile-modify-virt-type",
		"update-profile-modify-virt-bridge",
		"update-profile-save",
	})
	p := NewProfile()
	p.Name = "Ubuntu-20.04-x86_64"
	p.Distro = "0000000000000000000000000000000e"
	p.Virt.Path = "/var/lib/libvirt/images"
	p.Virt.Cpus = Value[int]{Data: 2}

	err := c.UpdateProfile(&p)
	FailOnError(t, err)
}

func TestSaveProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-profile")
	err := c.SaveProfile("00000000000000000000000000000013", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-profile")
	err := c.CopyProfile("00000000000000000000000000000013", "testprof2")
	FailOnError(t, err)
}

func TestRenameProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-profile")
	err := c.RenameProfile("00000000000000000000000000000014", "testprof1")
	FailOnError(t, err)
}

func TestGetProfileHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-handle")
	res, err := c.GetProfileHandle("testprof")
	FailOnError(t, err)

	if res != "00000000000000000000000000000013" {
		t.Error("Wrong object id returned.")
	}
}

func TestGetValidProfileBootLoaders(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-profile-boot-loaders")
	res, err := c.GetValidProfileBootLoaders("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if len(res) < 1 {
		t.Error("Expected at least one boot loader.")
	}
}

func TestGetProfileAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-profile-as-rendered")
	res, err := c.GetProfileAsRendered("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if res["name"] != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong profile name returned: %v", res["name"])
	}
}

func TestNewSubprofile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "new-subprofile")
	res, err := c.NewSubprofile()
	FailOnError(t, err)
	if res == "" {
		t.Error("Expected a non-empty object ID from NewSubprofile.")
	}
}
