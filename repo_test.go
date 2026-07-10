/*
Copyright 2017 HomeAway, Inc

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
	"testing"
	"time"
)

func TestNewRepo(t *testing.T) {
	// Arrange, Act & Assert
	_ = NewRepo()
}

func TestGetRepos(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repos")
	repos, err := c.GetRepos()
	FailOnError(t, err)
	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetRepo(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-repo")

	// Act
	repo, err := c.GetRepo("rhel-7-x86_64", false, false)

	// Assert
	FailOnError(t, err)
	if repo.Name != "rhel-7-x86_64" {
		t.Errorf("Wrong repo returned.")
	}
	if !repo.CreateRepoFlags.IsInherited {
		t.Errorf("Expected CreateRepoFlags.IsInherited to be true, got false")
	}
}

// TestGetRepoExplicitStringValue verifies that a non-inherited Value[string] field (e.g.
// CreateRepoFlags, Proxy) is decoded with Data set to the actual string — not left as the
// zero value "". This is a regression test for the cobblerDataHacks bug where the string
// case set RawData but omitted Data for non-inherited values.
func TestGetRepoExplicitStringValue(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-repo-explicit-string")

	// Act
	repo, err := c.GetRepo("rhel-7-x86_64", false, false)

	// Assert
	FailOnError(t, err)
	if repo.CreateRepoFlags.IsInherited {
		t.Errorf("Expected CreateRepoFlags.IsInherited to be false")
	}
	if repo.CreateRepoFlags.Data != "--no-database" {
		t.Errorf("Expected CreateRepoFlags.Data to be \"--no-database\", got %q", repo.CreateRepoFlags.Data)
	}
	if repo.Proxy.IsInherited {
		t.Errorf("Expected Proxy.IsInherited to be false")
	}
	if repo.Proxy.Data != "http://proxy.example.com:8080" {
		t.Errorf("Expected Proxy.Data to be \"http://proxy.example.com:8080\", got %q", repo.Proxy.Data)
	}
}

func TestDeleteRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-repo")
	err := c.DeleteRepo("test")
	FailOnError(t, err)
}

func TestDeleteRepoRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-repo")
	err := c.DeleteRepoRecursive("test", false)
	FailOnError(t, err)
}

func TestListRepoNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-repo")
	repos, err := c.ListRepoNames()
	FailOnError(t, err)

	if len(repos) != 0 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetReposSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-since")
	repos, err := c.GetReposSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestFindRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-repo")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	repos, err := c.FindRepo(criteria)
	FailOnError(t, err)

	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestFindRepoNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-repo-names")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	repos, err := c.FindRepoNames(criteria)
	FailOnError(t, err)

	if len(repos) != 1 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestSaveRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-repo")
	err := c.SaveRepo("repo::testrepo", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-repo")
	err := c.CopyRepo("repo::testrepo", "testrepo2")
	FailOnError(t, err)
}

func TestRenameRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-repo")
	err := c.RenameRepo("repo::testrepo2", "testrepo1")
	FailOnError(t, err)
}

func TestGetRepoHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-handle")
	res, err := c.GetRepoHandle("testrepo")
	FailOnError(t, err)

	if res != "repo::testrepo" {
		t.Error("Wrong object id returned.")
	}
}

func TestGetRepoConfigForProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-config-for-profile")
	res, err := c.GetRepoConfigForProfile("testprof")
	FailOnError(t, err)

	if res == "" {
		t.Error("Expected a non-empty repo config.")
	}
}

func TestGetRepoConfigForSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-config-for-system")
	res, err := c.GetRepoConfigForSystem("testsys")
	FailOnError(t, err)

	if res == "" {
		t.Error("Expected a non-empty repo config.")
	}
}

func TestGetRepoAsRendered(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-as-rendered")
	res, err := c.GetRepoAsRendered("rhel-7-x86_64")
	FailOnError(t, err)

	if res["name"] != "rhel-7-x86_64" {
		t.Errorf("Wrong repo name returned: %v", res["name"])
	}
}

/*
 * NOTE: We're skipping the testing of CREATE, UPDATE, DELETE methods for now because
 *       the current implementation of the StubHTTPClient does not allow
 *       buffered mock responses so as soon as the method makes the second
 *       call to Cobbler it'll fail.
 *       This is a system test, so perhaps we can run Cobbler in a Docker container
 *       and take it from there.
 */
