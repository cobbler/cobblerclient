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

func TestNewSystem(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewSystem()
}

// TestNewInterface removed - use TestNewNetworkInterface instead (4.0.0)

func TestGetSystems(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-systems")
	systems, err := c.GetSystems()
	FailOnError(t, err)

	if len(systems) != 6 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystem(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-system")

	// Act
	system, err := c.GetSystem("00000000000000000000000000000060", false, false)

	// Assert
	FailOnError(t, err)
	if system.Name != "test" {
		t.Errorf("Wrong system returned.")
	}
}

func TestSystemCreate(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"create-system-name-check",
		"new-system",
		"set-system-profile",
		"new-system-modify-image",
		"set-system-name",
		"new-system-modify-comment",
		"new-system-modify-kernel-options",
		"new-system-modify-kernel-options-post",
		"new-system-modify-autoinstall-meta",
		"new-system-modify-template-files",
		"new-system-modify-owners",
		"new-system-modify-autoinstall",
		"new-system-modify-boot-loaders",
		"set-system-nameservers",
		"new-system-modify-name-servers-search",
		"new-system-modify-enable-ipxe",
		"new-system-modify-filename",
		"new-system-modify-gateway",
		"set-system-hostname",
		"new-system-modify-ipv6-default-device",
		"new-system-modify-netboot-enabled",
		"new-system-modify-power-type",
		"new-system-modify-power-identity-file",
		"new-system-modify-power-options",
		"new-system-modify-power-user",
		"new-system-modify-power-password",
		"new-system-modify-power-address",
		"new-system-modify-power-id",
		"new-system-modify-proxy",
		"new-system-modify-redhat-management-key",
		"new-system-modify-redhat-management-org",
		"new-system-modify-redhat-management-user",
		"new-system-modify-redhat-management-password",
		"new-system-modify-serial-baud-rate",
		"new-system-modify-serial-device",
		"new-system-modify-server",
		"new-system-modify-status",
		"new-system-modify-next-server-v4",
		"new-system-modify-next-server-v6",
		"new-system-modify-virt-auto-boot",
		"new-system-modify-virt-cpus",
		"new-system-modify-virt-disk-driver",
		"new-system-modify-virt-file-size",
		"new-system-modify-virt-path",
		"new-system-modify-virt-pxe-boot",
		"new-system-modify-virt-ram",
		"new-system-modify-virt-type",
		"new-system-modify-virt-uefi",
		"new-system-modify-virt-pxe-boot-flag",
		"new-system-save",
		"new-system-get",
	})
	sys := NewSystem()
	sys.Name = "mytestsystem"
	sys.Hostname = "blahhost"
	sys.DNS.NameServers.Data = []string{"8.8.8.8", "8.8.4.4"}
	// Profile must be the referenced profile's real uid (centos7-x86_64), not its name.
	sys.Profile = "00000000000000000000000000000016"

	// Act
	newSys, err := c.CreateSystem(sys)

	// Assert
	FailOnError(t, err)
	if newSys.Name != "mytestsystem" {
		t.Errorf("Wrong system name returned.")
	}
	if newSys.Hostname != "blahhost" {
		t.Errorf("Wrong system hostname returned.")
	}
	if len(newSys.DNS.NameServers.Data) != 2 || newSys.DNS.NameServers.Data[0] != "8.8.8.8" {
		t.Errorf("Wrong system name servers returned.")
	}
	if newSys.Profile != "00000000000000000000000000000016" {
		t.Errorf("Wrong system profile returned.")
	}
}

func TestUpdateSystem(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-system-handle",
		"update-system-modify-profile",
		"update-system-modify-image",
		"update-system-modify-name",
		"update-system-modify-comment",
		"update-system-modify-kernel-options",
		"update-system-modify-kernel-options-post",
		"update-system-modify-autoinstall-meta",
		"update-system-modify-template-files",
		"update-system-modify-owners",
		"update-system-modify-autoinstall",
		"update-system-modify-boot-loaders",
		"update-system-modify-name-servers",
		"update-system-modify-name-servers-search",
		"update-system-modify-enable-ipxe",
		"update-system-modify-filename",
		"update-system-modify-gateway",
		"update-system-modify-hostname",
		"update-system-modify-ipv6-default-device",
		"update-system-modify-netboot-enabled",
		"update-system-modify-power-type",
		"update-system-modify-power-identity-file",
		"update-system-modify-power-options",
		"update-system-modify-power-user",
		"update-system-modify-power-password",
		"update-system-modify-power-address",
		"update-system-modify-power-id",
		"update-system-modify-proxy",
		"update-system-modify-redhat-management-key",
		"update-system-modify-redhat-management-org",
		"update-system-modify-redhat-management-user",
		"update-system-modify-redhat-management-password",
		"update-system-modify-serial-baud-rate",
		"update-system-modify-serial-device",
		"update-system-modify-server",
		"update-system-modify-status",
		"update-system-modify-next-server-v4",
		"update-system-modify-next-server-v6",
		"update-system-modify-virt-auto-boot",
		"update-system-modify-virt-cpus",
		"update-system-modify-virt-disk-driver",
		"update-system-modify-virt-file-size",
		"update-system-modify-virt-path",
		"update-system-modify-virt-pxe-boot",
		"update-system-modify-virt-ram",
		"update-system-modify-virt-type",
		"update-system-modify-virt-uefi",
		"update-system-modify-virt-pxe-boot-flag",
	})
	sys := NewSystem()
	sys.Name = "mytestsystem"
	sys.Hostname = "blahhost"
	sys.DNS.NameServers.Data = []string{"8.8.8.8", "8.8.4.4"}
	// Profile must be the referenced profile's real uid (centos7-x86_64), not its name.
	sys.Profile = "00000000000000000000000000000016"
	// UpdateSystem, unlike CreateSystem, does not apply default values for
	// these two fields, so they must be set explicitly to match the fixtures.
	sys.Power.Type = "ipmilanplus"
	sys.Status = "production"

	err := c.UpdateSystem(&sys)
	FailOnError(t, err)
}

func TestDeleteSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-system")
	err := c.DeleteSystem("00000000000000000000000000000061")
	FailOnError(t, err)
}

func TestDeleteSystemRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-system")
	err := c.DeleteSystemRecursive("00000000000000000000000000000061", false)
	FailOnError(t, err)
}

func TestListSystemNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-system")
	sytems, err := c.ListSystemNames()
	FailOnError(t, err)

	if len(sytems) != 6 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystemsSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-since")
	systems, err := c.GetSystemsSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(systems) != 6 {
		t.Errorf("Wrong number of profiles returned.")
	}
}

func TestFindSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-system")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	_, err := c.FindSystem(criteria, false)
	FailOnError(t, err)
}

func TestFindSystemNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-system-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	_, err := c.FindSystem(criteria, false)
	FailOnError(t, err)
}

func TestSaveSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-system")
	err := c.SaveSystem("0000000000000000000000000000001a", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopySystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-system")
	err := c.CopySystem("0000000000000000000000000000001a", "testsys2")
	FailOnError(t, err)
}

func TestRenameSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-system")
	err := c.RenameSystem("0000000000000000000000000000001b", "testsys1")
	FailOnError(t, err)
}

func TestGetSystemHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-handle")
	res, err := c.GetSystemHandle("testsys")
	FailOnError(t, err)

	if res != "0000000000000000000000000000001a" {
		t.Error("Wrong object id returned.")
	}
}

func TestGetValidSystemBootLoaders(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-system-boot-loaders")
	res, err := c.GetValidSystemBootLoaders("test")
	FailOnError(t, err)

	if len(res) < 1 {
		t.Error("Expected at least one boot loader.")
	}
}

func TestGetSystemAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-system-as-rendered")
	res, err := c.GetSystemAsRendered("test")
	FailOnError(t, err)

	if res["name"] != "test" {
		t.Errorf("Wrong system name returned: %v", res["name"])
	}
}

func TestGetInterfaces(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"get-interfaces-get-system",
	})
	testsys, err := c.GetSystem("0000000000000000000000000000001a", false, false)
	FailOnError(t, err)

	// Assert - in Cobbler 4.0.0 network interfaces are separate items and are
	// no longer embedded in the get_system response, so Interfaces should
	// come back as a non-nil, empty map.
	if testsys.Interfaces == nil {
		t.Fatal("Interfaces map should not be nil")
	}
	if len(testsys.Interfaces) != 0 {
		t.Fatalf("expected no embedded interfaces, got %d", len(testsys.Interfaces))
	}
}

func TestGetInterface(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"get-interfaces-get-system",
	})
	testsys, err := c.GetSystem("0000000000000000000000000000001a", false, false)
	FailOnError(t, err)

	// Assert - in Cobbler 4.0.0 network interfaces are separate items and are
	// no longer embedded in the get_system response, so a lookup by name
	// should come back empty.
	if _, ok := testsys.Interfaces["default"]; ok {
		t.Fatal("did not expect an embedded interface named \"default\"")
	}
}

func TestDisableNetboot(t *testing.T) {
	c := createStubHTTPClientSingle(t, "disable-netboot")
	err := c.DisableNetboot("testsys")
	FailOnError(t, err)
}

func TestUploadLogData(t *testing.T) {
	c := createStubHTTPClientSingle(t, "upload-log-data")
	res, err := c.UploadLogData("testsys", "/var/log/cobbler/testsys.log", 12, 0, "hello world!")
	FailOnError(t, err)
	if res {
		t.Error("Expected false result from UploadLogData.")
	}
}

func TestClearSystemLogs(t *testing.T) {
	c := createStubHTTPClientSingle(t, "clear-system-logs")
	res, err := c.ClearSystemLogs("0000000000000000000000000000001a")
	FailOnError(t, err)
	if !res {
		t.Error("Expected true result from ClearSystemLogs.")
	}
}
