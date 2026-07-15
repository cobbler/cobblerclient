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

func TestNewNetworkInterface(t *testing.T) {
	ni := NewNetworkInterface()
	if ni.InterfaceType != NetworkInterfaceTypeNA {
		t.Errorf("default InterfaceType = %v, want NA", ni.InterfaceType)
	}
	if ni.IPv4.StaticRoutes == nil {
		t.Errorf("IPv4.StaticRoutes should be non-nil empty slice")
	}
	if ni.IPv6.StaticRoutes == nil || ni.IPv6.Secondaries == nil {
		t.Errorf("IPv6 slices should be non-nil empty slices")
	}
	if ni.DNS.CNames == nil {
		t.Errorf("DNS.CNames should be non-nil empty slice")
	}
	if !ni.VirtBridge.IsInherited {
		t.Errorf("VirtBridge should default to inherited")
	}
}

func TestGetNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interface")
	ni, err := c.GetNetworkInterface("eth0-server1", false, false)
	FailOnError(t, err)

	if ni.Name != "eth0-server1" {
		t.Errorf("wrong name: %q", ni.Name)
	}
	if ni.MacAddress != "" {
		t.Errorf("wrong MAC: %q", ni.MacAddress)
	}
	if ni.IPv4.Address != "" {
		t.Errorf("wrong IPv4 address: %q", ni.IPv4.Address)
	}
	if ni.DNS.Name != "" {
		t.Errorf("wrong DNS name: %q", ni.DNS.Name)
	}
}

func TestGetNetworkInterfaces(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interfaces")
	ifaces, err := c.GetNetworkInterfaces()
	FailOnError(t, err)

	if len(ifaces) != 4 {
		t.Errorf("expected 4 ifaces, got %d", len(ifaces))
	}
}

func TestGetNetworkInterfaceHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interface-handle")
	handle, err := c.GetNetworkInterfaceHandle("eth0-server1")
	FailOnError(t, err)

	if handle != "000000000000000000000000000000d6" {
		t.Errorf("wrong handle: %q", handle)
	}
}

func TestListNetworkInterfaceNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-network-interface")
	names, err := c.ListNetworkInterfaceNames()
	FailOnError(t, err)

	if len(names) != 4 {
		t.Errorf("expected 4 names, got %d", len(names))
	}
}

func TestGetNetworkInterfacesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interfaces-since")
	ifaces, err := c.GetNetworkInterfacesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(ifaces) != 4 {
		t.Errorf("expected 4 ifaces, got %d", len(ifaces))
	}
}

func TestFindNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-network-interface")
	criteria := make(map[string]interface{}, 1)
	criteria["mac_address"] = "00:11:22:33:44:55"
	ifaces, err := c.FindNetworkInterface(criteria, false)
	FailOnError(t, err)

	// No seeded interface actually carries this MAC, so the search legitimately matches nothing.
	if len(ifaces) != 0 {
		t.Errorf("expected 0 ifaces, got %d", len(ifaces))
	}
}

func TestFindNetworkInterfaceNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-network-interface-names")
	criteria := make(map[string]interface{}, 1)
	criteria["mac_address"] = "00:11:22:33:44:55"
	names, err := c.FindNetworkInterfaceNames(criteria)
	FailOnError(t, err)

	// No seeded interface actually carries this MAC, so the search legitimately matches nothing.
	if len(names) != 0 {
		t.Errorf("expected 0 names, got %d", len(names))
	}
}

func TestCreateNetworkInterface(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"create-network-interface-name-check",
		"new-network-interface",
		"new-network-interface-modify-name",
		"new-network-interface-modify-comment",
		"new-network-interface-modify-kernel-options",
		"new-network-interface-modify-kernel-options-post",
		"new-network-interface-modify-autoinstall-meta",
		"new-network-interface-modify-template-files",
		"new-network-interface-modify-owners",
		"new-network-interface-modify-mac-address",
		"new-network-interface-modify-interface-type",
		"new-network-interface-modify-interface-master",
		"new-network-interface-modify-bonding-opts",
		"new-network-interface-modify-bridge-opts",
		"new-network-interface-modify-connected-mode",
		"new-network-interface-modify-management",
		"new-network-interface-modify-static",
		"new-network-interface-modify-dhcp-tag",
		"new-network-interface-modify-if-gateway",
		"new-network-interface-modify-mtu",
		"new-network-interface-modify-ipv6-default-gateway",
		"new-network-interface-modify-ipv6-static-routes",
		"new-network-interface-modify-virt-bridge",
		"new-network-interface-modify-ipv4-address",
		"new-network-interface-modify-ipv4-mtu",
		"new-network-interface-modify-ipv4-netmask",
		"new-network-interface-modify-ipv4-static-routes",
		"new-network-interface-modify-ipv6-address",
		"new-network-interface-modify-ipv6-prefix",
		"new-network-interface-modify-ipv6-secondaries",
		"new-network-interface-modify-ipv6-mtu",
		"new-network-interface-modify-ipv6-option-static-routes",
		"new-network-interface-modify-dns-name",
		"new-network-interface-modify-dns-cnames",
		"new-network-interface-save",
		"new-network-interface-get",
	})
	ni := NewNetworkInterface()
	ni.Name = "eth0-mytestsystem"
	ni.MacAddress = "AA:BB:CC:DD:EE:00"

	result, err := c.CreateNetworkInterface("00000000000000000000000000000019", ni)
	FailOnError(t, err)
	if result.Name != "eth0-mytestsystem" {
		t.Errorf("Wrong network interface name returned: %v", result.Name)
	}
}

func TestSaveNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-network-interface")
	err := c.SaveNetworkInterface("000000000000000000000000000000d6", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-network-interface")
	err := c.CopyNetworkInterface("000000000000000000000000000000d6", "eth1-server1")
	FailOnError(t, err)
}

func TestDeleteNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-network-interface")
	err := c.DeleteNetworkInterface("eth0-server1")
	FailOnError(t, err)
}

func TestDeleteNetworkInterfaceRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-network-interface")
	err := c.DeleteNetworkInterfaceRecursive("eth0-server1", false)
	FailOnError(t, err)
}

func TestRenameNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-network-interface")
	err := c.RenameNetworkInterface("000000000000000000000000000000d7", "eth2-server1")
	FailOnError(t, err)
}
