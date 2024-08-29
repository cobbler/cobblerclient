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
	"github.com/go-test/deep"
	"testing"
)

func TestSetCachedVersion(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "extended-version")
	expectedVersion := CobblerVersion{
		Major: 3,
		Minor: 4,
		Patch: 0,
	}

	// Act
	err := c.setCachedVersion()

	// Assert
	FailOnError(t, err)
	deep.Equal(c.CachedVersion, expectedVersion)
}

func TestInvalidateCachedVersion(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "extended-version")
	_ = c.setCachedVersion()

	// Act
	c.invalidateCachedVersion()

	// Assert
	deep.Equal(c.CachedVersion, CobblerVersion{})
}

func TestGenerateAutoinstall(t *testing.T) {
	c := createStubHTTPClientSingle(t, "generate-autoinstall")

	res, err := c.GenerateAutoinstall("", "")
	FailOnError(t, err)
	if res == "" {
		t.Fatalf("Expected a non-empty string.")
	}
}

func TestLastModifiedTime(t *testing.T) {
	c := createStubHTTPClientSingle(t, "last-modified-time")

	res, err := c.LastModifiedTime()
	FailOnError(t, err)
	if res < 0.0 {
		t.Fatalf("Expected the float to be greater or equal to zero.")
	}
}

func TestPing(t *testing.T) {
	c := createStubHTTPClientSingle(t, "ping")

	res, err := c.Ping()
	FailOnError(t, err)
	if res == false {
		t.Fatalf("Expected ping to return true")
	}
}

func TestAutoAddRepos(t *testing.T) {
	c := createStubHTTPClientSingle(t, "auto-add-repos")

	err := c.AutoAddRepos()
	FailOnError(t, err)
}

func TestGetAutoinstallTemplates(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"get-autoinstall-templates",
	)

	err := c.GetAutoinstallTemplates()
	FailOnError(t, err)
}

func TestGetAutoinstallSnippets(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"get-autoinstall-snippets",
	)

	err := c.GetAutoinstallSnippets()
	FailOnError(t, err)
}

func TestIsAutoinstallInUse(t *testing.T) {
	c := createStubHTTPClientSingle(t, "is-autoinstall-in-use")

	err := c.IsAutoinstallInUse("")
	FailOnError(t, err)
}

func TestGenerateIPxe(t *testing.T) {
	c := createStubHTTPClientSingle(t, "generate-ipxe")

	err := c.GenerateIPxe("", "", "")
	FailOnError(t, err)
}

func TestGenerateBootCfg(t *testing.T) {
	c := createStubHTTPClientSingle(t, "generate-boot-cfg")

	err := c.GenerateBootCfg("testprof", "")
	FailOnError(t, err)
}

func TestGenerateScript(t *testing.T) {
	c := createStubHTTPClientSingle(t, "generate-script")

	err := c.GenerateScript("testprof", "", "preseed_early_default")
	FailOnError(t, err)
}

func TestGetBlendedData(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-blended-data")

	result, err := c.GetBlendedData("testprof", "")
	FailOnError(t, err)
	if len(result) != 184 {
		t.Fatalf("Expected a map with 184 entries, got %d.", len(result))
	}
}

func TestRegisterNewSystem(t *testing.T) {
	// Skip for now as the XML appears to have a different order.
	t.Skip("XML has different order. Needs to be fixed at a later point!")

	c := createStubHTTPClientSingle(t, "register-new-system")

	err := c.RegisterNewSystem(
		map[string]interface{}{
			"name":    "test",
			"profile": "testprof",
			"interfaces": map[string]interface{}{
				"default": map[string]interface{}{
					"mac_address": "AA:BB:CC:DD:EE:FF",
					"ip_address":  "192.168.1.1",
					"netmask":     "255.255.255.0",
				},
			}})
	FailOnError(t, err)
}

func TestRunInstallTriggers(t *testing.T) {
	c := createStubHTTPClientSingle(t, "run-install-triggers")

	err := c.RunInstallTriggers("", "", "", "")
	FailOnError(t, err)
}

func TestGetReposCompatibleWithProfile(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"get-repos-compatible-with-profile",
	)

	err := c.GetReposCompatibleWithProfile("testprof")
	FailOnError(t, err)
}

func TestFindSystemByDnsName(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"find-system-by-dns-name",
	)

	err := c.FindSystemByDnsName("testname")
	FailOnError(t, err)
}

func TestGetRandomMac(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-random-mac")

	err := c.GetRandomMac()
	FailOnError(t, err)
}

func TestXmlRpcHacks(t *testing.T) {
	c := createStubHTTPClientSingle(t, "xmlrpc-hacks")

	err := c.XmlRpcHacks(map[string]interface{}{"test": true})
	FailOnError(t, err)
}

func TestGetStatus(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-status")

	err := c.GetStatus("normal")
	FailOnError(t, err)
}

func TestSyncDhcp(t *testing.T) {
	c := createStubHTTPClientSingle(t, "sync-dhcp")

	err := c.SyncDhcp()
	FailOnError(t, err)
}

func TestGetConfigData(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-config-data")

	err := c.GetConfigData("testsys")
	FailOnError(t, err)
}

func TestClient_IsValueInherit(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "inherit-string", args: struct{ value interface{} }{value: "<<inherit>>"}, want: true},
		{name: "non-inherit-string", args: struct{ value interface{} }{value: "garbage"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := NewStubHTTPClient(t)
			c := NewClient(hc, config)
			c.Token = "securetoken99"
			if got := c.IsValueInherit(tt.args.value); got != tt.want {
				t.Errorf("IsValueInherit() = %v, want %v", got, tt.want)
			}
		})
	}
}
