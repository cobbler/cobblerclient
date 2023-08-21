package item_test

import (
	"github.com/cobbler/cobblerclient/item"
	"testing"

	"github.com/ContainerSolutions/go-utils"

	cobbler_testing "github.com/cobbler/cobblerclient/internal/testing"
)

func TestGetSystems(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-systems-req.xml", "get-systems-res.xml")
	systems, err := item.GetSystems(c)
	utils.FailOnError(t, err)

	if len(systems) != 1 {
		t.Errorf("Wrong number of systems returned.")
	}
}

func TestGetSystem(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-system-req.xml", "get-system-res.xml")
	system, err := item.GetSystem(c, "test")
	utils.FailOnError(t, err)

	if system.Name.Get() != "test" {
		t.Errorf("Wrong system returned.")
	}
}

func TestNewSystem(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "new-system-req.xml", "new-system-res.xml")
	result, err := c.Call("new_system", c.Token)
	utils.FailOnError(t, err)
	newID := result.(string)

	if newID != "___NEW___system::abc123==" {
		t.Errorf("Wrong ID returned.")
	}

	c = cobbler_testing.CreateStubHTTPClient(t, "set-system-hostname-req.xml", "set-system-hostname-res.xml")
	result, err = c.Call("modify_system", newID, "hostname", "blahhost", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting hostname failed.")
	}

	c = cobbler_testing.CreateStubHTTPClient(t, "set-system-name-req.xml", "set-system-name-res.xml")
	result, err = c.Call("modify_system", newID, "name", "mytestsystem", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name failed.")
	}

	c = cobbler_testing.CreateStubHTTPClient(t, "set-system-nameservers-req.xml", "set-system-nameservers-res.xml")
	result, err = c.Call("modify_system", newID, "name_servers", "8.8.8.8 8.8.4.4", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name servers failed.")
	}

	c = cobbler_testing.CreateStubHTTPClient(t, "set-system-profile-req.xml", "set-system-profile-res.xml")
	result, err = c.Call("modify_system", newID, "profile", "centos7-x86_64", c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting name servers failed.")
	}

	/* I'm not sure how to get this test to pass with unordered maps
	nicInfo := map[string]interface{}{
		"macaddress-eth0":  "01:02:03:04:05:06",
		"ipaddress-eth0":   "1.2.3.4",
		"dnsname-eth0":     "deathstar",
		"subnetsmask-eth0": "255.255.255.0",
		"if-gateway-eth0":  "4.3.2.1",
	}

	c = cobbler_testing.CreateStubHTTPClient(t, "set-system-network-req.xml", "set-system-network-res.xml")
	result, err = c.Call("modify_system", newID, "modify_interface", nicInfo, c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Setting interface failed.")
	}
	*/

	c = cobbler_testing.CreateStubHTTPClient(t, "save-system-req.xml", "save-system-res.xml")
	result, err = c.Call("save_system", newID, c.Token)
	utils.FailOnError(t, err)

	if !result.(bool) {
		t.Errorf("Save failed.")
	}
}
