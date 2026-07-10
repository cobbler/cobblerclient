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
	ni, err := c.GetNetworkInterface("eth0@server1", false, false)
	FailOnError(t, err)

	if ni.Name != "eth0@server1" {
		t.Errorf("wrong name: %q", ni.Name)
	}
	if ni.MacAddress != "00:11:22:33:44:55" {
		t.Errorf("wrong MAC: %q", ni.MacAddress)
	}
	if ni.IPv4.Address != "192.168.1.10" {
		t.Errorf("wrong IPv4 address: %q", ni.IPv4.Address)
	}
	if ni.DNS.Name != "server1.example.com" {
		t.Errorf("wrong DNS name: %q", ni.DNS.Name)
	}
}

func TestGetNetworkInterfaces(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interfaces")
	ifaces, err := c.GetNetworkInterfaces()
	FailOnError(t, err)

	if len(ifaces) != 1 {
		t.Errorf("expected 1 iface, got %d", len(ifaces))
	}
}

func TestGetNetworkInterfaceHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interface-handle")
	handle, err := c.GetNetworkInterfaceHandle("eth0@server1")
	FailOnError(t, err)

	if handle != "network_interface::eth0@server1" {
		t.Errorf("wrong handle: %q", handle)
	}
}

func TestListNetworkInterfaceNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-network-interface")
	names, err := c.ListNetworkInterfaceNames()
	FailOnError(t, err)

	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestGetNetworkInterfacesSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-network-interfaces-since")
	ifaces, err := c.GetNetworkInterfacesSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(ifaces) != 1 {
		t.Errorf("expected 1 iface, got %d", len(ifaces))
	}
}

func TestFindNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-network-interface")
	criteria := make(map[string]interface{}, 1)
	criteria["mac_address"] = "00:11:22:33:44:55"
	ifaces, err := c.FindNetworkInterface(criteria)
	FailOnError(t, err)

	if len(ifaces) != 1 {
		t.Errorf("expected 1 iface, got %d", len(ifaces))
	}
	if ifaces[0].MacAddress != "00:11:22:33:44:55" {
		t.Errorf("wrong MAC: %q", ifaces[0].MacAddress)
	}
}

func TestFindNetworkInterfaceNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-network-interface-names")
	criteria := make(map[string]interface{}, 1)
	criteria["mac_address"] = "00:11:22:33:44:55"
	names, err := c.FindNetworkInterfaceNames(criteria)
	FailOnError(t, err)

	if len(names) != 1 {
		t.Errorf("expected 1 name, got %d", len(names))
	}
	if names[0] != "eth0@server1" {
		t.Errorf("wrong name: %q", names[0])
	}
}

func TestSaveNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-network-interface")
	err := c.SaveNetworkInterface("network_interface::eth0@server1", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-network-interface")
	err := c.CopyNetworkInterface("network_interface::eth0@server1", "eth1@server1")
	FailOnError(t, err)
}

func TestDeleteNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-network-interface")
	err := c.DeleteNetworkInterface("eth0@server1")
	FailOnError(t, err)
}

func TestDeleteNetworkInterfaceRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-network-interface")
	err := c.DeleteNetworkInterfaceRecursive("eth0@server1", false)
	FailOnError(t, err)
}

func TestRenameNetworkInterface(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-network-interface")
	err := c.RenameNetworkInterface("network_interface::eth0@server1", "eth2@server1")
	FailOnError(t, err)
}
