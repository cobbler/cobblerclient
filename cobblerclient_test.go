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
	"reflect"
	"regexp"
	"testing"
)

var config = ClientConfig{
	URL:      "http://localhost:8081/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

// createStubHTTPClient ...
func createStubHTTPClient(t *testing.T, fixtures []string) Client {
	hc := NewStubHTTPClient(t)

	for _, fixture := range fixtures {
		if fixture != "" {
			rawRequest, err := Fixture(fixture + "-req.xml")
			FailOnError(t, err)
			response, err := Fixture(fixture + "-res.xml")
			FailOnError(t, err)

			// flatten the request so it matches the kolo generated xml
			r := regexp.MustCompile(`\s+<`)
			expectedReq := []byte(r.ReplaceAllString(string(rawRequest), "<"))
			hc.answers = append(hc.answers, APIResponsePair{
				Expected: expectedReq,
				Response: response,
			})
		}
	}

	c := NewClient(hc, config)
	c.Token = "securetoken99"
	return c
}

// createStubHTTPClientSingle ...
func createStubHTTPClientSingle(t *testing.T, fixture string) Client {
	return createStubHTTPClient(t, []string{fixture})
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

func TestVersion(t *testing.T) {
	c := createStubHTTPClientSingle(t, "version")

	res, err := c.Version()
	FailOnError(t, err)
	if res != 3.4 {
		t.Errorf("Wrong version returned.")
	}
}

func TestExtendedVersion(t *testing.T) {
	c := createStubHTTPClientSingle(t, "extended-version")
	expectedResult := ExtendedVersion{
		Gitdate:      "Mon Jun 13 16:13:33 2022 +0200",
		Gitstamp:     "0e20f01b",
		Builddate:    "Mon Jun 27 06:34:23 2022",
		Version:      "3.4.0",
		VersionTuple: []int{3, 4, 0},
	}

	result, err := c.ExtendedVersion()
	FailOnError(t, err)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Result from 'extended_version' did not match expected result.")
	}
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
