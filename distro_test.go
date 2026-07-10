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

	if len(distros) != 4 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro")
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

	if len(distros) != 4 {
		t.Errorf("Wrong number of distros returned.")
	}
}

func TestGetDistrosSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distros-since")
	distros, err := c.GetDistrosSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(distros) != 4 {
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

func TestCreateDistro(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"create-distro-name-check",
		"new-distro",
		"new-distro-modify-name",
		"new-distro-modify-comment",
		"new-distro-modify-kernel-options",
		"new-distro-modify-kernel-options-post",
		"new-distro-modify-autoinstall-meta",
		"new-distro-modify-template-files",
		"new-distro-modify-owners",
		"new-distro-modify-arch",
		"new-distro-modify-boot-loaders",
		"new-distro-modify-breed",
		"new-distro-modify-initrd",
		"new-distro-modify-remote-boot-initrd",
		"new-distro-modify-kernel",
		"new-distro-modify-remote-boot-kernel",
		"new-distro-modify-redhat-management-key",
		"new-distro-modify-redhat-management-org",
		"new-distro-modify-redhat-management-user",
		"new-distro-modify-redhat-management-password",
		"new-distro-modify-os-version",
		"new-distro-save",
		"new-distro-get",
	})
	d := NewDistro()
	d.Name = "Ubuntu-20.04-x86_64"
	d.Kernel = "/code/system-tests/images/fake/vmlinuz"
	d.Initrd = "/code/system-tests/images/fake/initramfs"

	newDistro, err := c.CreateDistro(d)
	FailOnError(t, err)

	if newDistro.Name != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong distro name returned.")
	}
}

func TestUpdateDistro(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-distro-handle",
		"update-distro-modify-name",
		"update-distro-modify-comment",
		"update-distro-modify-kernel-options",
		"update-distro-modify-kernel-options-post",
		"update-distro-modify-autoinstall-meta",
		"update-distro-modify-template-files",
		"update-distro-modify-owners",
		"update-distro-modify-arch",
		"update-distro-modify-boot-loaders",
		"update-distro-modify-breed",
		"update-distro-modify-initrd",
		"update-distro-modify-remote-boot-initrd",
		"update-distro-modify-kernel",
		"update-distro-modify-remote-boot-kernel",
		"update-distro-modify-redhat-management-key",
		"update-distro-modify-redhat-management-org",
		"update-distro-modify-redhat-management-user",
		"update-distro-modify-redhat-management-password",
		"update-distro-modify-os-version",
		"update-distro-save",
	})
	d := NewDistro()
	d.Name = "Ubuntu-20.04-x86_64"
	d.Kernel = "/code/system-tests/images/fake/vmlinuz"
	d.Initrd = "/code/system-tests/images/fake/initramfs"

	err := c.UpdateDistro(&d)
	FailOnError(t, err)
}

func TestSaveDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-distro")
	err := c.SaveDistro("0000000000000000000000000000000f", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-distro")
	err := c.CopyDistro("0000000000000000000000000000000f", "test2")
	FailOnError(t, err)
}

func TestRenameDistro(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-distro")
	err := c.RenameDistro("00000000000000000000000000000010", "test1")
	FailOnError(t, err)
}

func TestGetDistroHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-handle")
	res, err := c.GetDistroHandle("test")
	FailOnError(t, err)

	if res != "0000000000000000000000000000000f" {
		t.Error("Wrong object id returned.")
	}
}

func TestGetValidDistroBootLoaders(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-distro-boot-loaders")
	res, err := c.GetValidDistroBootLoaders("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if len(res) < 1 {
		t.Error("Expected at least one boot loader.")
	}
}

func TestGetDistroAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-distro-as-rendered")
	res, err := c.GetDistroAsRendered("Ubuntu-20.04-x86_64")
	FailOnError(t, err)

	if res["name"] != "Ubuntu-20.04-x86_64" {
		t.Errorf("Wrong distro name returned: %v", res["name"])
	}
}
