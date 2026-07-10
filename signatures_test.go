package cobblerclient

import (
	"testing"
)

func TestGetSignatures(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-signatures")

	result, err := c.GetSignatures()
	FailOnError(t, err)
	if result.Breeds["redhat"]["rhel4"].VersionFile != `(redhat|sl|centos)-release-4(AS|WS|ES)[\.-]+(.*)\.rpm` {
		t.Fatalf("Expected a different regex!")
	}
}

func TestGetValidBreeds(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-breeds")

	res, err := c.GetValidBreeds()
	FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidOsVersionsForBreed(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"get-valid-os-verions-for-breed",
	)

	res, err := c.GetValidOsVersionsForBreed("redhat")
	FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidOsVersions(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-os-versions")

	res, err := c.GetValidOsVersions()
	FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestGetValidArchs(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-valid-archs")

	res, err := c.GetValidArchs()
	FailOnError(t, err)
	if len(res) <= 1 {
		t.Fatalf("Expected a length of greater then one.")
	}
}

func TestBackgroundSignatureUpdate(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"background-signature-update",
	)

	res, err := c.BackgroundSignatureUpdate()
	FailOnError(t, err)
	if res != "2000-01-01_000000_Updating Signatures_0000000000000000000000000000000c" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundSignatureReload(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-signature-reload")

	res, err := c.BackgroundSignatureReload()
	FailOnError(t, err)
	if res != "2000-01-01_000000_Reloading Signatures_0000000000000000000000000000000d" {
		t.Fatalf("Expected a different Event-ID!")
	}
}
