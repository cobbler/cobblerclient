package cobblerclient

import (
	"testing"
)

func TestGetSettings(t *testing.T) {
	c := createStubHTTPClient(t, "get-settings-req.xml", "get-settings-res.xml")

	result, err := c.GetSettings()
	FailOnError(t, err)
	if result.AuthTokenExpiration != 3600 {
		t.Errorf("Expected AuthTokenExpiration to be 3600, instead got %d", result.AuthTokenExpiration)
	}
}

func TestModifySettings(t *testing.T) {
	c := createStubHTTPClient(t, "modify-settings-req.xml", "modify-settings-res.xml")

	result, err := c.ModifySetting("auth_token_expiration", 7200)
	FailOnError(t, err)
	if result != 1 {
		t.Fatalf("Expected 1 but got %d", result)
	}
}
