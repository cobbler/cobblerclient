package cobblerclient

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"
)

func TestGetSignatures(t *testing.T) {
	c := createStubHTTPClient(t, "get-signatures-req.xml", "get-signatures-res.xml")

	result, err := c.GetSignatures()
	utils.FailOnError(t, err)
	if result.Breeds["redhat"]["rhel4"].VersionFile != `(redhat|sl|centos)-release-4(AS|WS|ES)[\.-]+(.*)\.rpm` {
		t.Fatalf("Expected a different regex!")
	}
}

func TestGetValidBreeds(t *testing.T) {
	c := createStubHTTPClient(t, "get-valid-breeds-req.xml", "get-valid-breeds-res.xml")

	res, err := c.GetValidBreeds()
	utils.FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidOsVersionsForBreed(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"get-valid-os-verions-for-breed-req.xml",
		"get-valid-os-verions-for-breed-res.xml",
	)

	res, err := c.GetValidOsVersionsForBreed("redhat")
	utils.FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidOsVersions(t *testing.T) {
	c := createStubHTTPClient(t, "get-valid-os-versions-req.xml", "get-valid-os-versions-res.xml")

	res, err := c.GetValidOsVersions()
	utils.FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidArchs(t *testing.T) {
	c := createStubHTTPClient(t, "get-valid-archs-req.xml", "get-valid-archs-res.xml")

	res, err := c.GetValidArchs()
	utils.FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestBackgroundSignatureUpdate(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"background-signature-update-req.xml",
		"background-signature-update-res.xml",
	)

	// FIXME: Test event-id return
	_, err := c.BackgroundSignatureUpdate()
	utils.FailOnError(t, err)
}
