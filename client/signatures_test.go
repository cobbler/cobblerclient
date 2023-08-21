package client_test

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"

	cobbler_testing "github.com/cobbler/cobblerclient/internal/testing"
)

func TestGetSignatures(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-signatures-req.xml", "get-signatures-res.xml")

	err := c.GetSignatures()
	utils.FailOnError(t, err)
}

func TestGetValidBreeds(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-valid-breeds-req.xml", "get-valid-breeds-res.xml")

	res, err := c.GetValidBreeds()
	utils.FailOnError(t, err)
	if len(res) <= 0 {
		t.Fatalf("Expected a lenght of greater then one.")
	}
}

func TestGetValidOsVersionsForBreed(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(
		t,
		"get-valid-os-verions-for-breed-req.xml",
		"get-valid-os-verions-for-breed-res.xml",
	)

	res, err := c.GetValidOsVersionsForBreed("redhat")
	utils.FailOnError(t, err)
	if len(res) <= 0 {
		t.Fatalf("Expected a lenght of greater then one.")
	}
}

func TestGetValidOsVersions(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-valid-os-versions-req.xml", "get-valid-os-versions-res.xml")

	res, err := c.GetValidOsVersions()
	utils.FailOnError(t, err)
	if len(res) <= 0 {
		t.Fatalf("Expected a lenght of greater then one.")
	}
}

func TestGetValidArchs(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-valid-archs-req.xml", "get-valid-archs-res.xml")

	res, err := c.GetValidArchs()
	utils.FailOnError(t, err)
	if len(res) <= 0 {
		t.Fatalf("Expected a lenght of greater then one.")
	}
}

func TestBackgroundSignatureUpdate(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(
		t,
		"background-signature-update-req.xml",
		"background-signature-update-res.xml",
	)

	// FIXME: Test event-id return
	_, err := c.BackgroundSignatureUpdate()
	utils.FailOnError(t, err)
}
