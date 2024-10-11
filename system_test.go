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
	"fmt"
	"github.com/go-test/deep"
	"testing"
	"time"
)

func TestNewSystem(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewSystem()
}

func TestNewInterface(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewInterface()
}

func TestGetSystems(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-systems")
	systems, err := c.GetSystems()
	FailOnError(t, err)

	if len(systems) != 1 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystem(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-system")
	c.CachedVersion = CobblerVersion{3, 3, 2}

	// Act
	system, err := c.GetSystem("test", false, false)

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
		"set-system-name",
		"new-system-modify-comment",
		"new-system-modify-kernel-options",
		"new-system-modify-kernel-options-post",
		"new-system-modify-autoinstall-meta",
		"new-system-modify-fetchable-files",
		"new-system-modify-boot-files",
		"new-system-modify-template-files",
		"new-system-modify-owners",
		"new-system-modify-mgmt-classes",
		"new-system-modify-mgmt-parameters",
		"set-system-profile",
		"new-system-modify-image",
		"new-system-modify-autoinstall",
		"new-system-modify-boot-loaders",
		"new-system-modify-enable-ipxe",
		"new-system-modify-filename",
		"new-system-modify-gateway",
		"set-system-hostname",
		"new-system-modify-ipv6-default-device",
		"set-system-nameservers",
		"new-system-modify-name-servers-search",
		"new-system-modify-netboot-enabled",
		"new-system-modify-next-server-v4",
		"new-system-modify-next-server-v6",
		"new-system-modify-power-address",
		"new-system-modify-power-id",
		"new-system-modify-power-identity-file",
		"new-system-modify-power-options",
		"new-system-modify-power-pass",
		"new-system-modify-power-type",
		"new-system-modify-power-user",
		"new-system-modify-proxy",
		"new-system-modify-redhat-management-key",
		"new-system-modify-serial-baud-rate",
		"new-system-modify-serial-device",
		"new-system-modify-server",
		"new-system-modify-status",
		"new-system-modify-virt-auto-boot",
		"new-system-modify-virt-cpus",
		"new-system-modify-virt-disk-driver",
		"new-system-modify-virt-file-size",
		"new-system-modify-virt-pxe-boot",
		"new-system-modify-virt-path",
		"new-system-modify-virt-ram",
		"new-system-modify-virt-type",
		"new-system-save",
		"new-system-get",
	})
	c.CachedVersion = CobblerVersion{3, 3, 2}
	sys := NewSystem()
	sys.Name = "mytestsystem"
	sys.Hostname = "blahhost"
	sys.NameServers = []string{"8.8.8.8", "8.8.4.4"}
	sys.Profile = "centos7-x86_64"

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
	if len(newSys.NameServers) != 2 || newSys.NameServers[0] != "8.8.8.8" {
		t.Errorf("Wrong system name servers returned.")
	}
	if newSys.Profile != "centos7-x86_64" {
		t.Errorf("Wrong system profile returned.")
	}
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

func TestCreateInterface(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"extended-version",
		"get-interfaces-get-system",
		"delete-interface-get-system-handle",
		"create-interface-create-interface",
		"delete-interface-save-system",
	})
	testsys, err := c.GetSystem("testsys", false, false)
	FailOnError(t, err)
	testinterface := NewInterface()

	// Act
	err = testsys.CreateInterface("eth0", testinterface)

	// Assert
	FailOnError(t, err)
}

func TestModifyInterface(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "modify-interface")
	testnic := NewInterface()
	testnic.IPAddress = "10.168.0.5"
	testnicmap := makeInterfaceOptionsMap("eth0", testnic)

	// Act
	err := c.ModifyInterface("system::testsys", testnicmap)

	// Assert
	FailOnError(t, err)
}

func TestGetInterfaces(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"extended-version",
		"get-interfaces-get-system",
	})
	testsys, err := c.GetSystem("testsys", false, false)
	FailOnError(t, err)

	// Act
	interfaces, err := testsys.GetInterfaces()

	// Assert
	if len(interfaces) < 1 {
		t.Fatal("there should be at least one interface")
	}
	if interfaces["default"].IPAddress != "10.1.0.1" {
		t.Fatal("incorrect IP address")
	}
	FailOnError(t, err)
}

func TestGetInterface(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"extended-version",
		"get-interfaces-get-system",
	})
	testsys, err := c.GetSystem("testsys", false, false)
	FailOnError(t, err)

	// Act
	nic, err := testsys.GetInterface("default")

	// Assert
	expectedInterface := NewInterface()
	expectedInterface.IPAddress = "10.1.0.1"
	differences := deep.Equal(nic, expectedInterface)
	if len(differences) > 0 {
		fmt.Println(differences)
		t.Fatal("interfaces non-equal")
	}
	FailOnError(t, err)
}

func TestDeleteInterface(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"extended-version",
		"get-interfaces-get-system",
		"delete-interface-get-system-handle",
		"delete-interface-delete-interface",
		"delete-interface-save-system",
	})
	testsys, err := c.GetSystem("testsys", false, false)
	FailOnError(t, err)

	// Act
	err = testsys.DeleteInterface("default")

	// Assert
	FailOnError(t, err)
}

func TestRenameInterface(t *testing.T) {
	// Arrange
	c := createStubHTTPClient(t, []string{
		"extended-version",
		"get-interfaces-get-system",
		"delete-interface-get-system-handle",
		"rename-interface-rename-interface",
	})
	testsys, err := c.GetSystem("testsys", false, false)
	FailOnError(t, err)

	// Act
	err = testsys.RenameInterface("", "")

	// Assert
	FailOnError(t, err)
}
