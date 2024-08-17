package cobblerclient

import (
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
	c := createStubHTTPClientSingle(t, "login")
	ok, err := c.Login()
	FailOnError(t, err)

	if !ok {
		t.Errorf("true expected; got false")
	}

	expected := "sa/1EWr40BWU+Pq3VEOOpD4cQtxkeMuFUw=="
	if c.Token != expected {
		t.Errorf(`"%s" expected; got "%s"`, expected, c.Token)
	}
}

func TestLoginWithError(t *testing.T) {
	c := createStubHTTPClientSingle(t, "login-err")
	expected := `Fault(1): <class 'cobbler.cexceptions.CX'>:'login failed (cobbler)'`

	ok, err := c.Login()
	if ok {
		t.Errorf("false expected; got true")
	}

	if err == nil || err.Error() != expected {
		t.Errorf("%s expected; got %s", expected, err)
	}
}

func TestLogout(t *testing.T) {
	c := createStubHTTPClientSingle(t, "logout")
	var expected = false

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
	if res == expected {
		t.Errorf(`"%t" expected; got "%t"`, expected, res)
	}
}

func TestGetUserFromToken(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-user-from-token")
	var expected = "testuser"

	res, err := c.GetUserFromToken("securetoken99")
	FailOnError(t, err)
	if res != expected {
		t.Errorf(`"%s" expected; got "%s"`, expected, res)
	}
}
