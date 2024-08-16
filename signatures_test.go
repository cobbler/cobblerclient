package cobblerclient

import (
	"testing"
)

func TestGetSignatures(t *testing.T) {
	c := createStubHTTPClient(t, "get-signatures-req.xml", "get-signatures-res.xml")

	result, err := c.GetSignatures()
	FailOnError(t, err)
	if result.Breeds["redhat"]["rhel4"].VersionFile != `(redhat|sl|centos)-release-4(AS|WS|ES)[\.-]+(.*)\.rpm` {
		t.Fatalf("Expected a different regex!")
	}
}

func TestGetValidBreeds(t *testing.T) {
	c := createStubHTTPClient(t, "get-valid-breeds-req.xml", "get-valid-breeds-res.xml")

	res, err := c.GetValidBreeds()
	FailOnError(t, err)
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
	FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidOsVersions(t *testing.T) {
	c := createStubHTTPClient(t, "get-valid-os-versions-req.xml", "get-valid-os-versions-res.xml")

	res, err := c.GetValidOsVersions()
	FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidArchs(t *testing.T) {
	c := createStubHTTPClient(t, "get-valid-archs-req.xml", "get-valid-archs-res.xml")

	res, err := c.GetValidArchs()
	FailOnError(t, err)
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

	res, err := c.BackgroundSignatureUpdate()
	FailOnError(t, err)
	if res != "2022-09-30_195846_Updating Signatures_6c5300d51c224984b4319fb536cc21d0" {
		t.Fatalf("Expected a different Event-ID!")
	}
}
