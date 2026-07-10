package cobblerclient

import (
	"errors"
	"testing"
)

func TestCheckAccessNoFail(t *testing.T) {
	c := createStubHTTPClientSingle(t, "check-access-no-fail")

	res, err := c.CheckAccessNoFail("", "", "")
	FailOnError(t, err)
	if res != false {
		t.Errorf(`"%t" expected; got "%t"`, false, res)
	}
}

func TestCheckAccess(t *testing.T) {
	c := createStubHTTPClientSingle(t, "check-access")

	res, err := c.CheckAccess("", "", "")
	FailOnError(t, err)
	if res < 0 || res > 1 {
		t.Errorf(`"0" or "1" expected; got "%d"`, res)
	}
}

func TestGetAuthnModuleName(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-authn-module-name")
	var expected = "authentication.configfile"

	res, err := c.GetAuthnModuleName()
	FailOnError(t, err)
	if res != expected {
		t.Errorf(`"%s" expected; got "%s"`, expected, res)
	}
}

func TestLogin(t *testing.T) {
	c := createStubHTTPClient(t, []string{"login", "extended-version"})
	c.CachedVersion = CobblerVersion{} // allow Login to fetch and validate the server version
	ok, err := c.Login()
	FailOnError(t, err)

	if !ok {
		t.Errorf("true expected; got false")
	}

	if c.Token != "securetoken99" {
		t.Errorf(`"securetoken99" expected; got "%s"`, c.Token)
	}
}

func TestLoginWithError(t *testing.T) {
	c := createStubHTTPClientSingle(t, "login-err")
	c.config.Username = "wrong"
	c.config.Password = "wrong"
	expected := `Fault(1): <class 'ValueError'>:login failed (wrong)`

	ok, err := c.Login()
	if ok {
		t.Errorf("false expected; got true")
	}

	if err == nil || err.Error() != expected {
		t.Errorf("%s expected; got %s", expected, err)
	}
}

func TestLoginWithOldServerVersion(t *testing.T) {
	c := createStubHTTPClient(t, []string{"login", "extended-version-old"})
	c.CachedVersion = CobblerVersion{} // allow Login to fetch and validate the server version

	ok, err := c.Login()
	if ok {
		t.Errorf("false expected; got true")
	}
	if err == nil {
		t.Fatal("expected an error for server version < 4, got nil")
	}
	var versionErr *UnsupportedServerVersionError
	if !errors.As(err, &versionErr) {
		t.Errorf("expected UnsupportedServerVersionError, got %T: %v", err, err)
	}
}

func TestLogout(t *testing.T) {
	c := createStubHTTPClientSingle(t, "logout")
	var expected = true

	res, err := c.Logout()
	FailOnError(t, err)
	if res != expected {
		t.Errorf(`"%t" expected; got "%t"`, expected, res)
	}
}

func TestTokenCheck(t *testing.T) {
	c := createStubHTTPClientSingle(t, "token-check")
	var expected = false

	res, err := c.TokenCheck("my_fake_token")
	FailOnError(t, err)
	if res != expected {
		t.Errorf(`"%t" expected; got "%t"`, expected, res)
	}
}

func TestGetUserFromToken(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-user-from-token")
	var expected = "cobbler"

	res, err := c.GetUserFromToken("securetoken99")
	FailOnError(t, err)
	if res != expected {
		t.Errorf(`"%s" expected; got "%s"`, expected, res)
	}
}
