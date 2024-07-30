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

	"github.com/ContainerSolutions/go-utils"
)

var config = ClientConfig{
	URL:      "http://localhost:8081/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

func createStubHTTPClient(t *testing.T, reqFixture string, resFixture string) Client {
	hc := utils.NewStubHTTPClient(t)

	if reqFixture != "" {
		rawRequest, err := utils.Fixture(reqFixture)
		utils.FailOnError(t, err)

		// flatten the request so it matches the kolo generated xml
		r := regexp.MustCompile(`\s+<`)
		expectedReq := []byte(r.ReplaceAllString(string(rawRequest), "<"))
		hc.Expected = expectedReq
	}

	if resFixture != "" {
		response, err := utils.Fixture(resFixture)
		utils.FailOnError(t, err)
		hc.Response = response
	}

	c := NewClient(hc, config)
	c.Token = "securetoken99"
	return c
}

func TestGenerateAutoinstall(t *testing.T) {
	c := createStubHTTPClient(t, "generate-autoinstall-req.xml", "generate-autoinstall-res.xml")

	res, err := c.GenerateAutoinstall("", "")
	utils.FailOnError(t, err)
	if res == "" {
		t.Fatalf("Expected a non-empty string.")
	}
}

func TestLastModifiedTime(t *testing.T) {
	c := createStubHTTPClient(t, "last-modified-time-req.xml", "last-modified-time-res.xml")

	res, err := c.LastModifiedTime()
	utils.FailOnError(t, err)
	if res < 0.0 {
		t.Fatalf("Expected the float to be greater or equal to zero.")
	}
}

func TestPing(t *testing.T) {
	c := createStubHTTPClient(t, "ping-req.xml", "ping-res.xml")

	res, err := c.Ping()
	utils.FailOnError(t, err)
	if res == false {
		t.Fatalf("Expected ping to return true")
	}
}

func TestAutoAddRepos(t *testing.T) {
	c := createStubHTTPClient(t, "auto-add-repos-req.xml", "auto-add-repos-res.xml")

	err := c.AutoAddRepos()
	utils.FailOnError(t, err)
}

func TestGetAutoinstallTemplates(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"get-autoinstall-templates-req.xml",
		"get-autoinstall-templates-res.xml",
	)

	err := c.GetAutoinstallTemplates()
	utils.FailOnError(t, err)
}

func TestGetAutoinstallSnippets(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"get-autoinstall-snippets-req.xml",
		"get-autoinstall-snippets-res.xml",
	)

	err := c.GetAutoinstallSnippets()
	utils.FailOnError(t, err)
}

func TestIsAutoinstallInUse(t *testing.T) {
	c := createStubHTTPClient(t, "is-autoinstall-in-use-req.xml", "is-autoinstall-in-use-res.xml")

	err := c.IsAutoinstallInUse("")
	utils.FailOnError(t, err)
}

func TestGenerateIPxe(t *testing.T) {
	c := createStubHTTPClient(t, "generate-ipxe-req.xml", "generate-ipxe-res.xml")

	err := c.GenerateIPxe("", "", "")
	utils.FailOnError(t, err)
}

func TestGenerateBootCfg(t *testing.T) {
	c := createStubHTTPClient(t, "generate-boot-cfg-req.xml", "generate-boot-cfg-res.xml")

	err := c.GenerateBootCfg("testprof", "")
	utils.FailOnError(t, err)
}

func TestGenerateScript(t *testing.T) {
	c := createStubHTTPClient(t, "generate-script-req.xml", "generate-script-res.xml")

	err := c.GenerateScript("testprof", "", "preseed_early_default")
	utils.FailOnError(t, err)
}

func TestGetBlendedData(t *testing.T) {
	c := createStubHTTPClient(t, "get-blended-data-req.xml", "get-blended-data-res.xml")

	result, err := c.GetBlendedData("testprof", "")
	utils.FailOnError(t, err)
	if len(result) != 184 {
		t.Fatalf("Expected a map with 184 entries, got %d.", len(result))
	}
}

func TestRegisterNewSystem(t *testing.T) {
	// Skip for now as the XML appears to have a different order.
	t.Skip("XML has different order. Needs to be fixed at a later point!")

	c := createStubHTTPClient(t, "register-new-system-req.xml", "register-new-system-res.xml")

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
	utils.FailOnError(t, err)
}

func TestRunInstallTriggers(t *testing.T) {
	c := createStubHTTPClient(t, "run-install-triggers-req.xml", "run-install-triggers-res.xml")

	err := c.RunInstallTriggers("", "", "", "")
	utils.FailOnError(t, err)
}

func TestVersion(t *testing.T) {
	c := createStubHTTPClient(t, "version-req.xml", "version-res.xml")

	res, err := c.Version()
	utils.FailOnError(t, err)
	if res != 3.4 {
		t.Errorf("Wrong version returned.")
	}
}

func TestExtendedVersion(t *testing.T) {
	c := createStubHTTPClient(t, "extended-version-req.xml", "extended-version-res.xml")
	expectedResult := ExtendedVersion{
		Gitdate:      "Mon Jun 13 16:13:33 2022 +0200",
		Gitstamp:     "0e20f01b",
		Builddate:    "Mon Jun 27 06:34:23 2022",
		Version:      "3.4.0",
		VersionTuple: []int{3, 4, 0},
	}

	result, err := c.ExtendedVersion()
	utils.FailOnError(t, err)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Result from 'extended_version' did not match expected result.")
	}
}

func TestGetReposCompatibleWithProfile(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"get-repos-compatible-with-profile-req.xml",
		"get-repos-compatible-with-profile-res.xml",
	)

	err := c.GetReposCompatibleWithProfile("testprof")
	utils.FailOnError(t, err)
}

func TestFindSystemByDnsName(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"find-system-by-dns-name-req.xml",
		"find-system-by-dns-name-res.xml",
	)

	err := c.FindSystemByDnsName("testname")
	utils.FailOnError(t, err)
}

func TestGetRandomMac(t *testing.T) {
	c := createStubHTTPClient(t, "get-random-mac-req.xml", "get-random-mac-res.xml")

	err := c.GetRandomMac()
	utils.FailOnError(t, err)
}

func TestXmlRpcHacks(t *testing.T) {
	c := createStubHTTPClient(t, "xmlrpc-hacks-req.xml", "xmlrpc-hacks-res.xml")

	err := c.XmlRpcHacks(map[string]interface{}{"test": true})
	utils.FailOnError(t, err)
}

func TestGetStatus(t *testing.T) {
	c := createStubHTTPClient(t, "get-status-req.xml", "get-status-res.xml")

	err := c.GetStatus("normal")
	utils.FailOnError(t, err)
}

func TestSyncDhcp(t *testing.T) {
	c := createStubHTTPClient(t, "sync-dhcp-req.xml", "sync-dhcp-res.xml")

	err := c.SyncDhcp()
	utils.FailOnError(t, err)
}

func TestGetConfigData(t *testing.T) {
	c := createStubHTTPClient(t, "get-config-data-req.xml", "get-config-data-res.xml")

	err := c.GetConfigData("testsys")
	utils.FailOnError(t, err)
}
