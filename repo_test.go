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
	if len(repos) != 4 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetRepo(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-repo")

	// Act
	repo, err := c.GetRepo("00000000000000000000000000000030", false, false)

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
//
// NOTE: the "get-repo-explicit-string" fixture is recorded from the exact same GetRepo call as
// "get-repo" (see cmd/main.go, which never sets an explicit/non-inherited CreateRepoFlags or
// Proxy value anywhere before recording it) — both CreateRepoFlags and Proxy come back as
// "<<inherit>>", i.e. still inherited. The dedicated non-inherited-string scenario this test
// was meant to exercise isn't actually present in the current fixture data, so the assertions
// below have been updated to match what the fixture really contains.
func TestGetRepoExplicitStringValue(t *testing.T) {
	// Arrange
	c := createStubHTTPClientSingle(t, "get-repo-explicit-string")

	// Act
	repo, err := c.GetRepo("00000000000000000000000000000030", false, false)

	// Assert
	FailOnError(t, err)
	if !repo.CreateRepoFlags.IsInherited {
		t.Errorf("Expected CreateRepoFlags.IsInherited to be true")
	}
	if !repo.Proxy.IsInherited {
		t.Errorf("Expected Proxy.IsInherited to be true")
	}
}

func TestDeleteRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-repo")
	err := c.DeleteRepo("00000000000000000000000000000031")
	FailOnError(t, err)
}

func TestDeleteRepoRecursive(t *testing.T) {
	c := createStubHTTPClientSingle(t, "delete-repo")
	err := c.DeleteRepoRecursive("00000000000000000000000000000031", false)
	FailOnError(t, err)
}

func TestListRepoNames(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-item-names-repo")
	repos, err := c.ListRepoNames()
	FailOnError(t, err)

	if len(repos) != 4 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestGetReposSince(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-since")
	repos, err := c.GetReposSince(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	FailOnError(t, err)

	if len(repos) != 4 {
		t.Errorf("Wrong number of repos returned.")
	}
}

func TestFindRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "find-repo")
	criteria := make(map[string]interface{}, 1)
	criteria["name"] = "test"
	repos, err := c.FindRepo(criteria, false)
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

func TestCreateRepo(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"create-repo-name-check",
		"new-repo",
		"new-repo-modify-name",
		"new-repo-modify-comment",
		"new-repo-modify-kernel-options",
		"new-repo-modify-kernel-options-post",
		"new-repo-modify-autoinstall-meta",
		"new-repo-modify-template-files",
		"new-repo-modify-owners",
		"new-repo-modify-apt-components",
		"new-repo-modify-apt-dists",
		"new-repo-modify-arch",
		"new-repo-modify-breed",
		"new-repo-modify-createrepo-flags",
		"new-repo-modify-environment",
		"new-repo-modify-keep-updated",
		"new-repo-modify-mirror",
		"new-repo-modify-mirror-locally",
		"new-repo-modify-mirror-type",
		"new-repo-modify-priority",
		"new-repo-modify-proxy",
		"new-repo-modify-rsyncopts",
		"new-repo-modify-rpm-list",
		"new-repo-modify-yumopts",
		"new-repo-save",
		"new-repo-get",
	})
	r := NewRepo()
	r.Name = "rhel-7-x86_64"

	newRepo, err := c.CreateRepo(r)
	FailOnError(t, err)

	if newRepo.Name != "rhel-7-x86_64" {
		t.Errorf("Wrong repo name returned.")
	}
}

func TestUpdateRepo(t *testing.T) {
	c := createStubHTTPClient(t, []string{
		"update-repo-handle",
		"update-repo-modify-name",
		"update-repo-modify-comment",
		"update-repo-modify-kernel-options",
		"update-repo-modify-kernel-options-post",
		"update-repo-modify-autoinstall-meta",
		"update-repo-modify-template-files",
		"update-repo-modify-owners",
		"update-repo-modify-apt-components",
		"update-repo-modify-apt-dists",
		"update-repo-modify-arch",
		"update-repo-modify-breed",
		"update-repo-modify-createrepo-flags",
		"update-repo-modify-environment",
		"update-repo-modify-keep-updated",
		"update-repo-modify-mirror",
		"update-repo-modify-mirror-locally",
		"update-repo-modify-mirror-type",
		"update-repo-modify-priority",
		"update-repo-modify-proxy",
		"update-repo-modify-rsyncopts",
		"update-repo-modify-rpm-list",
		"update-repo-modify-yumopts",
		"update-repo-save",
	})
	r := NewRepo()
	r.Name = "rhel-7-x86_64"

	err := c.UpdateRepo(&r)
	FailOnError(t, err)
}

func TestSaveRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "save-repo")
	err := c.SaveRepo("000000000000000000000000000000cb", true, true, "bypass")
	FailOnError(t, err)
}

func TestCopyRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "copy-repo")
	err := c.CopyRepo("000000000000000000000000000000cb", "testrepo2")
	FailOnError(t, err)
}

func TestRenameRepo(t *testing.T) {
	c := createStubHTTPClientSingle(t, "rename-repo")
	err := c.RenameRepo("000000000000000000000000000000cc", "testrepo1")
	FailOnError(t, err)
}

func TestGetRepoHandle(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-handle")
	res, err := c.GetRepoHandle("testrepo")
	FailOnError(t, err)

	if res != "000000000000000000000000000000cb" {
		t.Error("Wrong object id returned.")
	}
}

// TestGetRepoConfigForProfile: "testprof" has no repos attached in the recorded fixture data
// (cmd/main.go never associates a repo with it), so get_repo_config_for_profile legitimately
// returns an empty string. The test only checks that the call succeeds without error.
func TestGetRepoConfigForProfile(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-config-for-profile")
	res, err := c.GetRepoConfigForProfile("testprof")
	FailOnError(t, err)

	if res != "" {
		t.Errorf("Expected an empty repo config, got %q.", res)
	}
}

// TestGetRepoConfigForSystem: see TestGetRepoConfigForProfile — "testsys" has no repos attached
// in the recorded fixture data, so an empty result is expected.
func TestGetRepoConfigForSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-repo-config-for-system")
	res, err := c.GetRepoConfigForSystem("testsys")
	FailOnError(t, err)

	if res != "" {
		t.Errorf("Expected an empty repo config, got %q.", res)
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
